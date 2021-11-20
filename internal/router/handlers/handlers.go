package handlers

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/reaper47/recipya/internal/templates"
)

// Index handles the "/" page.
func Index(wr http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := templates.Render(wr, "index.gohtml", nil)
	if err != nil {
		log.Println(err)
	}
}
