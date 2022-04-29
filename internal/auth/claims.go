package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type customClaims struct {
	jwt.RegisteredClaims
	SID string
}

// IsValid verifies whether the custom claim is valid.
func (c *customClaims) IsValid() bool {
	return c.VerifyExpiresAt(time.Now(), true) && c.SID != ""
}
