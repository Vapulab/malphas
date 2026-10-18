package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/chromedp/chromedp/kb"
)

// --- Drawer: the modal half of the pair ------------------------------------

func TestDrawerIsModal(t *testing.T) {
	ctx := openPage(t, "/components/drawer")

	click(t, ctx, `[x-data] button`)
	settle(t, ctx)
	if !eval[bool](t, ctx, `!!document.querySelector('[role="dialog"]')`) {
		t.Fatal("the drawer did not open")
	}

	for i := 0; i < 12; i++ {
		press(t, ctx, kb.Tab)
	}
	settle(t, ctx)
	if !eval[bool](t, ctx, `document.querySelector('[role="dialog"]').contains(document.activeElement)`) {
		where := eval[string](t, ctx, `document.activeElement.outerHTML.slice(0, 80)`)
		t.Errorf("focus escaped the drawer; it landed on: %s", where)
	}

	if got := eval[string](t, ctx, `getComputedStyle(document.documentElement).overflow`); got != "hidden" {
		t.Errorf("page still scrolls behind the drawer: overflow=%s", got)
	}
	if !eval[bool](t, ctx, `
		[...document.body.children]
			.filter(el => !el.contains(document.querySelector('[role="dialog"]')))
			.every(el => el.hasAttribute('inert'))`) {
		t.Error("content behind the drawer is not inert")
	}
}

// --- Dropdown: the non-modal half ------------------------------------------
//
// The mirror assertion, and the one worth having: focus must NOT be trapped.
// Getting this backwards is easy, and a trapped menu is broken in a way that
// looks fine in a screenshot.

func TestDropdownArrowsMoveBetweenItems(t *testing.T) {
	ctx := openPage(t, "/components/dropdown")

	click(t, ctx, `#dd-acoes-menu, [aria-controls="dd-acoes-menu"]`)
	settle(t, ctx)
	if !eval[bool](t, ctx, `!!document.querySelector('[role="menu"]')`) {
		t.Fatal("the menu did not open")
	}

	first := eval[string](t, ctx, `document.activeElement.textContent.trim()`)
	press(t, ctx, kb.ArrowDown)
	settle(t, ctx)
	second := eval[string](t, ctx, `document.activeElement.textContent.trim()`)

	if first == second {
		t.Errorf("ArrowDown did not move between items; still on %q", first)
	}
	if !eval[bool](t, ctx, `document.querySelector('[role="menu"]').contains(document.activeElement)`) {
		t.Error("arrow navigation left the menu")
	}
}

func TestDropdownDoesNotTrapFocus(t *testing.T) {
	ctx := openPage(t, "/components/dropdown")

	click(t, ctx, `[aria-controls="dd-acoes-menu"]`)
	settle(t, ctx)

	press(t, ctx, kb.Tab)
	settle(t, ctx)

	// Tab must both leave the menu and close it. A menu that holds focus is
	// as wrong as a dialog that does not.
	if eval[bool](t, ctx, `!!document.querySelector('[role="menu"]')`) {
		t.Error("Tab did not close the menu")
	}
	if eval[bool](t, ctx, `getComputedStyle(document.documentElement).overflow === 'hidden'`) {
		t.Error("a dropdown must not lock page scroll — it is not modal")
	}
	if eval[int](t, ctx, `[...document.body.children].filter(el => el.hasAttribute('inert')).length`) != 0 {
		t.Error("a dropdown must not make the page inert")
	}
}

func TestDropdownEscapeReturnsFocusToTrigger(t *testing.T) {
	ctx := openPage(t, "/components/dropdown")

	focus(t, ctx, `[aria-controls="dd-acoes-menu"]`)
	trigger := eval[string](t, ctx, `document.activeElement.textContent.trim()`)

	press(t, ctx, kb.Enter)
	settle(t, ctx)
	if eval[bool](t, ctx, `document.activeElement === document.querySelector('[aria-controls="dd-acoes-menu"]')`) {
		t.Fatal("focus never entered the menu; this test would prove nothing")
	}

	press(t, ctx, kb.Escape)
	settle(t, ctx)
	back := eval[string](t, ctx, `document.activeElement.textContent.trim()`)
	if !strings.Contains(back, strings.TrimSpace(strings.Split(trigger, "\n")[0])) {
		t.Errorf("focus returned to %q, expected the trigger %q", back, trigger)
	}
}

// --- OTP --------------------------------------------------------------------

func TestOTPAdvancesBackspacesAndPastes(t *testing.T) {
	ctx := openPage(t, "/components/otp")

	const boxes = `document.querySelectorAll('[data-otp-index]')`
	const hidden = `document.querySelector('input[name="codigo"]').value`

	focus(t, ctx, `[data-otp-index="0"]`)
	press(t, ctx, "1")
	press(t, ctx, "2")
	settle(t, ctx)

	if got := eval[string](t, ctx, `document.activeElement.dataset.otpIndex`); got != "2" {
		t.Errorf("typing should advance to box 2, focus is on %q", got)
	}
	if got := eval[string](t, ctx, hidden); got != "12" {
		t.Errorf("hidden field is %q, want \"12\"", got)
	}

	// Backspace on an empty box steps back *and* clears, so one press deletes
	// one digit — the whole row behaves like a single field.
	press(t, ctx, kb.Backspace)
	settle(t, ctx)
	if got := eval[string](t, ctx, hidden); got != "1" {
		t.Errorf("one backspace should delete one digit; hidden field is %q, want \"1\"", got)
	}
	press(t, ctx, kb.Backspace)
	settle(t, ctx)
	if got := eval[string](t, ctx, hidden); got != "" {
		t.Errorf("a second backspace should clear the last digit; hidden field is %q", got)
	}

	// Pasting the whole code is how most people enter it.
	eval[bool](t, ctx, `(() => {
		const box = document.querySelector('[data-otp-index="0"]');
		box.focus();
		const dt = new DataTransfer();
		dt.setData('text', '987654');
		box.dispatchEvent(new ClipboardEvent('paste', { clipboardData: dt, bubbles: true, cancelable: true }));
		return true;
	})()`)
	settle(t, ctx)
	if got := eval[string](t, ctx, hidden); got != "987654" {
		t.Errorf("paste filled %q, want \"987654\"", got)
	}
}

// --- Select -----------------------------------------------------------------

func TestSelectSubmitsItsValue(t *testing.T) {
	ctx := openPage(t, "/components/select")

	value := eval[string](t, ctx, `document.querySelector('select[name="plano"]').value`)
	if value != "pro" {
		t.Errorf("expected the server-rendered selection, got %q", value)
	}
	// It is a real select, so the browser's own form serialisation applies.
	if !eval[bool](t, ctx, `document.querySelector('select[name="plano"]') instanceof HTMLSelectElement`) {
		t.Error("expected a native select element")
	}
}

// --- Combobox ---------------------------------------------------------------
//
// Focus stays in the text field throughout; the listbox is only pointed at.
// Every assertion that matters reads either the hidden input — what the form
// sends — or aria-activedescendant — what a screen reader announces.

const (
	cbEstado     = `#cb-estado`
	cbEstadoSent = `document.querySelector('input[name="estado"]').value`
	cbPlano      = `#cb-plano`
	cbPlanoSent  = `document.querySelector('input[name="plano"]').value`
)

// retype replaces the field's text the way a user would: select it all, type.
func retype(t *testing.T, ctx context.Context, selector, text string) {
	t.Helper()
	focus(t, ctx, selector)
	eval[bool](t, ctx, `(document.querySelector('`+selector+`').select(), true)`)
	press(t, ctx, text)
	settle(t, ctx)
}

func TestComboboxFiltersIgnoringAccentsAndEnterPicks(t *testing.T) {
	ctx := openPage(t, "/components/combobox")

	retype(t, ctx, cbEstado, "para")
	visible := eval[[]string](t, ctx,
		`[...document.querySelectorAll('#cb-estado-listbox [role="option"]')].filter(o => o.style.display !== 'none').map(o => o.dataset.label)`)
	want := []string{"Pará", "Paraíba", "Paraná"}
	if strings.Join(visible, ",") != strings.Join(want, ",") {
		t.Fatalf("filter for \"para\" showed %v, want %v", visible, want)
	}

	press(t, ctx, kb.Enter)
	settle(t, ctx)
	if got := eval[string](t, ctx, cbEstadoSent); got != "pa" {
		t.Errorf("Enter should pick the top match; form would send %q", got)
	}
	if got := eval[string](t, ctx, `document.querySelector('`+cbEstado+`').value`); got != "Pará" {
		t.Errorf("the field should show the picked label, shows %q", got)
	}
	if got := eval[string](t, ctx, `document.querySelector('`+cbEstado+`').getAttribute('aria-expanded')`); got != "false" {
		t.Errorf("picking should close the list; aria-expanded=%s", got)
	}
}

// The disabled option is skipped, and focus never leaves the field — the
// listbox is reached through aria-activedescendant alone.
func TestComboboxArrowsSkipDisabledAndKeepFocus(t *testing.T) {
	ctx := openPage(t, "/components/combobox")
	const activeLabel = `document.getElementById(document.querySelector('#cb-plano').getAttribute('aria-activedescendant'))?.dataset.label ?? ''`

	focus(t, ctx, cbPlano)
	var seen []string
	for range 3 {
		press(t, ctx, kb.ArrowDown)
		settle(t, ctx)
		seen = append(seen, eval[string](t, ctx, activeLabel))
	}
	if strings.Join(seen, ",") != "Grátis,Pro,Grátis" {
		t.Errorf("↓ ↓ ↓ highlighted %v; want to wrap past the disabled Enterprise", seen)
	}
	if got := eval[string](t, ctx, `document.activeElement.id`); got != "cb-plano" {
		t.Errorf("focus moved to %q; it must stay in the text field", got)
	}
}

// Free text never reaches the form: leaving the field reverts it to the last
// real choice, and emptying it clears the value.
func TestComboboxBlurRevertsFreeTextAndEmptyClears(t *testing.T) {
	ctx := openPage(t, "/components/combobox")

	retype(t, ctx, cbEstado, "xyz")
	focus(t, ctx, cbPlano)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.querySelector('`+cbEstado+`').value`); got != "São Paulo" {
		t.Errorf("free text should revert to the selection, field shows %q", got)
	}
	if got := eval[string](t, ctx, cbEstadoSent); got != "sp" {
		t.Errorf("free text changed the submitted value to %q", got)
	}

	// An exact label, typed without accents, counts as choosing it.
	retype(t, ctx, cbEstado, "goias")
	press(t, ctx, kb.Escape)
	focus(t, ctx, cbPlano)
	settle(t, ctx)
	if got := eval[string](t, ctx, cbEstadoSent); got != "go" {
		t.Errorf("typing a full label should select it; form would send %q", got)
	}

	retype(t, ctx, cbEstado, kb.Backspace)
	focus(t, ctx, cbPlano)
	settle(t, ctx)
	if got := eval[string](t, ctx, cbEstadoSent); got != "" {
		t.Errorf("an emptied field should clear the value, still %q", got)
	}
}

func TestComboboxEscapeClosesThenClears(t *testing.T) {
	ctx := openPage(t, "/components/combobox")

	focus(t, ctx, cbEstado)
	press(t, ctx, kb.ArrowDown)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.querySelector('`+cbEstado+`').getAttribute('aria-expanded')`); got != "true" {
		t.Fatal("↓ did not open the list; this test would prove nothing")
	}

	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.querySelector('`+cbEstado+`').getAttribute('aria-expanded')`); got != "false" {
		t.Error("the first Escape should close the list")
	}
	if got := eval[string](t, ctx, cbEstadoSent); got != "sp" {
		t.Errorf("closing must not touch the value, now %q", got)
	}

	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if got := eval[string](t, ctx, cbEstadoSent); got != "" {
		t.Errorf("the second Escape should clear, value is %q", got)
	}
}

// Retyping the current label is not a search: the list must not jump to the
// first option, or Enter would silently change the selection.
func TestComboboxRetypingTheSelectionKeepsIt(t *testing.T) {
	ctx := openPage(t, "/components/combobox")

	focus(t, ctx, cbEstado)
	eval[bool](t, ctx, `(document.querySelector('`+cbEstado+`').setSelectionRange(99, 99), true)`)
	press(t, ctx, kb.Backspace)
	press(t, ctx, "o")
	settle(t, ctx)
	press(t, ctx, kb.Enter)
	settle(t, ctx)
	if got := eval[string](t, ctx, cbEstadoSent); got != "sp" {
		t.Errorf("retyping \"São Paulo\" and pressing Enter changed the value to %q", got)
	}
}

// Tab after arrowing to an option takes it, the same as Enter would.
func TestComboboxTabTakesTheHighlightedOption(t *testing.T) {
	ctx := openPage(t, "/components/combobox")

	retype(t, ctx, cbEstado, "para")
	press(t, ctx, kb.ArrowDown)
	press(t, ctx, kb.ArrowDown)
	settle(t, ctx)
	press(t, ctx, kb.Tab)
	settle(t, ctx)
	if got := eval[string](t, ctx, cbEstadoSent); got != "pr" {
		t.Errorf("Tab on the highlighted Paraná should select it; form would send %q", got)
	}
}

// Enter with nothing highlighted submits the form, so the field has to be
// reconciled first: what is on screen must be what is sent.
func TestComboboxEnterWithoutHighlightCommitsBeforeSubmit(t *testing.T) {
	ctx := openPage(t, "/components/combobox")

	retype(t, ctx, cbEstado, "xyz")
	press(t, ctx, kb.Enter)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.querySelector('`+cbEstado+`').value`); got != "São Paulo" {
		t.Errorf("Enter should revert unmatched text before submitting; field shows %q", got)
	}
}

// Inside a Modal, the Escape that closes the list must stop there: the modal
// listens on window and would close with the form in it.
func TestComboboxEscapeDoesNotReachWindow(t *testing.T) {
	ctx := openPage(t, "/components/combobox")
	eval[bool](t, ctx, `(window.__esc = 0, window.addEventListener('keydown', e => { if (e.key === 'Escape') window.__esc++ }), true)`)

	focus(t, ctx, cbEstado)
	press(t, ctx, kb.ArrowDown)
	settle(t, ctx)
	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if got := eval[int](t, ctx, `window.__esc`); got != 0 {
		t.Errorf("Escape that closed the list reached window %d time(s)", got)
	}
}

// --- Carousel ---------------------------------------------------------------

func TestCarouselMovesWithArrowKeys(t *testing.T) {
	ctx := openPage(t, "/components/carousel")

	const pos = `document.querySelector('[aria-roledescription="carousel"] [tabindex="0"]').scrollLeft`

	before := eval[float64](t, ctx, pos)
	focus(t, ctx, `[aria-roledescription="carousel"] [tabindex="0"]`)
	press(t, ctx, kb.ArrowRight)
	sleep(t, ctx, 700)

	if after := eval[float64](t, ctx, pos); after <= before {
		t.Errorf("ArrowRight did not scroll the track: %v -> %v", before, after)
	}
}
