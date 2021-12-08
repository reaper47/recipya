package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type customClaims struct {
	jwt.StandardClaims
	SID string
}

func (u *customClaims) IsValid() bool {
	return u.VerifyExpiresAt(time.Now().Unix(), true) && u.SID != ""
}
