package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
)

func TestInputGroupPlacesAddonsAroundTheInput(t *testing.T) {
	html := renderToString(t, components.InputGroup(components.InputGroupProps{
		ID: "preco", Name: "preco", Label: "Preço",
		Prefix: components.InputGroupText("R$"),
		Suffix: components.InputGroupText(",00"),
		Attrs:  templ.Attributes{"inputmode": "numeric"},
	}))
	prefix := strings.Index(html, "R$")
	input := strings.Index(html, `<input`)
	suffix := strings.Index(html, ",00")
	if !(prefix >= 0 && prefix < input && input < suffix) {
		t.Errorf("expected prefix, input, suffix in order; got: %s", html)
	}
	// Attrs land on the input, as in Input, so Field.ControlAttrs() fits.
	if !strings.Contains(html, `inputmode="numeric"`) || !strings.Contains(html, `type="text"`) {
		t.Errorf("expected attrs and default type on the input, got: %s", html)
	}
	if !strings.Contains(html, `for="preco"`) {
		t.Errorf("expected the label bound, got: %s", html)
	}
}

func TestInputGroupWithoutAddons(t *testing.T) {
	html := renderToString(t, components.InputGroup(components.InputGroupProps{ID: "q", Name: "q"}))
	if strings.Contains(html, "border-r") || strings.Contains(html, "border-l") {
		t.Errorf("no addons, no addon separators: %s", html)
	}
}
