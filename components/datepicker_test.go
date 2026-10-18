package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/vapulab/malphas/components"
)

// The platform picker is the point: on a phone it is the system's own wheel or
// calendar, which no custom widget matches for touch and accessibility.
func TestDatePickerIsNativeWithBounds(t *testing.T) {
	html := renderToString(t, components.DatePicker(components.DatePickerProps{
		ID: "entrega", Name: "entrega", Label: "Entrega", Value: "2026-10-20",
		Min: "2026-10-10", Max: "2026-12-31", Required: true,
		Attrs: templ.Attributes{"aria-describedby": "x"},
	}))
	for _, want := range []string{
		`type="date"`, `name="entrega"`, `value="2026-10-20"`,
		`min="2026-10-10"`, `max="2026-12-31"`, `for="entrega"`,
		`required`, `aria-describedby="x"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}

// An empty min would be an invalid bound; the attribute must simply be absent.
func TestDatePickerOmitsEmptyBounds(t *testing.T) {
	html := renderToString(t, components.DatePicker(components.DatePickerProps{ID: "d", Name: "d"}))
	if strings.Contains(html, "min=") || strings.Contains(html, "max=") {
		t.Errorf("no bounds given, none rendered: %s", html)
	}
}
