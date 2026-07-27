package components_test

import (
	"strings"
	"testing"

	"github.com/darkanum/vapulab/malphas/components"
)

// Range is the one control in this slice that keeps the native input visible.
// Hiding it would mean reimplementing pointer drag, capture outside the
// element, touch and aria-valuenow — a lot of code for a worse result.
func TestRangeUsesTheNativeInput(t *testing.T) {
	html := renderToString(t, components.Range(components.RangeProps{
		ID: "vol", Name: "vol", Value: 40,
	}))
	if !strings.Contains(html, `type="range"`) {
		t.Errorf("expected a native range input, got: %s", html)
	}
	if strings.Contains(html, "sr-only") {
		t.Errorf("the range input must stay visible, got: %s", html)
	}
}

func TestRangeEmitsBoundsAndValue(t *testing.T) {
	html := renderToString(t, components.Range(components.RangeProps{
		ID: "vol", Name: "vol", Min: 10, Max: 20, Step: 2, Value: 14,
	}))
	for _, want := range []string{`min="10"`, `max="20"`, `step="2"`, `value="14"`, `name="vol"`} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %s, got: %s", want, html)
		}
	}
}

// Without an explicit max the browser assumes 100; emitting it keeps the
// rendered markup honest about the range being offered.
func TestRangeDefaultsMaxTo100(t *testing.T) {
	html := renderToString(t, components.Range(components.RangeProps{ID: "v", Name: "v"}))
	if !strings.Contains(html, `max="100"`) {
		t.Errorf("expected a default max of 100, got: %s", html)
	}
}

func TestRangeShowValueRendersCurrentValue(t *testing.T) {
	html := renderToString(t, components.Range(components.RangeProps{
		ID: "v", Name: "v", Value: 42, ShowValue: true,
	}))
	if !strings.Contains(html, "42") {
		t.Errorf("expected the value displayed, got: %s", html)
	}
	if !strings.Contains(html, "x-text") {
		t.Errorf("the displayed value must track the input, got: %s", html)
	}
}
