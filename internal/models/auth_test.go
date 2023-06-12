package models_test

import (
	"github.com/reaper47/recipya/internal/models"
	"testing"
	"time"
)

func TestAuthToken_IsExpired(t *testing.T) {
	t.Run("is expired", func(t *testing.T) {
		token := models.NewAuthToken(1, "", "", time.Now().Add(-1*time.Second).Unix(), 1)

		if !token.IsExpired() {
			t.Fail()
		}
	})

	t.Run("is not expired", func(t *testing.T) {
		token := models.NewAuthToken(1, "", "", time.Now().Add(1*time.Hour).Unix(), 1)

		if token.IsExpired() {
			t.Fail()
		}
	})
}
