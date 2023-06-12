package web

import "embed"

//go:embed templates/**/*.gohtml emails/*.mjml
var FS embed.FS
