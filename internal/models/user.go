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
func (m User) GetInitials() string {
	var initials string
	if m.Username != "" {
		initials = string(m.Username[0])
	} else {
		initials = string(m.Email[0])
	}
	return strings.ToUpper(initials)
}

// Populate populates the fields of a User.
func (m *User) Populate(id int64, username, email, hash string) {
	m.ID = id
	m.Username = username
	m.Email = email
	m.HashedPassword = hash
}
