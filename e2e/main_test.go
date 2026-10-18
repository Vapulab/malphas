// Package e2e drives a real browser against the catalog, to cover the one
// thing render tests cannot reach: the behavior that lives in assets/malphas.js.
//
// Arrow keys, focus, Escape and the countdown tick have no Go representation —
// they only exist once Alpine runs. The project takes no Node dependency, so
// there is no jest or vitest either; chromedp speaks CDP directly from Go and
// keeps that constraint intact.
//
// Every test skips, rather than fails, when no Chrome or Chromium is present,
// so a bare CI box still gets a green build from the rest of the suite. Set
// CHROME_PATH to point at a specific binary.
package e2e_test

import (
	"context"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/vapulab/malphas/e2e"
)

var (
	serverURL string
	allocCtx  context.Context
	skipWhy   string
)

// staticDir holds app.css (built from tailwind.config.js) and the vendors.
const staticDir = "static"

func TestMain(m *testing.M) {
	// x-cloak and every utility class live in the built stylesheet. Without it
	// the page still works, but visibility assertions would be meaningless, so
	// say so plainly instead of failing somewhere confusing.
	if _, err := os.Stat(filepath.Join(staticDir, "app.css")); err != nil {
		skipWhy = "e2e/static/app.css is missing — build it with Tailwind (see README)"
	}

	srv := httptest.NewServer(e2e.NewServer(staticDir))
	serverURL = srv.URL

	var cancelAlloc context.CancelFunc
	if skipWhy == "" {
		if path := findChrome(); path == "" {
			skipWhy = "no Chrome or Chromium found — set CHROME_PATH to run browser behaviour tests"
		} else {
			// Chrome's *new* headless, not the old one. Under old headless
			// document.hasFocus() is false, and a browser does not dispatch
			// focus events to an unfocused document — so focus-driven
			// components (Tooltip) would look broken when they are fine.
			opts := append(chromedp.DefaultExecAllocatorOptions[:],
				chromedp.ExecPath(path),
				chromedp.Flag("headless", "new"),
				chromedp.Flag("disable-gpu", true),
			)
			allocCtx, cancelAlloc = chromedp.NewExecAllocator(context.Background(), opts...)
		}
	}

	code := m.Run()

	if cancelAlloc != nil {
		cancelAlloc()
	}
	srv.Close()
	os.Exit(code)
}

func findChrome() string {
	if p := os.Getenv("CHROME_PATH"); p != "" {
		if _, err := os.Stat(p); err == nil {
			return p
		}
		return ""
	}

	var candidates []string
	switch runtime.GOOS {
	case "windows":
		for _, base := range []string{
			os.Getenv("ProgramFiles"),
			os.Getenv("ProgramFiles(x86)"),
			os.Getenv("LOCALAPPDATA"),
		} {
			if base == "" {
				continue
			}
			candidates = append(candidates,
				filepath.Join(base, `Google\Chrome\Application\chrome.exe`),
				filepath.Join(base, `Microsoft\Edge\Application\msedge.exe`),
			)
		}
	case "darwin":
		candidates = []string{
			"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
			"/Applications/Chromium.app/Contents/MacOS/Chromium",
		}
	default:
		candidates = []string{
			"/usr/bin/google-chrome",
			"/usr/bin/chromium",
			"/usr/bin/chromium-browser",
		}
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	for _, name := range []string{"google-chrome", "chromium", "chromium-browser", "chrome"} {
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
	}
	return ""
}

// assetsLoadedJS is true when every same-origin script and stylesheet arrived.
// A failed download shows up in Resource Timing with responseStatus 0.
const assetsLoadedJS = `(() => {
	const failed = performance.getEntriesByType('resource').filter(r =>
		r.name.startsWith(location.origin) && /\.(js|css)(\?|$)/.test(r.name) && r.responseStatus === 0);
	return failed.length === 0 && typeof window.Alpine !== 'undefined';
})()`

// openPage returns a browser context already showing path, with Alpine
// initialised. Every behaviour test starts here.
func openPage(t *testing.T, path string) context.Context {
	t.Helper()
	if skipWhy != "" {
		t.Skip(skipWhy)
	}

	ctx, cancel := chromedp.NewContext(allocCtx)
	t.Cleanup(cancel)
	ctx, cancelTimeout := context.WithTimeout(ctx, 45*time.Second)
	t.Cleanup(cancelTimeout)

	// On Windows, Chrome occasionally fails a localhost download with
	// net::ERR_NO_BUFFER_SPACE (a socket-buffer condition of the OS, not of
	// this app — CI on Linux never sees it). Whichever asset it hits, the page
	// is broken in a way no test is about: no alpine.min.js and nothing
	// initialises (a 15s timeout), no malphas.js and every component is inert
	// (an instant, baffling assertion failure), no app.css and every
	// visibility check lies. A page with a failed asset is reloaded instead;
	// any other failure still fails the test.
	for attempt := 1; ; attempt++ {
		var assetsLoaded bool
		if err := chromedp.Run(ctx,
			chromedp.Navigate(serverURL+path),
			chromedp.WaitReady("body"),
			// Focus events only fire in a focused document; without this the
			// tab stays in the background and Tooltip never opens.
			page.BringToFront(),
			chromedp.Poll(`document.readyState === 'complete'`, nil, chromedp.WithPollingTimeout(15*time.Second)),
			chromedp.Evaluate(assetsLoadedJS, &assetsLoaded),
		); err != nil {
			t.Fatalf("open %s: %v", path, err)
		}
		if assetsLoaded {
			break
		}
		if attempt == 3 {
			t.Fatalf("open %s: an asset failed to load %d times in a row", path, attempt)
		}
		t.Logf("open %s: an asset failed to load (attempt %d), reloading", path, attempt)
	}
	if err := chromedp.Run(ctx, waitForAlpine()); err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	return ctx
}

// waitForAlpine blocks until Alpine has processed the page. Asserting before
// that would test the server-rendered markup, which the render tests already
// cover, and would flake depending on how fast the script loaded.
func waitForAlpine() chromedp.Action {
	const ready = `
		(() => {
			if (!window.Alpine) return false;
			const root = document.querySelector('[x-data]');
			if (!root) return true;             // page has no behaviour to wait for
			try { return Alpine.$data(root) !== undefined; } catch (e) { return false; }
		})()`
	var ok bool
	return chromedp.Poll(ready, &ok, chromedp.WithPollingTimeout(15*time.Second))
}

// navigate loads path inside an already-open context, unlike openPage which
// always starts a fresh tab. Same tab means same origin storage (localStorage,
// cookies) carries over — the point when a test needs to prove something set
// on one page still holds after moving to another.
func navigate(t *testing.T, ctx context.Context, path string) {
	t.Helper()
	if err := chromedp.Run(ctx,
		chromedp.Navigate(serverURL+path),
		chromedp.WaitReady("body"),
		waitForAlpine(),
	); err != nil {
		t.Fatalf("navigate to %s: %v", path, err)
	}
}

// eval runs a JS expression and decodes the result.
func eval[T any](t *testing.T, ctx context.Context, expr string) T {
	t.Helper()
	var out T
	if err := chromedp.Run(ctx, chromedp.Evaluate(expr, &out)); err != nil {
		t.Fatalf("evaluate %s: %v", expr, err)
	}
	return out
}

// press sends a real key event to whatever currently has focus. Dispatching a
// synthetic KeyboardEvent from JS would bypass exactly the browser behaviour
// these tests exist to check.
func press(t *testing.T, ctx context.Context, key string) {
	t.Helper()
	if err := chromedp.Run(ctx, chromedp.KeyEvent(key)); err != nil {
		t.Fatalf("press %q: %v", key, err)
	}
}

func focus(t *testing.T, ctx context.Context, selector string) {
	t.Helper()
	if err := chromedp.Run(ctx, chromedp.Focus(selector, chromedp.ByQuery)); err != nil {
		t.Fatalf("focus %q: %v", selector, err)
	}
}

func click(t *testing.T, ctx context.Context, selector string) {
	t.Helper()
	if err := chromedp.Run(ctx, chromedp.Click(selector, chromedp.ByQuery)); err != nil {
		t.Fatalf("click %q: %v", selector, err)
	}
}

func sleep(t *testing.T, ctx context.Context, ms int) {
	t.Helper()
	if err := chromedp.Run(ctx, chromedp.Sleep(time.Duration(ms)*time.Millisecond)); err != nil {
		t.Fatalf("sleep: %v", err)
	}
}

// settle gives Alpine a tick to apply reactive updates after an interaction.
func settle(t *testing.T, ctx context.Context) {
	t.Helper()
	if err := chromedp.Run(ctx, chromedp.Sleep(150*time.Millisecond)); err != nil {
		t.Fatalf("settle: %v", err)
	}
}
