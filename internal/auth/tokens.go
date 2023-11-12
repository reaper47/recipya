package auth

import (
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"time"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	secretUUID, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	tokenAuth = jwtauth.New("HS512", []byte(secretUUID.String()), nil)
}

// CreateToken creates a JWT token.
func CreateToken(claims map[string]any, expiry time.Duration) (string, error) {
	jwtauth.SetIssuedNow(claims)
	jwtauth.SetExpiry(claims, time.Now().Add(expiry))

	_, token, err := tokenAuth.Encode(claims)
	return token, err
}

// ParseToken parses the JWT token for the userID.
func ParseToken(token string) (int64, error) {
	jwtToken, err := jwtauth.VerifyToken(tokenAuth, token)
	if err != nil {
		return -1, err
	}

	userIDAny, ok := jwtToken.Get("userID")
	if !ok {
		return -1, errors.New("custom claims does not have the user field")
	}

	return parseUserID(userIDAny)
}

func parseUserID(userIDAny any) (int64, error) {
	var userID int64
	switch x := userIDAny.(type) {
	case int64:
		userID = x
	case float64:
		userID = int64(x)
	default:
		userID = userIDAny.(int64)
	}
	return userID, nil
}
