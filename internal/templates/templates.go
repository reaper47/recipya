package templates

import (
	"bytes"
	"github.com/reaper47/recipya/web"
	"html/template"
	"io/fs"
)

var templatesEmail map[string]*template.Template

func init() {
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
