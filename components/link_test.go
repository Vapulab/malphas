package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/vapulab/malphas/components"
)

func TestLinkRendersHref(t *testing.T) {
	html := renderWithChildren(t,
		components.Link(components.LinkProps{Href: "/todos"}),
		templ.Raw("Todos"))
	if !strings.Contains(html, `href="/todos"`) {
		t.Errorf("expected the href, got: %s", html)
	}
	if !strings.Contains(html, "Todos") {
		t.Errorf("expected children rendered, got: %s", html)
	}
}

// templ sanitises URLs it considers unsafe; the component must not defeat that.
func TestLinkSanitisesJavascriptURL(t *testing.T) {
	html := renderWithChildren(t,
		components.Link(components.LinkProps{Href: "javascript:alert(1)"}),
		templ.Raw("x"))
	if strings.Contains(html, "javascript:alert") {
		t.Errorf("javascript: URL must not survive into the output, got: %s", html)
	}
}

func TestLinkExternalGetsSafeRel(t *testing.T) {
	html := renderWithChildren(t,
		components.Link(components.LinkProps{Href: "https://go.dev", External: true}),
		templ.Raw("Go"))
	if !strings.Contains(html, `target="_blank"`) {
		t.Errorf("expected target=_blank, got: %s", html)
	}
	// Without noopener the opened page can reach back through window.opener.
	if !strings.Contains(html, "noopener") {
		t.Errorf("external link must set rel=noopener, got: %s", html)
	}
}
