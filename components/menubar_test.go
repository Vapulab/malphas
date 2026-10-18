package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

func menubarFixture() components.MenubarProps {
	return components.MenubarProps{
		Label: "Editor",
		Menus: []components.MenubarMenu{
			{Label: "Arquivo", Entries: []components.MenuEntry{{Label: "Novo", OnClick: "x"}, {Label: "Abrir", OnClick: "x"}}},
			{Label: "Editar", Entries: []components.MenuEntry{{Label: "Desfazer", OnClick: "x"}, components.MenuDivider(), {Label: "Copiar", OnClick: "x"}}},
		},
	}
}

func TestMenubarRolesAndSingleTabStop(t *testing.T) {
	html := renderToString(t, components.Menubar(menubarFixture()))
	for _, want := range []string{
		`role="menubar"`,
		`aria-label="Editor"`,
		`aria-haspopup="menu"`,
		`aria-expanded="false"`,
		`x-data="menubar(2)"`,
		`role="separator"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	// One tab stop for the whole bar: the first top item. Submenu items are
	// reached with the arrows only.
	if n := strings.Count(html, `tabindex="0"`); n != 1 {
		t.Errorf("expected exactly one tab stop, got %d in: %s", n, html)
	}
	if n := strings.Count(html, `role="menu"`); n != 2 {
		t.Errorf("expected one submenu per top item, got %d", n)
	}
	if n := strings.Count(html, `role="menuitem"`); n != 2+4 {
		t.Errorf("expected 2 top items and 4 submenu items, got %d", n)
	}
}
