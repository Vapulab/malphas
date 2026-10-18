// Package e2e serves the Malphas catalog to the browser tests, the way a host
// app would: the catalog under /components, the embedded assets under
// /malphas-assets and the stylesheet and vendors under /static.
package e2e

import (
	"net/http"

	"github.com/vapulab/malphas"
	"github.com/vapulab/malphas/showcase"
)

// NewServer serves the catalog. staticDir holds app.css, built by Tailwind
// from tailwind.config.js, and the vendored htmx and Alpine.
func NewServer(staticDir string) http.Handler {
	mux := http.NewServeMux()
	catalog := showcase.Handler(Frame)
	mux.Handle("/components", catalog)
	mux.Handle("/components/", catalog)
	mux.Handle("/malphas-assets/", http.StripPrefix("/malphas-assets", http.FileServer(http.FS(malphas.Assets()))))
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(staticDir))))
	mux.Handle("GET /{$}", http.RedirectHandler("/components", http.StatusFound))
	return mux
}
