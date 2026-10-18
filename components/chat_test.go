package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
)

// --- Bubble -------------------------------------------------------------------

func TestBubbleVariantsAndAlign(t *testing.T) {
	start := renderWithChildren(t, components.Bubble(components.BubbleProps{}), templ.Raw("oi"))
	end := renderWithChildren(t, components.Bubble(components.BubbleProps{Variant: "primary", Align: "end"}), templ.Raw("oi"))
	if !strings.Contains(start, "bg-muted/10") || !strings.Contains(start, "self-start") {
		t.Errorf("default bubble should be muted and start-aligned: %s", start)
	}
	if !strings.Contains(end, "bg-primary text-primary-fg") || !strings.Contains(end, "self-end") {
		t.Errorf("primary end bubble: %s", end)
	}
}

// Consecutive messages from one author read as one block: the corner on the
// author's side flattens where bubbles meet.
func TestBubbleGroupCorners(t *testing.T) {
	cases := map[string][]string{
		"first":  {"rounded-bl-md"},
		"middle": {"rounded-tl-md", "rounded-bl-md"},
		"last":   {"rounded-tl-md"},
	}
	for group, want := range cases {
		html := renderWithChildren(t, components.Bubble(components.BubbleProps{Group: group}), templ.Raw("x"))
		for _, c := range want {
			if !strings.Contains(html, c) {
				t.Errorf("%s: missing %s in %s", group, c, html)
			}
		}
	}
	end := renderWithChildren(t, components.Bubble(components.BubbleProps{Group: "first", Align: "end"}), templ.Raw("x"))
	if !strings.Contains(end, "rounded-br-md") {
		t.Errorf("end-aligned groups flatten the right side: %s", end)
	}
}

// An emoji alone reads badly ("thumbs up sign"); each reaction says its count
// and meaning.
func TestBubbleReactionsAreLabelled(t *testing.T) {
	html := renderWithChildren(t, components.Bubble(components.BubbleProps{
		Reactions: []components.BubbleReaction{{Emoji: "👍", Label: "Curtir", Count: 3}},
	}), templ.Raw("x"))
	// The meaning is real text (sr-only), not an aria-label on a list item,
	// which NVDA and JAWS do not reliably read in browse mode.
	for _, want := range []string{`aria-label="Reactions"`, `<span class="sr-only">Curtir: 3</span>`, `aria-hidden="true">👍`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}

// --- Marker -------------------------------------------------------------------

func TestMarkerStatusIsLive(t *testing.T) {
	html := renderWithChildren(t, components.Marker(components.MarkerProps{Status: true, Icon: "●"}), templ.Raw("Ana entrou na conversa"))
	if !strings.Contains(html, `role="status"`) || !strings.Contains(html, `aria-hidden="true">●`) {
		t.Errorf("a status marker is a live region with a decorative icon: %s", html)
	}
}

// Not role="separator": a separator's children are presentational, so the
// label ("Ontem") would vanish for screen readers. The text reads as text;
// only the rules beside it are hidden.
func TestMarkerSeparatorKeepsItsLabelReadable(t *testing.T) {
	html := renderWithChildren(t, components.Marker(components.MarkerProps{Variant: "separator"}), templ.Raw("Ontem"))
	if strings.Contains(html, `role="separator"`) {
		t.Errorf("role=separator would hide the label: %s", html)
	}
	if strings.Count(html, `aria-hidden="true"`) != 2 || !strings.Contains(html, "Ontem") {
		t.Errorf("expected two decorative rules around the readable label: %s", html)
	}
}

// --- Message ------------------------------------------------------------------

func messageFixture(group string) components.MessageProps {
	return components.MessageProps{Author: "Ana Lima", Initials: "AL", Time: "14:02", Datetime: "2026-10-10T14:02", Group: group}
}

func TestMessageHeaderOnlyOnFirst(t *testing.T) {
	first := renderWithChildren(t, components.Message(messageFixture("first")), templ.Raw("oi"))
	middle := renderWithChildren(t, components.Message(messageFixture("middle")), templ.Raw("tudo bem?"))
	if !strings.Contains(first, ">AL<") || !strings.Contains(first, `<time datetime="2026-10-10T14:02"`) {
		t.Errorf("the first message shows avatar and time: %s", first)
	}
	if strings.Contains(middle, ">AL<") || strings.Contains(middle, "<time") {
		t.Errorf("a middle message hides the repeated avatar and header: %s", middle)
	}
}

// Hidden visually is not hidden from screen readers: every message still says
// who sent it and when.
func TestMessageKeepsAuthorAndTimeForScreenReaders(t *testing.T) {
	middle := renderWithChildren(t, components.Message(messageFixture("middle")), templ.Raw("tudo bem?"))
	if !strings.Contains(middle, `aria-label="Ana Lima, 14:02"`) || !strings.Contains(middle, "<article") {
		t.Errorf("every message keeps its author and time: %s", middle)
	}
}

func TestMessageEndAlignmentMirrors(t *testing.T) {
	p := messageFixture("single")
	p.Align = "end"
	html := renderWithChildren(t, components.Message(p), templ.Raw("oi"))
	if !strings.Contains(html, "flex-row-reverse") || !strings.Contains(html, "self-end") {
		t.Errorf("an end-aligned message mirrors: %s", html)
	}
}

// --- Message Scroller ---------------------------------------------------------

func TestMessageScrollerIsALog(t *testing.T) {
	html := renderWithChildren(t, components.MessageScroller(components.MessageScrollerProps{ID: "chat", Label: "Conversa do protocolo"}), templ.Raw("<article>oi</article>"))
	for _, want := range []string{
		`role="log"`,
		`aria-live="polite"`,
		`aria-label="Conversa do protocolo"`,
		`tabindex="0"`,
		`id="chat-content"`,
		`x-data="messageScroller"`,
		"[overflow-anchor:none]",
		"Jump to the end",
		`style="display:none"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}

// Messages carry a hook the scroller counts by; markers and other rows do not.
func TestMessageHasAMessageHook(t *testing.T) {
	html := renderWithChildren(t, components.Message(messageFixture("single")), templ.Raw("oi"))
	if !strings.Contains(html, "data-message") {
		t.Errorf("expected data-message on the article: %s", html)
	}
	marker := renderWithChildren(t, components.Marker(components.MarkerProps{}), templ.Raw("x"))
	if strings.Contains(marker, "data-message") {
		t.Errorf("a marker is not a message: %s", marker)
	}
}
