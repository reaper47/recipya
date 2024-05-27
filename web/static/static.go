package static

import "embed"

// FS is a read-only embedded file system containing static web assets.
// It includes CSS, images, and other static files needed by the application.
//
//go:embed css img favicon.ico robots.txt *.png site.webmanifest
var FS embed.FS
