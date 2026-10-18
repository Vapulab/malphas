package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/vapulab/malphas/components"
)

func TestResizableSeparatorARIA(t *testing.T) {
	html := renderToString(t, components.Resizable(components.ResizableProps{
		ID: "rz", Label: "Redimensionar editor", DefaultSize: 30, Min: 20, Max: 70,
		First: templ.Raw("<p>lista</p>"), Second: templ.Raw("<p>editor</p>"),
	}))
	for _, want := range []string{
		`role="separator"`,
		`tabindex="0"`,
		`aria-orientation="vertical"`, // side-by-side panels: the bar stands upright
		`aria-valuenow="30"`,
		`aria-valuemin="20"`,
		`aria-valuemax="70"`,
		`aria-controls="rz-first"`,
		`aria-label="Redimensionar editor"`,
		`id="rz-first"`,
		`style="flex-basis: 30%;"`,
		`x-data="resizable(&#34;rz&#34;, 30, 20, 70, false)"`,
		"lista", "editor",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}

func TestResizableDefaultsAndClamp(t *testing.T) {
	html := renderToString(t, components.Resizable(components.ResizableProps{ID: "a", First: templ.Raw(""), Second: templ.Raw("")}))
	for _, want := range []string{`aria-valuenow="50"`, `aria-valuemin="10"`, `aria-valuemax="90"`} {
		if !strings.Contains(html, want) {
			t.Errorf("defaults: missing %s in: %s", want, html)
		}
	}
	html = renderToString(t, components.Resizable(components.ResizableProps{ID: "b", DefaultSize: 95, Max: 80, First: templ.Raw(""), Second: templ.Raw("")}))
	if !strings.Contains(html, `aria-valuenow="80"`) {
		t.Errorf("a default past Max should be clamped: %s", html)
	}
	html = renderToString(t, components.Resizable(components.ResizableProps{ID: "c", Vertical: true, First: templ.Raw(""), Second: templ.Raw("")}))
	if !strings.Contains(html, `aria-orientation="horizontal"`) || !strings.Contains(html, "flex-col") {
		t.Errorf("stacked panels: horizontal bar, column layout: %s", html)
	}
}

// Without an ID there is no id to point at: two ID-less instances would both
// claim "-first" and aria-controls would cross over.
func TestResizableWithoutIDHasNoDuplicateID(t *testing.T) {
	html := renderToString(t, components.Resizable(components.ResizableProps{First: templ.Raw(""), Second: templ.Raw("")}))
	if strings.Contains(html, `id="-first"`) || strings.Contains(html, "aria-controls") {
		t.Errorf("no ID, no id or aria-controls: %s", html)
	}
}

// A 4px bar is too small to grab on a phone; a pseudo-element widens the hit
// area to 24px without changing the look.
func TestResizableHasATouchSizedHitArea(t *testing.T) {
	html := renderToString(t, components.Resizable(components.ResizableProps{ID: "r", First: templ.Raw(""), Second: templ.Raw("")}))
	if !strings.Contains(html, "after:-inset-x-2.5") {
		t.Errorf("expected a widened hit area on the upright bar: %s", html)
	}
}

// Bounds stay inside the component: Max 100 would push the bar out of reach,
// and Min above 100 is nonsense.
func TestResizableBoundsStayReachable(t *testing.T) {
	html := renderToString(t, components.Resizable(components.ResizableProps{ID: "r", Min: 150, Max: 100, DefaultSize: 120, First: templ.Raw(""), Second: templ.Raw("")}))
	if !strings.Contains(html, `aria-valuemax="95"`) || !strings.Contains(html, `aria-valuenow="95"`) {
		t.Errorf("bounds should clamp to at most 95: %s", html)
	}
}
