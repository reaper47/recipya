package handlers

import (
	"log"
	"net/http"

	"github.com/reaper47/recipya/internal/config"
	"github.com/reaper47/recipya/internal/templates"
)

// Index handles the "/" page.
func Index(wr http.ResponseWriter, req *http.Request) {
	recipes, err := config.App().Repo.GetAllRecipes()
	if err != nil {
		showErrorPage(wr, "Cannot retrieve all recipes.", err)
		return
	}
	log.Println(len(recipes))

	data := templates.IndexData{}
	err = templates.Render(wr, "index.gohtml", data)
	if err != nil {
		log.Println(err)
	}
}
