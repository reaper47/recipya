package server

import (
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/templates"
	"github.com/reaper47/recipya/web/components"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log/slog"
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
		About:           templates.NewAboutData(),
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
			About: templates.NewAboutData(),
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
		userIDAttr := slog.Int64("userID", userID)

		cookbookID, err := s.Repository.AddCookbook(title, userID)
		if err != nil {
			msg := "Could not create cookbook."
			slog.Error(msg, userIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		settings, err := s.Repository.UserSettings(userID)
		if err != nil {
			msg := "Could not get user settings."
			slog.Error(msg, userIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
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
			msg := "Error updating pagination."
			slog.Error(msg, userIDAttr, "page", page, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if p.NumResults == 1 {
			s.cookbooksHandler().ServeHTTP(w, r)
			return
		}

		slog.Info("Created cookbook", userIDAttr, "cookbookID", cookbookID, "title", title)

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

		pageStr := query.Get("page")
		page, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			page = 1
		}

		userID := getUserID(r)
		userIDAttr := slog.Int64("userID", userID)
		cookbookIDAttr := slog.Int64("cookbookID", cookbookID)
		pageAttr := slog.Uint64("page", page)

		err = s.Repository.DeleteCookbook(cookbookID, userID)
		if err != nil {
			msg := "Error deleting cookbook."
			slog.Error(msg, userIDAttr, cookbookIDAttr, pageAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		p, err := newCookbooksPagination(s, w, userID, page, true)
		if err != nil {
			msg := "Could not update pagination."
			slog.Error(msg, userIDAttr, cookbookIDAttr, pageAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorGeneralToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if p.NumResults == 0 {
			w.Header().Set("HX-Refresh", "true")
			return
		}

		slog.Info("Deleted cookbook", userIDAttr, cookbookIDAttr)
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
		userID := getUserID(r)
		userIDAttr := slog.Int64("userID", userID)

		cookbookID, err := parsePathPositiveID(r.PathValue("id"))
		if err != nil {
			msg := "Cookbook ID must be >= 0."
			slog.Error(msg, userIDAttr)
			w.Header().Set("HX-Trigger", models.NewErrorReqToast(msg).Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recipeID, err := parsePathPositiveID(r.PathValue("recipeID"))
		if err != nil {
			msg := "Recipe ID must be >= 0."
			slog.Error(msg, userIDAttr)
			w.Header().Set("HX-Trigger", models.NewErrorReqToast(msg).Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cookbookIDAttr := slog.Int64("cookbookID", cookbookID)
		recipeIDAttr := slog.Int64("recipeID", recipeID)

		numRecipes, err := s.Repository.DeleteRecipeFromCookbook(recipeID, cookbookID, userID)
		if err != nil {
			msg := "Error deleting recipe from cookbook."
			slog.Error(msg, userIDAttr, cookbookIDAttr, recipeIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		slog.Info("Removed recipe from cookbook", userIDAttr, cookbookIDAttr, recipeIDAttr)

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

		mode := query.Get("mode")
		if mode == "" {
			mode = "full"
		}

		sorts := query.Get("sort")
		if sorts == "" {
			sorts = "default"
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
			Searchbar: templates.SearchbarData{
				Mode: mode,
				Sort: sorts,
			},
			Title: cookbook.Title,
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

		userID := getUserID(r)
		userIDAttr := slog.Int64("userID", userID)
		cookbookIDAttr := slog.Int64("cookbookID", cookbookID)

		r.Body = http.MaxBytesReader(w, r.Body, 128<<20)

		err = r.ParseMultipartForm(128 << 20)
		form := r.MultipartForm
		if err != nil || form == nil || form.File == nil {
			msg := "Could not parse the uploaded image."
			slog.Error(msg, userIDAttr, cookbookIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorFormToast(msg).Render())
			return
		}

		imageFile, ok := form.File["image"]
		if !ok {
			msg := "Could not retrieve the image from the form."
			slog.Error(msg, userIDAttr, cookbookIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorFormToast(msg).Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		f, err := imageFile[0].Open()
		if err != nil {
			msg := "Could not open the image from the form."
			slog.Error(msg, userIDAttr, cookbookIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorFormToast(msg).Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer f.Close()

		imageUUID, err := s.Files.UploadImage(f)
		if err != nil {
			msg := "Error uploading image."
			slog.Error(msg, userIDAttr, cookbookIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorFilesToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		imageUUIDAttr := slog.String("imageUUID", imageUUID.String())

		err = s.Repository.UpdateCookbookImage(cookbookID, imageUUID, userID)
		if err != nil {
			msg := "Error updating image."
			slog.Error(msg, userIDAttr, cookbookIDAttr, imageUUIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		slog.Info("Updated cookbook image", userIDAttr, cookbookIDAttr, imageUUIDAttr)
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

		userIDAttr := slog.Int64("userID", userID)
		recipeIDAttr := slog.Int64("recipeID", recipeID)
		cookbookIDAttr := slog.Int64("cookbookID", cookbookID)

		err = s.Repository.AddCookbookRecipe(cookbookID, recipeID, userID)
		if err != nil {
			msg := "Could not add recipe to cookbook."
			slog.Error(msg, userIDAttr, cookbookIDAttr, recipeIDAttr)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		slog.Info("Added recipe to cookbook", userIDAttr, cookbookIDAttr, recipeIDAttr)
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
		sorts := query.Get("sort")

		opts := models.NewSearchOptionsRecipe(mode, sorts, page)
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

		params := "q=" + q + "&mode=" + mode + "&sort=" + sorts
		htmx := templates.PaginationHtmx{IsSwap: isHxReq, Target: "#search-results"}
		p := templates.NewPagination(page, numPages, totalCount, templates.ResultsPerPage, "/cookbooks/"+idStr+"/recipes/search", params, htmx)
		p.Search.CurrentPage = page

		_ = components.CookbookSearchRecipes(templates.Data{
			About: templates.NewAboutData(),
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
			Searchbar: templates.SearchbarData{
				Mode: mode,
				Sort: sorts,
				Term: q,
			},
		}).Render(r.Context(), w)
	}
}

func (s *Server) cookbooksPostCookbookReorderHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)
		userIDAttr := slog.Int64("userID", userID)

		cookbookID, err := parsePathPositiveID(r.PathValue("id"))
		if err != nil {
			msg := "Missing cookbook ID in body."
			slog.Error(msg, userIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorReqToast(msg).Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = r.ParseForm()
		if err != nil {
			msg := "Form could not be parsed."
			slog.Error(msg, userIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorReqToast(msg).Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		recipeIDsStr := r.Form["recipe-id"]
		if len(recipeIDsStr) == 0 {
			msg := "Missing recipe IDs in body."
			slog.Error(msg, userIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorFormToast(msg).Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cookbookIDAttr := slog.Int64("cookbookID", cookbookID)

		recipeIDs := make([]uint64, len(recipeIDsStr))
		for i, id := range recipeIDsStr {
			id, err := strconv.ParseUint(id, 10, 64)
			if err != nil {
				msg := "Recipe ID could not be parsed."
				slog.Error(msg, userIDAttr, cookbookIDAttr, "id", id, "error", err)
				w.Header().Set("HX-Trigger", models.NewErrorFormToast(msg).Render())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			recipeIDs[i] = id
		}

		err = s.Repository.ReorderCookbookRecipes(cookbookID, recipeIDs, userID)
		if err != nil {
			msg := "Failed to update indices."
			slog.Error(msg, userIDAttr, cookbookIDAttr, "recipeIDs", recipeIDs, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		slog.Info("Reordered recipes in cookbook", userIDAttr, cookbookIDAttr, "recipeIDs", recipeIDs)
		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *Server) cookbookSharePostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := getUserID(r)
		userIDAttr := slog.Int64("user_id", userID)

		cookbookID, err := parsePathPositiveID(r.PathValue("id"))
		if err != nil {
			msg := "Cookbook ID must be positive."
			slog.Error(msg, userIDAttr, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorReqToast(msg).Render())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		share := models.Share{CookbookID: cookbookID, RecipeID: -1, UserID: userID}

		link, err := s.Repository.AddShareLink(share)
		if err != nil {
			msg := "Failed to create share link."
			slog.Error(msg, userIDAttr, "share", share, "error", err)
			w.Header().Set("HX-Trigger", models.NewErrorDBToast(msg).Render())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		slog.Info("Cookbook shared", userIDAttr, "share", share, "link", link)

		_ = components.ShareLink(templates.Data{
			Content: r.Host + link,
		}).Render(r.Context(), w)
	}
}
