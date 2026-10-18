package e2e_test

import (
	"context"
	"testing"

	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

// --- Sidebar ------------------------------------------------------------------

const sbAside = `#sb-app-aside`

// The default test window is 762px wide, just under md (768px), where the
// sidebar is the mobile drawer by design. Desktop tests ask for a wide one.
func desktop(t *testing.T, ctx context.Context) {
	t.Helper()
	if err := chromedp.Run(ctx, chromedp.EmulateViewport(1100, 800)); err != nil {
		t.Fatal(err)
	}
	settle(t, ctx)
}

const sbToggle = `[aria-controls="sb-app-aside"]`

func TestSidebarCollapsesAndKeepsNames(t *testing.T) {
	ctx := openPage(t, "/components/sidebar")
	desktop(t, ctx)
	eval[bool](t, ctx, `(localStorage.removeItem('malphas-sidebar:sb-app'), true)`)
	wide := eval[float64](t, ctx, `document.querySelector('`+sbAside+`').getBoundingClientRect().width`)

	click(t, ctx, sbToggle)
	sleep(t, ctx, 400) // width transition
	narrow := eval[float64](t, ctx, `document.querySelector('`+sbAside+`').getBoundingClientRect().width`)
	if narrow >= wide/2 {
		t.Errorf("collapsing should shrink the aside: %.0f → %.0f", wide, narrow)
	}
	if got := eval[string](t, ctx, `document.querySelector('`+sbToggle+`').getAttribute('aria-expanded')`); got != "false" {
		t.Errorf("aria-expanded=%s after collapsing", got)
	}
	// Labels are hidden visually, not removed: the link's name is intact.
	label := eval[string](t, ctx, `(() => { const s = document.querySelector('`+sbAside+` a[aria-current] span:not([aria-hidden])'); return getComputedStyle(s).position + '|' + s.closest('a').textContent.replace(/\s+/g, ' ').trim() })()`)
	if label != "absolute|◆ Protocolos 12" {
		t.Errorf("collapsed label should be sr-only (absolute) with the name intact, got %q", label)
	}
}

func TestSidebarRemembersCollapse(t *testing.T) {
	ctx := openPage(t, "/components/sidebar")
	eval[bool](t, ctx, `(localStorage.removeItem('malphas-sidebar:sb-app'), true)`)
	desktop(t, ctx)
	navigate(t, ctx, "/components/sidebar")

	click(t, ctx, sbToggle)
	settle(t, ctx)
	navigate(t, ctx, "/components/sidebar")
	if got := eval[string](t, ctx, `document.querySelector('`+sbAside+`').dataset.collapsed`); got != "true" {
		t.Errorf("the collapsed choice was not remembered across a reload: %q", got)
	}
	eval[bool](t, ctx, `(localStorage.removeItem('malphas-sidebar:sb-app'), true)`)
}

func TestSidebarIsADrawerOnMobile(t *testing.T) {
	ctx := openPage(t, "/components/sidebar")
	if err := chromedp.Run(ctx, chromedp.EmulateViewport(400, 700)); err != nil {
		t.Fatal(err)
	}
	settle(t, ctx)
	if eval[bool](t, ctx, visible(sbAside)) {
		t.Error("the aside should give way on a phone")
	}
	click(t, ctx, `[aria-haspopup="dialog"]`)
	settle(t, ctx)
	links := eval[[]string](t, ctx, `[...document.querySelectorAll('[role="dialog"][aria-modal="true"] nav a')].map(a => a.textContent.replace(/\s+/g, ' ').trim())`)
	if len(links) != 5 || links[0] != "◆ Protocolos 12" {
		t.Errorf("the drawer should hold the same navigation, got %v", links)
	}
}

// --- Resizable ----------------------------------------------------------------

const rzBar = `#rz-editor-first + [role="separator"]`
const rzNow = `+document.querySelector('` + rzBar + `').getAttribute('aria-valuenow')`

func clearRz(t *testing.T, ctx context.Context) {
	t.Helper()
	eval[bool](t, ctx, `(localStorage.removeItem('malphas-resizable:rz-editor'), true)`)
}

func TestResizableKeyboardMovesAndClamps(t *testing.T) {
	ctx := openPage(t, "/components/resizable")
	clearRz(t, ctx)
	navigate(t, ctx, "/components/resizable")

	focus(t, ctx, rzBar)
	for _, step := range []struct {
		key  string
		want int
	}{
		{kb.ArrowRight, 40},
		{kb.End, 70},
		{kb.ArrowRight, 70}, // clamped at Max
		{kb.Home, 20},
		{kb.ArrowLeft, 20}, // clamped at Min
	} {
		press(t, ctx, step.key)
		settle(t, ctx)
		if got := eval[int](t, ctx, rzNow); got != step.want {
			t.Fatalf("after %q aria-valuenow=%d, want %d", step.key, got, step.want)
		}
	}
	// The width follows the value, not only the ARIA.
	share := eval[float64](t, ctx, `(() => { const f = document.querySelector('#rz-editor-first'); return f.getBoundingClientRect().width / f.parentElement.getBoundingClientRect().width * 100 })()`)
	if share < 19 || share > 21 {
		t.Errorf("the first panel should be ~20%% wide, is %.1f%%", share)
	}
	clearRz(t, ctx)
}

// The drag continues past the component's edge (pointer capture) and stops at
// Max; releasing the button out there still ends the drag.
func TestResizableDragWithPointerCapture(t *testing.T) {
	ctx := openPage(t, "/components/resizable")
	clearRz(t, ctx)
	navigate(t, ctx, "/components/resizable")
	// Below the fold in the test window: a press outside the viewport hits
	// nothing at all.
	eval[bool](t, ctx, `(document.querySelector('`+rzBar+`').scrollIntoView({block: 'center'}), true)`)
	settle(t, ctx)
	bar := rectOf(t, ctx, rzBar)
	root := rectOf(t, ctx, `#rz-editor-first`)
	x, y := bar.Left+bar.Width/2, bar.Top+bar.Height/2

	mouse := func(kind input.MouseType, x, y float64, buttons int64) chromedp.Action {
		return input.DispatchMouseEvent(kind, x, y).WithButton(input.Left).WithButtons(buttons).WithClickCount(1)
	}
	if err := chromedp.Run(ctx,
		mouse(input.MouseMoved, x, y, 0),
		mouse(input.MousePressed, x, y, 1),
		mouse(input.MouseMoved, x+40, y, 1),
		mouse(input.MouseMoved, root.Left+2000, y+300, 1),
		mouse(input.MouseReleased, root.Left+2000, y+300, 0),
	); err != nil {
		t.Fatal(err)
	}
	settle(t, ctx)
	if got := eval[int](t, ctx, rzNow); got != 70 {
		t.Errorf("dragging far past the edge should stop at Max (70), got %d", got)
	}
	if got := eval[string](t, ctx, `document.querySelector('`+rzBar+`').dataset.dragging`); got != "false" {
		t.Errorf("releasing outside the component should end the drag, dragging=%q", got)
	}
	clearRz(t, ctx)
}

func TestResizableRemembersSize(t *testing.T) {
	ctx := openPage(t, "/components/resizable")
	clearRz(t, ctx)
	navigate(t, ctx, "/components/resizable")

	focus(t, ctx, rzBar)
	press(t, ctx, kb.ArrowRight)
	press(t, ctx, kb.ArrowRight)
	settle(t, ctx)
	navigate(t, ctx, "/components/resizable")
	if got := eval[int](t, ctx, rzNow); got != 45 {
		t.Errorf("the size should survive a reload as 45, got %d", got)
	}
	clearRz(t, ctx)
}

func TestResizableIgnoresCorruptStorage(t *testing.T) {
	ctx := openPage(t, "/components/resizable")
	for _, bad := range []string{"abc", "500", "-3", ""} {
		eval[bool](t, ctx, `(localStorage.setItem('malphas-resizable:rz-editor', '`+bad+`'), true)`)
		navigate(t, ctx, "/components/resizable")
		if got := eval[int](t, ctx, rzNow); got != 35 {
			t.Errorf("stored %q should fall back to the default 35, got %d", bad, got)
		}
	}
	clearRz(t, ctx)
}

// The drawer sits deep inside the page here, not at the top of <body>; it must
// still make everything outside it unreachable, not just the body children
// that happen not to contain it.
func TestSidebarDrawerInertsTheRestOfThePage(t *testing.T) {
	ctx := openPage(t, "/components/sidebar")
	if err := chromedp.Run(ctx, chromedp.EmulateViewport(400, 700)); err != nil {
		t.Fatal(err)
	}
	settle(t, ctx)
	click(t, ctx, `[aria-haspopup="dialog"]`)
	settle(t, ctx)
	if !eval[bool](t, ctx, `!!document.querySelector('[aria-haspopup="dialog"]').closest('[inert]')`) {
		t.Error("the Menu button behind the drawer should be inert")
	}
	if !eval[bool](t, ctx, `!!document.querySelector('header').closest('[inert]')`) {
		t.Error("the page header should be inert")
	}
	if eval[bool](t, ctx, `!!document.querySelector('[role="dialog"][aria-modal="true"]').closest('[inert]')`) {
		t.Fatal("the drawer itself must not be inert")
	}
	press(t, ctx, kb.Escape)
	settle(t, ctx)
	if n := eval[int](t, ctx, `document.querySelectorAll('[inert]').length`); n != 0 {
		t.Errorf("closing should release every inert element, %d left", n)
	}
}

// In a right-to-left page the first panel is on the right: the drag and the
// arrows mirror.
func TestResizableMirrorsInRTL(t *testing.T) {
	ctx := openPage(t, "/components/resizable")
	clearRz(t, ctx)
	navigate(t, ctx, "/components/resizable")
	eval[bool](t, ctx, `(document.querySelector('#rz-editor-first').parentElement.setAttribute('dir', 'rtl'), true)`)

	focus(t, ctx, rzBar)
	press(t, ctx, kb.ArrowLeft) // toward the second panel, which is on the left now
	settle(t, ctx)
	if got := eval[int](t, ctx, rzNow); got != 40 {
		t.Errorf("in RTL, ← should grow the first (right-hand) panel to 40, got %d", got)
	}
	clearRz(t, ctx)
}
