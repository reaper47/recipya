package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/models"
	"github.com/reaper47/recipya/internal/scraper"
	"github.com/reaper47/recipya/internal/templates"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

func (s *Server) recipesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query == nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse query.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pageStr := query.Get("page")
	pageNumber, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		pageNumber = 1
	}

	userID := getUserID(r)
	p, err := newRecipesPagination(s, w, userID, pageNumber, false)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error updating pagination.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isHxRequest := r.Header.Get("HX-Request") == "true"

	baseData := templates.Data{
		About: templates.AboutData{
			Version: app.Version,
		},
		Functions:   templates.NewFunctionsData(),
		IsHxRequest: isHxRequest,
		Pagination:  p,
		Recipes:     s.Repository.Recipes(userID, pageNumber),
	}

	if isHxRequest {
		templates.RenderComponent(w, "recipes", "recipes-index", baseData)
		return
	}

	page := templates.HomePage
	baseData.IsAuthenticated = true
	baseData.Title = page.Title()
	templates.Render(w, page, baseData)
}

func newRecipesPagination(srv *Server, w http.ResponseWriter, userID int64, page uint64, isSwap bool) (templates.Pagination, error) {
	counts, err := srv.Repository.Counts(userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Error getting counts.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return templates.Pagination{}, err
	}

	numPages := counts.Recipes / templates.ResultsPerPage
	if numPages == 0 {
		numPages = 1
	}

	return templates.NewPagination(page, numPages, counts.Recipes, templates.ResultsPerPage, "/recipes", isSwap), nil
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

	userID := getUserID(r)
	_, found := s.Brokers[userID]
	if !found {
		w.Header().Set("HX-Trigger", makeToast("Connection lost. Please reload page.", warningToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.Brokers[userID].SendProgressStatus("Preparing...", true, 0, -1)

	err := r.ParseMultipartForm(128 << 20)
	form := r.MultipartForm
	if err != nil || form == nil || form.File == nil {
		s.Brokers[userID].HideNotification()
		w.Header().Set("HX-Trigger", makeToast("Could not parse the uploaded files.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	files, filesOk := form.File["files"]
	if !filesOk {
		s.Brokers[userID].HideNotification()
		w.Header().Set("HX-Trigger", makeToast("Could not retrieve the files or the directory from the form.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	settings, err := s.Repository.UserSettings(userID)
	if err != nil {
		s.Brokers[userID].HideNotification()
		w.Header().Set("HX-Trigger", makeToast("Error getting user settings.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go func() {
		var (
			count     atomic.Int64
			processed int
			progress  = make(chan models.Progress)
			recipes   = s.Files.ExtractRecipes(files)
			total     = len(recipes)
			recipeIDs = make([]int64, 0, total)
		)

		if total == 0 {
			s.Brokers[userID].HideNotification()
			s.Brokers[userID].SendToast("No recipes found.", "bg-orange-500")
			return
		}

		s.Brokers[userID].SendProgress(fmt.Sprintf("Importing 1/%d", total), 1, total)

		for i, recipe := range recipes {
			go func(index int, r models.Recipe) {
				defer func() {
					progress <- models.Progress{Total: total, Value: index}
				}()

				id, err := s.Repository.AddRecipe(&r, userID, settings)
				if err != nil {
					fmt.Println(err)
					return
				}

				recipeIDs = append(recipeIDs, id)
				count.Add(1)
			}(i, recipe)
		}

		for p := range progress {
			processed++
			s.Brokers[userID].SendProgress(fmt.Sprintf("Importing %d/%d", processed, p.Total), processed, p.Total)
			if processed == total {
				close(progress)
			}
		}

		s.Repository.CalculateNutrition(userID, recipeIDs, settings)
		s.Brokers[userID].HideNotification()
		s.Brokers[userID].SendToast(fmt.Sprintf("Imported %d recipes. %d skipped", count.Load(), int64(len(recipes))-count.Load()), "bg-blue-500")
	}()

	w.WriteHeader(http.StatusAccepted)
}

func (s *Server) recipeAddManualHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
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
	multipartForm := r.MultipartForm
	form := r.Form
	if err != nil || form == nil || multipartForm == nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var imageUUID uuid.UUID
	imageFile, ok := multipartForm.File["image"]
	if ok {
		f, err := imageFile[0].Open()
		if err != nil {
			w.Header().Set("HX-Trigger", makeToast("Could not open the image from the form.", errorToast))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer func() {
			_ = f.Close()
		}()

		imageUUID, err = s.Files.UploadImage(f)
		if err != nil {
			w.Header().Set("HX-Trigger", makeToast("Error uploading image.", errorToast))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	ingredients := make([]string, 0)
	i := 1
	for {
		key := fmt.Sprintf("ingredient-%d", i)
		if form.Has(key) {
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
		if form.Has(key) {
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

	userID := getUserID(r)
	settings, err := s.Repository.UserSettings(userID)
	if err != nil {
		log.Printf("recipeAddManualPostHandler.UserSettings error: %q", err)
		w.Header().Set("HX-Trigger", makeToast("Error getting user settings.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	recipeNumber, err := s.Repository.AddRecipe(recipe, userID, settings)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not add recipe.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	s.Repository.CalculateNutrition(userID, []int64{recipeNumber}, settings)
	w.Header().Set("HX-Redirect", "/recipes/"+strconv.FormatInt(recipeNumber, 10))
	w.WriteHeader(http.StatusCreated)
}

func recipeAddManualIngredientHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	form := r.Form
	if err != nil || form == nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	i := 1
	for {
		if !form.Has("ingredient-" + strconv.Itoa(i)) {
			break
		}

		i++
	}

	if form.Get(fmt.Sprintf("ingredient-%d", i-1)) == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	templates.RenderComponent(w, "recipes", "add-ingredient", i)
}

func recipeAddManualIngredientDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	form := r.Form
	if err != nil || form == nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	count := 0
	i := 1
	for {
		if !form.Has("ingredient-" + strconv.Itoa(i)) {
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
		if !form.Has(key) {
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
			`<label><input type="text" name="ingredient-` + currStr + `" placeholder="Ingredient #` + currStr + `" required class="w-8/12 py-1 pl-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400 dark:bg-gray-900 dark:border-none dark:text-gray-200" ` + `value="` + form.Get(key) + `" _="on keydown if event.key is 'Enter' then halt the event then get next <button/> from the parentElement of me then call htmx.trigger(it, 'click')"></label>`,
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

	_, _ = fmt.Fprint(w, sb.String())
}

func recipeAddManualInstructionHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	form := r.Form
	if err != nil || form == nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	i := 1
	for {
		if !form.Has("instruction-" + strconv.Itoa(i)) {
			break
		}

		i++
	}

	if form.Get(fmt.Sprintf("instruction-%d", i-1)) == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	templates.RenderComponent(w, "recipes", "add-instruction", i)
}

func recipeAddManualInstructionDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	form := r.Form
	if err != nil || form == nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	count := 0
	i := 1
	for {
		if !form.Has("instruction-" + strconv.Itoa(i)) {
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
		if !form.Has(key) {
			break
		}

		n, _ := strconv.Atoi(entry)
		if n == i {
			i++
			continue
		}

		currStr := strconv.Itoa(curr)
		value := form.Get(key)

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

	_, _ = fmt.Fprint(w, sb.String())
}

func (s *Server) recipesAddOCRHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<24)
	err := r.ParseMultipartForm(1 << 24)
	form := r.MultipartForm
	if err != nil || form == nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	images, ok := form.File["image"]
	if !ok {
		w.Header().Set("HX-Trigger", makeToast("Could not retrieve the image from the form.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	f, err := images[0].Open()
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not open the image from the form.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer func() {
		_ = f.Close()
	}()

	recipe, err := s.Integrations.ProcessImageOCR(f)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not process OCR.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := getUserID(r)
	settings, err := s.Repository.UserSettings(userID)
	if err != nil {
		log.Printf("recipesAddOCRHandler.UserSettings error: %q", err)
		w.Header().Set("HX-Trigger", makeToast("Error getting user settings.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := s.Repository.AddRecipe(&recipe, userID, settings)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Recipe could not be added.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.Repository.CalculateNutrition(userID, []int64{id}, settings)
	w.Header().Set("HX-Trigger", makeToast("Recipe scanned and uploaded.", infoToast))
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) recipesAddRequestWebsiteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[info] Website Request: Please support %s", r.FormValue("website"))
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

	userID := getUserID(r)
	settings, err := s.Repository.UserSettings(userID)
	if err != nil {
		log.Printf("recipesAddWebsiteHandler.UserSettings error: %q", err)
		w.Header().Set("HX-Trigger", makeToast("Error getting user settings.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := s.Repository.AddRecipe(recipe, userID, settings)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Recipe could not be added.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.Repository.CalculateNutrition(userID, []int64{id}, settings)
	w.Header().Set("HX-Redirect", "/recipes/"+strconv.FormatInt(id, 10))
	w.WriteHeader(http.StatusSeeOther)
}

func (s *Server) recipeDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID := getUserID(r)

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

	userID := getUserID(r)
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
	form := r.Form
	multipartForm := r.MultipartForm
	if err != nil || form == nil || multipartForm == nil {
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

	imageFile, ok := multipartForm.File["image"]
	if ok {
		f, err := imageFile[0].Open()
		if err != nil {
			w.Header().Set("HX-Trigger", makeToast("Could not open the image from the form.", errorToast))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer func() {
			_ = f.Close()
		}()

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
		if !form.Has(s) {
			break
		}
		updatedRecipe.Ingredients = append(updatedRecipe.Ingredients, r.FormValue(s))
		i++
	}

	i = 1
	for {
		s := "instruction-" + strconv.Itoa(i)
		if !form.Has(s) {
			break
		}
		updatedRecipe.Instructions = append(updatedRecipe.Instructions, r.FormValue(s))
		i++
	}

	yield, err := strconv.ParseInt(r.FormValue("yield"), 10, 16)
	if err == nil {
		updatedRecipe.Yield = int16(yield)
	}

	userID := getUserID(r)
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
	userID, isLoggedIn := s.findUserID(r)

	share, err := s.Repository.RecipeShared(r.URL.String())
	if err != nil {
		notFoundHandler(w, r)
		return
	}

	recipe, err := s.Repository.Recipe(share.RecipeID, share.UserID)
	if err != nil {
		notFoundHandler(w, r)
		return
	}

	data := templates.Data{
		IsAuthenticated: isLoggedIn,
		Title:           recipe.Name,
		View:            templates.NewViewRecipeData(share.RecipeID, recipe, userID == share.UserID, true),
	}

	if r.Header.Get("Hx-Request") == "true" {
		templates.RenderComponent(w, "recipes", "view-recipe", data)
	} else {
		templates.Render(w, templates.ViewRecipePage, data)
	}
}

func (s *Server) recipeScaleHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query == nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse query.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	yieldStr := query.Get("yield")
	yield, err := strconv.ParseInt(yieldStr, 10, 16)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("No yield in the query.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if yield <= 0 {
		w.Header().Set("HX-Trigger", makeToast("Yield must be greater than zero.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID := getUserID(r)

	recipe, err := s.Repository.Recipe(id, userID)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Recipe not found.", errorToast))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	recipe.Scale(int16(yield))

	templates.RenderComponent(w, "recipes", "ingredients-instructions", recipe)
}

func (s *Server) recipeSharePostHandler(w http.ResponseWriter, r *http.Request) {
	recipeID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID := getUserID(r)
	share := models.Share{CookbookID: -1, RecipeID: recipeID, UserID: userID}

	link, err := s.Repository.AddShareLink(share)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Failed to create share link.", errorToast))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	templates.RenderComponent(w, "recipes", "share-link", templates.Data{
		Content: r.Host + link,
	})
}

func (s *Server) recipesSearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query == nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse query.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pageStr := query.Get("page")
	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		page = 1
	}

	userID := getUserID(r)

	var recipes models.Recipes
	q := r.FormValue("q")
	if q == "" {
		recipes = s.Repository.Recipes(userID, page)
	} else {
		q = strings.ReplaceAll(q, ",", " ")
		q = strings.Join(strings.Fields(q), " ")

		var err error
		recipes, err = s.Repository.SearchRecipes(q, models.SearchOptionsRecipes{FullSearch: true}, userID)
		if err != nil {
			w.Header().Set("HX-Trigger", makeToast("Error searching recipes.", errorToast))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if len(recipes) == 0 {
		templates.RenderComponent(w, "search", "no-result", nil)
		return
	}

	templates.RenderComponent(w, "recipes", "list-recipes", templates.Data{
		Functions: templates.NewFunctionsData(),
		Recipes:   recipes,
	})
}

func (s *Server) recipesSupportedWebsitesHandler(w http.ResponseWriter, _ *http.Request) {
	websites := s.Repository.Websites()
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprint(w, websites.TableHTML())
}

func (s *Server) recipesViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID := getUserID(r)
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

func (s *Server) recipesViewShareHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()
	if query == nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse query.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cookbookIDStr := query.Get("cookbook")
	cookbookID, err := strconv.ParseInt(cookbookIDStr, 10, 64)
	if err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse cookbookID query parameter.", errorToast))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recipe, cookbookUserID, err := s.Repository.CookbookRecipe(id, cookbookID)
	if err != nil {
		notFoundHandler(w, r)
		return
	}

	userID, isLoggedIn := s.findUserID(r)

	data := templates.Data{
		IsAuthenticated: isLoggedIn,
		Title:           recipe.Name,
		View:            templates.NewViewRecipeData(id, recipe, cookbookUserID == userID, true),
	}

	if r.Header.Get("Hx-Request") == "true" {
		templates.RenderComponent(w, "recipes", "view-recipe", data)
	} else {
		templates.Render(w, templates.ViewRecipePage, data)
	}
}
