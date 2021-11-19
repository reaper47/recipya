package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/reaper47/recipya/internal/templates"
)

// Index handles the "/" page.
func Index(wr http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	templates.ExecuteTemplate(wr, "index.gohtml", nil)
}
