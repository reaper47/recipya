package server

import (
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/internal/utils/regex"
	"github.com/reaper47/recipya/web/components"
	"net/http"
)

func (s *Server) adminHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = components.Admin(templates.Data{
			IsAdmin:         true,
			IsAuthenticated: true,
			IsHxRequest:     r.Header.Get("Hx-Request") == "true",
			Admin: templates.AdminData{
				Users: s.Repository.Users(),
			},
		}).Render(r.Context(), w)
	}
}

func (s *Server) adminUsersPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")
		if !regex.Email.MatchString(email) || password == "" {
			w.Header().Set("HX-Trigger", makeToast("Email and/or password is invalid.", errorToast))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID := s.Repository.UserID(email)
		if userID != -1 {
			w.Header().Set("HX-Trigger", makeToast("User exists.", errorToast))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hashPassword, err := auth.HashPassword(password)
		if err != nil {
			w.Header().Set("HX-Trigger", makeToast("Error encoding your password.", errorToast))
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		userID, err = s.Repository.Register(email, hashPassword)
		if err != nil {
			w.Header().Set("HX-Trigger", makeToast("Failed to add user.", errorToast))
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = components.AdminUserRow(email, true, true).Render(r.Context(), w)
	}
}

func (s *Server) adminUsersDeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := s.Repository.UserID(r.PathValue("email"))
		if userID == -1 {
			return
		}

		if app.Config.Server.IsDemo && s.Repository.UserID("demo@demo.com") == userID {
			w.Header().Set("HX-Trigger", makeToast("Don't look up.", errorToast))
			w.WriteHeader(http.StatusTeapot)
			return
		}

		if userID == 1 {
			w.Header().Set("HX-Trigger", makeToast("Cannot delete admin.", errorToast))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := s.Repository.DeleteUser(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
