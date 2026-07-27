package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/darkanum/vapulab/malphas/components"
)

func TestIndicatorOverlaysLabelOnChildren(t *testing.T) {
	html := renderWithChildren(t,
		components.Indicator(components.IndicatorProps{Label: "9"}),
		templ.Raw(`<button>Caixa</button>`))
	if !strings.Contains(html, "<button>Caixa</button>") {
		t.Errorf("expected children rendered, got: %s", html)
	}
	if !strings.Contains(html, "9") {
		t.Errorf("expected the label, got: %s", html)
	}
	if !strings.Contains(html, "absolute") {
		t.Errorf("the label must be positioned over the child, got: %s", html)
	}
}

func TestIndicatorWithoutLabelRendersNoBadge(t *testing.T) {
	html := renderWithChildren(t,
		components.Indicator(components.IndicatorProps{}),
		templ.Raw(`<button>Caixa</button>`))
	if strings.Contains(html, "absolute") {
		t.Errorf("no label should mean no badge, got: %s", html)
	}
}
