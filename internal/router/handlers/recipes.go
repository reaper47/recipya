package handlers

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/reaper47/recipya/internal/templates"
)

// RecipesAdd handles the "/recipes/new" page.
func RecipesAdd(wr http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := templates.Render(wr, "recipes-new.gohtml", nil)
	if err != nil {
		log.Println(err)
	}
}

// RecipesAddManual handles the "/recipes/new/manual" page.
func RecipesAddManual(wr http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	err := templates.Render(wr, "recipes-new-manual.gohtml", nil)
	if err != nil {
		log.Println(err)
	}
}
