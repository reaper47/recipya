package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/reaper47/recipya/internal/router/handlers"
	"github.com/reaper47/recipya/static"
)

// New creates a new, fully-configured router.
func New() *mux.Router {
	GET := http.MethodGet
	POST := http.MethodPost
	PATCH := http.MethodPatch
	DELETE := http.MethodDelete

	r := mux.NewRouter()

	r.HandleFunc("/favicon.ico", handlers.Favicon).Methods(GET)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.FS(static.FS))))

	r.HandleFunc("/", handlers.Index).Methods(GET)
	r.HandleFunc("/recipes", handlers.Index).Methods(GET)
	r.HandleFunc("/recipes/{id:[0-9]+}", handlers.Recipe).Methods(GET, DELETE)
	r.HandleFunc("/recipes/{id:[0-9]+}/edit", handlers.EditRecipe).Methods(GET, PATCH)

	r.HandleFunc("/recipes/new", handlers.RecipesAdd).Methods(GET)
	r.HandleFunc("/recipes/new/manual", handlers.GetRecipesNewManual).Methods(GET)
	r.HandleFunc("/recipes/new/manual", handlers.PostRecipesNewManual).Methods(POST)

	return r
}
