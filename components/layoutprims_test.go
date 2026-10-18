package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/vapulab/malphas/components"
)

func TestAspectRatioClasses(t *testing.T) {
	cases := map[string]string{
		"16/9":   "aspect-video",
		"1/1":    "aspect-square",
		"4/3":    "aspect-[4/3]",
		"":       "aspect-video", // default
		"banana": "aspect-video", // unknown falls back rather than collapsing to 0 height
	}
	for ratio, want := range cases {
		html := renderWithChildren(t,
			components.AspectRatio(components.AspectRatioProps{Ratio: ratio}), templ.Raw(`<img src="/x.png" alt="">`))
		if !strings.Contains(html, want) {
			t.Errorf("ratio %q: expected %s, got: %s", ratio, want, html)
		}
	}
}

// A keyboard user can only scroll a container that can take focus, and a
// focusable region with no name is announced as just "region".
func TestScrollAreaIsAFocusableNamedRegion(t *testing.T) {
	html := renderWithChildren(t,
		components.ScrollArea(components.ScrollAreaProps{Label: "Termos de uso", Class: "max-h-64"}),
		templ.Raw("<p>texto</p>"))
	for _, want := range []string{`role="region"`, `aria-label="Termos de uso"`, `tabindex="0"`, `max-h-64`, `overflow-y-auto`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}

func TestScrollAreaHorizontal(t *testing.T) {
	html := renderWithChildren(t,
		components.ScrollArea(components.ScrollAreaProps{Label: "Tags", Horizontal: true}), templ.Raw("x"))
	if !strings.Contains(html, "overflow-x-auto") || strings.Contains(html, "overflow-y-auto") {
		t.Errorf("expected horizontal scrolling only, got: %s", html)
	}
}

// Without a Label there is nothing to name the region with, and a focusable
// unnamed region is announced as just "region". It falls back to a plain
// scroll container, which browsers already make keyboard-scrollable.
func TestScrollAreaWithoutLabelClaimsNoRegion(t *testing.T) {
	html := renderWithChildren(t, components.ScrollArea(components.ScrollAreaProps{}), templ.Raw("x"))
	for _, unwanted := range []string{`role="region"`, `aria-label`, `tabindex`} {
		if strings.Contains(html, unwanted) {
			t.Errorf("no Label, so no %s; got: %s", unwanted, html)
		}
	}
}
