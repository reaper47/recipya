package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/reaper47/recipya/internal/config"
	"github.com/reaper47/recipya/internal/repository"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/static"
)

// Index handles the GET / page.
func Index(w http.ResponseWriter, req *http.Request) {
	if repository.IsAuthenticated(w, req) {
		handleIndexAuthenticated(w, req)
	} else {
		handleIndexUnauthenticated(w, req)
	}
}

func handleIndexAuthenticated(w http.ResponseWriter, req *http.Request) {
	recipesCount, err := config.App().Repo.GetRecipesCount()
	if err != nil {
		showErrorPage(w, "Could not retrieve total number of recipes", err)
		return
	}
	pg := templates.Pagination{
		NumResults: recipesCount,
		NumPages:   recipesCount / 12,
	}

	qpage := req.URL.Query().Get("page")
	var page int
	page, err = strconv.Atoi(qpage)
	if err != nil || page <= 0 {
		page = 1
	} else if page > pg.NumPages+1 {
		page = pg.NumPages + 1
	}

	recipes, err := config.App().Repo.GetAllRecipes(page)
	if err != nil {
		showErrorPage(w, "Cannot retrieve all recipes.", err)
		return
	}

	pg.Init(page)

	err = templates.Render(w, "index.gohtml", templates.Data{
		RecipesData: templates.RecipesData{
			Recipes:    recipes,
			Pagination: pg,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

func handleIndexUnauthenticated(w http.ResponseWriter, req *http.Request) {
	err := templates.Render(w, "landing.gohtml", templates.Data{
		HideSidebar:       true,
		IsUnauthenticated: true,
	})
	if err != nil {
		log.Println(err)
	}
}

// Favicon serves the favicon.ico file.
func Favicon(w http.ResponseWriter, req *http.Request) {
	serveFile(w, "favicon.ico", "image/x-icon")
}

// Robots serves the robots.txt file.
func Robots(w http.ResponseWriter, req *http.Request) {
	serveFile(w, "robots.txt", "text/plain")
}

func serveFile(w http.ResponseWriter, fname, contentType string) {
	f, err := static.FS.ReadFile(fname)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", contentType)
	w.Write(f)
}
