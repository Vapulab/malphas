package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/darkanum/vapulab/malphas/components"
)

// Join's whole job is squaring off the inner corners so the children read as
// one control; without that it is just a flex row.
func TestJoinSquaresInnerCorners(t *testing.T) {
	html := renderWithChildren(t,
		components.Join(components.JoinProps{}),
		templ.Raw(`<button>a</button><button>b</button>`))
	if !strings.Contains(html, "rounded-none") {
		t.Errorf("expected inner corners squared, got: %s", html)
	}
	if !strings.Contains(html, "first-child") || !strings.Contains(html, "last-child") {
		t.Errorf("expected the outer corners kept round, got: %s", html)
	}
}

func TestJoinVerticalStacks(t *testing.T) {
	html := renderWithChildren(t,
		components.Join(components.JoinProps{Vertical: true}),
		templ.Raw(`<button>a</button>`))
	if !strings.Contains(html, "flex-col") {
		t.Errorf("expected a column layout, got: %s", html)
	}
}
