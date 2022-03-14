package repository

import (
	"net/http"
	"strings"

	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
)

// Sessions stores the session ID associated of each authenticated user.
// They are only stored in memory; they will be wiped when the server is closed.
var Sessions = map[string]models.Session{}

// IsAuthenticated verifies whether the user is authenticated.
//
// It returns the user's session and whether he or she is authenticated.
func IsAuthenticated(w http.ResponseWriter, req *http.Request) (models.Session, bool) {
	c, err := req.Cookie("session")
	if err != nil {
		c = &http.Cookie{Name: "session"}
	}

	sid, err := auth.ParseToken(c.Value)
	if err != nil && !strings.HasSuffix(err.Error(), "token contains an invalid number of segments") {
		return models.Session{}, false
	}

	if sid == "" {
		return models.Session{}, false
	}

	s, found := Sessions[sid]
	return s, found
}
