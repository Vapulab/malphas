package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

func breadcrumbFixture() components.BreadcrumbsProps {
	return components.BreadcrumbsProps{Items: []components.BreadcrumbItem{
		{Label: "Início", Href: "/"},
		{Label: "Componentes", Href: "/components"},
		{Label: "Breadcrumbs"},
	}}
}

func TestBreadcrumbsIsANavigationLandmark(t *testing.T) {
	html := renderToString(t, components.Breadcrumbs(breadcrumbFixture()))
	if !strings.Contains(html, "<nav") || !strings.Contains(html, "aria-label") {
		t.Errorf("expected a labelled nav landmark, got: %s", html)
	}
	if !strings.Contains(html, "<ol") {
		t.Errorf("breadcrumbs are an ordered list, got: %s", html)
	}
}

// The last crumb is where you already are: it must not be a link, and it must
// be marked as the current page.
func TestBreadcrumbsLastItemIsCurrentAndNotALink(t *testing.T) {
	html := renderToString(t, components.Breadcrumbs(breadcrumbFixture()))
	if !strings.Contains(html, `aria-current="page"`) {
		t.Errorf("expected aria-current on the last crumb, got: %s", html)
	}
	if strings.Count(html, "<a ") != 2 {
		t.Errorf("expected exactly 2 links for 3 crumbs, got: %s", html)
	}
}

// An item with no Href is not clickable even in the middle of the trail.
func TestBreadcrumbsItemWithoutHrefIsNotALink(t *testing.T) {
	html := renderToString(t, components.Breadcrumbs(components.BreadcrumbsProps{
		Items: []components.BreadcrumbItem{{Label: "Sem link"}, {Label: "Atual"}},
	}))
	if strings.Contains(html, "<a ") {
		t.Errorf("expected no links at all, got: %s", html)
	}
}
