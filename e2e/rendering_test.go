package e2e_test

import (
	"testing"

	"github.com/chromedp/chromedp/kb"
)

// These tests assert *computed style*, not attributes. Everything else in the
// suite checked ARIA and properties, which is why a checkbox whose mixed state
// rendered as a white dash on a white box passed every test while being
// invisible on screen.
//
// The underlying trap: Alpine's x-bind:class only removes classes it added, so
// a bound utility coexists with a conflicting static one and CSS picks a winner
// by stylesheet order. Attribute variants decide by specificity instead. Two of
// the three components below happened to win that ordering lottery; one lost.

func TestCheckboxMixedStateIsVisible(t *testing.T) {
	ctx := openPage(t, "/components/checkbox")

	const mixedBox = `label[for="c-todos"] [role="checkbox"]`
	const plainBox = `label[for="c-termos"] [role="checkbox"]`

	mixed := eval[string](t, ctx, `getComputedStyle(document.querySelector('`+mixedBox+`')).backgroundColor`)
	plain := eval[string](t, ctx, `getComputedStyle(document.querySelector('`+plainBox+`')).backgroundColor`)

	if mixed == plain {
		t.Errorf("the mixed checkbox is painted like an empty one (%s); its dash is invisible", mixed)
	}

	// The dash must contrast with whatever the box is painted.
	glyph := eval[string](t, ctx,
		`getComputedStyle(document.querySelector('`+mixedBox+` span')).color`)
	if glyph == mixed {
		t.Errorf("dash colour %s equals the box background — invisible", glyph)
	}
}

func TestRadioSelectedDotIsVisible(t *testing.T) {
	ctx := openPage(t, "/components/radio")

	on := eval[string](t, ctx,
		`getComputedStyle(document.querySelector('[role="radio"][aria-checked="true"] span')).opacity`)
	off := eval[string](t, ctx,
		`getComputedStyle(document.querySelector('[role="radio"][aria-checked="false"] span')).opacity`)

	if on != "1" {
		t.Errorf("selected radio dot has opacity %s, want 1", on)
	}
	if off != "0" {
		t.Errorf("unselected radio dot has opacity %s, want 0", off)
	}
}

// Interaction is where the class conflict actually bites: only after Alpine
// swaps state does the bound class join the static one.
func TestRatingColoursSurviveInteraction(t *testing.T) {
	ctx := openPage(t, "/components/rating")

	const starColour = `[...document.querySelectorAll('[role="radiogroup"]:first-of-type [role="radio"]')].map(e => getComputedStyle(e).color)`

	before := eval[[]string](t, ctx, starColour)
	if len(before) < 5 {
		t.Fatalf("expected 5 stars, got %d", len(before))
	}
	filled, empty := before[0], before[4]
	if filled == empty {
		t.Fatalf("filled and empty stars look identical (%s)", filled)
	}

	// Rate it 5: stars 4 and 5 must actually turn the filled colour.
	focus(t, ctx, `[role="radiogroup"] [tabindex="0"]`)
	press(t, ctx, kb.ArrowRight)
	press(t, ctx, kb.ArrowRight)
	// Longer than settle(): the stars carry transition-colors, and sampling
	// mid-transition reads an interpolated colour rather than the final one.
	sleep(t, ctx, 600)

	after := eval[[]string](t, ctx, starColour)
	for i, got := range after {
		if got != filled {
			t.Errorf("star %d is %s after rating 5; expected the filled colour %s", i+1, got, filled)
		}
	}
}
