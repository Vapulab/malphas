package e2e_test

import (
	"strings"
	"testing"
)

// A confirmation opened over a Modal must be usable with the mouse. The
// Modal had made everything outside itself inert — the confirmation
// included — so its buttons could not be clicked and only Escape got out.
func TestAlertOverModalIsClickable(t *testing.T) {
	ctx := openPage(t, "/components/alert-dialog")
	const modal = `[role="dialog"][aria-labelledby="md-editar-title"]`
	const alert = `[role="alertdialog"][aria-labelledby="ad-descartar-title"]`

	click(t, ctx, "#abrir-md-editar")
	settle(t, ctx)
	click(t, ctx, "#abrir-ad-descartar")
	settle(t, ctx)
	if eval[bool](t, ctx, `!!document.querySelector('`+alert+`').closest('[inert]')`) {
		t.Fatal("the confirmation sits inside an inert ancestor")
	}
	if got := eval[string](t, ctx, `document.activeElement.textContent.trim()`); got != "Cancel" {
		t.Errorf("focus should start on Cancel, is on %q", got)
	}

	click(t, ctx, alert+` button[data-alert-close]`)
	settle(t, ctx)
	if eval[bool](t, ctx, `!!document.querySelector('`+alert+`')`) {
		t.Fatal("clicking Cancel did not close the confirmation")
	}
	if eval[bool](t, ctx, `!!document.querySelector('`+modal+`').closest('[inert]')`) {
		t.Fatal("the modal underneath was left inert")
	}
	click(t, ctx, modal+` [aria-label="Close"]`)
	settle(t, ctx)
	if eval[bool](t, ctx, `!!document.querySelector('`+modal+`')`) {
		t.Error("the modal's close button did not close it")
	}
	if n := eval[int](t, ctx, `document.querySelectorAll('[inert]').length`); n != 0 {
		t.Errorf("%d element(s) left inert after closing everything", n)
	}
}

// Every scroll container in Malphas uses the same thin, themed scrollbar as
// Scroll Area, instead of the platform's default bar.
func TestScrollContainersUseTheThemedScrollbar(t *testing.T) {
	ctx := openPage(t, "/components/todos")
	click(t, ctx, "#abrir-modal-popover")
	settle(t, ctx)
	for _, sel := range []string{
		`[role="dialog"][aria-labelledby="modal-popover-title"]`, // Modal panel
		`#chat-demo-content`, // Message Scroller viewport (its parent)
		`[aria-roledescription="carousel"] [tabindex="0"]`, // Carousel track
	} {
		style := eval[string](t, ctx, `(el => { const s = getComputedStyle(el.id === 'chat-demo-content' ? el.parentElement : el); return s.scrollbarWidth + '|' + s.scrollbarColor })(document.querySelector('`+sel+`'))`)
		if !strings.HasPrefix(style, "thin|") || strings.HasSuffix(style, "|auto") {
			t.Errorf("%s uses the default scrollbar (%s)", sel, style)
		}
	}
}

// Text around an inline component keeps its spaces.
func TestHoverCardSentenceKeepsItsSpaces(t *testing.T) {
	ctx := openPage(t, "/components/hover-card")
	got := eval[string](t, ctx, `document.querySelector('[x-data^="hoverCard"]').closest('p').innerText.replace(/\s+/g, ' ').trim()`)
	if !strings.HasPrefix(got, "Protocolo revisado por @lincoln ontem às 18h.") {
		t.Errorf("the sentence lost its spaces: %q", got)
	}
}

// Bubbles are shown as a conversation, inside a chat box.
func TestBubbleDemoSitsInAConversation(t *testing.T) {
	ctx := openPage(t, "/components/bubble")
	if !eval[bool](t, ctx, `(b => !!b && getComputedStyle(b).borderTopStyle !== 'none')(document.querySelector('[aria-label="Conversa de exemplo"]'))`) {
		t.Error("the bubbles should sit inside a bordered conversation box")
	}
}
