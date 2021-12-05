package handlers

import (
	"log"
	"net/http"

	"github.com/reaper47/recipya/internal/config"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/static"
)

// Index handles the "/" page.
func Index(w http.ResponseWriter, req *http.Request) {
	recipes, err := config.App().Repo.GetAllRecipes()
	if err != nil {
		showErrorPage(w, "Cannot retrieve all recipes.", err)
		return
	}

	err = templates.Render(w, "index.gohtml", templates.RecipesData{Recipes: recipes})
	if err != nil {
		log.Println(err)
	}
}

// Favicon serves the favicon.ico file.
func Favicon(w http.ResponseWriter, req *http.Request) {
	f, err := static.FS.ReadFile("favicon.ico")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/x-icon")
	w.Write(f)
}
