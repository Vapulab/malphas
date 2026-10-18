package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/vapulab/malphas/components"
)

func TestAlertUsesVariantColorAndAlertRole(t *testing.T) {
	html := renderWithChildren(t,
		components.Alert(components.AlertProps{Variant: "warning"}),
		templ.Raw("Cuidado"))
	if !strings.Contains(html, `role="alert"`) {
		t.Errorf("alert must carry role=alert for screen readers, got: %s", html)
	}
	if !strings.Contains(html, "warning") {
		t.Errorf("expected the warning color, got: %s", html)
	}
	if !strings.Contains(html, "Cuidado") {
		t.Errorf("expected children rendered, got: %s", html)
	}
}

// "error" is the public name; danger is the token behind it.
func TestAlertErrorMapsToDangerToken(t *testing.T) {
	html := renderWithChildren(t,
		components.Alert(components.AlertProps{Variant: "error"}),
		templ.Raw("x"))
	if !strings.Contains(html, "danger") {
		t.Errorf("error variant should use the danger token, got: %s", html)
	}
}

func TestAlertDefaultsToInfo(t *testing.T) {
	html := renderWithChildren(t, components.Alert(components.AlertProps{}), templ.Raw("x"))
	if !strings.Contains(html, "info") {
		t.Errorf("expected info as the default variant, got: %s", html)
	}
}
