package e2e_test

import (
	"context"
	"testing"

	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

// --- Popover ----------------------------------------------------------------

const popPanel = `#pop-filtros-panel`

func visible(sel string) string {
	return `(() => { const el = document.querySelector('` + sel + `'); return !!el && getComputedStyle(el).display !== 'none' })()`
}

func TestPopoverFocusesContentAndEscapeReturns(t *testing.T) {
	ctx := openPage(t, "/components/popover")

	click(t, ctx, `[aria-controls="pop-filtros-panel"]`)
	settle(t, ctx)
	if !eval[bool](t, ctx, visible(popPanel)) {
		t.Fatal("the popover did not open")
	}
	if got := eval[string](t, ctx, `document.activeElement.id`); got != "pop-autor" {
		t.Errorf("focus should move to the first field, is on %q", got)
	}

	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if eval[bool](t, ctx, visible(popPanel)) {
		t.Error("Escape did not close the popover")
	}
	if got := eval[string](t, ctx, `document.activeElement.getAttribute('aria-controls')`); got != "pop-filtros-panel" {
		t.Errorf("focus should return to the trigger, is on %q", got)
	}
}

// Non-modal: nothing behind it is blocked, and Tab out of it closes it.
func TestPopoverIsNotModal(t *testing.T) {
	ctx := openPage(t, "/components/popover")

	click(t, ctx, `[aria-controls="pop-filtros-panel"]`)
	settle(t, ctx)
	if eval[bool](t, ctx, `getComputedStyle(document.documentElement).overflow === 'hidden'`) {
		t.Error("a popover must not lock page scroll")
	}
	if n := eval[int](t, ctx, `[...document.body.children].filter(el => el.hasAttribute('inert')).length`); n != 0 {
		t.Errorf("a popover must not make the page inert (%d inert)", n)
	}

	focus(t, ctx, `#pop-filtros-panel input[type=checkbox]`)
	press(t, ctx, kb.Tab)
	settle(t, ctx)
	if eval[bool](t, ctx, visible(popPanel)) {
		t.Error("Tab out of the popover should close it")
	}
	if eval[bool](t, ctx, `document.querySelector('`+popPanel+`').contains(document.activeElement)`) {
		t.Error("focus is still inside a closed popover")
	}
}

// The Escape that closes a popover inside a Modal must not reach window, where
// the Modal would close too.
func TestPopoverEscapeInsideModalClosesOnlyThePopover(t *testing.T) {
	ctx := openPage(t, "/components/popover")

	click(t, ctx, "#abrir-modal-popover")
	settle(t, ctx)
	click(t, ctx, `[aria-controls="pop-no-modal-panel"]`)
	settle(t, ctx)
	if !eval[bool](t, ctx, visible(`#pop-no-modal-panel`)) {
		t.Fatal("the popover inside the modal did not open")
	}

	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if eval[bool](t, ctx, visible(`#pop-no-modal-panel`)) {
		t.Error("Escape did not close the popover")
	}
	if !eval[bool](t, ctx, `!!document.querySelector('[role="dialog"][aria-modal="true"]')`) {
		t.Error("Escape closed the modal too; it should have stopped at the popover")
	}
}

// --- Hover Card -------------------------------------------------------------

// Keyboard focus opens it too, but only after the delay: a card that flashed
// open on every Tab through a paragraph of links would be noise.
func TestHoverCardOpensOnFocusAfterDelay(t *testing.T) {
	ctx := openPage(t, "/components/hover-card")

	focus(t, ctx, hcTrigger)
	sleep(t, ctx, 100)
	if eval[bool](t, ctx, visible(`#hc-lincoln-card`)) {
		t.Error("the card opened before its delay")
	}
	sleep(t, ctx, 600)
	if !eval[bool](t, ctx, visible(`#hc-lincoln-card`)) {
		t.Error("the card did not open on keyboard focus")
	}
}

func TestHoverCardEscapeCloses(t *testing.T) {
	ctx := openPage(t, "/components/hover-card")

	focus(t, ctx, hcTrigger)
	sleep(t, ctx, 700)
	if !eval[bool](t, ctx, visible(`#hc-lincoln-card`)) {
		t.Fatal("the card never opened; this test would prove nothing")
	}
	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if eval[bool](t, ctx, visible(`#hc-lincoln-card`)) {
		t.Error("Escape did not close the card")
	}
}

// Scoped to the component: the showcase navigation links to the same page.
const hcTrigger = `[x-data^="hoverCard"] a`

// --- Alert Dialog -----------------------------------------------------------

const adOpen = `!!document.querySelector('[role="alertdialog"]')`

func TestAlertDialogFocusStartsOnCancel(t *testing.T) {
	ctx := openPage(t, "/components/alert-dialog")

	click(t, ctx, "#abrir-ad")
	settle(t, ctx)
	if !eval[bool](t, ctx, adOpen) {
		t.Fatal("the alert dialog did not open")
	}
	if got := eval[string](t, ctx, `document.activeElement.textContent.trim()`); got != "Cancelar" {
		t.Errorf("initial focus is on %q; it must be on the safe choice", got)
	}
}

// A stray click outside must not count as an answer.
func TestAlertDialogBackdropDoesNotClose(t *testing.T) {
	ctx := openPage(t, "/components/alert-dialog")

	click(t, ctx, "#abrir-ad")
	settle(t, ctx)
	eval[bool](t, ctx, `(document.elementFromPoint(5, 5).click(), true)`)
	settle(t, ctx)
	if !eval[bool](t, ctx, adOpen) {
		t.Error("clicking the backdrop closed the alert dialog")
	}
}

func TestAlertDialogConfirmRunsAndCloses(t *testing.T) {
	ctx := openPage(t, "/components/alert-dialog")

	click(t, ctx, "#abrir-ad")
	settle(t, ctx)
	click(t, ctx, `[role="alertdialog"] button.bg-danger`)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.getElementById('ad-resultado').textContent`); got != "excluído" {
		t.Errorf("the confirm action did not run: %q", got)
	}
	if eval[bool](t, ctx, adOpen) {
		t.Error("confirming should close the dialog")
	}
	if got := eval[string](t, ctx, `document.activeElement.id`); got != "abrir-ad" {
		t.Errorf("focus should return to the trigger, is on %q", got)
	}
}

// --- Context Menu -----------------------------------------------------------

const cmMenu = `#cm-arquivo-menu`

func rightClickAt(t *testing.T, ctx context.Context, x, y float64) {
	t.Helper()
	if err := chromedp.Run(ctx, chromedp.MouseClickXY(x, y, chromedp.ButtonRight)); err != nil {
		t.Fatalf("right-click at %v,%v: %v", x, y, err)
	}
}

type rect struct{ Left, Top, Right, Bottom, Width, Height float64 }

func rectOf(t *testing.T, ctx context.Context, sel string) rect {
	t.Helper()
	return eval[rect](t, ctx, `(() => { const r = document.querySelector('`+sel+`').getBoundingClientRect(); return {Left: r.left, Top: r.top, Right: r.right, Bottom: r.bottom, Width: r.width, Height: r.height} })()`)
}

func TestContextMenuOpensAtThePointer(t *testing.T) {
	ctx := openPage(t, "/components/context-menu")
	eval[bool](t, ctx, `(document.querySelector('#cm-arquivo-hint').previousElementSibling.scrollIntoView({block: 'center'}), true)`)
	target := rectOf(t, ctx, `[x-ref="target"]`)

	x, y := target.Left+20, target.Top+20
	rightClickAt(t, ctx, x, y)
	settle(t, ctx)
	if !eval[bool](t, ctx, visible(cmMenu)) {
		t.Fatal("right-click did not open the menu")
	}
	menu := rectOf(t, ctx, cmMenu)
	if menu.Left < x-2 || menu.Left > x+2 || menu.Top < y-2 || menu.Top > y+2 {
		t.Errorf("menu at %.0f,%.0f; expected at the pointer %.0f,%.0f", menu.Left, menu.Top, x, y)
	}
	if got := eval[string](t, ctx, `document.activeElement.textContent.trim()`); got != "Renomear" {
		t.Errorf("focus should be on the first item, is on %q", got)
	}
}

// Near the bottom-right corner, opening at the pointer would push the menu
// off screen. The viewport is narrowed and the target scrolled to the bottom
// so the click really lands in that corner.
func TestContextMenuStaysInsideTheViewport(t *testing.T) {
	ctx := openPage(t, "/components/context-menu")
	if err := chromedp.Run(ctx, chromedp.EmulateViewport(420, 600)); err != nil {
		t.Fatal(err)
	}
	// The padding pushes the target below the fold, so scrolling it to the
	// bottom edge works however little the host renders above it. Without it,
	// a target already on screen would need a scroll past the top.
	eval[bool](t, ctx, `(document.body.style.paddingTop = '100vh', document.querySelector('[x-ref="target"]').scrollIntoView({block: 'end'}), true)`)
	settle(t, ctx)
	target := rectOf(t, ctx, `[x-ref="target"]`)
	vw :=eval[float64](t, ctx, `window.innerWidth`)
	vh := eval[float64](t, ctx, `window.innerHeight`)

	x, y := target.Right-5, target.Bottom-5
	if x < vw-200 || y < vh-150 {
		t.Fatalf("fixture: the click at %.0f,%.0f is not near the corner of a %.0fx%.0f viewport", x, y, vw, vh)
	}
	rightClickAt(t, ctx, x, y)
	settle(t, ctx)
	menu := rectOf(t, ctx, cmMenu)
	if menu.Left < 0 || menu.Top < 0 || menu.Right > vw || menu.Bottom > vh {
		t.Errorf("menu %+v spills out of the %.0fx%.0f viewport", menu, vw, vh)
	}
}

func TestContextMenuOpensFromTheKeyboard(t *testing.T) {
	ctx := openPage(t, "/components/context-menu")

	focus(t, ctx, `[x-ref="target"]`)
	if err := chromedp.Run(ctx, chromedp.KeyEvent(kb.F10, chromedp.KeyModifiers(input.ModifierShift))); err != nil {
		t.Fatal(err)
	}
	settle(t, ctx)
	if !eval[bool](t, ctx, visible(cmMenu)) {
		t.Fatal("Shift+F10 did not open the menu")
	}
	if got := eval[string](t, ctx, `document.activeElement.textContent.trim()`); got != "Renomear" {
		t.Errorf("focus should be on the first item, is on %q", got)
	}

	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if eval[bool](t, ctx, visible(cmMenu)) {
		t.Error("Escape did not close the menu")
	}
	if !eval[bool](t, ctx, `document.activeElement === document.querySelector('[x-ref="target"]')`) {
		t.Error("focus should return to the target")
	}
}

// --- Menubar ----------------------------------------------------------------

func mbTop(label string) string {
	return `[role="menubar"] [data-menubar-top]` + `:is(` + `[data-menubar-top="` + map[string]string{"Arquivo": "0", "Editar": "1", "Exibir": "2"}[label] + `"])`
}

const mbActive = `document.activeElement.textContent.trim()`
const mbOpenMenus = `[...document.querySelectorAll('[role="menubar"] [role="menu"]')].filter(m => getComputedStyle(m).display !== 'none').map(m => m.getAttribute('aria-label'))`

func TestMenubarArrowsMoveBetweenMenus(t *testing.T) {
	ctx := openPage(t, "/components/menubar")

	focus(t, ctx, mbTop("Arquivo"))
	press(t, ctx, kb.ArrowRight)
	settle(t, ctx)
	if got := eval[string](t, ctx, mbActive); got != "Editar" {
		t.Fatalf("→ on the bar should focus Editar, focus is on %q", got)
	}
	if got := eval[[]string](t, ctx, mbOpenMenus); len(got) != 0 {
		t.Errorf("moving across a closed bar must not open menus, open: %v", got)
	}

	press(t, ctx, kb.ArrowDown)
	settle(t, ctx)
	if got := eval[string](t, ctx, mbActive); got != "Desfazer" {
		t.Errorf("↓ should open Editar and focus its first item, focus is on %q", got)
	}

	// Inside an open menu, → carries the open state to the next menu.
	press(t, ctx, kb.ArrowRight)
	settle(t, ctx)
	if got := eval[[]string](t, ctx, mbOpenMenus); len(got) != 1 || got[0] != "Exibir" {
		t.Errorf("→ inside Editar should open Exibir, open: %v", got)
	}
	if got := eval[string](t, ctx, mbActive); got != "Modo foco" {
		t.Errorf("focus should be on Exibir's first item, is on %q", got)
	}
}

func TestMenubarEscapeReturnsToTheTopItem(t *testing.T) {
	ctx := openPage(t, "/components/menubar")

	focus(t, ctx, mbTop("Arquivo"))
	press(t, ctx, kb.ArrowDown)
	settle(t, ctx)
	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if got := eval[[]string](t, ctx, mbOpenMenus); len(got) != 0 {
		t.Errorf("Escape should close the menu, open: %v", got)
	}
	if got := eval[string](t, ctx, mbActive); got != "Arquivo" {
		t.Errorf("focus should return to Arquivo, is on %q", got)
	}
}

func TestMenubarTabClosesEverything(t *testing.T) {
	ctx := openPage(t, "/components/menubar")

	focus(t, ctx, mbTop("Arquivo"))
	press(t, ctx, kb.ArrowDown)
	settle(t, ctx)
	press(t, ctx, kb.Tab)
	settle(t, ctx)
	if got := eval[[]string](t, ctx, mbOpenMenus); len(got) != 0 {
		t.Errorf("Tab should close every menu, open: %v", got)
	}
	if eval[bool](t, ctx, `document.querySelector('[role="menubar"]').contains(document.activeElement)`) {
		t.Error("Tab should let focus leave the menu bar")
	}
}

func TestMenubarItemRunsAndCloses(t *testing.T) {
	ctx := openPage(t, "/components/menubar")

	focus(t, ctx, mbTop("Arquivo"))
	press(t, ctx, kb.ArrowDown)
	settle(t, ctx) // the menu opens and takes focus on the next tick
	press(t, ctx, kb.ArrowDown)
	settle(t, ctx)
	press(t, ctx, kb.Enter)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.getElementById('mb-resultado').textContent`); got != "abrir" {
		t.Errorf("Enter on Abrir should run it, result is %q", got)
	}
	if got := eval[[]string](t, ctx, mbOpenMenus); len(got) != 0 {
		t.Errorf("running an item should close the menu, open: %v", got)
	}
	if got := eval[string](t, ctx, mbActive); got != "Arquivo" {
		t.Errorf("focus should return to Arquivo, is on %q", got)
	}
}

// --- Command ----------------------------------------------------------------

const cmdOpen = `!!document.querySelector('[aria-label="Paleta de comandos"]')`
const cmdVisibleOptions = `[...document.querySelectorAll('[role="listbox"] [role="option"]')].filter(o => getComputedStyle(o).display !== 'none').map(o => o.dataset.label)`

// pressEnterNative sends Enter as one keyDown carrying its text, which is what
// a physical keyboard produces: the browser derives the keypress from it, and
// preventDefault on keydown cancels that keypress. chromedp's KeyEvent sends a
// separate "char" event instead, which lands on whatever has focus afterwards
// — here the trigger the palette just returned focus to, reopening it.
func pressEnterNative(t *testing.T, ctx context.Context) {
	t.Helper()
	down := input.DispatchKeyEvent(input.KeyDown).WithKey("Enter").WithCode("Enter").
		WithWindowsVirtualKeyCode(13).WithNativeVirtualKeyCode(13).WithText(string(rune(13))) // carriage return, the text of Enter
	up := input.DispatchKeyEvent(input.KeyUp).WithKey("Enter").WithCode("Enter").
		WithWindowsVirtualKeyCode(13).WithNativeVirtualKeyCode(13)
	if err := chromedp.Run(ctx, down, up); err != nil {
		t.Fatalf("press Enter: %v", err)
	}
}

func ctrlK(t *testing.T, ctx context.Context) {
	t.Helper()
	if err := chromedp.Run(ctx, chromedp.KeyEvent("k", chromedp.KeyModifiers(input.ModifierCtrl))); err != nil {
		t.Fatal(err)
	}
}

func TestCommandOpensWithCtrlK(t *testing.T) {
	ctx := openPage(t, "/components/command")

	ctrlK(t, ctx)
	settle(t, ctx)
	if !eval[bool](t, ctx, cmdOpen) {
		t.Fatal("Ctrl+K did not open the palette")
	}
	if got := eval[string](t, ctx, `document.activeElement.getAttribute('role')`); got != "combobox" {
		t.Errorf("focus should be in the search field, is on role=%q", got)
	}
}

func TestCommandFiltersAndRunsWithEnter(t *testing.T) {
	ctx := openPage(t, "/components/command")

	click(t, ctx, "#abrir-cmd")
	settle(t, ctx)
	press(t, ctx, "tema")
	settle(t, ctx)
	got := eval[[]string](t, ctx, cmdVisibleOptions)
	if len(got) != 2 || got[0] != "Alternar tema" || got[1] != "Tema do sistema" {
		t.Fatalf("filtering \"tema\" shows %v", got)
	}
	// The group with no match disappears with its heading.
	if eval[bool](t, ctx, visible(`[data-command-group="1"]`)) {
		t.Error("the Páginas group should hide when none of its items match")
	}

	press(t, ctx, kb.ArrowDown)
	settle(t, ctx)
	pressEnterNative(t, ctx)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.getElementById('cmd-resultado').textContent`); got != "tema do sistema" {
		t.Errorf("Enter should run the highlighted command, result is %q", got)
	}
	if eval[bool](t, ctx, cmdOpen) {
		t.Error("running a command should close the palette")
	}
}

func TestCommandEnterWithNoResultsDoesNothing(t *testing.T) {
	ctx := openPage(t, "/components/command")
	eval[bool](t, ctx, `(window.__errors = [], window.addEventListener('error', e => window.__errors.push(e.message)), true)`)

	click(t, ctx, "#abrir-cmd")
	settle(t, ctx)
	press(t, ctx, "zzz")
	settle(t, ctx)
	press(t, ctx, kb.Enter)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.getElementById('cmd-resultado').textContent`); got != "nenhum" {
		t.Errorf("Enter with no results ran %q", got)
	}
	if !eval[bool](t, ctx, cmdOpen) {
		t.Error("Enter with no results should leave the palette open")
	}
	if errs := eval[[]string](t, ctx, `window.__errors`); len(errs) != 0 {
		t.Errorf("Enter with no results raised: %v", errs)
	}
}

func TestCommandIsModal(t *testing.T) {
	ctx := openPage(t, "/components/command")

	click(t, ctx, "#abrir-cmd")
	settle(t, ctx)
	if !eval[bool](t, ctx, `getComputedStyle(document.documentElement).overflow === 'hidden'`) {
		t.Error("the palette is modal: page scroll should be locked")
	}
	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if eval[bool](t, ctx, cmdOpen) {
		t.Error("Escape should close the palette")
	}
	if got := eval[string](t, ctx, `document.activeElement.id`); got != "abrir-cmd" {
		t.Errorf("focus should return to the trigger, is on %q", got)
	}
}

// --- Drawer, bottom ---------------------------------------------------------

func TestDrawerBottomSitsAtTheBottom(t *testing.T) {
	ctx := openPage(t, "/components/drawer")

	click(t, ctx, "#abrir-drawer-bottom")
	settle(t, ctx)
	panel := rectOf(t, ctx, `[role="dialog"][aria-labelledby="demo-drawer-bottom-title"]`)
	vw := eval[float64](t, ctx, `window.innerWidth`)
	vh := eval[float64](t, ctx, `window.innerHeight`)
	if panel.Bottom < vh-1 || panel.Bottom > vh+1 {
		t.Errorf("the sheet should touch the bottom edge: bottom=%.0f, viewport=%.0f", panel.Bottom, vh)
	}
	if panel.Width < vw*0.95 {
		t.Errorf("the sheet should span the width: %.0f of %.0f", panel.Width, vw)
	}
	if panel.Top <= 0 {
		t.Error("the sheet should not cover the whole page")
	}
}

// Walking the open menu around the bar must always show exactly the menu that
// is open. Three sibling x-show panels toggled this way stopped updating after
// a few switches — the state said Editar while Arquivo stayed on screen — and
// the shorter tests above never switched enough times to see it.
func TestMenubarSwitchingMenusKeepsTheScreenInStep(t *testing.T) {
	ctx := openPage(t, "/components/menubar")

	focus(t, ctx, mbTop("Arquivo"))
	press(t, ctx, kb.ArrowDown)
	settle(t, ctx)
	for _, step := range []struct {
		key  string
		want string
	}{
		{kb.ArrowRight, "Editar"},
		{kb.ArrowRight, "Exibir"},
		{kb.ArrowRight, "Arquivo"},
		{kb.ArrowLeft, "Exibir"},
		{kb.ArrowLeft, "Editar"},
		{kb.ArrowLeft, "Arquivo"},
	} {
		press(t, ctx, step.key)
		settle(t, ctx)
		got := eval[[]string](t, ctx, mbOpenMenus)
		if len(got) != 1 || got[0] != step.want {
			t.Fatalf("after %q expected only %s open, screen shows %v", step.key, step.want, got)
		}
		if !eval[bool](t, ctx, `document.querySelector('[role="menu"][aria-label="`+step.want+`"]').contains(document.activeElement)`) {
			t.Fatalf("after %q focus is not inside %s", step.key, step.want)
		}
	}
}

// --- Frame timing -------------------------------------------------------------
//
// Alpine's x-show reveals a panel inside requestAnimationFrame (after its first
// toggle), while $nextTick is a timer. Chrome happens to run the frame first;
// Firefox usually runs the timer first. slowFrames reproduces that order in
// Chrome by pushing every frame callback behind a timer, so work that needs
// the panel visible — focusing into it, measuring it — is tested against the
// order it cannot rely on.
func slowFrames(t *testing.T, ctx context.Context) {
	t.Helper()
	eval[bool](t, ctx, `(() => { const raf = window.requestAnimationFrame.bind(window); window.requestAnimationFrame = cb => setTimeout(() => raf(cb), 40); return true })()`)
}

func TestPopoverFocusesInWhenFramesAreSlow(t *testing.T) {
	ctx := openPage(t, "/components/popover")
	slowFrames(t, ctx)

	// Twice: x-show only defers to the frame from its second toggle on.
	for i := 0; i < 2; i++ {
		click(t, ctx, `[aria-controls="pop-filtros-panel"]`)
		sleep(t, ctx, 300)
		if got := eval[string](t, ctx, `document.activeElement.id`); got != "pop-autor" {
			t.Fatalf("open #%d: focus should move into the panel, is on %q", i+1, got)
		}
		press(t, ctx, kb.Escape)
		sleep(t, ctx, 300)
	}
}

func TestMenubarFocusesTheNextMenuWhenFramesAreSlow(t *testing.T) {
	ctx := openPage(t, "/components/menubar")
	slowFrames(t, ctx)

	focus(t, ctx, mbTop("Arquivo"))
	for _, step := range []struct{ key, want string }{
		{kb.ArrowDown, "Novo protocolo"},
		{kb.ArrowRight, "Desfazer"},
		{kb.ArrowRight, "Modo foco"},
		{kb.ArrowRight, "Novo protocolo"},
	} {
		press(t, ctx, step.key)
		sleep(t, ctx, 300)
		if got := eval[string](t, ctx, mbActive); got != step.want {
			t.Fatalf("after %q focus should be on %q, is on %q", step.key, step.want, got)
		}
	}
}

func TestContextMenuClampsWhenFramesAreSlow(t *testing.T) {
	ctx := openPage(t, "/components/context-menu")
	if err := chromedp.Run(ctx, chromedp.EmulateViewport(420, 600)); err != nil {
		t.Fatal(err)
	}
	eval[bool](t, ctx, `(document.querySelector('[x-ref="target"]').scrollIntoView({block: 'end'}), true)`)
	slowFrames(t, ctx)
	settle(t, ctx)
	target := rectOf(t, ctx, `[x-ref="target"]`)
	vw := eval[float64](t, ctx, `window.innerWidth`)
	vh := eval[float64](t, ctx, `window.innerHeight`)

	// First open and close, so the next open goes through the deferred path.
	rightClickAt(t, ctx, target.Left+20, target.Top+20)
	sleep(t, ctx, 300)
	press(t, ctx, kb.Escape)
	sleep(t, ctx, 300)

	rightClickAt(t, ctx, target.Right-5, target.Bottom-5)
	sleep(t, ctx, 300)
	menu := rectOf(t, ctx, cmMenu)
	if menu.Right > vw || menu.Bottom > vh {
		t.Errorf("menu %+v spills out of the %.0fx%.0f viewport", menu, vw, vh)
	}
	if got := eval[string](t, ctx, `document.activeElement.textContent.trim()`); got != "Renomear" {
		t.Errorf("focus should be on the first item, is on %q", got)
	}
}

// A second right-click while the menu is open must not make the menu its own
// focus-return target: Escape would then focus a hidden item and drop to body.
func TestContextMenuSecondRightClickStillReturnsFocus(t *testing.T) {
	ctx := openPage(t, "/components/context-menu")
	focus(t, ctx, `[x-ref="target"]`)
	target := rectOf(t, ctx, `[x-ref="target"]`)

	rightClickAt(t, ctx, target.Left+20, target.Top+20)
	settle(t, ctx)
	rightClickAt(t, ctx, target.Left+120, target.Top+40)
	settle(t, ctx)
	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if !eval[bool](t, ctx, `document.activeElement === document.querySelector('[x-ref="target"]')`) {
		t.Errorf("focus should return to the target, is on %q", eval[string](t, ctx, `document.activeElement.tagName`))
	}
}

// The menu must not reach into the content it wraps: roving items inside the
// area are not menu items, a left click on the area closes the menu, and a
// right-click elsewhere closes it too.
func TestContextMenuStaysOutOfTheWrappedContent(t *testing.T) {
	ctx := openPage(t, "/components/context-menu")
	// A roving item inside the area, as a ToggleGroup or Tabs there would add.
	eval[bool](t, ctx, `(() => { const b = document.createElement('button'); b.dataset.rovingItem = '0'; b.id = 'intruso'; b.textContent = 'intruso'; document.querySelector('[x-ref="target"]').appendChild(b); return true })()`)
	target := rectOf(t, ctx, `[x-ref="target"]`)

	rightClickAt(t, ctx, target.Left+20, target.Top+20)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.activeElement.textContent.trim()`); got != "Renomear" {
		t.Fatalf("opening focused %q; a roving item in the wrapped area is not a menu item", got)
	}

	if err := chromedp.Run(ctx, chromedp.MouseClickXY(target.Right-10, target.Bottom-10)); err != nil {
		t.Fatal(err)
	}
	settle(t, ctx)
	if eval[bool](t, ctx, visible(cmMenu)) {
		t.Error("a left click on the area (outside the menu) should close it")
	}

	rightClickAt(t, ctx, target.Left+20, target.Top+20)
	settle(t, ctx)
	heading := rectOf(t, ctx, `h1`)
	rightClickAt(t, ctx, heading.Left+5, heading.Top+5)
	settle(t, ctx)
	if eval[bool](t, ctx, visible(cmMenu)) {
		t.Error("a right-click elsewhere should close the menu")
	}
}

// Modal overlays stack — a confirmation over an edit form, a palette over a
// dialog — and every one of them listens for Escape at window. One Escape must
// close only the top one, or cancelling a confirmation discards the form.
func TestEscapeClosesOnlyTheTopModalOverlay(t *testing.T) {
	ctx := openPage(t, "/components/alert-dialog")
	const modalOpen = `!!document.querySelector('[role="dialog"][aria-labelledby="md-editar-title"]')`
	const alertOpen = `!!document.querySelector('[role="alertdialog"][aria-labelledby="ad-descartar-title"]')`

	click(t, ctx, "#abrir-md-editar")
	settle(t, ctx)
	click(t, ctx, "#abrir-ad-descartar")
	settle(t, ctx)
	if !eval[bool](t, ctx, alertOpen) || !eval[bool](t, ctx, modalOpen) {
		t.Fatal("fixture: expected the alert dialog open over the modal")
	}

	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if eval[bool](t, ctx, alertOpen) {
		t.Error("Escape should close the confirmation")
	}
	if !eval[bool](t, ctx, modalOpen) {
		t.Fatal("Escape closed the modal underneath too; the user's form is gone")
	}
	if got := eval[string](t, ctx, `document.activeElement.id`); got != "abrir-ad-descartar" {
		t.Errorf("focus should return to the button inside the modal, is on %q", got)
	}

	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if eval[bool](t, ctx, modalOpen) {
		t.Error("a second Escape should close the modal")
	}
}

// Pointer and focus each hold the card open on their own: moving the mouse
// away must not close a card the keyboard is still on.
func TestHoverCardStaysOpenWhileFocusedAfterThePointerLeaves(t *testing.T) {
	ctx := openPage(t, "/components/hover-card")
	link := rectOf(t, ctx, hcTrigger)

	if err := chromedp.Run(ctx, chromedp.MouseEvent(input.MouseMoved, link.Left+5, link.Top+5)); err != nil {
		t.Fatal(err)
	}
	focus(t, ctx, hcTrigger)
	sleep(t, ctx, 700)
	if !eval[bool](t, ctx, visible(`#hc-lincoln-card`)) {
		t.Fatal("fixture: the card did not open")
	}
	if err := chromedp.Run(ctx, chromedp.MouseEvent(input.MouseMoved, 2, 2)); err != nil {
		t.Fatal(err)
	}
	sleep(t, ctx, 500)
	if !eval[bool](t, ctx, visible(`#hc-lincoln-card`)) {
		t.Error("the pointer left, but focus is still on the link: the card should stay open")
	}
}
