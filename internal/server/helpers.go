package server

import (
	"encoding/json"
	"github.com/reaper47/recipya/internal/models"
	"net/http"
	"strconv"
)

type toast string

const (
	infoToast    toast = "showInfoToast"
	warningToast toast = "showWarningToast"
	errorToast   toast = "showErrorToast"
)

func makeToast(message string, toastType toast) string {
	var backgroundColor string
	switch toastType {
	case infoToast:
		backgroundColor = "alert-info"
	case warningToast:
		backgroundColor = "alert-warning"
	case errorToast:
		backgroundColor = "alert-error"
	default:
	}

	xb, _ := json.Marshal(map[string]string{
		"showToast": `{"message":"` + message + `","backgroundColor":"` + backgroundColor + `"}`,
	})
	return string(xb)
}

func isAuthenticated(r *http.Request, getAuthToken func(selector string, validator string) (models.AuthToken, error)) bool {
	return getUserIDFromSessionCookie(r) != -1 || getUserIDFromRememberMeCookie(r, getAuthToken) != -1
}

func getUserID(r *http.Request) int64 {
	return r.Context().Value(UserIDKey).(int64)
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

func parsePathPositiveID(value string) (int64, error) {
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, err
	}
	return id, nil
}
