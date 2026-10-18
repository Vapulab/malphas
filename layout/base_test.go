package layout_test

import (
	"context"
	"io/fs"
	"strings"
	"testing"

	"github.com/vapulab/malphas"
	"github.com/vapulab/malphas/layout"
)

func renderBase(t *testing.T, props layout.BaseProps) string {
	t.Helper()
	var sb strings.Builder
	if err := layout.Base(props).Render(context.Background(), &sb); err != nil {
		t.Fatalf("render: %v", err)
	}
	return sb.String()
}

// ui.js registers Alpine.data() from an alpine:init listener. Alpine fires that
// event the moment it loads, so a listener attached afterwards never runs and
// every design-system behavior silently disappears. Order is load-bearing.
func TestBaseLoadsUIJSBeforeAlpine(t *testing.T) {
	html := renderBase(t, layout.BaseProps{Title: "t"})
	uiJS := strings.Index(html, "/malphas-assets/malphas.js")
	alpine := strings.Index(html, "alpine.min.js")
	if uiJS < 0 || alpine < 0 {
		t.Fatalf("expected both scripts in head, got: %s", html)
	}
	if uiJS > alpine {
		t.Fatal("ui.js must be loaded before alpine.min.js, otherwise its alpine:init listener never fires")
	}
}

// The theme script must run before the stylesheet is applied, otherwise the
// page paints in the wrong theme and visibly snaps to the right one.
func TestBaseAppliesThemeBeforeStylesheet(t *testing.T) {
	html := renderBase(t, layout.BaseProps{Title: "t"})
	script := strings.Index(html, "malphas-theme")
	css := strings.Index(html, `href="/static/app.css"`)
	if script < 0 {
		t.Fatalf("expected the inline theme script in head, got: %s", html)
	}
	if script > css {
		t.Fatal("the inline theme script must come before the stylesheet, otherwise the page flashes the wrong theme")
	}
}

// The inline script and ui.js must agree on the storage key, or a reload loses
// the user's choice. Nothing else couples these two files.
func TestBaseThemeKeyMatchesUIJS(t *testing.T) {
	html := renderBase(t, layout.BaseProps{Title: "t"})
	js, err := fs.ReadFile(malphas.Assets(), "malphas.js")
	if err != nil {
		t.Fatalf("read ui.js: %v", err)
	}
	const key = `"malphas-theme"`
	if !strings.Contains(html, key) || !strings.Contains(string(js), key) {
		t.Fatalf("storage key %s must appear in both the inline script and ui.js", key)
	}
}

// An undeclared icon makes every page load probe /favicon.ico and log a 404.
func TestBaseDeclaresAFavicon(t *testing.T) {
	html := renderBase(t, layout.BaseProps{Title: "t"})
	if !strings.Contains(html, `rel="icon"`) {
		t.Fatalf("expected a declared favicon, got: %s", html)
	}
}

func TestBaseRendersTitleAndStylesheet(t *testing.T) {
	html := renderBase(t, layout.BaseProps{Title: "Vapulab"})
	for _, want := range []string{"<!doctype html>", "<title>Vapulab</title>", `href="/static/app.css"`} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected %q in base layout, got: %s", want, html)
		}
	}
}
