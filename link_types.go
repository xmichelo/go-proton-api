package proton

import (
	"encoding/base64"
	"errors"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

type LinkWalkFunc func([]string, Link, *crypto.KeyRing) error

// Link holds the tree structure, for the clients, they represent the files and folders of a given volume.
// They have a ParentLinkID that points to parent folders.
// Links also hold the file name (encrypted) and a hash of the name for name collisions.
// Link data is encrypted with its owning Share keyring.
type Link struct {
	LinkID       string // Encrypted file/folder ID
	ParentLinkID string // Encrypted parent folder ID (LinkID). Root link has null ParentLinkID.

	Type               LinkType
	Name               string // Encrypted file name
	NameSignatureEmail string // Link name signature email
	Hash               string // HMAC of name encrypted with parent hash key
	State              LinkState
	TotalSize          int64 // Encrypted size of Node (all active and obsolete revisions for files)

	MIMEType string

	CreateTime              int64  // Link creation time
	ModifyTime              int64  // Link modification time (on API, real modify date is stored in XAttr)
	Trashed                 *int64 // Time at which the file was trashed, null if file is not trashed.
	NodeKey                 string // The private NodeKey, used to decrypt any file/folder content.
	NodePassphrase          string // The passphrase used to unlock the NodeKey, encrypted by the owning Link/Share keyring.
	NodePassphraseSignature string // Node passphrase signature

	Attributes       int64   // ?
	XAttr            *string // Extended attributes encrypted with link key
	Permissions      int64   // ?
	FileProperties   *FileProperties
	FolderProperties *FolderProperties

	SignatureEmail string // Signature email address used for passphrase, should be the user's address associated with the Share.

}

type LinkState int

const (
	LinkStateDraft LinkState = iota
	LinkStateActive
	LinkStateTrashed
	LinkStateDeleted
	LinkStateRestoring
)

func (l LinkState) String() string {
	switch l {
	case LinkStateDraft:
		return "draft"
	case LinkStateActive:
		return "active"
	case LinkStateTrashed:
		return "trashed"
	case LinkStateDeleted:
		return "deleted"
	case LinkStateRestoring:
		return "restoring"
	default:
		return "unknown"
	}
}

func (l Link) GetName(parentNodeKR, addrKR *crypto.KeyRing) (string, error) {
	encName, err := crypto.NewPGPMessageFromArmored(l.Name)
	if err != nil {
		return "", err
	}

	decName, err := parentNodeKR.Decrypt(encName, addrKR, crypto.GetUnixTime())
	if err != nil {
		return "", err
	}

	return decName.GetString(), nil
}

func (l Link) GetKeyRing(parentNodeKR, addrKR *crypto.KeyRing) (*crypto.KeyRing, error) {
	enc, err := crypto.NewPGPMessageFromArmored(l.NodePassphrase)
	if err != nil {
		return nil, err
	}

	dec, err := parentNodeKR.Decrypt(enc, nil, crypto.GetUnixTime())
	if err != nil {
		return nil, err
	}

	sig, err := crypto.NewPGPSignatureFromArmored(l.NodePassphraseSignature)
	if err != nil {
		return nil, err
	}

	if err := addrKR.VerifyDetached(dec, sig, crypto.GetUnixTime()); err != nil {
		return nil, err
	}

	lockedKey, err := crypto.NewKeyFromArmored(l.NodeKey)
	if err != nil {
		return nil, err
	}

	unlockedKey, err := lockedKey.Unlock(dec.GetBinary())
	if err != nil {
		return nil, err
	}

	return crypto.NewKeyRing(unlockedKey)
}

func (l Link) GetHashKey(nodeKR, addrKR *crypto.KeyRing) ([]byte, error) {
	if l.Type != LinkTypeFolder {
		return nil, errors.New("link is not a folder")
	}

	enc, err := crypto.NewPGPMessageFromArmored(l.FolderProperties.NodeHashKey)
	if err != nil {
		return nil, err
	}

	dec, err := nodeKR.Decrypt(enc, nodeKR, crypto.GetUnixTime())
	if err != nil {
		var sigError crypto.SignatureVerificationError
		if errors.As(err, &sigError) {
			// nodeHashKey is supposed to be signed with the node key ring, however some legacy applications
			// signed it with the share address key.
			if addrKR == nil {
				return nil, err
			}

			dec, err = nodeKR.Decrypt(enc, addrKR, crypto.GetUnixTime())
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return dec.GetBinary(), nil
}

func (l Link) GetSessionKey(nodeKR *crypto.KeyRing) (*crypto.SessionKey, error) {
	if l.Type != LinkTypeFile {
		return nil, errors.New("link is not a file")
	}

	dec, err := base64.StdEncoding.DecodeString(l.FileProperties.ContentKeyPacket)
	if err != nil {
		return nil, err
	}

	key, err := nodeKR.DecryptSessionKey(dec)
	if err != nil {
		return nil, err
	}

	sig, err := crypto.NewPGPSignatureFromArmored(l.FileProperties.ContentKeyPacketSignature)
	if err != nil {
		return nil, err
	}

	if err := nodeKR.VerifyDetached(crypto.NewPlainMessage(key.Key), sig, crypto.GetUnixTime()); err != nil {
		return nil, err
	}

	return key, nil
}

type FileProperties struct {
	ContentKeyPacket          string           // The block's key packet, encrypted with the node key.
	ContentKeyPacketSignature string           // Signature of the content key packet. Signature of the session key, signed with the NodeKey.
	ActiveRevision            RevisionMetadata // The active revision of the file.
}

type FolderProperties struct {
	NodeHashKey string // HMAC key used to hash the folder's children names.
}

type LinkType int

const (
	LinkTypeFolder LinkType = iota + 1
	LinkTypeFile
)

func (t LinkType) String() string {
	switch t {
	case LinkTypeFolder:
		return "folder"
	case LinkTypeFile:
		return "file"
	default:
		return "unknown"
	}
}

type RevisionMetadata struct {
	ID                string        // Encrypted Revision ID
	CreateTime        int64         // Unix timestamp of the revision creation time
	Size              int64         // Size of the revision in bytes
	ManifestSignature string        // Signature of the revision manifest, signed with user's address key of the share.
	SignatureEmail    string        // Email of the user that signed the revision.
	State             RevisionState // State of revision
	Thumbnail         Bool          // Whether the revision has a thumbnail
	ThumbnailHash     string        // Hash of the thumbnail
}

// Revision Revisions are only for files, they represent “versions” of files.
// Each file can have 1 active revision and n obsolete revisions.
type Revision struct {
	RevisionMetadata

	Blocks []Block
}

type RevisionState int

const (
	RevisionStateDraft RevisionState = iota
	RevisionStateActive
	RevisionStateObsolete
	RevisionStateDeleted
)

func (r RevisionState) String() string {
	switch r {
	case RevisionStateDraft:
		return "draft"
	case RevisionStateActive:
		return "active"
	case RevisionStateObsolete:
		return "obsolete"
	case RevisionStateDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}
