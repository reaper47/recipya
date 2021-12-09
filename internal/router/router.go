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

	r.HandleFunc("/favicon.ico", handlers.Favicon).Methods(GET)
	r.HandleFunc("/robots.txt", handlers.Robots).Methods(GET)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.FS(static.FS))))
	r.PathPrefix("/data/img/").Handler(http.StripPrefix("/data/img/", http.FileServer(http.Dir("data/img"))))

	r.HandleFunc("/auth/register", handlers.Register).Methods(GET, POST)
	r.HandleFunc("/auth/signin", handlers.SignIn).Methods(GET, POST)
	r.HandleFunc("/auth/signout", handlers.SignOut).Methods(POST)

	r.HandleFunc("/", handlers.Index).Methods(GET)
	r.HandleFunc("/recipes", amw.Middleware(handlers.Recipes)).Methods(GET)
	r.HandleFunc("/recipes/{id:[0-9]+}", handlers.Recipe).Methods(GET)
	r.HandleFunc("/recipes/{id:[0-9]+}", amw.Middleware(handlers.Recipe)).Methods(DELETE)
	r.HandleFunc("/recipes/{id:[0-9]+}/edit", amw.Middleware(handlers.EditRecipe)).Methods(GET, POST)

	r.HandleFunc("/recipes/new", amw.Middleware(handlers.RecipesAdd)).Methods(GET)
	r.HandleFunc("/recipes/new/manual", amw.Middleware(handlers.GetRecipesNewManual)).Methods(GET)
	r.HandleFunc("/recipes/new/manual", amw.Middleware(handlers.PostRecipesNewManual)).Methods(POST)

	return r
}
