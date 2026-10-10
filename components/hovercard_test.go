package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/darkanum/vapulab/malphas/components"
)

// The card is supplemental: the trigger is a real link, so everything the card
// shows is also reachable without hovering.
func TestHoverCardTriggerIsALink(t *testing.T) {
	html := renderWithChildren(t,
		components.HoverCard(components.HoverCardProps{ID: "hc", Label: "@lincoln", Href: "/u/lincoln"}),
		templ.Raw("<p>Pesquisador</p>"))
	for _, want := range []string{
		`<a href="/u/lincoln"`,
		`x-data="hoverCard(400, 200)"`,
		`id="hc-card"`,
		`style="display:none"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}

// A hover card sits inside running text, usually a <p>. A <div> there makes
// the HTML parser close the paragraph and hoist the card out of the component,
// where Alpine never binds it: open becomes true and nothing shows.
func TestHoverCardIsPhrasingContent(t *testing.T) {
	html := renderWithChildren(t,
		components.HoverCard(components.HoverCardProps{ID: "hc", Label: "@x", Href: "/x"}), templ.Raw("y"))
	if strings.Contains(html, "<div") {
		t.Errorf("a <div> cannot live inside a paragraph, got: %s", html)
	}
}
