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
	userID := getUserID(r)

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

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		page = 1
	}

	cookbooks, err := s.Repository.Cookbooks(userID, page)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error getting cookbooks.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p, err := newCookbooksPagination(s, w, userID, page, false)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error updating pagination.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isHxRequest := r.Header.Get("Hx-Request") == "true"

	data := templates.Data{
		CookbookFeature: templates.CookbookFeature{
			Cookbooks: cookbooks,
			MakeCookbook: func(index int64, cookbook models.Cookbook) templates.CookbookView {
				return templates.CookbookView{
					ID:          cookbook.ID,
					Image:       cookbook.Image,
					IsUUIDValid: cookbook.Image != uuid.Nil,
					NumRecipes:  cookbook.Count,
					PageItemID:  index + 1,
					Title:       cookbook.Title,
				}
			},
			ViewMode: settings.CookbooksViewMode,
		},
		IsAuthenticated: true,
		IsHxRequest:     isHxRequest,
		Title:           "Cookbooks",
		Pagination:      p,
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

	userID := getUserID(r)
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

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		page = 1
	}

	p, err := newCookbooksPagination(s, w, userID, page, true)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error updating pagination.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if p.NumResults == 1 {
		s.cookbooksHandler(w, r)
		return
	}

	tmpl := "cookbook-grid"
	if settings.CookbooksViewMode == models.ListViewMode {
		tmpl = "cookbook-list"
	}

	w.WriteHeader(http.StatusCreated)
	templates.RenderComponent(w, "cookbooks", tmpl+"-add", templates.Data{
		CookbookFeature: templates.CookbookFeature{
			Cookbook: templates.CookbookView{
				ID:         cookbookID,
				PageItemID: int64(p.NumResults),
				Title:      title,
			},
		},
		Pagination: p,
	})
}

func (s *Server) cookbooksDeleteHandler(w http.ResponseWriter, r *http.Request) {
	cookbookIDStr := chi.URLParam(r, "id")
	cookbookID, err := strconv.ParseInt(cookbookIDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := getUserID(r)
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		page = 1
	}

	err = s.Repository.DeleteCookbook(cookbookID, userID, page)
	if err != nil {
		// TODO: Log it with slog
		w.Header().Set("HX-Trigger", makeToast("Error deleting cookbook.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	p, err := newCookbooksPagination(s, w, userID, page, true)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error updating pagination.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if p.NumResults == 0 {
		w.Header().Set("HX-Refresh", "true")
		return
	}

	templates.RenderComponent(w, "cookbooks", "pagination", p)
}

func newCookbooksPagination(srv *Server, w http.ResponseWriter, userID int64, page uint64, isSwap bool) (templates.Pagination, error) {
	counts, err := srv.Repository.Counts(userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error getting counts.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return templates.Pagination{}, err
	}

	numPages := counts.Cookbooks / templates.ResultsPerPage
	if numPages == 0 {
		numPages = 1
	}

	return templates.NewPagination(page, numPages, counts.Cookbooks, templates.ResultsPerPage, "/cookbooks", isSwap), nil
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

	userID := getUserID(r)
	err = s.Repository.UpdateCookbookImage(cookbookID, imageUUID, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error updating the cookbook's image.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
