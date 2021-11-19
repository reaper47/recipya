package config

import (
	"html/template"
	"sync"
)

var (
	tpl  *template.Template
	once sync.Once
)

// Tpl stores a reference to the templates.
func Tpl() *template.Template {
	once.Do(initTemplates)
	return tpl
}

func initTemplates() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	template.Must(tpl.ParseGlob("templates/layout/*.gohtml"))
	template.Must(tpl.ParseGlob("templates/partials/*.gohtml"))
}
