package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
)

func TestEmptyRendersTitleDescriptionAndActions(t *testing.T) {
	html := renderWithChildren(t,
		components.Empty(components.EmptyProps{Title: "No projects", Description: "Create the first one."}),
		templ.Raw(`<button>New project</button>`))
	for _, want := range []string{"No projects", "Create the first one.", "New project"} {
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
		Title: "Report", Href: "/relatorios/1",
		Actions: templ.Raw(`<button id="acao">Download</button>`),
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
	html := renderToString(t, components.Item(components.ItemProps{Title: "Report"}))
	if strings.Contains(html, "<a") {
		t.Errorf("no Href, no link: %s", html)
	}
}

func TestItemGroupIsAList(t *testing.T) {
	html := renderToString(t, components.ItemGroup(components.ItemGroupProps{
		Label: "Files",
		Items: []components.ItemProps{{Title: "a.pdf"}, {Title: "b.pdf"}},
	}))
	if !strings.Contains(html, `<ul aria-label="Files"`) || strings.Count(html, "<li") != 2 {
		t.Errorf("expected a labelled list of 2, got: %s", html)
	}
}

func TestProseAppliesTheClass(t *testing.T) {
	html := renderWithChildren(t, components.Prose(components.ProseProps{}), templ.Raw("<h2>Title</h2>"))
	if !strings.Contains(html, `class="prose-malphas"`) {
		t.Errorf("expected the prose class, got: %s", html)
	}
}

// A link with no text has no accessible name. Without a Title there is
// nothing to put in the link, so the item is not a link.
func TestItemWithHrefButNoTitleIsNotALink(t *testing.T) {
	html := renderToString(t, components.Item(components.ItemProps{Href: "/x", Description: "description only"}))
	if strings.Contains(html, "<a") {
		t.Errorf("an empty link would have no name, got: %s", html)
	}
}
