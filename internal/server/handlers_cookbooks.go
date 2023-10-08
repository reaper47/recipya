package server

import (
	"github.com/reaper47/recipya/internal/templates"
	"net/http"
)

func (s *Server) cookbooksHandler(w http.ResponseWriter, r *http.Request) {
	isHxRequest := r.Header.Get("Hx-Request") == "true"
	data := templates.Data{
		IsAuthenticated: true,
		IsHxRequest:     isHxRequest,
		Title:           "Cookbooks",
	}

	if isHxRequest {
		templates.RenderComponent(w, "cookbooks", "cookbooks-index", data)
	} else {
		templates.Render(w, templates.CookbooksPage, data)
	}
}
