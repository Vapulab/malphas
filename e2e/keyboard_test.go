package e2e_test

import (
	"testing"

	"github.com/chromedp/chromedp/kb"
)

// --- Tabs ------------------------------------------------------------------

func TestTabsArrowMovesFocusSelectionAndPanel(t *testing.T) {
	ctx := openPage(t, "/components/tabs")
	focus(t, ctx, "#conta-tab-0")
	press(t, ctx, kb.ArrowRight)
	settle(t, ctx)

	if got := eval[string](t, ctx, `document.activeElement.id`); got != "conta-tab-1" {
		t.Errorf("focus did not follow the arrow: got %q", got)
	}
	if got := eval[string](t, ctx, `document.querySelector('[aria-selected="true"]').id`); got != "conta-tab-1" {
		t.Errorf("selection did not follow focus: got %q", got)
	}
	visible := eval[[]string](t, ctx,
		`[...document.querySelectorAll('[role="tabpanel"]')].filter(p => !p.hidden).map(p => p.id)`)
	if len(visible) != 1 || visible[0] != "conta-panel-1" {
		t.Errorf("wrong panel shown: %v", visible)
	}
}

// The third tab is disabled, so the arrow must wrap past it to the first.
func TestTabsArrowSkipsDisabledTab(t *testing.T) {
	ctx := openPage(t, "/components/tabs")
	focus(t, ctx, "#conta-tab-0")
	press(t, ctx, kb.ArrowRight)
	press(t, ctx, kb.ArrowRight)
	settle(t, ctx)

	if got := eval[string](t, ctx, `document.activeElement.id`); got != "conta-tab-0" {
		t.Errorf("expected to wrap past the disabled tab back to the first, got %q", got)
	}
}

// End goes to the last *enabled* tab, not simply the last one.
func TestTabsEndGoesToLastEnabledTab(t *testing.T) {
	ctx := openPage(t, "/components/tabs")
	focus(t, ctx, "#conta-tab-0")
	press(t, ctx, kb.End)
	settle(t, ctx)

	if got := eval[string](t, ctx, `document.activeElement.id`); got != "conta-tab-1" {
		t.Errorf("End landed on %q; the last tab is disabled, so it should stop at conta-tab-1", got)
	}
}

// --- Radio -----------------------------------------------------------------

// The point of the hidden-input design: an arrow key has to change what the
// form would actually submit, not only what is painted.
func TestRadioArrowChangesSubmittedValue(t *testing.T) {
	ctx := openPage(t, "/components/radio")
	const selected = `document.querySelector('[role="radiogroup"] input[type=radio]:checked').value`

	if got := eval[string](t, ctx, selected); got != "pro" {
		t.Fatalf("fixture changed: expected to start at pro, got %q", got)
	}

	focus(t, ctx, `[role="radiogroup"] [tabindex="0"]`)
	press(t, ctx, kb.ArrowLeft)
	settle(t, ctx)

	if got := eval[string](t, ctx, selected); got != "free" {
		t.Errorf("arrow did not move the real input: still %q", got)
	}
	if got := eval[int](t, ctx,
		`[...document.querySelectorAll('[role="radio"]')].findIndex(r => r.getAttribute('aria-checked') === 'true')`); got != 0 {
		t.Errorf("aria-checked did not follow: index %d", got)
	}
}

// --- Rating ----------------------------------------------------------------

func TestRatingArrowChangesSubmittedValue(t *testing.T) {
	ctx := openPage(t, "/components/rating")
	const value = `document.querySelector('input[name="nota"]').value`

	before := eval[string](t, ctx, value)
	focus(t, ctx, `[role="radiogroup"] [tabindex="0"]`)
	press(t, ctx, kb.ArrowLeft)
	settle(t, ctx)

	after := eval[string](t, ctx, value)
	if after == before {
		t.Errorf("arrow did not change the submitted rating: still %q", after)
	}
	if after != "2" {
		t.Errorf("expected the rating to step down from 3 to 2, got %q", after)
	}
}

// --- Tooltip ---------------------------------------------------------------

// Hover-only tooltips are unreachable by keyboard, and one that cannot be
// dismissed covers content for anyone using magnification.
func TestTooltipOpensOnFocusAndClosesOnEscape(t *testing.T) {
	ctx := openPage(t, "/components/tooltip")
	const visible = `getComputedStyle(document.getElementById('tip-teclado')).display !== 'none'`

	if eval[bool](t, ctx, visible) {
		t.Fatal("tooltip should start hidden")
	}

	focus(t, ctx, `[aria-describedby="tip-teclado"] button`)
	settle(t, ctx)
	if !eval[bool](t, ctx, visible) {
		t.Error("tooltip did not open on focus")
	}

	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if eval[bool](t, ctx, visible) {
		t.Error("Escape did not dismiss the tooltip")
	}
	// Escape must not steal focus, or the user loses their place.
	if got := eval[string](t, ctx, `document.activeElement.tagName`); got != "BUTTON" {
		t.Errorf("focus left the trigger after Escape: now on %s", got)
	}
}

// --- Accordion -------------------------------------------------------------

func TestAccordionExclusiveClosesTheOthers(t *testing.T) {
	ctx := openPage(t, "/components/accordion")
	const open = `[...document.querySelectorAll('[id^="faq-x-panel-"]')].filter(p => getComputedStyle(p).display !== 'none').length`

	click(t, ctx, `[aria-controls="faq-x-panel-0"]`)
	settle(t, ctx)
	if got := eval[int](t, ctx, open); got != 1 {
		t.Fatalf("expected 1 open panel, got %d", got)
	}

	click(t, ctx, `[aria-controls="faq-x-panel-1"]`)
	settle(t, ctx)
	if got := eval[int](t, ctx, open); got != 1 {
		t.Errorf("exclusive mode left %d panels open", got)
	}
	if got := eval[string](t, ctx,
		`document.querySelector('[aria-controls="faq-x-panel-0"]').getAttribute('aria-expanded')`); got != "false" {
		t.Errorf("the first trigger still reports aria-expanded=%q", got)
	}
}

// Non-exclusive is the default and must genuinely allow several open at once.
func TestAccordionIndependentKeepsPanelsOpen(t *testing.T) {
	ctx := openPage(t, "/components/accordion")
	const open = `[...document.querySelectorAll('[id^="faq-panel-"]')].filter(p => getComputedStyle(p).display !== 'none').length`

	click(t, ctx, `[aria-controls="faq-panel-0"]`)
	settle(t, ctx)
	if got := eval[int](t, ctx, open); got != 2 {
		t.Errorf("expected 2 panels open (one starts open), got %d", got)
	}
}

// --- Collapse --------------------------------------------------------------

// Regression guard for the x-show mechanism: hiding a panel with a class or the
// hidden attribute keeps applying after x-show "shows" it, so the panel opens
// once and never again. This bug shipped three times during development and no
// render test could see it.
func TestCollapseReopensAfterClosing(t *testing.T) {
	ctx := openPage(t, "/components/collapse")
	const visible = `getComputedStyle(document.getElementById('col-1')).display !== 'none'`

	click(t, ctx, `[aria-controls="col-1"]`)
	settle(t, ctx)
	if !eval[bool](t, ctx, visible) {
		t.Fatal("panel did not open")
	}

	click(t, ctx, `[aria-controls="col-1"]`)
	settle(t, ctx)
	if eval[bool](t, ctx, visible) {
		t.Fatal("panel did not close")
	}

	click(t, ctx, `[aria-controls="col-1"]`)
	settle(t, ctx)
	if !eval[bool](t, ctx, visible) {
		t.Error("panel never reopened — x-show is fighting a class or attribute again")
	}
}

// --- Checkbox --------------------------------------------------------------

// indeterminate is a JS-only property, so it can only be verified in a browser.
func TestCheckboxIndeterminateIsSetAndResolvesOnInteraction(t *testing.T) {
	ctx := openPage(t, "/components/checkbox")

	if !eval[bool](t, ctx, `document.getElementById('c-todos').indeterminate`) {
		t.Fatal("the indeterminate property was never applied")
	}
	if got := eval[string](t, ctx,
		`document.querySelector('label[for="c-todos"] [role="checkbox"]').getAttribute('aria-checked')`); got != "mixed" {
		t.Errorf("expected aria-checked=mixed, got %q", got)
	}

	click(t, ctx, `label[for="c-todos"]`)
	settle(t, ctx)

	if eval[bool](t, ctx, `document.getElementById('c-todos').indeterminate`) {
		t.Error("the third state survived a deliberate interaction, so the box lies about what it submits")
	}
}

// --- Countdown -------------------------------------------------------------

func TestCountdownTicksDown(t *testing.T) {
	ctx := openPage(t, "/components/countdown")
	const first = `document.querySelector('[role="timer"]').textContent.trim()`

	before := eval[string](t, ctx, first)
	sleep(t, ctx, 2200)
	after := eval[string](t, ctx, first)

	if before == after {
		t.Errorf("the clock never advanced: still %q", after)
	}
}

// --- Form integration ------------------------------------------------------

// The whole reason these controls carry a hidden native input. Interacts only
// with the visible parts, then reads what the browser would actually send.
func TestFormControlsSubmitTheirValues(t *testing.T) {
	ctx := openPage(t, "/components/toggle")

	click(t, ctx, `form label[for="f-toggle"]`)
	click(t, ctx, `form label[for="f-check"]`)
	settle(t, ctx)

	data := eval[map[string]string](t, ctx,
		`Object.fromEntries(new FormData(document.querySelector('form')).entries())`)

	for field, want := range map[string]string{
		"toggle": "on",
		"check":  "yes",
		"radio":  "b",
		"rating": "3",
	} {
		if data[field] != want {
			t.Errorf("form field %q = %q, want %q (full payload: %v)", field, data[field], want, data)
		}
	}
}
