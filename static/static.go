package static

import "embed"

//go:embed css img js favicon.ico robots.txt
var FS embed.FS
