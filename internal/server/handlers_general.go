package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/reaper47/recipya/internal/templates"
	"net/http"
	"strconv"
)

func (s *Server) downloadHandler(w http.ResponseWriter, r *http.Request) {
	file := chi.URLParam(r, "tmpFile")
	data, err := s.Files.ReadTempFile(file)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", http.DetectContentType(data))
	w.Header().Set("Content-Disposition", `attachment; filename="`+file+`"`)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
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

func notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	templates.Render(w, templates.Simple, templates.PageNotFound)
}

func (s *Server) userInitialsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.Write([]byte(s.Repository.UserInitials(userID.(int64))))
}
