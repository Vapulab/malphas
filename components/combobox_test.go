package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

func comboboxFixture() components.ComboboxProps {
	return components.ComboboxProps{
		ID: "cidade", Name: "cidade", Label: "Cidade", Value: "sp",
		Options: []components.SelectOption{
			{Label: "Rio de Janeiro", Value: "rj"},
			{Label: "São Paulo", Value: "sp"},
			{Label: "Brasília", Value: "df", Disabled: true},
		},
	}
}

// The text field must not carry the name: free text would then be submitted
// as if it were a value.
func TestComboboxSubmitsThroughAHiddenInput(t *testing.T) {
	html := renderToString(t, components.Combobox(comboboxFixture()))
	if !strings.Contains(html, `<input type="hidden" name="cidade" value="sp"`) {
		t.Errorf("expected a hidden input carrying the value, got: %s", html)
	}
	if strings.Count(html, `name="cidade"`) != 1 {
		t.Errorf("only the hidden input may be named, got: %s", html)
	}
}

func TestComboboxWiresTheARIAPattern(t *testing.T) {
	html := renderToString(t, components.Combobox(comboboxFixture()))
	for _, want := range []string{
		`role="combobox"`,
		`aria-autocomplete="list"`,
		`aria-expanded="false"`,
		`aria-controls="cidade-listbox"`,
		`id="cidade-listbox" role="listbox"`,
		`id="cidade-opt-0" role="option"`,
		`for="cidade"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	if strings.Count(html, `role="option"`) != 3 {
		t.Errorf("expected 3 options, got: %s", html)
	}
}

// The server paints the selection before Alpine loads: the label in the field
// and aria-selected on the matching option.
func TestComboboxRendersTheSelection(t *testing.T) {
	html := renderToString(t, components.Combobox(comboboxFixture()))
	if !strings.Contains(html, `value="São Paulo"`) {
		t.Errorf("expected the selected label in the field, got: %s", html)
	}
	if strings.Count(html, `aria-selected="true"`) != 1 {
		t.Errorf("expected exactly one selected option, got: %s", html)
	}
}

// A value that matches no option, or a disabled one, is not a selection the
// user could have made; rendering it would submit something unpickable.
func TestComboboxIgnoresAnUnpickableValue(t *testing.T) {
	for _, value := range []string{"xx", "df"} {
		props := comboboxFixture()
		props.Value = value
		html := renderToString(t, components.Combobox(props))
		if !strings.Contains(html, `<input type="hidden" name="cidade" value=""`) {
			t.Errorf("value %q: expected an empty submission, got: %s", value, html)
		}
		if strings.Contains(html, `aria-selected="true"`) {
			t.Errorf("value %q: no option should be selected, got: %s", value, html)
		}
	}
}

// A disabled native control is left out of the submission; the hidden input
// has to follow, or a field the user could not edit still sends a value.
func TestComboboxDisabledSubmitsNothing(t *testing.T) {
	props := comboboxFixture()
	props.Disabled = true
	html := renderToString(t, components.Combobox(props))
	if !strings.Contains(html, `<input type="hidden" name="cidade" value="sp" disabled`) {
		t.Errorf("expected the hidden input disabled too, got: %s", html)
	}
}

func TestComboboxDisabledOption(t *testing.T) {
	html := renderToString(t, components.Combobox(comboboxFixture()))
	if strings.Count(html, `aria-disabled="true"`) != 1 {
		t.Errorf("expected the disabled option marked, got: %s", html)
	}
}

// Closed until Alpine opens it — otherwise the list flashes open on load.
func TestComboboxListboxStartsClosed(t *testing.T) {
	html := renderToString(t, components.Combobox(comboboxFixture()))
	i := strings.Index(html, `role="listbox"`)
	end := strings.Index(html[i:], ">")
	if !strings.Contains(html[i:i+end], `display:none`) {
		t.Errorf("expected the listbox hidden on first paint, got: %s", html[i:i+end])
	}
}

// An error is tied to the text field the way Field ties one to its control:
// announced on focus, and marked invalid.
func TestComboboxLinksItsError(t *testing.T) {
	p := comboboxFixture()
	p.Error = "Escolha uma cidade."
	html := renderToString(t, components.Combobox(p))
	for _, want := range []string{`aria-invalid="true"`, `aria-describedby="cidade-error"`, `id="cidade-error"`, "Escolha uma cidade."} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q in: %s", want, html)
		}
	}
	if html := renderToString(t, components.Combobox(comboboxFixture())); strings.Contains(html, "aria-invalid") || strings.Contains(html, "-error") {
		t.Error("no error, no invalid state")
	}
}
