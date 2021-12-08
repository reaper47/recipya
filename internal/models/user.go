package models

import (
	"fmt"

	"github.com/reaper47/recipya/internal/auth"
)

type User struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"password"`
}

// NewUser creates a new user.
func NewUser(username, email, password string) (User, error) {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return User{}, fmt.Errorf("could not create user: %s", err)
	}

	return User{
		Username:       username,
		Email:          email,
		HashedPassword: string(hash),
	}, nil
}

// IsPasswordOk checks whether the password is ok.
func (u User) IsPasswordOk(password string) bool {
	return auth.ComparePassword(password, []byte(u.HashedPassword)) == nil
}
