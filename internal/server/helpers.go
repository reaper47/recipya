package server

import (
	"errors"
	"github.com/reaper47/recipya/internal/models"
	"net/http"
	"strconv"
)

func isAuthenticated(r *http.Request, getAuthToken func(selector string, validator string) (models.AuthToken, error)) bool {
	return getUserIDFromSessionCookie(r) != -1 || getUserIDFromRememberMeCookie(r, getAuthToken) != -1
}

func (s *Server) findUserID(r *http.Request) (int64, bool) {
	isLoggedIn := true
	userID := getUserIDFromSessionCookie(r)
	if userID == -1 {
		userID = getUserIDFromRememberMeCookie(r, s.Repository.GetAuthToken)
		if userID == -1 {
			isLoggedIn = false
		}
	}
	return userID, isLoggedIn
}

func getUserID(r *http.Request) int64 {
	return r.Context().Value(UserIDKey).(int64)
}

func parsePathPositiveID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}

	if id <= 0 {
		return 0, errors.New("value must be > 0")
	}

	return id, nil
}
