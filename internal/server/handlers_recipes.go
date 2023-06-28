package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/scraper"
	"github.com/reaper47/recipya/internal/templates"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

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

	if err := r.ParseMultipartForm(128 << 20); err != nil {
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
		if err := s.Repository.AddRecipe(&r, userID); err != nil {
			continue
		}
		count += 1
	}

	msg := fmt.Sprintf("Imported %d recipes. %d failed", count, len(recipes)-count)
	w.Header().Set("HX-Trigger", makeToast(msg, infoToast))
	w.WriteHeader(http.StatusCreated)
}

func recipeAddManualHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Hx-Request") == "true" {
		templates.RenderComponent(w, "recipes", "add-recipe-manual", nil)
	} else {
		page := templates.AddRecipeManualPage
		templates.Render(w, page, templates.Data{
			IsAuthenticated: true,
			Title:           page.Title(),
		})
	}
}

func recipeAddManualIngredientHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	i := 1
	num := strconv.Itoa(i)
	for {
		if !r.Form.Has("ingredient-" + num) {
			break
		}

		i++
		num = strconv.Itoa(i)
	}

	if r.Form.Get(fmt.Sprintf("ingredient-%d", i-1)) == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var sb strings.Builder
	xs := []string{
		`<li class="pb-2 pl-2">`,
		`<label><input autofocus type="text" name="ingredient-` + num + `" placeholder="Ingredient #` + num + `" required class="w-8/12 py-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400" onkeydown="handleKeyDownIngredient(event)"></label>`,
		`&nbsp;<button type="button" class="w-10 h-10 text-center duration-300 bg-green-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-green-600 hover:text-white center" title="Shortcut: Enter" hx-post="/recipes/add/manual/ingredient" hx-target="#ingredients-list" hx-swap="beforeend" hx-include="[name^='ingredient']">+</button>`,
		`&nbsp;<button type="button" class="w-10 h-10 duration-300 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/` + num + `" hx-include="[name^='ingredient']">-</button>`,
		`&nbsp;<div class="inline-block h-4 cursor-move handle"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/></svg></div>`,
	}
	for _, x := range xs {
		sb.WriteString(x)
	}
	fmt.Fprintf(w, sb.String())
}

func recipeAddManualIngredientDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
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
			`<li class="pb-2 pl-2">`,
			`<label><input type="text" name="ingredient-` + currStr + `" placeholder="Ingredient #` + currStr + `" required class="w-8/12 py-1 text-gray-600 placeholder-gray-400 bg-white border border-gray-400" ` + `value="` + r.Form.Get(key) + `"></label>`,
			`&nbsp;<button type="button" class="w-10 h-10 text-center duration-300 bg-green-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-green-600 hover:text-white center" title="Shortcut: Enter" hx-post="/recipes/add/manual/ingredient" hx-target="#ingredients-list" hx-swap="beforeend" hx-include="[name^='ingredient']">+</button>`,
			`&nbsp;<button type="button" class="w-10 h-10 duration-300 bg-red-300 border border-gray-800 rounded-lg md:w-7 md:h-7 hover:bg-red-600 hover:text-white center" hx-target="#ingredients-list" hx-post="/recipes/add/manual/ingredient/` + currStr + `" hx-include="[name^='ingredient']">-</button>`,
			`&nbsp;<div class="inline-block h-4 cursor-move handle"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/></svg></div>`,
		}
		for _, x := range xs {
			sb.WriteString(x)
		}

		i++
		curr++
	}

	fmt.Fprintf(w, sb.String())
}

func recipeAddManualInstructionHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Header().Set("HX-Trigger", makeToast("Could not parse the form.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	i := 1
	num := fmt.Sprintf("%d", i)
	for {
		if !r.Form.Has("instruction-" + num) {
			break
		}

		i++
		num = fmt.Sprintf("%d", i)
	}

	if r.Form.Get(fmt.Sprintf("instruction-%d", i-1)) == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var sb strings.Builder
	xs := []string{
		`<li class="pb-2 pl-2 md:pl-0">`,
		`<label><textarea autofocus required name="instruction-` + num + `" rows="3" class="w-9/12 border border-gray-300 md:w-5/6 xl:w-11/12" placeholder="Instruction #` + num + `" onkeydown="handleKeyDownInstruction(event)"></textarea>&nbsp;</label>`,
		`<div class="inline-flex flex-col-reverse">`,
		`<button type="button" class="mt-4 md:flex-initial w-10 h-10 right-0.5 md:w-7 md:h-7 md:right-auto duration-300 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/` + num + `" hx-include="[name^='instruction']">-</button>`,
		`<button type="button" class="md:flex-initial bottom-0 right-0.5 md:w-7 md:h-7 md:right-auto w-10 h-10 text-center duration-300 bg-green-300 border border-gray-800 rounded-lg hover:bg-green-600 hover:text-white center" title="Shortcut: CTRL + Enter" hx-post="/recipes/add/manual/instruction" hx-target="#instructions-list" hx-swap="beforeend" hx-include="[name^='instruction']">+</button>`,
		`</div>&nbsp;<div class="inline-block h-4 cursor-move handle"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/></svg></div>`,
	}
	for _, x := range xs {
		sb.WriteString(x)
	}
	fmt.Fprintf(w, sb.String())
}

func recipeAddManualInstructionDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
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
			`<li class="pb-2 pl-2 md:pl-0">`,
			`<label><textarea required name="instruction-` + currStr + `" rows="3" class="w-9/12 border border-gray-300 md:w-5/6 xl:w-11/12" placeholder="Instruction #` + currStr + `" onkeydown="handleKeyDownInstruction(event)">` + value + `</textarea>&nbsp;</label>`,
			`<div class="inline-flex flex-col-reverse">`,
			`<button type="button" class="mt-4 md:flex-initial w-10 h-10 right-0.5 md:w-7 md:h-7 md:right-auto duration-300 bg-red-300 border border-gray-800 rounded-lg top-3 hover:bg-red-600 hover:text-white center" hx-target="#instructions-list" hx-post="/recipes/add/manual/instruction/` + currStr + `" hx-include="[name^='instruction']">-</button>`,
			`<button type="button" class="md:flex-initial bottom-0 right-0.5 md:w-7 md:h-7 md:right-auto w-10 h-10 text-center duration-300 bg-green-300 border border-gray-800 rounded-lg hover:bg-green-600 hover:text-white center" title="Shortcut: CTRL + Enter" hx-post="/recipes/add/manual/instruction" hx-target="#instructions-list" hx-swap="beforeend" hx-include="[name^='instruction']">+</button>`,
			`</div>&nbsp;<div class="inline-block h-4 cursor-move handle"><svg xmlns="http://www.w3.org/2000/svg" class="md:w-4 md:h-4 w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/></svg></div>`,
		}
		for _, x := range xs {
			sb.WriteString(x)
		}

		i++
		curr++
	}

	fmt.Fprintf(w, sb.String())
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
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if _, err := url.ParseRequestURI(rawURL); err != nil {
		w.Header().Set("HX-Trigger", makeToast("Invalid URI.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	rs, err := scraper.Scrape(rawURL)
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
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := s.Repository.AddRecipe(recipe, r.Context().Value("userID").(int64)); err != nil {
		w.Header().Set("HX-Trigger", makeToast("Recipe could not be added.", errorToast))
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusSeeOther)
}

func (s *Server) recipesSupportedWebsitesHandler(w http.ResponseWriter, _ *http.Request) {
	websites := s.Repository.Websites()
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprintf(w, websites.TableHTML())
}

func (s *Server) recipesSupportedWebsitesPostHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("search")
	if query == "" {
		s.recipesSupportedWebsitesHandler(w, r)
		return
	}

	websites := s.Repository.WebsitesSearch(query)
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprintf(w, websites.TableHTML())
}
