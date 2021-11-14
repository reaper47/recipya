package config

import "text/template"

var Tpl *template.Template

func init() {
	Tpl = template.Must(template.ParseGlob("templates/*gohtml"))
}
