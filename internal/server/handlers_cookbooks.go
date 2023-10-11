package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"net/http"
	"strconv"
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
			MakeCookbook: func(index int64, cookbook models.Cookbook) templates.CookbookView {
				return templates.CookbookView{
					ID:          index + 1,
					Image:       cookbook.Image,
					IsUUIDValid: cookbook.Image != uuid.Nil,
					NumRecipes:  cookbook.Count,
					Title:       cookbook.Title,
				}
			},
			ViewMode: settings.CookbooksViewMode,
		},
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

	settings, err := s.Repository.UserSettings(userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not get user settings.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl := "cookbook-grid"
	if settings.CookbooksViewMode == models.ListViewMode {
		tmpl = "cookbook-list"
	}

	w.WriteHeader(http.StatusCreated)
	templates.RenderComponent(w, "cookbooks", tmpl, templates.CookbookView{
		ID:          cookbookID,
		IsUUIDValid: false,
		Title:       title,
	})
}

func (s *Server) cookbooksImagePostHandler(w http.ResponseWriter, r *http.Request) {
	cookbookIDStr := chi.URLParam(r, "id")
	cookbookID, err := strconv.ParseInt(cookbookIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 128<<20)

	err = r.ParseMultipartForm(128 << 20)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the uploaded image.", errorToast))
		return
	}

	imageFile, ok := r.MultipartForm.File["image"]
	if !ok {
		w.Header().Set("HX-Trigger", makeToast("Could not retrieve the image from the form.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	f, err := imageFile[0].Open()
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could open the image from the form.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer f.Close()

	imageUUID, err := s.Files.UploadImage(f)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error uploading image.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value("userID").(int64)
	err = s.Repository.UpdateCookbookImage(cookbookID, imageUUID, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error updating the cookbook's image.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
