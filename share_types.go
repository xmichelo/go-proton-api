package proton

import (
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

type ShareMetadata struct {
	ShareID  string // Encrypted share ID
	LinkID   string // Encrypted link ID to which the share points (root of share).
	VolumeID string // Encrypted volume ID on which the share is mounted

	Type  ShareType  // Type of share
	State ShareState // The state of the share (active, deleted)

	CreationTime int64 // Creation time of the share in Unix time
	ModifyTime   int64 // Last modification time of the share in Unix time

	Creator           string     // Creator email address
	Flags             ShareFlags // The flag bitmap
	Locked            bool       // Whether the share is locked
	VolumeSoftDeleted bool       // Was the volume soft deleted
}

// Share is an entry point to a location in the file structure (Volume).
// It points to a file or folder anywhere in the tree and holds a key called the ShareKey.
// To access a file or folder in Drive, a user must be a member of a share.
// The membership information is tied to a specific address, and key.
// This key then allows the user to decrypt the share key, giving access to the file system rooted at that share.
type Share struct {
	ShareMetadata

	AddressID    string // Encrypted address ID
	AddressKeyID string // Encrypted address key ID

	Key                 string // The private ShareKey, encrypted with a passphrase
	Passphrase          string // The encrypted passphrase
	PassphraseSignature string // The signature of the passphrase
}

// CreateShareReq is the request body for the share POST request used to create a share
type CreateShareReq struct {
	AddressID                string // The encrypted addressID.
	RootLinkID               string // The linkID for the share root node.
	ShareKey                 string // The armored PGP private key of the share.
	SharePassphrase          string // The armored PGP message containing the share key passphrase, encrypted with the address key.
	SharePassphraseSignature string // The armored PGP signature of the share passphrase.
	PassphraseKeyPacket      string // The key packet for the root link passphrase, in base64 encoding.
	NameKeyPacket            string // The key packet for the root link name, in base64 encoding.
}

// CreateShareRes is the response body for a share creation request
type CreateShareRes struct {
	ID string // The shareID.
}

func (s Share) GetKeyRing(addrKR *crypto.KeyRing) (*crypto.KeyRing, error) {
	enc, err := crypto.NewPGPMessageFromArmored(s.Passphrase)
	if err != nil {
		return nil, err
	}

	dec, err := addrKR.Decrypt(enc, nil, crypto.GetUnixTime())
	if err != nil {
		return nil, err
	}

	sig, err := crypto.NewPGPSignatureFromArmored(s.PassphraseSignature)
	if err != nil {
		return nil, err
	}

	if err := addrKR.VerifyDetached(dec, sig, crypto.GetUnixTime()); err != nil {
		return nil, err
	}

	lockedKey, err := crypto.NewKeyFromArmored(s.Key)
	if err != nil {
		return nil, err
	}

	unlockedKey, err := lockedKey.Unlock(dec.GetBinary())
	if err != nil {
		return nil, err
	}

	return crypto.NewKeyRing(unlockedKey)
}

type ShareType int

const (
	ShareTypeMain     ShareType = 1
	ShareTypeStandard ShareType = 2
	ShareTypeDevice   ShareType = 3
)

func (t ShareType) String() string {
	switch t {
	case ShareTypeMain:
		return "main"
	case ShareTypeStandard:
		return "standard"
	case ShareTypeDevice:
		return "device"
	default:
		return fmt.Sprintf("unknown (%v)", int(t))
	}
}

type ShareState int

const (
	ShareStateActive  ShareState = 1
	ShareStateDeleted ShareState = 2
)

func (s ShareState) String() string {
	switch s {
	case ShareStateActive:
		return "active"
	case ShareStateDeleted:
		return "deleted"
	default:
		return fmt.Sprintf("unknown (%v)", int(s))
	}
}

type ShareFlags int

const (
	NoFlags ShareFlags = iota
	PrimaryShare
)

func (f ShareFlags) String() string {
	switch f {
	case NoFlags:
		return "none"
	case PrimaryShare:
		return "primary"
	default:
		return fmt.Sprintf("unknown (%v)", int(f))
	}
}
