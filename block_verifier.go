package proton

import (
	"encoding/base64"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

type BlockVerifier struct {
	verificationCode []byte
	sessionKey       *crypto.SessionKey
}

// NewBlockVerifier return a new block verifier for file upload.
func NewBlockVerifier(data VerificationData, kr *crypto.KeyRing) (BlockVerifier, error) {
	code, err := base64.StdEncoding.DecodeString(data.VerificationCode)
	if err != nil {
		return BlockVerifier{}, err
	}

	keyPacket, err := base64.StdEncoding.DecodeString(data.ContentKeyPacket)
	if err != nil {
		return BlockVerifier{}, err
	}

	sessionKey, err := kr.DecryptSessionKey(keyPacket)
	if err != nil {
		return BlockVerifier{}, err
	}

	return BlockVerifier{
		verificationCode: code,
		sessionKey:       sessionKey,
	}, nil
}

// GetVerificationToken return the verification token for an encrypted data packet.
func (v BlockVerifier) GetVerificationToken(data []byte) (string, error) {
	// Requirement: check that we can decrypt the data packet
	// Note: we can optimize this as the spec only require use to decrypt the 16 first bytes of the message, and check against original message,
	// as this method does not include the checksum verification performed by the full decryption.
	_, err := v.sessionKey.Decrypt(data)
	if err != nil {
		return "", err
	}

	// zero-pad the data packet if it is less than 32 bytes in size.
	if len(data) < 32 {
		pad := make([]byte, 32-len(data))
		data = append(data, pad...)
	}

	// compute the token (xor of the verification code and the 32 first bytes of the data packet).
	result := make([]byte, 32)
	for i := 0; i < 32; i++ {
		result[i] = v.verificationCode[i] ^ data[i]
	}

	return base64.StdEncoding.EncodeToString(result), nil
}
