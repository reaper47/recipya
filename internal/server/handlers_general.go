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
	_, _ = w.Write(data)
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	if isAuthenticated(r, s.Repository.GetAuthToken) {
		middleware := s.mustBeLoggedInMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.recipesHandler(w, r)
		}))
		middleware.ServeHTTP(w, r)
	} else {
		page := templates.LandingPage
		templates.Render(w, page, templates.Data{
			IsAuthenticated: false,
			Title:           page.Title(),
		})
	}
}

func notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	templates.Render(w, templates.Simple, templates.PageNotFound)
}

func (s *Server) userInitialsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey)
	if userID == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	_, _ = w.Write([]byte(s.Repository.UserInitials(userID.(int64))))
}
