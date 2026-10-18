package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

func dropdownFixture() components.DropdownProps {
	return components.DropdownProps{
		ID:    "acoes",
		Label: "Ações",
		Entries: []components.MenuEntry{
			{Label: "Editar", OnClick: "$dispatch('edit')"},
			{Label: "Duplicar", Href: "/duplicar"},
			components.MenuDivider(),
			{Label: "Excluir", OnClick: "$dispatch('del')", Disabled: true},
		},
	}
}

// Unlike Menu, a dropdown really is an application menu, so it claims the role.
func TestDropdownUsesMenuRoles(t *testing.T) {
	html := renderToString(t, components.Dropdown(dropdownFixture()))
	if !strings.Contains(html, `role="menu"`) {
		t.Errorf("expected role=menu, got: %s", html)
	}
	if strings.Count(html, `role="menuitem"`) != 3 {
		t.Errorf("expected 3 menuitems (the divider is not one), got: %s", html)
	}
	if !strings.Contains(html, `role="separator"`) {
		t.Errorf("expected the divider as a separator, got: %s", html)
	}
}

// The trigger has to say what it opens and whether it is open.
func TestDropdownTriggerIsWired(t *testing.T) {
	html := renderToString(t, components.Dropdown(dropdownFixture()))
	for _, want := range []string{`aria-haspopup="menu"`, `aria-expanded="false"`, `aria-controls="acoes-menu"`, `id="acoes-menu"`} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %s, got: %s", want, html)
		}
	}
}

// Roving tabindex: the menu is reached once, then arrows move within it.
func TestDropdownHasOneTabbableItem(t *testing.T) {
	html := renderToString(t, components.Dropdown(dropdownFixture()))
	if got := strings.Count(html, `tabindex="0"`); got != 1 {
		t.Errorf("exactly one item may be tabbable, got %d: %s", got, html)
	}
}

func TestDropdownDisabledItemOptsOutOfRoving(t *testing.T) {
	html := renderToString(t, components.Dropdown(dropdownFixture()))
	if !strings.Contains(html, "data-roving-disabled") {
		t.Errorf("a disabled item must be skipped by the arrows, got: %s", html)
	}
}

// This is the assertion that separates Dropdown from Drawer: a menu is not
// modal, so it must not use the overlay primitive.
func TestDropdownIsNotModal(t *testing.T) {
	html := renderToString(t, components.Dropdown(dropdownFixture()))
	for _, forbidden := range []string{"openOverlay", "closeOverlay", `aria-modal`, "onTab"} {
		if strings.Contains(html, forbidden) {
			t.Errorf("a dropdown must not trap focus; found %s in: %s", forbidden, html)
		}
	}
	if !strings.Contains(html, "onMenuKeydown") {
		t.Errorf("expected the non-modal key handler, got: %s", html)
	}
}

func TestDropdownPosition(t *testing.T) {
	props := dropdownFixture()
	props.Position = "top-end"
	html := renderToString(t, components.Dropdown(props))
	if !strings.Contains(html, "bottom-full") || !strings.Contains(html, "right-0") {
		t.Errorf("expected the menu placed above and end-aligned, got: %s", html)
	}
}
