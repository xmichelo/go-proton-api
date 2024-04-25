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

// CommitRevisionReq holds the request body for a revision commit request.
type CommitRevisionReq struct {
	ManifestSignature string  // Signature of the manifest.
	SignatureAddress  string  // Address used to sign the manifest.
	BlockNumber       *int    // The index of the last block to keep when creating a revision while preserving partial content from a previous revision.
	XAttr             *string // File extended attributes encrypted with the link key.
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

// XAttr holds the extended attributes for a file revision.
type XAttr struct {
	Common    *XAttrCommon    // Common attributes.
	Location  *XAttrLocation  // Location attributes.
	Camera    *XAttrCamera    // Camera attributes.
	Media     *XAttrMedia     // Media attributes.
	IOSPhotos *XAttrIOSPhotos `json:"iOS.photos"` // iOS photos attributes.
}

// XAttrCommon contains the common attributes for file revisions.
type XAttrCommon struct {
	ModificationTime string       // UTC time in ISO 8601 format.
	Size             *int64       // Size in bytes of the unencrypted content.
	BlockSizes       []int64      // array containing the size of each unencrypted block.
	Digest           *XAttrDigest // The digests.
}

// XAttrDigest contain the digests for a file revision.
type XAttrDigest struct {
	SHA1 string // SHA1 hash, in lower-case hex format.
}

// XAttrLocation contains the location attributes for a file revision.
type XAttrLocation struct {
	Latitude  float64
	Longitude float64
}

// XAttrCamera contains the camera-related attributes for a file revision.
type XAttrCamera struct {
	CaptureTime        *string                  // UTC time in ISO 8601 format.
	Device             *string                  // The name of the camera device.
	Orientation        *int                     // EXIF orientation index in the range [1-8].
	SubjectCoordinates *XAttrSubjectCoordinates // The subject coordinates.
}

// XAttrSubjectCoordinates holds a photo's EXIF subject coordinates.
type XAttrSubjectCoordinates struct {
	Top    int // Top coordinate.
	Left   int // Left coordinate.
	Bottom int // Bottom coordinate.
	Right  int // Right coordinate.
}

// XAttrMedia contains the media-related attributes.
type XAttrMedia struct {
	Width    *int // Photo or video width.
	Height   *int // Photo or video height.
	Duration *int // Media duration.
}

type XAttrIOSPhotos struct {
	ICloudID         string // iCloud library ID.
	ModificationTime string // UTC time in ISO 8601 format.
}
