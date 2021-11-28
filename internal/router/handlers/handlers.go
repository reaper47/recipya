package handlers

import (
	"log"
	"net/http"

	"github.com/reaper47/recipya/internal/templates"
)

// Index handles the "/" page.
func Index(wr http.ResponseWriter, req *http.Request) {
	data := templates.IndexData{}
	err := templates.Render(wr, "index.gohtml", data)
	if err != nil {
		log.Println(err)
	}
}
