package e2e_test

import (
	"context"
	"math"
	"strconv"
	"strings"
	"testing"
)

// relativeLuminance implements the WCAG 2.x relative luminance formula:
// https://www.w3.org/TR/WCAG21/#dfn-relative-luminance
func relativeLuminance(r, g, b int) float64 {
	channel := func(c int) float64 {
		cs := float64(c) / 255
		if cs <= 0.03928 {
			return cs / 12.92
		}
		return math.Pow((cs+0.055)/1.055, 2.4)
	}
	return 0.2126*channel(r) + 0.7152*channel(g) + 0.0722*channel(b)
}

// contrastRatio follows the same spec: (L1+0.05)/(L2+0.05), L1 being the
// lighter color. It is symmetric on purpose — (text, background) gives the same
// ratio as (background, text), so it does not matter which color is which.
func contrastRatio(a, b [3]int) float64 {
	la, lb := relativeLuminance(a[0], a[1], a[2]), relativeLuminance(b[0], b[1], b[2])
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

// parseRGB reads the format getComputedStyle returns: "rgb(r, g, b)" or
// "rgba(r, g, b, a)". The alpha channel, when present, is ignored on purpose —
// the elements this file measures (text, solid fills) are opaque; handling
// alpha would mean compositing against the real background, which is more than
// the test needs to prove.
func parseRGB(t *testing.T, s string) [3]int {
	t.Helper()
	inner := strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(s, "rgba("), "rgb("), ")")
	parts := strings.Split(inner, ",")
	if len(parts) < 3 {
		t.Fatalf("could not parse the color %q", s)
	}
	var out [3]int
	for i := 0; i < 3; i++ {
		n, err := strconv.Atoi(strings.TrimSpace(parts[i]))
		if err != nil {
			t.Fatalf("could not parse the color %q: %v", s, err)
		}
		out[i] = n
	}
	return out
}

// openPageForcedDark opens path with the theme explicitly dark, written to
// localStorage before navigating — not with an emulated prefers-color-scheme.
// An explicit value in localStorage beats the OS both ways (see
// TestExplicitThemeBeatsTheOperatingSystem in theme_test.go), so this works
// whatever another test left in the same Chrome profile, and whatever the
// color scheme of the machine running the suite.
// t.Cleanup removes the key afterwards: without it, every test that ran later
// would silently inherit the dark theme.
func openPageForcedDark(t *testing.T, path string) context.Context {
	t.Helper()
	ctx := openPage(t, "/components")
	t.Cleanup(func() {
		eval[bool](t, ctx, `(() => { localStorage.removeItem('malphas-theme'); return true; })()`)
	})
	eval[bool](t, ctx, `(() => { localStorage.setItem('malphas-theme', 'dark'); return true; })()`)
	navigate(t, ctx, path)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.documentElement.dataset.theme || ''`); got != "dark" {
		t.Fatalf("could not force the dark theme: data-theme=%q", got)
	}
	return ctx
}

// minAA is the WCAG 2.x floor for normal text (level AA). Large text and UI
// boundaries have a 3:1 floor, but measuring against the stricter one covers
// both.
const minAA = 4.5

// TestDarkThemePrimaryTextMeetsAA guards the contrast rendering_test.go never
// covered: never in the dark theme, and never as an actual ratio — only colors
// compared with each other (mixed !== plain), which does not catch a color that
// is technically different but still unreadable. --color-primary had no
// override in the dark block of tokens.css before this fix: the light value
// (indigo-700) gives ~2.2:1 on the dark surface, well below the AA floor. The
// catalog index has a text-primary link right at the top, so it is enough.
func TestDarkThemePrimaryTextMeetsAA(t *testing.T) {
	ctx := openPageForcedDark(t, "/components")

	textColor := eval[string](t, ctx, `getComputedStyle(document.querySelector('.text-primary')).color`)
	bgColor := eval[string](t, ctx, `getComputedStyle(document.body).backgroundColor`)

	ratio := contrastRatio(parseRGB(t, textColor), parseRGB(t, bgColor))
	if ratio < minAA {
		t.Errorf("text-primary on the dark surface: %.2f:1, below the AA floor of %.1f:1 (text %s, background %s)",
			ratio, minAA, textColor, bgColor)
	}
}

// TestDarkThemePrimaryFillMeetsAA covers the other side of the same trap:
// raising --color-primary alone fixes text-primary and breaks every
// bg-primary text-primary-fg fill (Button, solid Badge, Pagination, done Steps,
// Checkbox, Toggle) — the original white would drop to ~2.9:1 on the lighter
// blue. --color-primary-fg has to move with it, in the same tokens.css blocks.
// /components/button renders a "primary" Button (bg-primary text-primary-fg)
// in the page's first demo.
func TestDarkThemePrimaryFillMeetsAA(t *testing.T) {
	ctx := openPageForcedDark(t, "/components/button")

	textColor := eval[string](t, ctx, `getComputedStyle(document.querySelector('.text-primary-fg')).color`)
	bgColor := eval[string](t, ctx, `getComputedStyle(document.querySelector('.text-primary-fg')).backgroundColor`)

	ratio := contrastRatio(parseRGB(t, textColor), parseRGB(t, bgColor))
	if ratio < minAA {
		t.Errorf("text-primary-fg on bg-primary: %.2f:1, below the AA floor of %.1f:1 (text %s, background %s)",
			ratio, minAA, textColor, bgColor)
	}
}

// The same gap as primary, now in danger: the dark block had no override, and
// the light red (red-600) gives ~3.9:1 on the dark surface. The Field error is
// text the person has to read to fix the field.
func TestDarkThemeDangerTextMeetsAA(t *testing.T) {
	ctx := openPageForcedDark(t, "/components/field")

	textColor := eval[string](t, ctx, `getComputedStyle(document.querySelector('#f-site-error')).color`)
	bgColor := eval[string](t, ctx, `getComputedStyle(document.body).backgroundColor`)

	ratio := contrastRatio(parseRGB(t, textColor), parseRGB(t, bgColor))
	if ratio < minAA {
		t.Errorf("text-danger on the dark surface: %.2f:1, below the AA floor of %.1f:1 (text %s, background %s)",
			ratio, minAA, textColor, bgColor)
	}
}

// And the other side, as with primary: lightening --color-danger alone would
// break the white on the danger Button. /components/button has a "danger"
// Button.
func TestDarkThemeDangerFillMeetsAA(t *testing.T) {
	ctx := openPageForcedDark(t, "/components/button")

	textColor := eval[string](t, ctx, `getComputedStyle(document.querySelector('.text-danger-fg')).color`)
	bgColor := eval[string](t, ctx, `getComputedStyle(document.querySelector('.text-danger-fg')).backgroundColor`)

	ratio := contrastRatio(parseRGB(t, textColor), parseRGB(t, bgColor))
	if ratio < minAA {
		t.Errorf("text-danger-fg on bg-danger: %.2f:1, below the AA floor of %.1f:1 (text %s, background %s)",
			ratio, minAA, textColor, bgColor)
	}
}
