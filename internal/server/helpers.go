package server

import (
	"encoding/json"
	"fmt"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
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

func sendErrorAdminEmail(sendFunc func(to string, template templates.EmailTemplate, data any), errFuncName string, err error) {
	sendFunc(app.Config.Email.From, templates.EmailErrorAdmin, templates.EmailData{
		Text: fmt.Sprintf("error in %s: %q", errFuncName, err),
	})
}
