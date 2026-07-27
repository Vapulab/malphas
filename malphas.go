// Package malphas is the root of the Malphas design system: the UI layer of
// the Vapulab stack.
package malphas

import (
	"embed"
	"io/fs"
)

//go:embed assets
var assetsFS embed.FS

// Assets returns the embedded design-system static assets, rooted at assets/
// (e.g. read "malphas.js", not "assets/ui.js").
func Assets() fs.FS {
	sub, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		panic(err)
	}
	return sub
}
