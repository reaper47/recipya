package server

import (
	"fmt"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/web/components"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"strconv"
)

func (s *Server) cookbookShareHandler(w http.ResponseWriter, r *http.Request) {
	userID, isLoggedIn := s.findUserID(r)

	share, err := s.Repository.CookbookShared(r.URL.String())
	if err != nil {
		notFoundHandler(w, r)
		return
	}

	cookbook, err := s.Repository.Cookbook(share.CookbookID, share.UserID)
	if err != nil {
		notFoundHandler(w, r)
		return
	}

	_ = components.CookbookIndex(templates.Data{
		About: templates.AboutData{
			Version: app.Version,
		},
		IsAdmin:         userID == 1,
		IsAuthenticated: isLoggedIn,
		IsHxRequest:     r.Header.Get("Hx-Request") == "true",
		Title:           cookbook.Title,
		Functions:       templates.NewFunctionsData[int64](),
		CookbookFeature: templates.CookbookFeature{
			Cookbook: templates.MakeCookbookView(cookbook, 1, 1),
			ShareData: templates.ShareData{
				IsFromHost: userID == share.UserID,
				IsShared:   true,
			},
		},
	}).Render(r.Context(), w)
}

func (s *Server) cookbooksHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query == nil {
			w.Header().Set("HX-Trigger", models.NewErrorFormToast("Could not parse query.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID := getUserID(r)
		view := query.Get("view")
		if view != "" {
			mode := models.ViewModeFromString(view)
			err := s.Repository.UpdateUserSettingsCookbooksViewMode(userID, mode)
			if err != nil {
				w.Header().Set("HX-Trigger", models.NewErrorDBToast("Error updating user settings.").Render())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		settings, err := s.Repository.UserSettings(userID)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Could not get user settings.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		page, err := strconv.ParseUint(query.Get("page"), 10, 64)
		if err != nil {
			page = 1
		}

		cookbooks, err := s.Repository.Cookbooks(userID, page)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Error getting cookbooks.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		p, err := newCookbooksPagination(s, w, userID, page, false)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorGeneralToast("Error updating pagination.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = components.CookbooksIndex(templates.Data{
			About: templates.AboutData{
				Version: app.Version,
			},
			CookbookFeature: templates.CookbookFeature{
				Cookbooks: cookbooks,
				MakeCookbook: func(index int64, cookbook models.Cookbook, page uint64) templates.CookbookView {
					return templates.MakeCookbookView(cookbook, index, page)
				},
				ShareData: templates.ShareData{IsFromHost: true},
				ViewMode:  settings.CookbooksViewMode,
			},
			IsAdmin:         userID == 1,
			IsAuthenticated: true,
			IsHxRequest:     r.Header.Get("Hx-Request") == "true",
			Title:           "Cookbooks",
			Pagination:      p,
		}).Render(r.Context(), w)
	}
}

func (s *Server) cookbooksPostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := []cases.Caser{cases.Title(language.AmericanEnglish, cases.NoLower)}[0].String(r.Header.Get("HX-Prompt"))
		if title == "" {
			w.Header().Set("HX-Trigger", models.NewErrorReqToast("Title must not be empty.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID := getUserID(r)
		cookbookID, err := s.Repository.AddCookbook(title, userID)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Could not create cookbook.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		settings, err := s.Repository.UserSettings(userID)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Could not get user settings.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		query := r.URL.Query()
		if query == nil {
			w.Header().Set("HX-Trigger", models.NewErrorReqToast("Query is empty.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		pageStr := query.Get("page")
		page, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			page = 1
		}

		p, err := newCookbooksPagination(s, w, userID, page, true)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Error updating pagination.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if p.NumResults == 1 {
			s.cookbooksHandler().ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusCreated)

		data := templates.Data{
			CookbookFeature: templates.CookbookFeature{
				Cookbook: templates.CookbookView{
					ID:         cookbookID,
					PageItemID: int64(p.NumResults),
					PageNumber: page,
					Title:      title,
				},
				ShareData: templates.ShareData{IsFromHost: true},
			},
			Pagination: p,
		}

		if settings.CookbooksViewMode == models.GridViewMode {
			_ = components.CookbookGridAdd(data).Render(r.Context(), w)
		} else {
			_ = components.CookbookListAdd(data).Render(r.Context(), w)
		}
	}
}

func (s *Server) cookbooksDeleteCookbookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookbookID, err := parsePathPositiveID(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		query := r.URL.Query()
		if query == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID := getUserID(r)
		pageStr := query.Get("page")
		page, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			page = 1
		}

		err = s.Repository.DeleteCookbook(cookbookID, userID)
		if err != nil {
			// TODO: Log it with slog
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Error deleting cookbook.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		p, err := newCookbooksPagination(s, w, userID, page, true)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorGeneralToast("Could not update pagination.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if p.NumResults == 0 {
			w.Header().Set("HX-Refresh", "true")
			return
		}

		_ = components.Pagination(p).Render(r.Context(), w)
	}
}

func newCookbooksPagination(srv *Server, w http.ResponseWriter, userID int64, page uint64, isSwap bool) (templates.Pagination, error) {
	counts, err := srv.Repository.Counts(userID)
	if err != nil {
		w.Header().Set("HX-Trigger", models.NewErrorDBToast("Error getting counts.").Render())
		w.WriteHeader(http.StatusInternalServerError)
		return templates.Pagination{}, err
	}

	numPages := counts.Cookbooks / templates.ResultsPerPage
	if numPages == 0 {
		numPages = 1
	}

	htmx := templates.PaginationHtmx{
		IsSwap: isSwap,
		Target: "#content",
	}
	return templates.NewPagination(page, numPages, counts.Cookbooks, templates.ResultsPerPage, "/cookbooks", "", htmx), nil
}

func (s *Server) cookbooksDeleteCookbookRecipeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookbookID, err := parsePathPositiveID(r.PathValue("id"))
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorReqToast("Cookbook ID must be >= 0.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recipeID, err := parsePathPositiveID(r.PathValue("recipeID"))
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorReqToast("Recipe ID must be >= 0.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID := getUserID(r)

		numRecipes, err := s.Repository.DeleteRecipeFromCookbook(recipeID, cookbookID, userID)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Error deleting recipe from cookbook.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if numRecipes == 0 {
			_ = components.CookbookIndexNoRecipes(true).Render(r.Context(), w)
		}
	}
}

func (s *Server) cookbooksDownloadCookbookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookbookID, err := parsePathPositiveID(r.PathValue("id"))
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorReqToast("Could not parse cookbook ID.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cookbook, err := s.Repository.Cookbook(cookbookID, getUserID(r))
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Could not fetch cookbook.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(cookbook.Recipes) == 0 {
			w.Header().Set("HX-Trigger", models.NewErrorReqToast("Cookbook is empty.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fileName, err := s.Files.ExportCookbook(cookbook, models.PDF)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorFilesToast("Failed to export cookbook.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/download/"+fileName)
		w.WriteHeader(http.StatusSeeOther)
	}
}

func (s *Server) cookbooksGetCookbookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query == nil {
			w.Header().Set("HX-Trigger", models.NewErrorReqToast("Could not parse query.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := parsePathPositiveID(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		isHxRequest := r.Header.Get("Hx-Request") == "true"

		page, err := strconv.ParseUint(query.Get("page"), 10, 64)
		if err != nil {
			page = 1
		}

		cookbook, err := s.Repository.Cookbook(id, getUserID(r))
		if err != nil {
			// TODO: Create an error type to differentiate between http.StatusNotFound and http.StatusInternalServerError
			http.NotFound(w, r)
			return
		}

		_ = components.CookbookIndex(templates.Data{
			About: templates.NewAboutData(),
			CookbookFeature: templates.CookbookFeature{
				Cookbook:  templates.MakeCookbookView(cookbook, id-1, page),
				ShareData: templates.ShareData{IsFromHost: true},
			},
			IsAdmin:         getUserID(r) == 1,
			IsAuthenticated: true,
			IsHxRequest:     isHxRequest,
			Functions:       templates.NewFunctionsData[int64](),
			Pagination:      templates.Pagination{IsHidden: true},
			Title:           cookbook.Title,
		}).Render(r.Context(), w)
	}
}

func (s *Server) cookbooksImagePostCookbookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookbookID, err := parsePathPositiveID(r.PathValue("id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 128<<20)

		err = r.ParseMultipartForm(128 << 20)
		form := r.MultipartForm
		if err != nil || form == nil || form.File == nil {
			w.Header().Set("HX-Trigger", models.NewErrorFormToast("Could not parse the uploaded image.").Render())
			return
		}

		imageFile, ok := form.File["image"]
		if !ok {
			w.Header().Set("HX-Trigger", models.NewErrorFormToast("Could not retrieve the image from the form.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		f, err := imageFile[0].Open()
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorFormToast("Could not open the image from the form.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer func() {
			_ = f.Close()
		}()

		imageUUID, err := s.Files.UploadImage(f)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorFilesToast("Error uploading image.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userID := getUserID(r)
		err = s.Repository.UpdateCookbookImage(cookbookID, imageUUID, userID)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Error updating image.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Server) cookbookPostCookbookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookbookID, err := parsePathPositiveID(r.PathValue("id"))
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorFormToast("Missing 'cookbookId' in body.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recipeIDStr := r.FormValue("recipeId")
		recipeID, err := strconv.ParseInt(recipeIDStr, 10, 64)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorFormToast("Missing 'recipeId' in body.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID := getUserID(r)
		err = s.Repository.AddCookbookRecipe(cookbookID, recipeID, userID)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Could not add recipe to cookbook.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (s *Server) cookbooksRecipesSearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		idStr := r.PathValue("id")
		id, err := parsePathPositiveID(idStr)
		if err != nil || id <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Cookbook ID in path must be > 0."))
			return
		}

		page, err := strconv.ParseUint(query.Get("page"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Missing page number in query."))
			return
		}

		if page == 0 {
			page = 1
		}

		userID := getUserID(r)

		cookbook, err := s.Repository.Cookbook(id, userID)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Error getting cookbooks.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		q := query.Get("q")
		if q == "" {
			w.Header().Set("HX-Redirect", "/cookbooks/"+strconv.FormatInt(id, 10))
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("Query parameter must not be 'q' empty."))
			return
		}

		mode := query.Get("mode")
		if mode == "" {
			mode = "name"
		}

		opts := models.NewSearchOptionsRecipe(mode, page)
		opts.CookbookID = id

		recipes, totalCount, err := s.Repository.SearchRecipes(q, page, opts, userID)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Error searching recipes.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(recipes) == 0 {
			_ = components.SearchNoResult().Render(r.Context(), w)
			return
		}

		numPages := totalCount / templates.ResultsPerPage
		if numPages == 0 {
			numPages = 1
		}

		isHxReq := r.Header.Get("HX-Request") == "true"

		params := "q=" + q + "&mode=" + mode
		htmx := templates.PaginationHtmx{IsSwap: isHxReq, Target: "#search-results"}
		p := templates.NewPagination(page, numPages, totalCount, templates.ResultsPerPage, "/cookbooks/"+idStr+"/recipes/search", params, htmx)
		p.Search.CurrentPage = page

		_ = components.CookbookSearchRecipes(templates.Data{
			About:   templates.NewAboutData(),
			Content: q,
			CookbookFeature: templates.CookbookFeature{
				Cookbook: templates.CookbookView{
					ID:         cookbook.ID,
					PageItemID: id,
					PageNumber: page,
					Recipes:    recipes,
					Title:      cookbook.Title,
				},
				ShareData: templates.ShareData{
					IsFromHost: true,
					IsShared:   false,
				},
			},
			IsAdmin:         userID == 1,
			IsAutologin:     app.Config.Server.IsAutologin,
			IsAuthenticated: true,
			IsHxRequest:     isHxReq,
			Functions:       templates.NewFunctionsData[int64](),
			Pagination:      p,
		}).Render(r.Context(), w)
	}
}

func (s *Server) cookbooksPostCookbookReorderHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookbookID, err := parsePathPositiveID(r.PathValue("id"))
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorReqToast("Missing cookbook ID in body.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = r.ParseForm()
		if err != nil {
			return
		}

		recipeIDsStr := r.Form["recipe-id"]
		if len(recipeIDsStr) == 0 {
			w.Header().Set("HX-Trigger", models.NewErrorFormToast("Missing recipe IDs in body.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recipeIDs := make([]uint64, len(recipeIDsStr))
		for i, s := range recipeIDsStr {
			id, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				w.Header().Set("HX-Trigger", models.NewErrorFormToast(fmt.Sprintf("Recipe ID %q is invalid.", s)).Render())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			recipeIDs[i] = id
		}

		err = s.Repository.ReorderCookbookRecipes(cookbookID, recipeIDs, getUserID(r))
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Failed to update indices.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) cookbookSharePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookbookID, err := parsePathPositiveID(r.PathValue("id"))
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorReqToast("Cookbook ID must be positive.").Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userID := getUserID(r)
		share := models.Share{CookbookID: cookbookID, RecipeID: -1, UserID: userID}

		link, err := s.Repository.AddShareLink(share)
		if err != nil {
			w.Header().Set("HX-Trigger", models.NewErrorDBToast("Failed to create share link.").Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = components.ShareLink(templates.Data{
			Content: r.Host + link,
		}).Render(r.Context(), w)
	}
}
