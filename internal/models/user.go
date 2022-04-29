package models

import (
	"strings"

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

// GetInitials gets the user's 1-character initials from either the email address of username.
func (u User) GetInitials() string {
	var initials string
	if u.Username != "" {
		initials = string(u.Username[0])
	} else {
		initials = string(u.Email[0])
	}
	return strings.ToUpper(initials)
}

// Populate populates the fields of a User.
func (u *User) Populate(id int64, username, email, hash string) {
	u.ID = id
	u.Username = username
	u.Email = email
	u.HashedPassword = hash
}
