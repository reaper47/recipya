package auth_test

import (
	"github.com/reaper47/recipya/internal/auth"
	"testing"
	"time"
)

func TestAuthTokens(t *testing.T) {
	t.Run("token is valid", func(t *testing.T) {
		userID := int64(1)
		claims := map[string]any{"userID": userID}
		token, err := auth.CreateToken(claims, 30*time.Minute)
		if err != nil {
			t.Fatal(err)
		}

		actual, err := auth.ParseToken(token)
		if err != nil {
			t.Fatal(err)
		}

		if actual != userID {
			t.Fatalf("wanted userID %d but got %d", userID, actual)
		}
	})

	t.Run("token is expired", func(t *testing.T) {
		claims := map[string]any{"userID": 1}
		jwtAuth, _ := auth.CreateToken(claims, -1*time.Second)

		_, err := auth.ParseToken(jwtAuth)
		if err == nil {
			t.Error(err)
		}
	})
}
