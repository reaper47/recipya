package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// GenerateSelectorAndValidator creates a selector and validator pair.
// The selector is a unique ID, and the validator is a random token.
func GenerateSelectorAndValidator() (string, string) {
	selector := make([]byte, 12)
	validator := make([]byte, 32)

	_, _ = rand.Read(selector)
	_, _ = rand.Read(validator)

	return base64.URLEncoding.EncodeToString(selector), base64.URLEncoding.EncodeToString(validator)
}

// HashValidator creates a hashed validator using SHA-256.
func HashValidator(validator string) string {
	hash := sha256.Sum256([]byte(validator))
	return base64.URLEncoding.EncodeToString(hash[:])
}
