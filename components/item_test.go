package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/darkanum/vapulab/malphas/components"
)

func TestEmptyRendersTitleDescriptionAndActions(t *testing.T) {
	html := renderWithChildren(t,
		components.Empty(components.EmptyProps{Title: "Nenhum projeto", Description: "Crie o primeiro."}),
		templ.Raw(`<button>Novo projeto</button>`))
	for _, want := range []string{"Nenhum projeto", "Crie o primeiro.", "Novo projeto"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	// The page decides the heading level; Empty cannot know it.
	if strings.Contains(html, "<h") {
		t.Errorf("Empty must not impose a heading level, got: %s", html)
	}
}

// The stretched-link pattern: the actions are siblings of the link, never
// inside it. A button nested in an <a> is invalid HTML and would navigate.
func TestItemLinkDoesNotWrapTheActions(t *testing.T) {
	html := renderToString(t, components.Item(components.ItemProps{
		Title: "Relatório", Href: "/relatorios/1",
		Actions: templ.Raw(`<button id="acao">Baixar</button>`),
	}))
	linkEnd := strings.Index(html, "</a>")
	action := strings.Index(html, `id="acao"`)
	if linkEnd < 0 || action < linkEnd {
		t.Errorf("the action must come after the link closes, got: %s", html)
	}
	if !strings.Contains(html, `href="/relatorios/1"`) || !strings.Contains(html, "after:absolute") {
		t.Errorf("expected a stretched title link, got: %s", html)
	}
}

func TestItemWithoutHrefHasNoLink(t *testing.T) {
	html := renderToString(t, components.Item(components.ItemProps{Title: "Relatório"}))
	if strings.Contains(html, "<a") {
		t.Errorf("no Href, no link: %s", html)
	}
}

func TestItemGroupIsAList(t *testing.T) {
	html := renderToString(t, components.ItemGroup(components.ItemGroupProps{
		Label: "Arquivos",
		Items: []components.ItemProps{{Title: "a.pdf"}, {Title: "b.pdf"}},
	}))
	if !strings.Contains(html, `<ul aria-label="Arquivos"`) || strings.Count(html, "<li") != 2 {
		t.Errorf("expected a labelled list of 2, got: %s", html)
	}
}
