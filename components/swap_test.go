package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
)

func swapFixture(active bool) components.SwapProps {
	return components.SwapProps{
		On:    templ.Raw("<span>pause</span>"),
		Off:   templ.Raw("<span>play</span>"),
		Label: "Reproduzir",
		Active: active,
	}
}

// A div that toggles is invisible to assistive tech. Swap is a real button
// carrying its pressed state.
func TestSwapIsAToggleButton(t *testing.T) {
	html := renderToString(t, components.Swap(swapFixture(false)))
	if !strings.Contains(html, "<button") {
		t.Errorf("expected a button, got: %s", html)
	}
	if !strings.Contains(html, `aria-pressed="false"`) {
		t.Errorf("expected the initial pressed state, got: %s", html)
	}
	if !strings.Contains(html, "Reproduzir") {
		t.Errorf("expected an accessible label, got: %s", html)
	}
}

func TestSwapRendersBothFacesOnce(t *testing.T) {
	html := renderToString(t, components.Swap(swapFixture(false)))
	if !strings.Contains(html, "play") || !strings.Contains(html, "pause") {
		t.Errorf("both faces must be in the DOM so swapping needs no fetch, got: %s", html)
	}
}

func TestSwapActiveStartsPressed(t *testing.T) {
	html := renderToString(t, components.Swap(swapFixture(true)))
	if !strings.Contains(html, `aria-pressed="true"`) {
		t.Errorf("expected pressed, got: %s", html)
	}
}
