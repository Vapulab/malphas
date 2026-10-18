// Package e2e serves the Malphas catalog to the browser tests, the way a host
// app would: the catalog under /components, the embedded assets under
// /malphas-assets and the stylesheet and vendors under /static.
package e2e

import (
	"net/http"
	"time"

	"github.com/vapulab/malphas"
	"github.com/vapulab/malphas/components"
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
	mux.HandleFunc("GET /fixtures/calendar", calendarFixture)
	return mux
}

// calendarFixture renders a lone Calendar in the language named by ?locale=,
// for the tests of what changes with the language: the text and the week's
// first day. The catalog itself is always English.
func calendarFixture(w http.ResponseWriter, r *http.Request) {
	l, _ := malphas.ParseLocale(r.URL.Query().Get("locale"))
	ctx := showcase.WithQuery(malphas.WithLocale(r.Context(), l), r)
	today := time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC)
	cal := components.Calendar(components.CalendarProps{
		ID: "fx-cal", Name: "day", Today: today,
		Month:  showcase.QueryDate(ctx, "month", today),
		Focus:  showcase.QueryDate(ctx, "focus", time.Time{}),
		NavURL: "/fixtures/calendar?locale=" + l.Segment(),
	})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = Frame("Calendar", false, cal).Render(ctx, w)
}
