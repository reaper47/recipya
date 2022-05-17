package handlers

import (
	"errors"
	"net/http"

	"github.com/reaper47/recipya/internal/logger"
	"github.com/reaper47/recipya/internal/templates"
)

func showErrorPage(w http.ResponseWriter, message string, err error) {
	if err == nil {
		err = errors.New("")
	}

	logger.Sanitize(message, err.Error())
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusInternalServerError)
	_ = templates.Render(w, "error-500", message)
}
