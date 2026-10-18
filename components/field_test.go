package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
)

func TestLabelMarksRequiredForSightedUsersOnly(t *testing.T) {
	html := renderWithChildren(t,
		components.Label(components.LabelProps{For: "email", Required: true}),
		templ.Raw("Email"))
	if !strings.Contains(html, `for="email"`) {
		t.Errorf("expected the label bound, got: %s", html)
	}
	// The control's own required attribute is what assistive tech announces;
	// reading the asterisk too would say it twice.
	if !strings.Contains(html, `<span aria-hidden="true" class="ml-0.5 text-danger">*</span>`) {
		t.Errorf("expected an aria-hidden asterisk, got: %s", html)
	}
}

func TestFieldRendersDescriptionAndErrorWithIDs(t *testing.T) {
	f := components.FieldProps{ID: "email", Label: "Email", Description: "Used to sign in.", Error: "Invalid email."}
	html := renderWithChildren(t, components.Field(f), templ.Raw(`<input id="email">`))
	for _, want := range []string{
		`for="email"`,
		`id="email-description"`,
		`id="email-error"`,
		`Invalid email.`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}

func TestFieldControlAttrsReferenceOnlyRenderedIDs(t *testing.T) {
	cases := []struct {
		name        string
		props       components.FieldProps
		describedBy string
		invalid     bool
	}{
		{"nothing", components.FieldProps{ID: "x"}, "", false},
		{"description", components.FieldProps{ID: "x", Description: "d"}, "x-description", false},
		{"error only", components.FieldProps{ID: "x", Error: "e"}, "x-error", true},
		{"both", components.FieldProps{ID: "x", Description: "d", Error: "e"}, "x-description x-error", true},
	}
	for _, c := range cases {
		attrs := c.props.ControlAttrs()
		got, has := attrs["aria-describedby"]
		if c.describedBy == "" && has {
			t.Errorf("%s: aria-describedby=%v points at nothing", c.name, got)
		}
		if c.describedBy != "" && got != c.describedBy {
			t.Errorf("%s: aria-describedby=%v, want %q", c.name, got, c.describedBy)
		}
		if _, has := attrs["aria-invalid"]; has != c.invalid {
			t.Errorf("%s: aria-invalid present=%v, want %v", c.name, has, c.invalid)
		}
	}
}

func TestFieldControlAttrsCarryRequired(t *testing.T) {
	attrs := components.FieldProps{ID: "x", Required: true}.ControlAttrs()
	if attrs["required"] != true {
		t.Errorf("expected required on the control, got %v", attrs)
	}
}
