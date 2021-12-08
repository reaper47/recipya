package repository

import (
	"net/http"
	"strings"

	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
)

var Users = map[string]models.User{}
var Sessions = map[string]string{}

// IsAuthenticated verifies whether the user is authenticated.
func IsAuthenticated(w http.ResponseWriter, req *http.Request) bool {
	c, err := req.Cookie("session")
	if err != nil {
		c = &http.Cookie{Name: "session"}
	}

	sid, err := auth.ParseToken(c.Value)
	if err != nil && !strings.HasSuffix(err.Error(), "token contains an invalid number of segments") {
		return false
	}

	if sid == "" {
		return false
	}

	_, found := Sessions[sid]
	return found
}
