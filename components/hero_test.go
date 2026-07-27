package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/darkanum/vapulab/malphas/components"
)

func TestHeroRendersChildrenInASection(t *testing.T) {
	html := renderWithChildren(t,
		components.Hero(components.HeroProps{}),
		templ.Raw("<h1>Vapulab</h1>"))
	if !strings.Contains(html, "<section") {
		t.Errorf("expected a section landmark, got: %s", html)
	}
	if !strings.Contains(html, "<h1>Vapulab</h1>") {
		t.Errorf("expected children rendered, got: %s", html)
	}
}

func TestHeroAcceptsExtraClasses(t *testing.T) {
	html := renderWithChildren(t,
		components.Hero(components.HeroProps{Class: "bg-primary/5"}),
		templ.Raw("x"))
	if !strings.Contains(html, "bg-primary/5") {
		t.Errorf("expected the caller's classes, got: %s", html)
	}
}
