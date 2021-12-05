package handlers

import (
	"log"
	"net/http"

	"github.com/reaper47/recipya/internal/templates"
)

func showErrorPage(w http.ResponseWriter, message string, err error) {
	log.Println(message, err)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusInternalServerError)
	templates.Render(w, "error-500", message)
}
