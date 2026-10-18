package e2e_test

import (
	"testing"

	"github.com/chromedp/chromedp/kb"
)

// Modality is the one thing in this design system that cannot be verified by
// looking at markup: whether focus can leave, whether the page still scrolls,
// whether a screen reader can still reach the content behind. All three only
// exist at runtime.
//
// The valuable assertion here is the negative one — focus must NOT escape a
// modal overlay, and must escape a non-modal one. The two patterns are mirror
// images and easy to swap by accident.

func TestModalTrapsFocus(t *testing.T) {
	ctx := openPage(t, "/components/modal")

	click(t, ctx, `[x-data] button`)
	settle(t, ctx)

	if !eval[bool](t, ctx, `!!document.querySelector('[role="dialog"]')`) {
		t.Fatal("the dialog did not open")
	}

	// Tab many more times than there are focusable elements: if the trap is
	// missing, focus walks out into the page behind and never comes back.
	for i := 0; i < 12; i++ {
		press(t, ctx, kb.Tab)
	}
	settle(t, ctx)

	inside := eval[bool](t, ctx,
		`document.querySelector('[role="dialog"]').contains(document.activeElement)`)
	if !inside {
		where := eval[string](t, ctx, `document.activeElement.outerHTML.slice(0, 80)`)
		t.Errorf("focus escaped the dialog after 12 tabs; it landed on: %s", where)
	}
}

func TestModalLocksScrollAndInertsTheBackground(t *testing.T) {
	ctx := openPage(t, "/components/modal")

	before := eval[string](t, ctx, `getComputedStyle(document.documentElement).overflow`)

	click(t, ctx, `[x-data] button`)
	settle(t, ctx)

	if got := eval[string](t, ctx, `getComputedStyle(document.documentElement).overflow`); got != "hidden" {
		t.Errorf("page still scrolls behind the dialog: overflow=%s", got)
	}
	// Everything outside the overlay must be unreachable, including to a
	// screen reader's virtual cursor.
	if !eval[bool](t, ctx, `
		[...document.body.children]
			.filter(el => !el.contains(document.querySelector('[role="dialog"]')))
			.every(el => el.hasAttribute('inert'))`) {
		t.Error("content behind the dialog is not inert; a screen reader still reads the whole page")
	}

	press(t, ctx, kb.Escape)
	settle(t, ctx)

	if got := eval[string](t, ctx, `getComputedStyle(document.documentElement).overflow`); got != before {
		t.Errorf("scroll was not restored after closing: overflow=%s, want %s", got, before)
	}
	if eval[int](t, ctx, `[...document.body.children].filter(el => el.hasAttribute('inert')).length`) != 0 {
		t.Error("inert was left behind after closing; the page is now permanently unreachable")
	}
}

// Being dumped at the top of the document on every dismissal is what makes
// keyboard navigation exhausting.
func TestModalReturnsFocusToItsTrigger(t *testing.T) {
	ctx := openPage(t, "/components/modal")

	focus(t, ctx, `[x-data] button`)
	trigger := eval[string](t, ctx, `document.activeElement.textContent.trim()`)

	press(t, ctx, kb.Enter)
	settle(t, ctx)

	// Move focus *into* the dialog first. Without this the assertion is
	// vacuous: focus never left the trigger, so it trivially "returns" there
	// even with no restore logic at all.
	press(t, ctx, kb.Tab)
	settle(t, ctx)
	if eval[bool](t, ctx, `document.activeElement === document.querySelector('[x-data] button')`) {
		t.Fatal("focus never entered the dialog; this test would prove nothing")
	}

	press(t, ctx, kb.Escape)
	settle(t, ctx)

	back := eval[string](t, ctx, `document.activeElement.textContent.trim()`)
	if back != trigger {
		t.Errorf("focus returned to %q, expected the trigger %q", back, trigger)
	}
}
