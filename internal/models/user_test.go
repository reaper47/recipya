package models

import (
	"testing"

	"github.com/reaper47/recipya/internal/auth"
)

func TestModelUser(t *testing.T) {
	u := User{
		ID:       1,
		Username: "adam jenkins",
		Email:    "user@name.com",
	}

	testcases1 := []struct {
		name string
		in   string
		want bool
	}{
		{
			name: "password is invalid",
			in:   "12345",
			want: false,
		},
		{
			name: "password is valid",
			in:   "12345",
			want: true,
		},
	}
	for i, tc := range testcases1 {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := auth.HashPassword(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			u.HashedPassword = string(hash)

			password := tc.in
			if !tc.want {
				password += "!"
			}
			if u.IsPasswordOk(password) != tc.want {
				t.Errorf("IsPasswordOk for test #%d: %#v, want %v", i, tc.in, tc.want)
			}
		})
	}

	testcases2 := []struct {
		name     string
		useEmail bool
		in       string
		want     string
	}{
		{
			name:     "get username from email",
			useEmail: true,
			in:       "macpoule@gmail.com",
			want:     "M",
		},
		{
			name:     "get initials from username",
			useEmail: false,
			in:       "adam jenkins",
			want:     "A",
		},
	}
	for i, tc := range testcases2 {
		t.Run(tc.name, func(t *testing.T) {
			if tc.useEmail {
				u.Email = tc.in
				u.Username = ""
			} else {
				u.Email = ""
				u.Username = tc.in
			}

			if u.GetInitials() != tc.want {
				t.Errorf("GetInitials for test #%d: %s, want %s", i, tc.in, tc.want)
			}
		})
	}
}
