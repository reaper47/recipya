package repository

import (
	"net/http"
	"strings"

	"github.com/reaper47/recipya/internal/auth"
)

// Sessions stores the session ID associated of each authenticated user.
// They are only stored in memory; they will be wiped when the server is closed.
var Sessions = map[string]int64{}

// IsAuthenticated verifies whether the user is authenticated.
//
// It returns the ID of the user and whether he or she is authenticated.
func IsAuthenticated(w http.ResponseWriter, req *http.Request) (int64, bool) {
	c, err := req.Cookie("session")
	if err != nil {
		c = &http.Cookie{Name: "session"}
	}

	sid, err := auth.ParseToken(c.Value)
	if err != nil && !strings.HasSuffix(err.Error(), "token contains an invalid number of segments") {
		return -1, false
	}

	if sid == "" {
		return -1, false
	}

	id, found := Sessions[sid]
	return id, found
}
