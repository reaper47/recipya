package templates

import (
	"io"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	template.Must(tpl.ParseGlob("templates/layout/*.gohtml"))
	template.Must(tpl.ParseGlob("templates/partials/*.gohtml"))
}

// ExecuteTemplates is a wrapper for template.ExecuteTemplate.
//
// This avoids the creation of an extra function in this file, i.e. Tpl(), to
// retrieve the `tpl` variable, and then calling templates.Tpl().ExecuteTemplate
// from a different package to execute a template.
func ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return tpl.ExecuteTemplate(wr, name, data)
}
