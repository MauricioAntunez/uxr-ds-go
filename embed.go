package ds

import (
	"embed"
	"io/fs"
)

//go:embed components
var componentsFS embed.FS

//go:embed css
var cssFS embed.FS

//go:embed js
var jsFS embed.FS

// Components returns the filesystem containing shared HTML component templates.
func Components() fs.FS {
	sub, _ := fs.Sub(componentsFS, "components") // err always nil for embedded dirs
	return sub
}

// CSS returns the filesystem containing the design system stylesheet.
func CSS() fs.FS {
	sub, _ := fs.Sub(cssFS, "css") // err always nil for embedded dirs
	return sub
}

// JS returns the filesystem containing the design system JavaScript.
func JS() fs.FS {
	sub, _ := fs.Sub(jsFS, "js") // err always nil for embedded dirs
	return sub
}
