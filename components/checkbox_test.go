package components_test

import (
	"strings"
	"testing"

	"github.com/darkanum/vapulab/malphas/components"
)

func TestCheckboxSubmitsThroughAHiddenInput(t *testing.T) {
	html := renderToString(t, components.Checkbox(components.CheckboxProps{
		ID: "aceito", Name: "aceito", Value: "sim", Label: "Aceito os termos",
	}))
	if !strings.Contains(html, `type="checkbox"`) {
		t.Errorf("expected a real checkbox, got: %s", html)
	}
	if !strings.Contains(html, `name="aceito"`) || !strings.Contains(html, `value="sim"`) {
		t.Errorf("expected name and value on the input, got: %s", html)
	}
	if !strings.Contains(html, "Aceito os termos") {
		t.Errorf("expected the label text, got: %s", html)
	}
}

// Indeterminate cannot be expressed in HTML — it is a JS-only property — so it
// has to be applied on init, and announced separately.
func TestCheckboxIndeterminateIsAppliedAndAnnounced(t *testing.T) {
	html := renderToString(t, components.Checkbox(components.CheckboxProps{
		ID: "todos", Name: "todos", Label: "Todos", Indeterminate: true,
	}))
	if !strings.Contains(html, "indeterminate") {
		t.Errorf("expected the indeterminate property to be set, got: %s", html)
	}
	if !strings.Contains(html, `aria-checked="mixed"`) {
		t.Errorf("expected aria-checked=mixed, got: %s", html)
	}
}

func TestCheckboxDisabled(t *testing.T) {
	html := renderToString(t, components.Checkbox(components.CheckboxProps{
		ID: "x", Name: "x", Label: "x", Disabled: true,
	}))
	if !strings.Contains(html, "disabled") {
		t.Errorf("expected the input disabled, got: %s", html)
	}
}
