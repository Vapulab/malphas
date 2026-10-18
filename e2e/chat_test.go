package e2e_test

import (
	"context"
	"strconv"
	"strings"
	"testing"
	"time"
)

// --- Message Scroller ---------------------------------------------------------

const (
	chatLog      = `document.querySelector('#chat-demo-content').parentElement`
	chatAtBottom = `(() => { const v = ` + chatLog + `; return v.scrollHeight - v.scrollTop - v.clientHeight <= 24 })()`
	chatCount    = `document.querySelectorAll('#chat-demo-content article').length`
	chatJump     = `[...document.querySelectorAll('button')].find(b => b.textContent.includes('Jump to latest'))`
	chatWatched  = `document.querySelectorAll('#chat-demo-content article')[5]`
	chatContent  = `document.querySelector('#chat-demo-content')`
)

// scrolledTo puts the reader mid-conversation and returns the watched
// message's position, for "did it move" checks.
// The watched message is kept by reference (window.__watched): inserting
// rows above it shifts its index, and re-querying by index would measure a
// different message. Its position is measured against the log, not the
// window: clicking a demo button scrolls the page itself to bring the button
// into view, which would read as movement in the window's coordinates.
const watchedTop = `(window.__watched.getBoundingClientRect().top - ` + chatLog + `.getBoundingClientRect().top)`

func scrolledTo(t *testing.T, ctx context.Context, top int) float64 {
	t.Helper()
	eval[bool](t, ctx, `(`+chatLog+`.scrollTop = `+strconv.Itoa(top)+`, window.__watched = `+chatWatched+`, true)`)
	settle(t, ctx)
	return eval[float64](t, ctx, watchedTop)
}

func movedBy(t *testing.T, ctx context.Context, before float64) float64 {
	t.Helper()
	return eval[float64](t, ctx, watchedTop) - before
}

func TestScrollerFollowsNewMessagesAtTheBottom(t *testing.T) {
	ctx := openPage(t, "/components/message-scroller")
	if !eval[bool](t, ctx, chatAtBottom) {
		t.Fatal("the conversation should open at its end")
	}
	before := eval[int](t, ctx, chatCount)

	click(t, ctx, "#chat-nova")
	sleep(t, ctx, 500)
	if got := eval[int](t, ctx, chatCount); got != before+1 {
		t.Fatalf("expected one new message, count %d → %d", before, got)
	}
	if !eval[bool](t, ctx, chatAtBottom) {
		t.Error("a reader at the end should follow the new message")
	}
}

// Reading higher up: new messages must not move the view, and the button
// counts them.
func TestScrollerHoldsPositionAndCountsWhenScrolledUp(t *testing.T) {
	ctx := openPage(t, "/components/message-scroller")
	// Mid-conversation, not at the very top: a view at 0 could never move by
	// an append, and the test would pass without the feature.
	before := scrolledTo(t, ctx, 150)

	click(t, ctx, "#chat-nova")
	sleep(t, ctx, 400)
	click(t, ctx, "#chat-nova")
	sleep(t, ctx, 400)
	if diff := movedBy(t, ctx, before); diff > 1 || diff < -1 {
		t.Errorf("the message being read moved by %.1fpx", diff)
	}
	if got := eval[string](t, ctx, `(b => b && getComputedStyle(b).display !== 'none' ? b.textContent.replace(/\s+/g, ' ').trim() : '')(`+chatJump+`)`); !strings.Contains(got, "(2 new)") {
		t.Errorf("the jump button should show 2 new messages, shows %q", got)
	}
}

func TestScrollerJumpButtonGoesToTheEnd(t *testing.T) {
	ctx := openPage(t, "/components/message-scroller")
	eval[bool](t, ctx, `(`+chatLog+`.scrollTop = 0, true)`)
	settle(t, ctx)
	click(t, ctx, "#chat-nova")
	sleep(t, ctx, 400)

	eval[bool](t, ctx, `(`+chatJump+`.click(), true)`)
	sleep(t, ctx, 900) // smooth scroll
	if !eval[bool](t, ctx, chatAtBottom) {
		t.Error("the button should bring the reader to the end")
	}
	if eval[bool](t, ctx, `getComputedStyle(`+chatJump+`).display !== 'none'`) {
		t.Error("the button should hide at the end")
	}
	if !eval[bool](t, ctx, `document.activeElement === `+chatLog) {
		t.Error("focus should stay in the conversation, not fall to <body>")
	}
}

// Older messages entering at the top: whatever the reader was looking at
// stays on the same pixel.
func TestScrollerKeepsPlaceWhenOlderMessagesArrive(t *testing.T) {
	ctx := openPage(t, "/components/message-scroller")
	before := scrolledTo(t, ctx, 150)
	count := eval[int](t, ctx, chatCount)

	click(t, ctx, "#chat-antigas")
	sleep(t, ctx, 500)
	if got := eval[int](t, ctx, chatCount); got != count+3 {
		t.Fatalf("expected three older messages, count %d → %d", count, got)
	}
	if diff := movedBy(t, ctx, before); diff > 1 || diff < -1 {
		t.Errorf("the message being read moved by %.1fpx when history loaded above it", diff)
	}
}

func TestScrollerFollowsStreaming(t *testing.T) {
	ctx := openPage(t, "/components/message-scroller")
	height := eval[float64](t, ctx, chatLog+`.scrollHeight`)

	click(t, ctx, "#chat-stream")
	sleep(t, ctx, 2400)
	if got := eval[float64](t, ctx, chatLog+`.scrollHeight`); got <= height {
		t.Fatalf("the streamed text did not grow the conversation (%v → %v)", height, got)
	}
	if !eval[bool](t, ctx, chatAtBottom) {
		t.Error("a reader at the end should keep seeing the streamed text")
	}
}

// The common htmx "load older" sentinel replaces itself (outerHTML) with the
// older messages plus a new sentinel: the old first element is gone, and the
// reader must still not move.
func TestScrollerKeepsPlaceWhenASentinelReplacesItself(t *testing.T) {
	ctx := openPage(t, "/components/message-scroller")
	before := scrolledTo(t, ctx, 150)
	eval[bool](t, ctx, `(() => { const first = `+chatContent+`.firstElementChild; const older = [...`+chatContent+`.querySelectorAll('article')].slice(0, 3).map(a => a.outerHTML).join(''); first.outerHTML = '<div id="sentinela">load</div>' + older; return true })()`)
	settle(t, ctx)
	if d := movedBy(t, ctx, before); d > 1 || d < -1 {
		t.Errorf("the message being read moved by %.1fpx when the sentinel swapped", d)
	}
}

// Prepend and append in one batch: the view stays, and only the appended
// message counts as new.
func TestScrollerPrependAndAppendInOneBatch(t *testing.T) {
	ctx := openPage(t, "/components/message-scroller")
	before := scrolledTo(t, ctx, 150)
	eval[bool](t, ctx, `(() => { const c = `+chatContent+`; const a = c.querySelector('article'); c.prepend(a.cloneNode(true)); c.append(a.cloneNode(true)); return true })()`)
	settle(t, ctx)
	if d := movedBy(t, ctx, before); d > 1 || d < -1 {
		t.Errorf("the message being read moved by %.1fpx", d)
	}
	if got := eval[string](t, ctx, chatJump+`.textContent.replace(/\s+/g, ' ')`); !strings.Contains(got, "(1 new)") {
		t.Errorf("only the appended message is new, button says %q", got)
	}
}

// Only messages count as new: a "Generating a reply…" marker does not.
func TestScrollerCountsOnlyMessages(t *testing.T) {
	ctx := openPage(t, "/components/message-scroller")
	scrolledTo(t, ctx, 150)
	eval[bool](t, ctx, `(() => { const m = document.createElement('div'); m.setAttribute('role', 'status'); m.textContent = 'Generating a reply…'; `+chatContent+`.append(m); return true })()`)
	settle(t, ctx)
	if eval[bool](t, ctx, `/\(\d/.test(`+chatJump+`.innerText)`) { // innerText: the hidden counter span always holds text
		t.Errorf("a marker was counted as a new message: %q", eval[string](t, ctx, chatJump+`.innerText`))
	}
}

// Growth without a DOM mutation (an image above finishing loading): the
// reader stays on what they were reading.
func TestScrollerKeepsPlaceWhenContentAboveGrows(t *testing.T) {
	ctx := openPage(t, "/components/message-scroller")
	before := scrolledTo(t, ctx, 300)
	// The first message sits wholly above a view scrolled to 300px.
	if !eval[bool](t, ctx, `(a => a.offsetTop + a.offsetHeight < `+chatLog+`.scrollTop)(document.querySelector('#chat-demo-content article'))`) {
		t.Fatal("fixture: the first message should be above the viewport")
	}
	eval[bool](t, ctx, `(document.querySelector('#chat-demo-content article').style.paddingTop = '120px', true)`)
	sleep(t, ctx, 300)
	if d := movedBy(t, ctx, before); d > 1 || d < -1 {
		t.Errorf("the message being read moved by %.1fpx when content above grew", d)
	}
}

// Older history must not be read out as if it just arrived: the log goes
// quiet for history loads and speaks again afterwards.
func TestScrollerHistoryLoadIsNotAnnounced(t *testing.T) {
	ctx := openPage(t, "/components/message-scroller")
	eval[bool](t, ctx, `(window.__live = [], new MutationObserver(() => window.__live.push(`+chatLog+`.getAttribute('aria-live'))).observe(`+chatLog+`, { attributes: true, attributeFilter: ['aria-live'] }), true)`)
	click(t, ctx, "#chat-antigas")
	sleep(t, ctx, 600)
	got := eval[[]string](t, ctx, `window.__live`)
	if len(got) < 2 || got[0] != "off" || got[len(got)-1] != "polite" {
		t.Errorf("the log should go quiet (off) for the history load and return to polite, saw %v", got)
	}
}

// A streamed answer is announced once, when it is complete: it is aria-busy
// while it grows.
func TestScrollerStreamingIsBusyWhileItGrows(t *testing.T) {
	ctx := openPage(t, "/components/message-scroller")
	click(t, ctx, "#chat-stream")
	sleep(t, ctx, 300)
	if got := eval[string](t, ctx, `[...document.querySelectorAll('#chat-demo-content article')].at(-1).getAttribute('aria-busy')`); got != "true" {
		t.Errorf("the streaming message should be aria-busy, got %q", got)
	}
	// The stream is 30 words at 60ms, 1.8s on a fast machine and twice that
	// on a slow one: wait for the end, not for a guess at when it comes.
	if !waitUntil(t, ctx, `![...document.querySelectorAll('#chat-demo-content article')].at(-1).hasAttribute('aria-busy')`, 15*time.Second) {
		t.Errorf("aria-busy should clear when streaming ends, got %q", eval[string](t, ctx, `[...document.querySelectorAll('#chat-demo-content article')].at(-1).getAttribute('aria-busy') || ''`))
	}
}

// A host whose page shell writes an invalid lang (pt_BR, with an underscore)
// must not lose the counter: Intl.PluralRules throws on such a tag.
func TestScrollerCounterSurvivesAnInvalidPageLang(t *testing.T) {
	ctx := openPage(t, "/components/message-scroller")
	eval[bool](t, ctx, `(document.documentElement.lang = "pt_BR", true)`)
	scrolledTo(t, ctx, 150)
	click(t, ctx, "#chat-nova")
	if !waitUntil(t, ctx, `(b => !!b && /\(1 new\)/.test(b.textContent))(`+chatJump+`)`, 5*time.Second) {
		t.Errorf("the counter should still read (1 new), shows %q", eval[string](t, ctx, `(b => b ? b.textContent.replace(/\s+/g, ' ').trim() : '')(`+chatJump+`)`))
	}
}
