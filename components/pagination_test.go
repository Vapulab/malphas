package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/vapulab/malphas/components"
)

func TestPaginationIsANavigationLandmark(t *testing.T) {
	html := renderToString(t, components.Pagination(components.PaginationProps{
		Current: 2, Total: 5, BaseHref: "/todos?page=",
	}))
	if !strings.Contains(html, "<nav") || !strings.Contains(html, "aria-label") {
		t.Errorf("expected a labelled nav, got: %s", html)
	}
	if !strings.Contains(html, `href="/todos?page=3"`) {
		t.Errorf("expected the page number appended to BaseHref, got: %s", html)
	}
}

// Linking to the page you are already on is a dead control; aria-current is
// what actually announces the position.
func TestPaginationCurrentPageIsNotALink(t *testing.T) {
	html := renderToString(t, components.Pagination(components.PaginationProps{
		Current: 2, Total: 5, BaseHref: "/p/",
	}))
	if strings.Count(html, `aria-current="page"`) != 1 {
		t.Errorf("expected exactly one current page, got: %s", html)
	}
	if strings.Contains(html, `href="/p/2"`) {
		t.Errorf("the current page must not be a link, got: %s", html)
	}
}

// A link that goes nowhere still takes focus and still looks clickable.
func TestPaginationDisabledStepsAreNotLinks(t *testing.T) {
	first := renderToString(t, components.Pagination(components.PaginationProps{
		Current: 1, Total: 5, BaseHref: "/p/",
	}))
	if !strings.Contains(first, `aria-disabled="true"`) {
		t.Errorf("on page 1, Previous must be disabled, got: %s", first)
	}
	if strings.Contains(first, `href="/p/0"`) {
		t.Errorf("Previous must not link to page 0, got: %s", first)
	}

	last := renderToString(t, components.Pagination(components.PaginationProps{
		Current: 5, Total: 5, BaseHref: "/p/",
	}))
	if strings.Contains(last, `href="/p/6"`) {
		t.Errorf("Next must not link past the last page, got: %s", last)
	}
}

// Short ranges show every page; long ones collapse around the current.
func TestPaginationCollapsesLongRanges(t *testing.T) {
	short := renderToString(t, components.Pagination(components.PaginationProps{
		Current: 3, Total: 7, BaseHref: "/p/",
	}))
	if strings.Contains(short, "&hellip;") {
		t.Errorf("7 pages fit without an ellipsis, got: %s", short)
	}

	long := renderToString(t, components.Pagination(components.PaginationProps{
		Current: 10, Total: 20, BaseHref: "/p/",
	}))
	if !strings.Contains(long, "&hellip;") {
		t.Errorf("expected the range collapsed, got: %s", long)
	}
	// First and last stay reachable so the ends are always one click away.
	if !strings.Contains(long, `>1</a>`) || !strings.Contains(long, `>20</a>`) {
		t.Errorf("first and last page must always be present, got: %s", long)
	}
}
