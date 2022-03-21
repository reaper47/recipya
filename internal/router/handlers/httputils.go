package handlers

import (
	"fmt"
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

func writeJson(w http.ResponseWriter, message string, code int) {
	j, err := models.NewErrorJSON(http.StatusBadRequest, message)
	if err != nil {
		fmt.Fprintf(w, constants.ErrDecodingJSON+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(constants.HeaderContentType, constants.ApplicationJSON)
	w.WriteHeader(code)
	w.Write(j)
}
