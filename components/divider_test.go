package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

func TestDividerWithLabel(t *testing.T) {
	html := renderToString(t, components.Divider(components.DividerProps{Label: "or continue with"}))
	if !strings.Contains(html, "or continue with") {
		t.Errorf("expected the label, got: %s", html)
	}
}

// A plain rule carries no meaning for a screen reader, so it must not be
// announced; a labelled one is a real separator.
func TestDividerWithoutLabelIsHiddenFromScreenReaders(t *testing.T) {
	html := renderToString(t, components.Divider(components.DividerProps{}))
	if !strings.Contains(html, "aria-hidden") {
		t.Errorf("an unlabelled divider should be aria-hidden, got: %s", html)
	}
}

func TestDividerVertical(t *testing.T) {
	html := renderToString(t, components.Divider(components.DividerProps{Vertical: true}))
	if !strings.Contains(html, "w-px") {
		t.Errorf("vertical divider should be a 1px-wide column, got: %s", html)
	}
}
