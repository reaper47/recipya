package handlers

import (
	"log"
	"net/http"

	"github.com/reaper47/recipya/internal/templates"
)

func showErrorPage(wr http.ResponseWriter, message string, err error) {
	log.Println(message, err)
	wr.Header().Set("Content-Type", "text/html")
	wr.WriteHeader(http.StatusInternalServerError)
	templates.Render(wr, "error-500", message)
}
