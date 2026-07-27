package components_test

import (
	"strings"
	"testing"

	"github.com/darkanum/vapulab/malphas/components"
)

// The visual is a styled div, so the switch role and state have to be declared
// explicitly — nothing here is a native checkbox from the browser's point of
// view except the hidden input.
func TestToggleHasSwitchSemantics(t *testing.T) {
	html := renderToString(t, components.Toggle(components.ToggleProps{
		ID: "notif", Name: "notif", Label: "Notificações",
	}))
	if !strings.Contains(html, `role="switch"`) {
		t.Errorf("expected role=switch, got: %s", html)
	}
	if !strings.Contains(html, "aria-checked") {
		t.Errorf("expected aria-checked, got: %s", html)
	}
}

// Without the hidden input the control looks right and submits nothing.
func TestToggleSubmitsThroughAHiddenInput(t *testing.T) {
	html := renderToString(t, components.Toggle(components.ToggleProps{
		ID: "notif", Name: "notif", Label: "Notificações",
	}))
	if !strings.Contains(html, `type="checkbox"`) {
		t.Errorf("expected a real checkbox behind the visual, got: %s", html)
	}
	if !strings.Contains(html, `name="notif"`) {
		t.Errorf("expected the field name on the input, got: %s", html)
	}
	if !strings.Contains(html, "sr-only") {
		t.Errorf("the real input should be visually hidden, not display:none, got: %s", html)
	}
}

func TestToggleCheckedInitialState(t *testing.T) {
	html := renderToString(t, components.Toggle(components.ToggleProps{
		ID: "x", Name: "x", Label: "x", Checked: true,
	}))
	if !strings.Contains(html, "checked") {
		t.Errorf("expected the input to start checked, got: %s", html)
	}
}

func TestToggleLabelIsBoundToInput(t *testing.T) {
	html := renderToString(t, components.Toggle(components.ToggleProps{
		ID: "notif", Name: "notif", Label: "Notificações",
	}))
	if !strings.Contains(html, `for="notif"`) {
		t.Errorf("expected the label bound to the input, got: %s", html)
	}
}
