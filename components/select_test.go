package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/vapulab/malphas/components"
)

func selectFixture() components.SelectProps {
	return components.SelectProps{
		ID: "plano", Name: "plano", Label: "Plano", Value: "pro",
		Options: []components.SelectOption{
			{Label: "Grátis", Value: "free"},
			{Label: "Pro", Value: "pro"},
			{Label: "Enterprise", Value: "ent", Disabled: true},
		},
	}
}

// The deliberate exception: a real <select>, so type-ahead, the platform
// picker and form submission come for free.
func TestSelectUsesTheNativeElement(t *testing.T) {
	html := renderToString(t, components.Select(selectFixture()))
	if !strings.Contains(html, "<select") {
		t.Errorf("expected a native select, got: %s", html)
	}
	if strings.Count(html, "<option") != 3 {
		t.Errorf("expected 3 options, got: %s", html)
	}
	if !strings.Contains(html, `name="plano"`) {
		t.Errorf("expected the field name, got: %s", html)
	}
}

func TestSelectMarksTheCurrentValue(t *testing.T) {
	html := renderToString(t, components.Select(selectFixture()))
	if strings.Count(html, "selected") != 1 {
		t.Errorf("expected exactly one selected option, got: %s", html)
	}
}

func TestSelectLabelIsBound(t *testing.T) {
	html := renderToString(t, components.Select(selectFixture()))
	if !strings.Contains(html, `for="plano"`) {
		t.Errorf("expected the label bound to the field, got: %s", html)
	}
}

// A placeholder that can be re-selected would let the user submit nothing by
// accident.
func TestSelectPlaceholderIsDisabled(t *testing.T) {
	props := selectFixture()
	props.Value = ""
	props.Placeholder = "Escolha…"
	html := renderToString(t, components.Select(props))
	if !strings.Contains(html, `<option value="" disabled`) {
		t.Errorf("expected a disabled placeholder option, got: %s", html)
	}
}

func TestSelectDisabledOption(t *testing.T) {
	html := renderToString(t, components.Select(selectFixture()))
	if !strings.Contains(html, `disabled`) {
		t.Errorf("expected the disabled option, got: %s", html)
	}
}
