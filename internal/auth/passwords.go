package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword hashes a plain text password using bcrypt.
func HashPassword(password string) ([]byte, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// ComparePassword verifies whether the plan
// text password and  the hashed password are equal.
func ComparePassword(password string, hashedPassword []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
