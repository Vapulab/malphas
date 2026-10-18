package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

func TestProgressUsesNativeElementWithValue(t *testing.T) {
	html := renderToString(t, components.Progress(components.ProgressProps{Value: 40}))
	if !strings.Contains(html, "<progress") {
		t.Errorf("expected the native <progress> element, got: %s", html)
	}
	if !strings.Contains(html, `value="40"`) {
		t.Errorf("expected value=40, got: %s", html)
	}
	// Max must always be emitted, otherwise the browser assumes 1 and a value
	// of 40 renders as a full bar.
	if !strings.Contains(html, `max="100"`) {
		t.Errorf("expected a default max of 100, got: %s", html)
	}
}

func TestProgressRespectsCustomMax(t *testing.T) {
	html := renderToString(t, components.Progress(components.ProgressProps{Value: 3, Max: 5}))
	if !strings.Contains(html, `max="5"`) {
		t.Errorf("expected max=5, got: %s", html)
	}
}

// An indeterminate <progress> is one with no value attribute at all.
func TestProgressIndeterminateOmitsValue(t *testing.T) {
	html := renderToString(t, components.Progress(components.ProgressProps{Indeterminate: true}))
	if strings.Contains(html, "value=") {
		t.Errorf("indeterminate progress must omit value, got: %s", html)
	}
}
