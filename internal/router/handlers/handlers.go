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
	_, isAuthenticated := repository.IsAuthenticated(w, req)
	if isAuthenticated {
		http.Redirect(w, req, "/recipes", http.StatusSeeOther)
	} else {
		handleIndexUnauthenticated(w, req)
	}
}

func handleIndexUnauthenticated(w http.ResponseWriter, req *http.Request) {
	err := templates.Render(w, "landing.gohtml", templates.Data{
		HideSidebar: true,
		HeaderData: templates.HeaderData{
			IsUnauthenticated: true,
		},
	})
	if err != nil {
		log.Println(err)
	}
}

// Recipes handles the GET /recipes page.
func Recipes(w http.ResponseWriter, req *http.Request) {
	count, err := config.App().Repo.RecipesCount()
	if err != nil {
		showErrorPage(w, "Could not retrieve total number of recipes", err)
		return
	}
	pg := templates.Pagination{
		NumResults: count,
		NumPages:   count / 12,
	}

	qpage := req.URL.Query().Get("page")
	var page int
	page, err = strconv.Atoi(qpage)
	if err != nil || page <= 0 {
		page = 1
	} else if page > pg.NumPages+1 {
		page = pg.NumPages + 1
	}

	s := getSession(req)
	recipes, err := config.App().Repo.Recipes(s.UserID, page)
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
		HeaderData: templates.HeaderData{
			AvatarInitials: s.UserInitials,
		},
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

// Settings serves the GET /settings endpoint.
func Settings(w http.ResponseWriter, req *http.Request) {
	xc := config.App().Repo.Categories(getSession(req).UserID)
	for i, c := range xc {
		if c == "uncategorized" {
			xc[i] = xc[len(xc)-1]
			xc = xc[:len(xc)-1]
			break
		}
	}

	err := templates.Render(w, "settings.gohtml", templates.Data{
		HeaderData: templates.HeaderData{
			AvatarInitials: getSession(req).UserInitials,
		},
		HideSidebar: true,
		HideGap:     true,
		Categories:  xc,
	})
	if err != nil {
		log.Println(err)
	}
}
