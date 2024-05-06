package proton

// ShareURL represent a share URL object
type ShareURL struct {
	Token                    string
	ShareURLID               string
	ShareID                  string
	PublicURL                string `json:"PublicUrl"`
	ExpirationTime           int64
	LastAccessTime           int64
	CreateTime               int64
	MaxAccesses              int
	NumAccesses              int
	Name                     string // armored PGP message
	CreatorEmail             string
	Permissions              ShareURLPermissions
	Flags                    int
	UrlPasswordSalt          string // Base64 encoded.
	SharePasswordSalt        string // Base64 encoded.
	SRPVerifier              string // Base64 encoded.
	SRPModulusID             string
	Password                 string // Armored PGP message
	SharePassphraseKeyPacket string // Base64 encoded
}

// CreateShareURLReq hold the request body for a CreateShareURL POST request.
type CreateShareURLReq struct {
	CreatorEmail             string
	Permissions              ShareURLPermissions
	UrlPasswordSalt          string
	SharePasswordSalt        string
	SRPVerifier              string
	SRPModulusID             string
	Flags                    ShareURLFlags
	SharePassphraseKeyPacket string
	Password                 string
	MaxAccesses              int
	ExpirationTime           *int64
	ExpirationDuration       *int64
	Name                     *string
}

type ShareURLPermissions int

const (
	ShareURLPermissionWrite ShareURLPermissions = 1 << 1 // Note: not yet supported.
	ShareURLPermissionRead  ShareURLPermissions = 1 << 2
)

type ShareURLFlags int

const (
	ShareURLFlagLegacyRandomPassword ShareURLFlags = iota
	ShareURLFlagLegacyCustomPassword
	ShareURLFlagRandomPassword
	ShareURLFlagCustomPassword
)
