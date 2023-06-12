package models

import "time"

// NewAuthToken creates a new AuthToken.
func NewAuthToken(id int64, selector, hashValidator string, expiresSeconds, userID int64) *AuthToken {
	return &AuthToken{
		ID:            id,
		Selector:      selector,
		HashValidator: hashValidator,
		Expires:       time.Unix(expiresSeconds, 0),
		UserID:        userID,
	}
}

// AuthToken holds details on an authentication token.
type AuthToken struct {
	ID            int64
	Selector      string
	HashValidator string
	Expires       time.Time
	UserID        int64
}

// IsExpired verifies whether the auth token is expired.
func (a *AuthToken) IsExpired() bool {
	return time.Now().After(a.Expires)
}
