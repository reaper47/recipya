package server

import (
	"github.com/reaper47/recipya/internal/templates"
	"net/http"
)

func avatarDropdownHandler(w http.ResponseWriter, _ *http.Request) {
	templates.RenderComponent(w, "core", "avatar-menu", nil)
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	page := templates.LandingPage
	isAuth := isAuthenticated(r, s.Repository.GetAuthToken)
	if isAuth {
		page = templates.HomePage
	}

	templates.Render(w, page, templates.Data{
		IsAuthenticated: isAuth,
		Title:           page.Title(),
	})
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	templates.Render(w, templates.Simple, templates.PageNotFound)
}

func (s *Server) settingsHandler(w http.ResponseWriter, r *http.Request) {
	panic("to implement")
}

func (s *Server) userInitialsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.Write([]byte(s.Repository.UserInitials(userID.(int64))))
}
