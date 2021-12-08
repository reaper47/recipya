package auth

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var key = []byte(os.Getenv("key"))

// ExpiresAt is the time at which a token or cookie expires when the session is not desired.
var ExpiresAt = time.Now().Add(14 * 24 * time.Hour)

// CreateToken creates a JWT token.
func CreateToken(sid string) (string, error) {
	cc := customClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: ExpiresAt.Unix(),
		},
		SID: sid,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS512, cc)
	st, err := t.SignedString(key)
	if err != nil {
		return "", err
	}
	return st, nil
}

// ParseToken parses the JWT token.
func ParseToken(st string) (string, error) {
	t, err := jwt.ParseWithClaims(st, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, errors.New("invalid signing algorithm")
		}
		return key, nil
	})
	if err != nil {
		return "", err
	}

	if !t.Valid {
		return "", err
	}

	return t.Claims.(*customClaims).SID, nil
}
