package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
)

func tooltipFixture() components.TooltipProps {
	return components.TooltipProps{ID: "tip-save", Text: "Salva sem fechar"}
}

// aria-describedby is what makes the tip reach a screen reader; hover styling
// alone reaches nobody who is not using a mouse.
func TestTooltipIsWiredWithDescribedby(t *testing.T) {
	html := renderWithChildren(t,
		components.Tooltip(tooltipFixture()),
		templ.Raw("<button>Salvar</button>"))
	if !strings.Contains(html, `aria-describedby="tip-save"`) {
		t.Errorf("expected aria-describedby on the trigger, got: %s", html)
	}
	if !strings.Contains(html, `id="tip-save"`) {
		t.Errorf("expected the tip to carry the referenced id, got: %s", html)
	}
	if !strings.Contains(html, `role="tooltip"`) {
		t.Errorf("expected role=tooltip, got: %s", html)
	}
}

// Hover-only tooltips are unreachable by keyboard, and a tooltip that cannot be
// dismissed traps a screen-magnifier user behind it.
func TestTooltipOpensOnFocusAndClosesOnEscape(t *testing.T) {
	html := renderWithChildren(t,
		components.Tooltip(tooltipFixture()),
		templ.Raw("<button>Salvar</button>"))
	if !strings.Contains(html, "focusin") {
		t.Errorf("tooltip must open on focus, not only hover, got: %s", html)
	}
	if !strings.Contains(html, "escape") {
		t.Errorf("tooltip must be dismissable with Escape, got: %s", html)
	}
}

func TestTooltipRendersTextAndChildren(t *testing.T) {
	html := renderWithChildren(t,
		components.Tooltip(tooltipFixture()),
		templ.Raw("<button>Salvar</button>"))
	if !strings.Contains(html, "Salva sem fechar") {
		t.Errorf("expected the tip text, got: %s", html)
	}
	if !strings.Contains(html, "<button>Salvar</button>") {
		t.Errorf("expected the trigger rendered, got: %s", html)
	}
}
