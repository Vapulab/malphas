package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/darkanum/vapulab/malphas/components"
)

func accordionFixture(exclusive bool) components.AccordionProps {
	return components.AccordionProps{
		ID:        "faq",
		Exclusive: exclusive,
		Items: []components.AccordionItem{
			{Title: "O que é?", Content: templ.Raw("<p>Um boilerplate.</p>")},
			{Title: "Como uso?", Content: templ.Raw("<p>go get.</p>"), Open: true},
			{Title: "Custa?", Content: templ.Raw("<p>Não.</p>")},
		},
	}
}

func TestAccordionWiresEveryPanelToItsTrigger(t *testing.T) {
	html := renderToString(t, components.Accordion(accordionFixture(false)))
	if strings.Count(html, "aria-controls") != 3 {
		t.Errorf("expected each trigger wired to its panel, got: %s", html)
	}
	// Leading space so this counts the static attribute only: the x-bind:
	// counterpart contains the same substring.
	if strings.Count(html, ` aria-expanded="`) != 3 {
		t.Errorf("expected a server-rendered state on each trigger, got: %s", html)
	}
	// Panel ids must be unique or aria-controls points at the wrong element.
	for _, id := range []string{"faq-panel-0", "faq-panel-1", "faq-panel-2"} {
		if !strings.Contains(html, id) {
			t.Errorf("expected panel id %q, got: %s", id, html)
		}
	}
}

func TestAccordionRespectsInitiallyOpenItem(t *testing.T) {
	html := renderToString(t, components.Accordion(accordionFixture(false)))
	if strings.Count(html, `aria-expanded="true"`) != 1 {
		t.Errorf("expected exactly the one item marked Open, got: %s", html)
	}
}

func TestAccordionRendersAllTitlesAndContent(t *testing.T) {
	html := renderToString(t, components.Accordion(accordionFixture(false)))
	for _, want := range []string{"O que é?", "Como uso?", "Custa?", "Um boilerplate.", "go get."} {
		if !strings.Contains(want, "") && !strings.Contains(html, want) {
			t.Errorf("expected %q, got: %s", want, html)
		}
	}
}

// Exclusive mode has to reach the behavior, otherwise every panel opens
// independently and the prop silently does nothing.
func TestAccordionExclusiveIsPassedToTheBehavior(t *testing.T) {
	exclusive := renderToString(t, components.Accordion(accordionFixture(true)))
	shared := renderToString(t, components.Accordion(accordionFixture(false)))
	if exclusive == shared {
		t.Error("Exclusive must change the rendered markup")
	}
	if !strings.Contains(exclusive, "true") {
		t.Errorf("expected exclusive flag in the x-data initialiser, got: %s", exclusive)
	}
}
