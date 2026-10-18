package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
)

func navbarFixture(drawerID string) components.NavbarProps {
	return components.NavbarProps{
		Brand:    "Vapulab",
		DrawerID: drawerID,
		Entries: []components.MenuEntry{
			{Label: "Components", Href: "/", Current: true},
			{Label: "All", Href: "/todos"},
		},
	}
}

func navbarHTML(t *testing.T, drawerID string) string {
	t.Helper()
	return renderWithChildren(t, components.Navbar(navbarFixture(drawerID)), templ.Raw(""))
}

func TestNavbarIsANavigationLandmark(t *testing.T) {
	html := navbarHTML(t, "")
	if !strings.Contains(html, "<header") || !strings.Contains(html, "<nav") {
		t.Errorf("expected a header with a nav landmark, got: %s", html)
	}
	if !strings.Contains(html, "<ul") {
		t.Errorf("navigation links belong in a list, got: %s", html)
	}
}

// Colour alone does not tell a screen reader which page you are on.
func TestNavbarMarksTheCurrentPage(t *testing.T) {
	html := navbarHTML(t, "")
	if strings.Count(html, `aria-current="page"`) != 1 {
		t.Errorf("expected exactly one current entry, got: %s", html)
	}
}

// Navbar knows how to open a Drawer but does not require one. Without a
// DrawerID it must not render a button that dispatches at nothing.
func TestNavbarMenuButtonOnlyWithADrawer(t *testing.T) {
	without := navbarHTML(t, "")
	if strings.Contains(without, "open-drawer") {
		t.Errorf("no DrawerID should mean no menu button, got: %s", without)
	}

	with := navbarHTML(t, "mobile-nav")
	if !strings.Contains(with, `open-drawer`) || !strings.Contains(with, `mobile-nav`) {
		t.Errorf("expected the button to target the given drawer, got: %s", with)
	}
	if !strings.Contains(with, `aria-label="Open menu"`) {
		t.Errorf("a hamburger glyph needs a label, got: %s", with)
	}
	if !strings.Contains(with, `aria-haspopup="dialog"`) {
		t.Errorf("expected aria-haspopup so the button announces what it opens, got: %s", with)
	}
}

// Dividers make sense in a dropdown, not in a horizontal bar.
func TestNavbarSkipsDividers(t *testing.T) {
	props := navbarFixture("")
	props.Entries = append(props.Entries, components.MenuDivider())
	html := renderWithChildren(t, components.Navbar(props), templ.Raw(""))
	if strings.Count(html, "<li") != 2 {
		t.Errorf("expected 2 items and no divider row, got: %s", html)
	}
}
