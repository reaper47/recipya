package templates

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/oxtoacart/bpool"
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

	layouts, err := filepath.Glob("templates/layouts/*.gohtml")
	if err != nil {
		log.Fatalln(err)
	}

	includes, err := filepath.Glob("templates/*.gohtml")
	if err != nil {
		log.Fatalln(err)
	}

	const mainTmpl = `{{define "main" }} {{ template "base" . }} {{ end }}`
	mainTemplate := template.New("main")
	mainTemplate, err = mainTemplate.Parse(mainTmpl)
	if err != nil {
		log.Fatalln(err)
	}

	for _, layout := range layouts {
		files := append(includes, layout)
		templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
	}

	for _, file := range includes {
		fname := filepath.Base(file)
		files := append(layouts, file)
		templates[fname], err = mainTemplate.Clone()
		if err != nil {
			log.Fatalln(err)
		}
		templates[fname] = template.Must(templates[fname].ParseFiles(files...))
	}
}

// Render is a wrapper for template.ExecuteTemplate.
func Render(wr http.ResponseWriter, name string, data interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		err := fmt.Sprintf("The template %s does not exist.", name)
		http.Error(wr, err, http.StatusInternalServerError)
		return fmt.Errorf(err)
	}

	buf := bufpool.Get()
	defer bufpool.Put(buf)

	err := tmpl.Execute(buf, data)
	if err != nil {
		http.Error(wr, err.Error(), http.StatusInternalServerError)
		return err
	}

	wr.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(wr)
	return nil
}
