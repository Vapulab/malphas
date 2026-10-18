package e2e_test

import (
	"context"
	"strings"
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

// --- Textarea ---------------------------------------------------------------

func TestTextareaCounterFollowsTyping(t *testing.T) {
	ctx := openPage(t, "/components/textarea")

	focus(t, ctx, "#ta-bio")
	press(t, ctx, "olá 😀")
	settle(t, ctx)

	if got := eval[string](t, ctx, `document.querySelector('#ta-bio').closest('[x-data]').querySelector('[x-text="count"]').textContent`); got != "6" {
		t.Errorf("counter shows %q after typing \"olá 😀\", want 6 (emoji counts 2, as maxlength does)", got)
	}
}

// The parser drops one newline right after <textarea>; the component pads for
// it, so a saved value that starts with a blank line keeps it.
func TestTextareaKeepsALeadingNewlineInTheBrowser(t *testing.T) {
	ctx := openPage(t, "/components/textarea")

	if got := eval[string](t, ctx, `document.querySelector('#ta-nota').value`); got != "\nComeça com uma linha em branco." {
		t.Errorf("the parsed value is %q; the leading newline was lost", got)
	}
}

// reset restores the default value without an input event, so a counter that
// only listens to input would keep showing the typed length.
func TestTextareaCounterFollowsAFormReset(t *testing.T) {
	ctx := openPage(t, "/components/textarea")
	const counter = `document.querySelector('#ta-bio').closest('[x-data]').querySelector('[x-text="count"]').textContent`

	focus(t, ctx, "#ta-bio")
	press(t, ctx, "abc")
	settle(t, ctx)
	if got := eval[string](t, ctx, counter); got != "3" {
		t.Fatalf("fixture: counter shows %q after typing abc", got)
	}

	click(t, ctx, `main form button[type="reset"]`)
	settle(t, ctx)
	if got := eval[string](t, ctx, counter); got != "0" {
		t.Errorf("after reset the field is empty but the counter shows %q", got)
	}
}

// --- ToggleButton -----------------------------------------------------------

const formData = `Object.fromEntries(new FormData(document.querySelector('main form')).entries())`

func TestToggleButtonSubmitsWhilePressed(t *testing.T) {
	ctx := openPage(t, "/components/toggle-button")

	if got := eval[map[string]string](t, ctx, formData); got["favorito"] != "" || got["negrito"] != "on" || got["bloqueado"] != "" {
		t.Fatalf("fixture changed: %v", got)
	}

	click(t, ctx, "#tb-fav")
	click(t, ctx, "#tb-bold")
	settle(t, ctx)

	got := eval[map[string]string](t, ctx, formData)
	if got["favorito"] != "on" {
		t.Errorf("pressing should submit the value; payload %v", got)
	}
	if _, has := got["negrito"]; has {
		t.Errorf("releasing should stop submitting; payload %v", got)
	}
	if _, has := got["bloqueado"]; has {
		t.Errorf("a disabled toggle must never submit; payload %v", got)
	}
	if p := eval[string](t, ctx, `document.querySelector('#tb-fav').getAttribute('aria-pressed')`); p != "true" {
		t.Errorf("aria-pressed=%s after pressing", p)
	}
}

// Pressed must be visible, not only announced: the checkbox once shipped a
// state that passed every ARIA assertion and was invisible on screen.
func TestToggleButtonPressedLooksDifferent(t *testing.T) {
	ctx := openPage(t, "/components/toggle-button")
	const bg = `getComputedStyle(document.querySelector('%s')).backgroundColor`

	on := eval[string](t, ctx, strings.Replace(bg, "%s", "#tb-bold", 1))
	off := eval[string](t, ctx, strings.Replace(bg, "%s", "#tb-fav", 1))
	if on == off {
		t.Errorf("pressed and released share a background (%s)", on)
	}
}

// --- ToggleGroup ------------------------------------------------------------

const groupValues = `name => new FormData(document.querySelector('main form')).getAll(name)`

func groupButton(label string) string {
	return `//button[normalize-space()="` + label + `"]`
}

func clickByText(t *testing.T, ctx context.Context, label string) {
	t.Helper()
	if err := chromedp.Run(ctx, chromedp.Click(groupButton(label), chromedp.BySearch)); err != nil {
		t.Fatalf("click %q: %v", label, err)
	}
}

func TestToggleGroupSingleSwitchesAndReleases(t *testing.T) {
	ctx := openPage(t, "/components/toggle-group")
	values := func() []string { return eval[[]string](t, ctx, `(`+groupValues+`)('alinhamento')`) }

	if got := values(); len(got) != 1 || got[0] != "centro" {
		t.Fatalf("fixture changed: %v", got)
	}

	clickByText(t, ctx, "Esquerda")
	settle(t, ctx)
	if got := values(); len(got) != 1 || got[0] != "esq" {
		t.Errorf("single mode should switch to esq, form sends %v", got)
	}

	// Pressing the pressed item releases it: nothing is submitted.
	clickByText(t, ctx, "Esquerda")
	settle(t, ctx)
	if got := values(); len(got) != 0 {
		t.Errorf("pressing the pressed item should release it, form sends %v", got)
	}
}

func TestToggleGroupMultipleSubmitsEveryPressedValue(t *testing.T) {
	ctx := openPage(t, "/components/toggle-group")

	clickByText(t, ctx, "Negrito")
	settle(t, ctx)
	got := eval[[]string](t, ctx, `(`+groupValues+`)('formato')`)
	if strings.Join(got, ",") != "negrito,italico" {
		t.Errorf("form sends formato=%v, want [negrito italico]", got)
	}
}

// Arrows move focus only and skip the disabled item; nothing gets pressed.
func TestToggleGroupArrowsMoveFocusNotState(t *testing.T) {
	ctx := openPage(t, "/components/toggle-group")

	focus(t, ctx, `[aria-label="Formatação"] [tabindex="0"]`)
	press(t, ctx, kb.ArrowRight)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.activeElement.textContent.trim()`); got != "Tachado" {
		t.Errorf("→ from Itálico should skip the disabled Sublinhado and land on Tachado, got %q", got)
	}
	if got := eval[[]string](t, ctx, `(`+groupValues+`)('formato')`); strings.Join(got, ",") != "italico" {
		t.Errorf("arrows changed the pressed state: form sends %v", got)
	}
}

// --- Scroll Area ------------------------------------------------------------

// The whole point of tabindex=0: the keyboard alone can scroll it.
func TestScrollAreaScrollsFromTheKeyboard(t *testing.T) {
	ctx := openPage(t, "/components/scroll-area")
	const area = `[aria-label="Termos de uso"]`

	focus(t, ctx, area)
	press(t, ctx, kb.PageDown)
	sleep(t, ctx, 400) // smooth scrolling animates

	if got := eval[float64](t, ctx, `document.querySelector('`+area+`').scrollTop`); got <= 0 {
		t.Errorf("PageDown on the focused area did not scroll it (scrollTop=%v)", got)
	}
}

// --- Item -------------------------------------------------------------------

// The action sits above the stretched link's overlay: clicking it must run the
// action and stay on the page.
func TestItemActionDoesNotFollowTheLink(t *testing.T) {
	ctx := openPage(t, "/components/item")
	eval[bool](t, ctx, `(history.replaceState(null, '', '/components/item?marca=1'), true)`)

	click(t, ctx, "#item-baixar")
	settle(t, ctx)

	if got := eval[string](t, ctx, `location.search`); got != "?marca=1" {
		t.Errorf("clicking the action navigated (location.search=%q)", got)
	}
	if got := eval[string](t, ctx, `document.querySelector('#item-baixar').textContent`); got != "Baixado" {
		t.Errorf("the action did not run: %q", got)
	}
}

// --- Typography -------------------------------------------------------------

// Preflight zeroes heading sizes and list markers; prose has to win anyway.
func TestProseSurvivesPreflight(t *testing.T) {
	ctx := openPage(t, "/components/typography")
	const style = `(sel, prop) => getComputedStyle(document.querySelector('.prose-malphas ' + sel))[prop]`

	h2 := eval[string](t, ctx, `(`+style+`)('h2', 'fontSize')`)
	p := eval[string](t, ctx, `(`+style+`)('p', 'fontSize')`)
	if h2 == p {
		t.Errorf("h2 and p share a font size (%s): preflight won", h2)
	}
	if got := eval[string](t, ctx, `(`+style+`)('ul', 'listStyleType')`); got != "disc" {
		t.Errorf("list markers are %q, want disc", got)
	}
}
