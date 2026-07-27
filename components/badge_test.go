package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/darkanum/vapulab/malphas/components"
)

func TestBadgeSolidUsesVariantBackground(t *testing.T) {
	html := renderWithChildren(t,
		components.Badge(components.BadgeProps{Variant: "success"}),
		templ.Raw("Ativo"))
	if !strings.Contains(html, "bg-success") {
		t.Errorf("expected a solid success background, got: %s", html)
	}
	if !strings.Contains(html, "Ativo") {
		t.Errorf("expected children rendered, got: %s", html)
	}
}

// Outline must not paint a solid background, or the two styles look identical.
func TestBadgeOutlineHasNoSolidBackground(t *testing.T) {
	html := renderWithChildren(t,
		components.Badge(components.BadgeProps{Variant: "success", Outline: true}),
		templ.Raw("Ativo"))
	if strings.Contains(html, "bg-success ") || strings.HasSuffix(html, "bg-success") {
		t.Errorf("outline badge must not use a solid background, got: %s", html)
	}
	if !strings.Contains(html, "border-success") {
		t.Errorf("expected a success border, got: %s", html)
	}
}
