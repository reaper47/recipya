package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/reaper47/recipya/internal/templates"
)

// RecipesAdd handles the GET /recipes/new URI.
func RecipesAdd(wr http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := templates.Render(wr, "recipes-new.gohtml", nil)
	if err != nil {
		log.Println(err)
	}
}

// GetRecipesNewManual handles the GET /recipes/new/manual URI.
func GetRecipesNewManual(wr http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := templates.Render(wr, "recipes-new-manual.gohtml", nil)
	if err != nil {
		log.Println(err)
	}
}

// PostRecipesNewManual handles the POST /recipes/new/manual URI.
func PostRecipesNewManual(wr http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	req.ParseForm()
	fmt.Println(req.Form)
}
