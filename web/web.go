package web

import "embed"

// FS is a read-only embedded file system containing static HTML assets.
// It includes CSS, images, and other static files needed by the application.
//
//go:embed templates/**/*.gohtml emails/transpiled/*.gohtml
var FS embed.FS
