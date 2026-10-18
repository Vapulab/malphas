package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/vapulab/malphas/components"
)

// The trigger must be a button: a clickable div is not focusable, not
// activatable with Enter or Space, and announced as nothing.
func TestCollapseTriggerIsAButtonWithState(t *testing.T) {
	html := renderWithChildren(t,
		components.Collapse(components.CollapseProps{ID: "c1", Title: "Detalhes"}),
		templ.Raw("<p>corpo</p>"))
	if !strings.Contains(html, "<button") {
		t.Errorf("expected a button trigger, got: %s", html)
	}
	if !strings.Contains(html, `aria-expanded="false"`) {
		t.Errorf("expected the initial collapsed state, got: %s", html)
	}
	if !strings.Contains(html, `aria-controls="c1"`) {
		t.Errorf("expected the trigger wired to the panel, got: %s", html)
	}
	if !strings.Contains(html, `id="c1"`) {
		t.Errorf("expected the panel to carry the referenced id, got: %s", html)
	}
}

func TestCollapseOpenStartsExpanded(t *testing.T) {
	html := renderWithChildren(t,
		components.Collapse(components.CollapseProps{ID: "c1", Title: "Detalhes", Open: true}),
		templ.Raw("<p>corpo</p>"))
	if !strings.Contains(html, `aria-expanded="true"`) {
		t.Errorf("expected expanded, got: %s", html)
	}
}

func TestCollapseRendersTitleAndChildren(t *testing.T) {
	html := renderWithChildren(t,
		components.Collapse(components.CollapseProps{ID: "c1", Title: "Detalhes"}),
		templ.Raw("<p>corpo</p>"))
	if !strings.Contains(html, "Detalhes") || !strings.Contains(html, "<p>corpo</p>") {
		t.Errorf("expected title and children, got: %s", html)
	}
}
