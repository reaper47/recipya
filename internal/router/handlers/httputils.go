package handlers

import (
	"net/http"

	"github.com/reaper47/recipya/internal/constants"
	"github.com/reaper47/recipya/internal/models"
)

func getSession(req *http.Request) models.Session {
	var s models.Session
	if m := req.Context().Value(constants.UserID); m != nil {
		if value, ok := m.(models.Session); ok {
			s = value
		}
	}
	return s
}
