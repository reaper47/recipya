package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/oxtoacart/bpool"
	"github.com/reaper47/recipya/views"
)

var templates map[string]*template.Template
var bufpool *bpool.BufferPool

func init() {
	templates = make(map[string]*template.Template)
	bufpool = bpool.NewBufferPool(64)

	tmpls, err := fs.ReadDir(views.FS, ".")
	if err != nil {
		panic(err)
	}

	for _, tmpl := range tmpls {
		if tmpl.IsDir() {
			continue
		}
		templates[tmpl.Name()] = template.Must(template.New("main").Funcs(fm).ParseFS(views.FS, tmpl.Name(), "layouts/*.gohtml"))
	}
}

// Render is a wrapper for template.ExecuteTemplate.
func Render(w http.ResponseWriter, name string, data any) error {
	tmpl, ok := templates[name]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", name)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := tmpl.Execute(buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
	return nil
}
