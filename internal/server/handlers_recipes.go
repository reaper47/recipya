package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/scraper"
	"github.com/reaper47/recipya/internal/templates"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (s *Server) recipesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	baseData := templates.Data{
		Functions: templates.FunctionsData{
			CutString: func(s string, numCharacters int) string {
				if len(s) < numCharacters {
					return s
				}
				return s[:numCharacters] + "…"
			},
			IsUUIDValid: func(u uuid.UUID) bool {
				return u != uuid.Nil
			},
		},
		Recipes: s.Repository.Recipes(userID),
	}

	if r.Header.Get("HX-Request") == "true" {
		templates.RenderComponent(w, "recipes", "list-recipes", baseData)
		return
	} else {
		page := templates.HomePage
		baseData.IsAuthenticated = true
		baseData.Title = page.Title()
		templates.Render(w, page, baseData)
		return
	}
}

func recipesAddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Hx-Request") == "true" {
		parsedURL, err := url.Parse(r.Header.Get("HX-Current-Url"))
		if err == nil && parsedURL.Path == "/recipes/add/unsupported-website" {
			w.Header().Set("HX-Trigger", makeToast("Website requested.", infoToast))
		}

		templates.RenderComponent(w, "recipes", "add-recipe", nil)
	} else {
		page := templates.AddRecipePage
		templates.Render(w, page, templates.Data{
			IsAuthenticated: true,
			Title:           page.Title(),
		})
	}
}

func (s *Server) recipesAddImportHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 128<<20)

	err := r.ParseMultipartForm(128 << 20)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the uploaded files.", errorToast))
		return
	}

	files, filesOk := r.MultipartForm.File["files"]
	if !filesOk {
		w.Header().Set("HX-Trigger", makeToast("Could not retrieve the files or the directory from the form.", errorToast))
		return
	}

	recipes := s.Files.ExtractRecipes(files)
	userID := r.Context().Value("userID").(int64)

	count := 0
	for _, r := range recipes {
		_, err := s.Repository.AddRecipe(&r, userID)
		if err != nil {
			continue
		}
		count += 1
	}

	msg := fmt.Sprintf("Imported %d recipes. %d skipped", count, len(recipes)-count)
	w.Header().Set("HX-Trigger", makeToast(msg, infoToast))
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) recipeAddManualHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)
	categories, err := s.Repository.Categories(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	viewData := &templates.ViewRecipeData{Categories: categories}

	if r.Header.Get("Hx-Request") == "true" {
		templates.RenderComponent(w, "recipes", "add-recipe-manual", templates.Data{
			View: viewData,
		})
	} else {
		page := templates.AddRecipeManualPage
		templates.Render(w, page, templates.Data{
			IsAuthenticated: true,
			Title:           page.Title(),
			View:            viewData,
		})
	}
}

func (s *Server) recipeAddManualPostHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 128<<20)
	err := r.ParseMultipartForm(128 << 20)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ingredients := make([]string, 0)
	i := 1
	for {
		key := fmt.Sprintf("ingredient-%d", i)
		if r.Form.Has(key) {
			ingredients = append(ingredients, r.FormValue(key))
			i++
		} else {
			break
		}
	}

	instructions := make([]string, 0)
	i = 1
	for {
		key := fmt.Sprintf("instruction-%d", i)
		if r.Form.Has(key) {
			instructions = append(instructions, r.FormValue(key))
			i++
		} else {
			break
		}
	}

	times, err := models.NewTimes(r.FormValue("time-preparation"), r.FormValue("time-cooking"))
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error parsing times.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	yield, err := strconv.ParseInt(r.FormValue("yield"), 10, 16)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error parsing yield.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recipe := &models.Recipe{
		Category:     strings.ToLower(r.FormValue("category")),
		CreatedAt:    time.Time{},
		Cuisine:      "",
		Description:  r.FormValue("description"),
		Image:        imageUUID,
		Ingredients:  ingredients,
		Instructions: instructions,
		Keywords:     nil,
		Name:         r.FormValue("title"),
		Nutrition: models.Nutrition{
			Calories:           r.FormValue("calories"),
			Cholesterol:        r.FormValue("cholesterol"),
			Fiber:              r.FormValue("fiber"),
			Protein:            r.FormValue("protein"),
			SaturatedFat:       r.FormValue("saturated-fat"),
			Sodium:             r.FormValue("sodium"),
			Sugars:             r.FormValue("sugars"),
			TotalCarbohydrates: r.FormValue("total-carbohydrates"),
			TotalFat:           r.FormValue("total-fat"),
			UnsaturatedFat:     "",
		},
		Times:     times,
		Tools:     nil,
		UpdatedAt: time.Time{},
		URL:       r.FormValue("source"),
		Yield:     int16(yield),
	}

	userID := r.Context().Value("userID").(int64)
	recipeNumber, err := s.Repository.AddRecipe(recipe, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not add recipe.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("HX-Redirect", "/recipes/"+strconv.FormatInt(recipeNumber, 10))
	w.WriteHeader(http.StatusCreated)
}

func recipeAddManualIngredientHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	i := 1
	for {
		if !r.Form.Has("ingredient-" + strconv.Itoa(i)) {
			break
		}

		i++
	}

	if r.Form.Get(fmt.Sprintf("ingredient-%d", i-1)) == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	templates.RenderComponent(w, "recipes", "add-ingredient", i)
}

func recipeAddManualIngredientDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	count := 0
	i := 1
	for {
		if !r.Form.Has("ingredient-" + strconv.Itoa(i)) {
			break
		}

		i++
		count++
	}

	if count == 1 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var sb strings.Builder
	entry := chi.URLParam(r, "entry")

	curr := 1
	i = 1
	for {
		key := "ingredient-" + strconv.Itoa(i)
		if !r.Form.Has(key) {
			break
		}

		n, _ := strconv.Atoi(entry)
		if n == i {
			i++
			continue
		}

		currStr := strconv.Itoa(curr)
		xs := []string{
			`<li class="pb-2">`,
			`<label><input type="text" name="ingredient-` + currStr + `" placeholder="Ingredient #` + currStr + `" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" ` + `value="` + r.Form.Get(key) + `" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')"></label>`,
			`&nbsp;<button type="button" class="w-10 h-10 text-center bg-green-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-green-600 hover:text-white center dark:bg-green-500" title="Shortcut: Enter" hx-post="/recipes/add/manual/ingredient" hx-target="#ingredients-list" hx-swap="beforeend" hx-include="[name^='ingredient']">+</button>`,
			`&nbsp;<button type="button" class="delete-button w-10 h-10 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/` + currStr + `" hx-include="[name^='ingredient']">-</button>`,
			`&nbsp;<div class="inline-block h-4 cursor-move handle"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/></svg></div>`,
			`</li>`,
		}
		for _, x := range xs {
			sb.WriteString(x)
		}

		i++
		curr++
	}

	fmt.Fprint(w, sb.String())
}

func recipeAddManualInstructionHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	i := 1
	for {
		if !r.Form.Has("instruction-" + strconv.Itoa(i)) {
			break
		}

		i++
	}

	if r.Form.Get(fmt.Sprintf("instruction-%d", i-1)) == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	templates.RenderComponent(w, "recipes", "add-instruction", i)
}

func recipeAddManualInstructionDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	count := 0
	i := 1
	for {
		if !r.Form.Has("instruction-" + strconv.Itoa(i)) {
			break
		}

		i++
		count++
	}

	if count == 1 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var sb strings.Builder
	entry := chi.URLParam(r, "entry")

	curr := 1
	i = 1
	for {
		key := "instruction-" + strconv.Itoa(i)
		if !r.Form.Has(key) {
			break
		}

		n, _ := strconv.Atoi(entry)
		if n == i {
			i++
			continue
		}

		currStr := strconv.Itoa(curr)
		value := r.Form.Get(key)

		xs := []string{
			`<li class="pt-2 md:pl-0"><div class="flex">`,
			`<label class="w-[95%]"><textarea required name="instruction-` + currStr + `" rows="3" class="w-[97%] border border-gray-300 pl-1 dark:bg-gray-900 dark:border-none" placeholder="Instruction #` + currStr + `" _="on keydown if event.ctrlKey and event.key === 'Enter' then halt the event then get next <div/> from the parentElement of me then get first <button/> in it then call htmx.trigger(it, 'click')">` + value + `</textarea>&nbsp;</label>`,
			`<div class="grid">`,
			`<div class="h-4 cursor-move handle grid place-content-center"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/></svg></div>`,
			`<button type="button" class="md:w-7 md:h-7 md:right-auto w-10 h-10 text-center bg-green-300 border border-gray-800 rounded-lg hover:bg-green-600 hover:text-white center dark:bg-green-500" title="Shortcut: CTRL + Enter" hx-post="/recipes/add/manual/instruction" hx-target="#instructions-list" hx-swap="beforeend" hx-include="[name^='instruction']">+</button>`,
			`<button type="button" class="delete-button w-10 h-10 md:w-7 md:h-7 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center dark:bg-red-500" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/` + currStr + `" hx-include="[name^='instruction']">-</button>`,
			`</div></div></li>`,
		}
		for _, x := range xs {
			sb.WriteString(x)
		}

		i++
		curr++
	}

	fmt.Fprint(w, sb.String())
}

func (s *Server) recipesAddRequestWebsiteHandler(w http.ResponseWriter, r *http.Request) {
	s.Email.Send(app.Config.Email.From, templates.EmailRequestWebsite, templates.EmailData{
		Text: r.FormValue("website"),
	})

	w.Header().Set("HX-Redirect", "/recipes/add")
	w.Header().Set("HX-Trigger", makeToast("I love chicken", infoToast))
	http.Redirect(w, r, "/recipes/add", http.StatusSeeOther)
}

func (s *Server) recipesAddWebsiteHandler(w http.ResponseWriter, r *http.Request) {
	rawURL := r.Header.Get("HX-Prompt")
	if rawURL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Invalid URI.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rs, err := scraper.Scrape(rawURL, s.Files)
	if err != nil {
		templates.RenderComponent(w, "recipes", "unsupported-website", templates.Data{
			IsAuthenticated: true,
			Scraper: templates.ScraperData{
				UnsupportedWebsite: rawURL,
			},
		})
		return
	}

	recipe, err := rs.Recipe()
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Recipe schema is invalid.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("userID").(int64)
	recipeNumber, err := s.Repository.AddRecipe(recipe, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Recipe could not be added.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/recipes/"+strconv.FormatInt(recipeNumber, 10))
	w.WriteHeader(http.StatusSeeOther)
}

func (s *Server) recipeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := r.Context().Value("userID").(int64)

	rowsAffected, err := s.Repository.DeleteRecipe(id, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Recipe could not be deleted.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		w.Header().Set("HX-Trigger", makeToast("Recipe not found.", errorToast))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) recipesEditHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value("userID").(int64)
	recipe, err := s.Repository.Recipe(id, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Failed to retrieve recipe.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templates.RenderComponent(w, "recipes", "edit-recipe", templates.Data{
		View: templates.NewViewRecipeData(id, recipe, true, false),
	})
}

func (s *Server) recipesEditPostHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 128<<20)
	err := r.ParseMultipartForm(128 << 20)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedRecipe := models.Recipe{
		Category:     r.FormValue("category"),
		Description:  r.FormValue("description"),
		Ingredients:  make([]string, 0),
		Instructions: make([]string, 0),
		Name:         r.FormValue("title"),
		Nutrition: models.Nutrition{
			Calories:           r.FormValue("calories"),
			Cholesterol:        r.FormValue("cholesterol"),
			Fiber:              r.FormValue("fiber"),
			Protein:            r.FormValue("protein"),
			SaturatedFat:       r.FormValue("saturated-fat"),
			Sodium:             r.FormValue("sodium"),
			Sugars:             r.FormValue("sugars"),
			TotalCarbohydrates: r.FormValue("total-carbohydrates"),
			TotalFat:           r.FormValue("total-fat"),
		},
		URL: r.FormValue("source"),
	}

	imageFile, ok := r.MultipartForm.File["image"]
	if ok {
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updatedRecipe.Image = imageUUID
	}

	times, err := models.NewTimes(r.FormValue("time-preparation"), r.FormValue("time-cooking"))
	if err == nil {
		updatedRecipe.Times = times
	}

	i := 1
	for {
		s := "ingredient-" + strconv.Itoa(i)
		if !r.Form.Has(s) {
			break
		}
		updatedRecipe.Ingredients = append(updatedRecipe.Ingredients, r.FormValue(s))
		i++
	}

	i = 1
	for {
		s := "instruction-" + strconv.Itoa(i)
		if !r.Form.Has(s) {
			break
		}
		updatedRecipe.Instructions = append(updatedRecipe.Instructions, r.FormValue(s))
		i++
	}

	yield, err := strconv.ParseInt(r.FormValue("yield"), 10, 16)
	if err == nil {
		updatedRecipe.Yield = int16(yield)
	}

	userID := r.Context().Value("userID").(int64)
	recipeNumStr := chi.URLParam(r, "id")
	recipeNum, err := strconv.ParseInt(recipeNumStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.Repository.UpdateRecipe(&updatedRecipe, userID, recipeNum)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error updating recipe.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Redirect", "/recipes/"+recipeNumStr)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) recipeShareHandler(w http.ResponseWriter, r *http.Request) {
	var userID int64
	isLoggedIn := true
	userID = getUserIDFromSessionCookie(r)
	if userID == -1 {
		userID = getUserIDFromRememberMeCookie(r, s.Repository.GetAuthToken)
		if userID == -1 {
			isLoggedIn = false
		}
	}

	share, err := s.Repository.RecipeShared(r.URL.String())
	if err != nil {
		notFoundHandler(w, r)
		return
	}

	userRecipeID := s.Repository.RecipeUser(share.RecipeID)
	recipe, err := s.Repository.Recipe(share.RecipeID, share.UserID)
	if err != nil {
		notFoundHandler(w, r)
		return
	}

	data := templates.Data{
		IsAuthenticated: isLoggedIn,
		Title:           recipe.Name,
		View:            templates.NewViewRecipeData(share.RecipeID, recipe, userRecipeID == share.UserID, true),
	}

	if r.Header.Get("Hx-Request") == "true" {
		templates.RenderComponent(w, "recipes", "view-recipe", data)
	} else {
		templates.Render(w, templates.ViewRecipePage, data)
	}
}

func (s *Server) recipeSharePostHandler(w http.ResponseWriter, r *http.Request) {
	recipeID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID := r.Context().Value("userID").(int64)
	share := models.ShareRecipe{RecipeID: recipeID, UserID: userID}

	link, err := s.Repository.AddShareLink(share)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Failed to create share link.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templates.RenderComponent(w, "recipes", "share-recipe", templates.Data{
		Content: r.Host + link,
	})
}

func (s *Server) recipesSupportedWebsitesHandler(w http.ResponseWriter, _ *http.Request) {
	websites := s.Repository.Websites()
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprintf(w, websites.TableHTML())
}

func (s *Server) recipesViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value("userID").(int64)
	recipe, err := s.Repository.Recipe(id, userID)
	if err != nil {
		notFoundHandler(w, r)
		return
	}

	data := templates.Data{
		IsAuthenticated: true,
		Title:           recipe.Name,
		View:            templates.NewViewRecipeData(id, recipe, true, false),
	}

	if r.Header.Get("Hx-Request") == "true" {
		templates.RenderComponent(w, "recipes", "view-recipe", data)
	} else {
		templates.Render(w, templates.ViewRecipePage, data)
	}
}
