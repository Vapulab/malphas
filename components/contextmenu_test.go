package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
)

func contextMenuFixture() components.ContextMenuProps {
	return components.ContextMenuProps{
		ID: "arquivo", Label: "File actions",
		Entries: []components.MenuEntry{
			{Label: "Rename", OnClick: "$dispatch('renomear')"},
			components.MenuDivider(),
			{Label: "Delete", OnClick: "$dispatch('excluir')"},
		},
	}
}

// A right-click menu is invisible to keyboard users unless the target can be
// focused and says how to open it.
func TestContextMenuTargetIsReachable(t *testing.T) {
	html := renderWithChildren(t, components.ContextMenu(contextMenuFixture()), templ.Raw("<p>report.pdf</p>"))
	for _, want := range []string{
		`tabindex="0"`,
		`aria-describedby="arquivo-hint"`,
		`id="arquivo-hint"`,
		`Shift+F10`,
		`x-data="contextMenu"`,
		`role="menu"`,
		`aria-label="File actions"`,
		`role="separator"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	if strings.Count(html, `role="menuitem"`) != 2 {
		t.Errorf("expected 2 items, got: %s", html)
	}
}
