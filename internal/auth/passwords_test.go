package auth

import "testing"

func FuzzAuthPasswords(f *testing.F) {
	testcases := []string{"password", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		hashed, err := HashPassword(orig)
		if err != nil {
			t.Errorf("Password cannot be hashed: %q", orig)
		}

		if string(hashed) == orig {
			t.Errorf("hashed password must be equal to the plain text password: %q", orig)
		}

		err = ComparePassword(orig, hashed)
		if err != nil {
			t.Errorf("passwords do not match: %s", err)
		}
	})
}
