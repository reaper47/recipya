package templates

import (
	"bytes"
	"fmt"
	"github.com/reaper47/recipya/web"
	"html/template"
	"io/fs"
	"net/http"
	"slices"
	"strings"
)

var (
	templates      map[string]*template.Template
	templatesEmail map[string]*template.Template
)

func init() {
	templates = make(map[string]*template.Template)

	err := fs.WalkDir(web.FS, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".gohtml") {
			authLayouts := []string{
				"templates/pages/forgot-password.gohtml",
				"templates/pages/forgot-password-reset.gohtml",
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
	})
	if err != nil {
		panic(err)
	}

	initEmailTemplates()
}

func initEmailTemplates() {
	templatesEmail = make(map[string]*template.Template)

	emailDir, err := fs.ReadDir(web.FS, "emails/transpiled")
	if err != nil {
		panic(err)
	}

	for _, entry := range emailDir {
		n := entry.Name()

		data, err := fs.ReadFile(web.FS, "emails/transpiled/"+n)
		if err != nil {
			panic(err)
		}
		data = bytes.ReplaceAll(data, []byte("[["), []byte("{{"))
		data = bytes.ReplaceAll(data, []byte("]]"), []byte("}}"))

		tmpl := template.Must(template.New(n).Parse(string(data)))
		if tmpl == nil || tmpl.Tree == nil || tmpl.Tree.Root == nil {
			panic("template or tree or root of " + entry.Name() + " is nil")
		}
		templatesEmail[n] = tmpl
	}
}

// Render renders a page to the response writer.
func Render(w http.ResponseWriter, page Page, data any) {
	err := templates["pages/"+page.String()+".gohtml"].Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RenderComponent renders a component to the response writer.
func RenderComponent(w http.ResponseWriter, component, name string, data any) {
	var buf bytes.Buffer
	_ = templates["components/"+component+".gohtml"].ExecuteTemplate(&buf, name, data)
	w.Header().Set("Content-Type", "text/html")
	_, _ = fmt.Fprint(w, buf.String())
}

// ToastWS returns the toast for use in conjunction with websockets for further processing.
func ToastWS() *template.Template {
	return templates["components/toast.gohtml"]
}
