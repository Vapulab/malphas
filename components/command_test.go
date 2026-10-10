package components_test

import (
	"strings"
	"testing"

	"github.com/darkanum/vapulab/malphas/components"
)

func commandFixture() components.CommandProps {
	return components.CommandProps{
		ID: "cmd", Placeholder: "Busque um comando…", Hotkey: true,
		Groups: []components.CommandGroup{
			{Heading: "Navegar", Items: []components.CommandItem{
				{Label: "Componentes", Href: "/components"},
				{Label: "Landing", Href: "/"},
			}},
			{Heading: "Ações", Items: []components.CommandItem{
				{Label: "Alternar tema", OnClick: "x()", Shortcut: "Ctrl T"},
				{Label: "Sair", OnClick: "y()", Disabled: true},
			}},
		},
	}
}

func TestCommandWiresComboboxToListbox(t *testing.T) {
	html := renderToString(t, components.Command(commandFixture()))
	for _, want := range []string{
		`x-data="command"`,
		`role="combobox"`,
		`aria-controls="cmd-listbox"`,
		`id="cmd-listbox" role="listbox"`,
		`x-on:keydown.window="hotkey($event)"`,
		`aria-modal="true"`,
		`data-href="/components"`,
		`aria-disabled="true"`,
		`Ctrl T`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	// One index space across groups, so the arrows walk straight through.
	for i, label := range []string{"Componentes", "Landing", "Alternar tema", "Sair"} {
		if !strings.Contains(html, `data-option-index="`+string(rune('0'+i))+`" data-label="`+label+`"`) {
			t.Errorf("option %q should have index %d, got: %s", label, i, html)
		}
	}
}

func TestCommandGroupsAreLabelled(t *testing.T) {
	html := renderToString(t, components.Command(commandFixture()))
	if !strings.Contains(html, `role="group" aria-labelledby="cmd-group-0"`) || !strings.Contains(html, `id="cmd-group-1"`) {
		t.Errorf("each group needs its heading as a label, got: %s", html)
	}
}

func TestCommandWithoutHotkeyRegistersNoShortcut(t *testing.T) {
	props := commandFixture()
	props.Hotkey = false
	html := renderToString(t, components.Command(props))
	if strings.Contains(html, "hotkey(") {
		t.Errorf("Hotkey false must not listen for Ctrl+K, got: %s", html)
	}
}
