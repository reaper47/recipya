package config

import "html/template"

var Tpl *template.Template

func init() {
	Tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	template.Must(Tpl.ParseGlob("templates/layout/*.gohtml"))
	template.Must(Tpl.ParseGlob("templates/partials/*.gohtml"))
}
