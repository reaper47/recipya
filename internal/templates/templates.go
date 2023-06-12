package templates

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Boostport/mjml-go"
	"github.com/reaper47/recipya/web"
	"golang.org/x/exp/slices"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	templates      map[string]*template.Template
	templatesEmail map[string]*template.Template
	emailsFuncMap  = template.FuncMap{
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.ReplaceAll(template.HTMLEscapeString(text), "\n", "<br />"))
		},
	}
)

func init() {
	templates = make(map[string]*template.Template)

	if err := fs.WalkDir(web.FS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".gohtml") {
			authLayouts := []string{
				"templates/pages/login.gohtml",
				"templates/pages/register.gohtml",
				"templates/pages/simple.gohtml",
			}
			var layouts []string
			if slices.Contains(authLayouts, path) {
				layouts, _ = fs.Glob(web.FS, "templates/layouts/auth.gohtml")
			} else {
				layouts, _ = fs.Glob(web.FS, "templates/layouts/main.gohtml")
			}

			components, _ := fs.Glob(web.FS, "templates/components/*.gohtml")

			files := append(layouts, components...)
			files = append(files, path)

			templates[strings.TrimPrefix(path, "templates/")] = template.Must(template.ParseFS(web.FS, files...))
		}
		return nil
	}); err != nil {
		panic(err)
	}

	initEmailTemplates()
}

func initEmailTemplates() {
	templatesEmail = make(map[string]*template.Template)

	emailDir, err := fs.ReadDir(web.FS, "emails")
	if err != nil {
		panic(err)
	}

	for _, entry := range emailDir {
		n := entry.Name()

		if filepath.Ext(n) == ".mjml" {
			tmpl := template.Must(template.New(n).ParseFS(web.FS, "emails/"+n))

			html, err := mjml.ToHTML(context.Background(), tmpl.Tree.Root.String(), mjml.WithMinify(true))
			if err != nil {
				log.Fatal(err)
			}
			html = strings.ReplaceAll(html, "[[", "{{")
			html = strings.ReplaceAll(html, "]]", "}}")

			templatesEmail[n] = template.Must(template.New(n).Funcs(emailsFuncMap).Parse(html))
		}
	}
}

// Render renders a page to the response writer.
func Render(w http.ResponseWriter, page Page, data any) {
	if err := templates["pages/"+page.String()+".gohtml"].Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RenderComponent renders a component to the response writer.
func RenderComponent(w http.ResponseWriter, component, name string, data any) {
	div, err := fs.ReadFile(web.FS, "templates/components/"+component+".gohtml")
	if err != nil {
		return
	}

	t, err := template.New("component").Parse(string(div))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	_ = t.ExecuteTemplate(&buf, name, data)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, buf.String())
}
