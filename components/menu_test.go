package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/vapulab/malphas/components"
)

func menuFixture() components.MenuProps {
	return components.MenuProps{
		Entries: []components.MenuEntry{
			{Label: "Início", Href: "/", Current: true},
			{Label: "Componentes", Href: "/components"},
			components.MenuDivider(),
			{Label: "Sair", OnClick: "$dispatch('logout')"},
		},
	}
}

// Site navigation is a list of links, not an application menu. role="menu"
// here would make a screen reader announce the site as a program menu and
// expect arrow-key navigation that does not exist.
func TestMenuUsesNavigationSemanticsNotMenuRole(t *testing.T) {
	html := renderToString(t, components.Menu(menuFixture()))
	if strings.Contains(html, `role="menu"`) || strings.Contains(html, `role="menuitem"`) {
		t.Errorf("site navigation must not claim the menu role, got: %s", html)
	}
	if !strings.Contains(html, "<ul") || !strings.Contains(html, "<li") {
		t.Errorf("expected a real list, got: %s", html)
	}
}

// The row lights up on hover, so the whole row must be clickable. Wrapping the
// link in a padded span made the target 303x40 to the eye and 35x18 to the
// pointer — the kind of bug that only shows up when someone tries to use it.
func TestMenuLinkFillsTheRow(t *testing.T) {
	html := renderToString(t, components.Menu(menuFixture()))
	if strings.Contains(html, `<span class="block rounded px-3 py-2`) {
		t.Errorf("the padding must be on the link, not on a wrapper around it: %s", html)
	}
	if !strings.Contains(html, `class="block w-full rounded px-3 py-2`) {
		t.Errorf("expected the link itself to fill the row, got: %s", html)
	}
}

func TestMenuMarksCurrentPage(t *testing.T) {
	html := renderToString(t, components.Menu(menuFixture()))
	if strings.Count(html, `aria-current="page"`) != 1 {
		t.Errorf("expected exactly one current entry, got: %s", html)
	}
}

// A divider is a separator, not a focusable list item.
func TestMenuRendersDividerAsSeparator(t *testing.T) {
	html := renderToString(t, components.Menu(menuFixture()))
	if !strings.Contains(html, `role="separator"`) {
		t.Errorf("expected the divider to be a separator, got: %s", html)
	}
	if strings.Count(html, "<a ") != 2 {
		t.Errorf("expected 2 links: the divider and the action entry are not links, got: %s", html)
	}
}

// Entries with no Href act rather than navigate, so they are buttons.
func TestMenuActionEntryIsAButton(t *testing.T) {
	html := renderToString(t, components.Menu(menuFixture()))
	// templ escapes the quotes inside the attribute, so match on the payload
	// rather than the literal expression.
	if !strings.Contains(html, "<button") || !strings.Contains(html, "logout") {
		t.Errorf("expected the action entry rendered as a button, got: %s", html)
	}
	if !strings.Contains(html, "x-on:click") {
		t.Errorf("expected the click handler wired, got: %s", html)
	}
}

// Nesting navigation landmarks adds noise; the wrapper is opt-in.
func TestMenuLandmarkIsOptional(t *testing.T) {
	plain := renderToString(t, components.Menu(menuFixture()))
	if strings.Contains(plain, "<nav") {
		t.Errorf("Menu should not be a landmark by default, got: %s", plain)
	}

	props := menuFixture()
	props.Landmark = true
	props.Label = "Principal"
	landmark := renderToString(t, components.Menu(props))
	if !strings.Contains(landmark, "<nav") || !strings.Contains(landmark, `aria-label="Principal"`) {
		t.Errorf("expected a labelled nav landmark, got: %s", landmark)
	}
}
