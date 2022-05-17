package repository_test

import (
	"net/http"
	"testing"

	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/constants"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/repository"
)

func TestAuth(t *testing.T) {
	t.Run("user has a session", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		sid, _ := auth.CreateToken("test")
		req.AddCookie(&http.Cookie{Name: constants.CookieSession, Value: sid})
		repository.Sessions["test"] = models.Session{}

		_, got := repository.IsAuthenticated(req)

		assertAuth(t, got, true)
	})

	t.Run("user has no session", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		_, got := repository.IsAuthenticated(req)

		assertAuth(t, got, false)
	})
}

func assertAuth(t testing.TB, got, want bool) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v but want %v", got, want)
	}
}
