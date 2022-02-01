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
	DELETE := http.MethodDelete

	r := mux.NewRouter()

	amw := authenticationMiddleware{}

	r.HandleFunc("/", handlers.Index).Methods(GET)
	r.HandleFunc("/favicon.ico", handlers.Favicon).Methods(GET)
	r.HandleFunc("/robots.txt", handlers.Robots).Methods(GET)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.FS(static.FS))))
	r.PathPrefix("/data/img/").Handler(http.StripPrefix("/data/img/", http.FileServer(http.Dir("data/img"))))

	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", handlers.Register).Methods(GET, POST)
	auth.HandleFunc("/signin", handlers.SignIn).Methods(GET, POST)
	auth.HandleFunc("/signout", handlers.SignOut).Methods(POST)

	recipes := r.PathPrefix("/recipes").Subrouter()
	recipes.HandleFunc("", amw.Middleware(handlers.Recipes)).Methods(GET)
	recipes.HandleFunc("/{id:[0-9]+}", handlers.Recipe).Methods(GET)
	recipes.HandleFunc("/{id:[0-9]+}", amw.Middleware(handlers.Recipe)).Methods(DELETE)
	recipes.HandleFunc("/{id:[0-9]+}/edit", amw.Middleware(handlers.EditRecipe)).Methods(GET, POST)
	recipes.HandleFunc("/new", amw.Middleware(handlers.RecipesAdd)).Methods(GET)
	recipes.HandleFunc("/new/manual", amw.Middleware(handlers.GetRecipesNewManual)).Methods(GET)
	recipes.HandleFunc("/new/manual", amw.Middleware(handlers.PostRecipesNewManual)).Methods(POST)
	recipes.HandleFunc("/categories", amw.Middleware(handlers.Categories)).Methods(POST)
	recipes.HandleFunc("/import", amw.Middleware(handlers.ImportRecipes)).Methods(POST)

	return r
}
