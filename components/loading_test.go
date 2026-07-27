package components_test

import (
	"strings"
	"testing"

	"github.com/darkanum/vapulab/malphas/components"
)

// A spinner that only animates is invisible to a screen reader; it needs a
// status role and a text label.
func TestLoadingIsAnnouncedToScreenReaders(t *testing.T) {
	html := renderToString(t, components.Loading(components.LoadingProps{}))
	if !strings.Contains(html, `role="status"`) {
		t.Errorf("expected role=status, got: %s", html)
	}
	if !strings.Contains(html, "sr-only") {
		t.Errorf("expected a visually hidden label, got: %s", html)
	}
}

func TestLoadingKinds(t *testing.T) {
	for _, kind := range []string{"spinner", "dots", "bars"} {
		html := renderToString(t, components.Loading(components.LoadingProps{Kind: kind}))
		if !strings.Contains(html, "animate-") {
			t.Errorf("kind %q renders nothing animated: %s", kind, html)
		}
	}
}

func TestLoadingCustomLabel(t *testing.T) {
	html := renderToString(t, components.Loading(components.LoadingProps{Label: "Salvando"}))
	if !strings.Contains(html, "Salvando") {
		t.Errorf("expected the custom label, got: %s", html)
	}
}
