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
func (u *customClaims) IsValid() bool {
	return u.VerifyExpiresAt(time.Now(), true) && u.SID != ""
}
