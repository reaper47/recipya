// Package docs provides the read-only collection of files of the documentation website.
package docs

import "embed"

// FS is a read-only embedded file system containing the static documentation website.
//
//go:embed website/public
var FS embed.FS
