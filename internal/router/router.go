package router

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/reaper47/recipya/internal/router/handlers"
	"github.com/reaper47/recipya/static"
)

// New creates a new, fully-configured router.
func New() *httprouter.Router {
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.FS(static.FS))

	router.GET("/", handlers.Index)

	return router
}
