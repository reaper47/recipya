package models

import (
	"github.com/reaper47/recipya/internal/auth"
)

// User holds data related to a user.
type User struct {
	ID             int64
	Username       string
	Email          string
	HashedPassword string
}

// IsPasswordOk checks whether the password is ok.
func (u User) IsPasswordOk(password string) bool {
	return auth.ComparePassword(password, []byte(u.HashedPassword)) == nil
}
