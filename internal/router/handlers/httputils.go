package handlers

import (
	"net/http"

	"github.com/reaper47/recipya/internal/constants"
)

func getUserID(req *http.Request) int64 {
	var userID int64
	if m := req.Context().Value(constants.UserID); m != nil {
		if value, ok := m.(int64); ok {
			userID = value
		}
	}
	return userID
}
