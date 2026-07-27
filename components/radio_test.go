package components_test

import (
	"strings"
	"testing"

	"github.com/darkanum/vapulab/malphas/components"
)

func radioFixture() components.RadioProps {
	return components.RadioProps{
		Name:  "plano",
		Value: "pro",
		Options: []components.RadioOption{
			{Label: "Grátis", Value: "free"},
			{Label: "Pro", Value: "pro"},
			{Label: "Enterprise", Value: "ent", Disabled: true},
		},
	}
}

func TestRadioIsARadiogroup(t *testing.T) {
	html := renderToString(t, components.Radio(radioFixture()))
	if !strings.Contains(html, `role="radiogroup"`) {
		t.Errorf("expected role=radiogroup, got: %s", html)
	}
	if strings.Count(html, `role="radio"`) != 3 {
		t.Errorf("expected 3 radio roles, got: %s", html)
	}
}

// Every option needs a real input sharing the group name, or the form submits
// nothing.
func TestRadioSubmitsThroughHiddenInputs(t *testing.T) {
	html := renderToString(t, components.Radio(radioFixture()))
	if strings.Count(html, `type="radio"`) != 3 {
		t.Errorf("expected one real radio per option, got: %s", html)
	}
	if strings.Count(html, `name="plano"`) != 3 {
		t.Errorf("all options must share the group name, got: %s", html)
	}
	if !strings.Contains(html, `value="pro"`) {
		t.Errorf("expected the option values, got: %s", html)
	}
}

// Roving tabindex: exactly one option is reachable with Tab, the arrows do the
// rest. The selected option is the one that gets it.
func TestRadioHasExactlyOneTabbableOption(t *testing.T) {
	html := renderToString(t, components.Radio(radioFixture()))
	if strings.Count(html, `tabindex="0"`) != 1 {
		t.Errorf("exactly one option may be tabbable, got: %s", html)
	}
	if strings.Count(html, `tabindex="-1"`) != 2 {
		t.Errorf("the other options must be removed from the tab order, got: %s", html)
	}
}

func TestRadioMarksSelectedOption(t *testing.T) {
	html := renderToString(t, components.Radio(radioFixture()))
	if !strings.Contains(html, `aria-checked="true"`) {
		t.Errorf("expected the selected option marked, got: %s", html)
	}
	if strings.Count(html, `aria-checked="true"`) != 1 {
		t.Errorf("only one option may be checked, got: %s", html)
	}
}

// A disabled option must be skipped by the arrow keys, which rovingList decides
// from this attribute.
func TestRadioDisabledOptionOptsOutOfRoving(t *testing.T) {
	html := renderToString(t, components.Radio(radioFixture()))
	if !strings.Contains(html, "data-roving-disabled") {
		t.Errorf("expected the disabled option to opt out of arrow navigation, got: %s", html)
	}
}
