package proton

type CreateFileReq struct {
	ParentLinkID string

	Name     string // The file name. Encrypted using the parent folder node key, signed using the address key.
	Hash     string // The file name hash as a hex string, hashed using the parent hash key.
	MIMEType string // The MIME Type of the file.

	ContentKeyPacket          string // The content key packet, encrypted with the node key.
	ContentKeyPacketSignature string // Unencrypted signature of the content key, signed with the NodeKey

	NodeKey                 string // The private NodeKey, used to decrypt any file/folder content.
	NodePassphrase          string // The passphrase used to unlock the NodeKey, encrypted by the owning Link/Share keyring.
	NodePassphraseSignature string // The signature of the NodePassphrase

	SignatureAddress string // Signature email address used to sign passphrase and name

	ClientUID string // The client unique ID.
}

type CreateFileRes struct {
	ID         string // Encrypted Link ID
	RevisionID string // Encrypted Revision ID
}

type UpdateRevisionReq struct {
	BlockList         []BlockToken
	State             RevisionState
	ManifestSignature string
	SignatureAddress  string
}

type BlockToken struct {
	Index int
	Token string
}

// ConflictErrorResponse holds the fields in the API error details when a conflict occurs (code proton.AlreadyExists).
type ConflictErrorResponse struct {
	ConflictLinkID          string
	ConflictRevisionID      string
	ConflictDraftRevisionID string
	ConflictDraftClientUID  string
}
