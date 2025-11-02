package static

import (
	"embed"
	"io/fs"
)

// EmbedFS embeds the frontend static files built from the web directory.
//
//go:embed dist
var embedFS embed.FS

// GetFS returns the embedded filesystem containing the built frontend files.
// The files are located under the "dist" subdirectory.
func GetFS() (fs.FS, error) {
	return fs.Sub(embedFS, "dist")
}

// GetRawFS returns the raw embedded filesystem (for direct access).
func GetRawFS() embed.FS {
	return embedFS
}
