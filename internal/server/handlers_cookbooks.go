package server

import (
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"net/http"
)

func (s *Server) cookbooksHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	view := r.URL.Query().Get("view")
	if view != "" {
		mode := models.ViewModeFromString(view)
		err := s.Repository.UpdateUserSettingsCookbooksViewMode(userID, mode)
		if err != nil {
			w.Header().Set("HX-Trigger", makeToast("Error updating user settings.", errorToast))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	settings, err := s.Repository.UserSettings(userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not get user settings.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookbooks, err := s.Repository.Cookbooks(userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error getting cookbooks.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isHxRequest := r.Header.Get("Hx-Request") == "true"

	data := templates.Data{
		CookbookFeature: templates.CookbookFeature{
			Cookbooks: cookbooks,
			ViewMode:  settings.CookbooksViewMode,
		},
		Functions: templates.FunctionsData{
			Inc: func(v int64) int64 {
				return v + 1
			},
			IsUUIDValid: func(u uuid.UUID) bool {
				return u != uuid.Nil
			},
		},
		IsAuthenticated: true,
		IsHxRequest:     isHxRequest,

		Title: "Cookbooks",
	}

	if isHxRequest {
		templates.RenderComponent(w, "cookbooks", "cookbooks-index", data)
	} else {
		templates.Render(w, templates.CookbooksPage, data)
	}
}

func (s *Server) cookbooksPostHandler(w http.ResponseWriter, r *http.Request) {
	title := r.Header.Get("HX-Prompt")
	if title == "" {
		w.Header().Set("HX-Trigger", makeToast("Title must not be empty.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int64)
	cookbookID, err := s.Repository.AddCookbook(title, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not create cookbook.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", makeToast("Cookbook created.", infoToast))
	w.WriteHeader(http.StatusCreated)

	templates.RenderComponent(w, "cookbooks", "cookbooks-grid", templates.Data{
		CookbookFeature: templates.CookbookFeature{
			Cookbooks: []models.Cookbook{{
				ID:    cookbookID,
				Title: title,
			}},
		},
		Functions: templates.FunctionsData{
			Inc: func(v int64) int64 {
				return v + 1
			},
			IsUUIDValid: func(u uuid.UUID) bool {
				return u != uuid.Nil
			},
		},
	})
}
