package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"net/http"
	"slices"
	"strconv"
	"strings"
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
			MakeCookbook: func(index int64, cookbook models.Cookbook, page uint64) models.CookbookView {
				return cookbook.MakeView(index, page)
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
			Cookbook: models.CookbookView{
				ID:         cookbookID,
				PageItemID: int64(p.NumResults),
				PageNumber: page,
				Title:      title,
			},
		},
		Pagination: p,
	})
}

func (s *Server) cookbooksDeleteCookbookHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) cookbooksDeleteCookbookRecipeHandler(w http.ResponseWriter, r *http.Request) {
	cookbookIDStr := chi.URLParam(r, "id")
	cookbookID, err := strconv.ParseUint(cookbookIDStr, 10, 64)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Cookbook ID in the URL must be >= 0.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recipeIDStr := chi.URLParam(r, "recipeID")
	recipeID, err := strconv.ParseUint(recipeIDStr, 10, 64)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Recipe ID in the URL must be >= 0.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := getUserID(r)

	numRecipes, err := s.Repository.DeleteRecipeFromCookbook(recipeID, cookbookID, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error deleting recipe from cookbook.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if numRecipes == 0 {
		templates.RenderComponent(w, "cookbooks", "cookbook-index-no-recipes", nil)
	}
}

func (s *Server) cookbooksGetCookbookHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isHxRequest := r.Header.Get("Hx-Request") == "true"

	userID := getUserID(r)
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		if !isHxRequest {
			http.Redirect(w, r, "/cookbooks", http.StatusSeeOther)
			return
		}

		w.Header().Set("HX-Trigger", makeToast("Missing page query parameter.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookbook, err := s.Repository.Cookbook(id, userID, page)
	if err != nil {
		// TODO: Create an error type to differentiate between http.StatusNotFound and http.StatusInternalServerError
		http.NotFound(w, r)
		return
	}

	data := templates.Data{
		IsAuthenticated: true,
		IsHxRequest:     isHxRequest,
		Functions:       templates.NewFunctionsData(),
		CookbookFeature: templates.CookbookFeature{
			Cookbook: cookbook.MakeView(id-1, page),
		},
		Title: cookbook.Title,
	}

	if isHxRequest {
		templates.RenderComponent(w, "cookbooks", "cookbook-index", data)
	} else {
		templates.Render(w, templates.CookbookPage, data)
	}
}

func (s *Server) cookbooksImagePostCookbookHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) cookbookPostCookbookHandler(w http.ResponseWriter, r *http.Request) {
	cookbookIDStr := r.FormValue("cookbookId")
	cookbookID, err := strconv.ParseInt(cookbookIDStr, 10, 64)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Missing 'cookbookId' in body.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recipeIDStr := r.FormValue("recipeId")
	recipeID, err := strconv.ParseInt(recipeIDStr, 10, 64)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Missing 'recipeId' in body.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := getUserID(r)
	err = s.Repository.AddCookbookRecipe(cookbookID, recipeID, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not add recipe to cookbook.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) cookbooksRecipesSearchPostHandler(w http.ResponseWriter, r *http.Request) {
	cookbookID := r.FormValue("id")
	id, err := strconv.ParseInt(cookbookID, 10, 64)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Missing cookbook ID in body.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pageNumber := r.FormValue("page")
	page, err := strconv.ParseUint(pageNumber, 10, 64)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Missing page number in body.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := getUserID(r)

	cookbook, err := s.Repository.Cookbook(id, userID, page)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error getting cookbooks.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	q := r.FormValue("q")
	q = strings.TrimSpace(q)
	if q == "" {
		templates.RenderComponent(w, "cookbooks", "cookbook-recipes", templates.Data{
			Functions: templates.NewFunctionsData(),
			CookbookFeature: templates.CookbookFeature{
				Cookbook: cookbook.MakeView(1, page),
			},
		})
		return
	}

	recipes, err := s.Repository.SearchRecipes(q, models.SearchOptionsRecipes{ByName: true}, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error searching recipes.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	recipes = slices.DeleteFunc(recipes, func(r1 models.Recipe) bool {
		return slices.ContainsFunc(cookbook.Recipes, func(r2 models.Recipe) bool {
			return r1.ID == r2.ID
		})
	})

	if len(recipes) == 0 {
		templates.RenderComponent(w, "search", "no-result", nil)
		return
	}

	templates.RenderComponent(w, "search", "cookbooks-search-results-recipes", templates.Data{
		CookbookFeature: templates.CookbookFeature{
			Cookbook: models.CookbookView{
				ID:         cookbook.ID,
				PageItemID: id,
				PageNumber: page,
			},
		},
		Recipes: recipes,
	})
}

func (s *Server) cookbooksPostCookbookReorderHandler(w http.ResponseWriter, r *http.Request) {
	cookbookIDStr := r.FormValue("cookbook-id")
	cookbookID, err := strconv.ParseInt(cookbookIDStr, 10, 64)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Missing cookbook ID in body.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := getUserID(r)

	recipeIDsStr := r.Form["recipe-id"]
	if len(recipeIDsStr) == 0 {
		w.Header().Set("HX-Trigger", makeToast("Missing recipe IDs in body.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var recipeIDs []uint64
	for _, s := range recipeIDsStr {
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			w.Header().Set("HX-Trigger", makeToast(fmt.Sprintf("Recipe ID %q is invalid.", s), errorToast))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		recipeIDs = append(recipeIDs, id)
	}

	err = s.Repository.ReorderCookbookRecipes(cookbookID, recipeIDs, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Failed to update indices.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
