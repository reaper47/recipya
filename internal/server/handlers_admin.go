package server

import (
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/auth"
	"github.com/reaper47/recipya/internal/models"
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
		if app.Config.Server.IsDemo {
			w.Header().Set("HX-Trigger", models.NewErrorToast("Every day is Christmas.", "", "OK").Render())
			w.WriteHeader(http.StatusTeapot)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		if !regex.Email.MatchString(email) || password == "" {
			w.Header().Set("HX-Trigger", models.NewErrorFormToast("Email and/or password is invalid.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID := s.Repository.UserID(email)
		if userID != -1 {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("User exists.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hashPassword, err := auth.HashPassword(password)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorAuthToast("Error encoding your password.").Render())
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		_, err = s.Repository.Register(email, hashPassword)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Failed to add user.").Render())
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

		if app.Config.Server.IsDemo {
			w.Header().Set("HX-Trigger", models.NewErrorGeneralToast("Who do you think you are, eh?").Render())
			w.WriteHeader(http.StatusTeapot)
			return
		}

		if userID == 1 {
			w.Header().Set("HX-Trigger", models.NewErrorGeneralToast("Cannot delete admin.").Render())
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
