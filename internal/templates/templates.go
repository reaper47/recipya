package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/oxtoacart/bpool"
	"github.com/reaper47/recipya/views"
)

var templates map[string]*template.Template
var bufpool *bpool.BufferPool

func init() {
	bufpool = bpool.NewBufferPool(64)
}

// Load loads the templates.
func Load() {
	templates = make(map[string]*template.Template)
	bufpool = bpool.NewBufferPool(64)

	mainTemplate := template.New("main")
	mainTemplate, err := mainTemplate.Parse(`{{define "main" }}{{ template "base" . }}{{ end }}`)
	if err != nil {
		log.Fatalln(err)
	}

	var layouts, includes []string
	fs.WalkDir(views.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if strings.HasPrefix(path, "layouts/") {
			layouts = append(layouts, path)
		} else {
			includes = append(includes, filepath.Base(path))
		}
		return nil
	})

	for _, layout := range layouts {
		files := append(includes, layout)
		templates[filepath.Base(layout)] = template.Must(template.New("").Funcs(fm).ParseFS(views.FS, files...))
	}

	for _, file := range includes {
		fname := filepath.Base(file)
		files := append(layouts, file)
		templates[fname], err = mainTemplate.Clone()
		if err != nil {
			log.Fatalln(err)
		}
		templates[fname] = template.Must(templates[fname].Funcs(fm).ParseFS(views.FS, files...))
	}
}

// Render is a wrapper for template.ExecuteTemplate.
func Render(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		err := fmt.Sprintf("The template %s does not exist.", name)
		http.Error(w, err, http.StatusInternalServerError)
		return fmt.Errorf(err)
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
