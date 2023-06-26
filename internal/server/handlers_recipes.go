package server

import (
	"fmt"
	"github.com/reaper47/recipya/internal/app"
	"github.com/reaper47/recipya/internal/scraper"
	"github.com/reaper47/recipya/internal/templates"
	"net/http"
	"net/url"
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
