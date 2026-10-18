package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/vapulab/malphas/components"
)

func TestRatingIsARadiogroupOfStars(t *testing.T) {
	html := renderToString(t, components.Rating(components.RatingProps{
		Name: "nota", Value: 3,
	}))
	if !strings.Contains(html, `role="radiogroup"`) {
		t.Errorf("expected role=radiogroup, got: %s", html)
	}
	// Default Max is 5.
	if strings.Count(html, `role="radio"`) != 5 {
		t.Errorf("expected 5 stars by default, got: %s", html)
	}
}

func TestRatingSubmitsAValue(t *testing.T) {
	html := renderToString(t, components.Rating(components.RatingProps{Name: "nota", Value: 3}))
	if !strings.Contains(html, `name="nota"`) {
		t.Errorf("expected the field name, got: %s", html)
	}
	if !strings.Contains(html, `value="3"`) {
		t.Errorf("expected the current value submitted, got: %s", html)
	}
}

func TestRatingRespectsMax(t *testing.T) {
	html := renderToString(t, components.Rating(components.RatingProps{Name: "n", Value: 2, Max: 3}))
	if strings.Count(html, `role="radio"`) != 3 {
		t.Errorf("expected 3 stars, got: %s", html)
	}
}

// Read-only ratings are output, not input: they must not be focusable and must
// not submit anything.
func TestRatingReadOnlyIsNotInteractive(t *testing.T) {
	html := renderToString(t, components.Rating(components.RatingProps{
		Name: "nota", Value: 4, ReadOnly: true,
	}))
	if strings.Contains(html, `tabindex="0"`) {
		t.Errorf("a read-only rating must not be focusable, got: %s", html)
	}
	if strings.Contains(html, "<input") {
		t.Errorf("a read-only rating must not submit a value, got: %s", html)
	}
	if !strings.Contains(html, "aria-label") {
		t.Errorf("a read-only rating still needs a text label, got: %s", html)
	}
}

func TestRatingInteractiveHasOneTabbableStar(t *testing.T) {
	html := renderToString(t, components.Rating(components.RatingProps{Name: "n", Value: 3}))
	if strings.Count(html, `tabindex="0"`) != 1 {
		t.Errorf("exactly one star may be tabbable, got: %s", html)
	}
}
