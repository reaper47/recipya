package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/reaper47/recipya/internal/router/handlers"
)

// New creates a new, fully-configured router.
func New() *httprouter.Router {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("static"))

	router.GET("/", handlers.Index)

	return router
}
