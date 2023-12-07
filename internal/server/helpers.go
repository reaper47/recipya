package server

import (
	"encoding/json"
	"github.com/reaper47/recipya/internal/models"
	"net/http"
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
		backgroundColor = "bg-blue-500"
	case warningToast:
		backgroundColor = "bg-orange-500"
	case errorToast:
		backgroundColor = "bg-red-500"
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
	var userID int64
	isLoggedIn := true
	userID = getUserIDFromSessionCookie(r)
	if userID == -1 {
		userID = getUserIDFromRememberMeCookie(r, s.Repository.GetAuthToken)
		if userID == -1 {
			isLoggedIn = false
		}
	}
	return userID, isLoggedIn
}
