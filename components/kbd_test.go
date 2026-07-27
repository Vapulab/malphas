package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/darkanum/vapulab/malphas/components"
)

func TestKbdRendersKbdElement(t *testing.T) {
	html := renderWithChildren(t, components.Kbd(components.KbdProps{}), templ.Raw("Ctrl"))
	if !strings.Contains(html, "<kbd") {
		t.Errorf("expected a <kbd> element, got: %s", html)
	}
	if !strings.Contains(html, "Ctrl") {
		t.Errorf("expected children rendered, got: %s", html)
	}
}
