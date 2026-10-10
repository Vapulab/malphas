package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/darkanum/vapulab/malphas/components"
)

func TestToggleButtonRendersPressedState(t *testing.T) {
	html := renderWithChildren(t,
		components.ToggleButton(components.ToggleButtonProps{ID: "b", Pressed: true, Label: "Negrito"}),
		templ.Raw("B"))
	for _, want := range []string{`type="button"`, `aria-pressed="true"`, `aria-label="Negrito"`, `x-data="toggleButton(true)"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}

// The value travels only while pressed, like a checked checkbox: a disabled
// hidden input is left out of the submission.
func TestToggleButtonSubmitsOnlyWhenPressed(t *testing.T) {
	off := renderWithChildren(t,
		components.ToggleButton(components.ToggleButtonProps{Name: "fav"}), templ.Raw("★"))
	if !strings.Contains(off, `<input type="hidden" name="fav" value="on" disabled`) {
		t.Errorf("unpressed: expected a disabled hidden input, got: %s", off)
	}
	on := renderWithChildren(t,
		components.ToggleButton(components.ToggleButtonProps{Name: "fav", Pressed: true}), templ.Raw("★"))
	if strings.Contains(on, `value="on" disabled`) {
		t.Errorf("pressed: the hidden input must be enabled, got: %s", on)
	}
}

// A disabled control submits nothing, pressed or not.
func TestToggleButtonDisabledSubmitsNothing(t *testing.T) {
	html := renderWithChildren(t,
		components.ToggleButton(components.ToggleButtonProps{Name: "fav", Pressed: true, Disabled: true}), templ.Raw("★"))
	if !strings.Contains(html, `value="on" disabled`) {
		t.Errorf("expected the hidden input disabled, got: %s", html)
	}
	if !strings.Contains(html, `x-bind:disabled="!pressed || true"`) {
		t.Errorf("Alpine must not re-enable it, got: %s", html)
	}
}
