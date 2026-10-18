package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
)

func TestPopoverWiresTriggerToADialog(t *testing.T) {
	html := renderWithChildren(t,
		components.Popover(components.PopoverProps{ID: "filtros", Label: "Filtros"}),
		templ.Raw(`<input id="q">`))
	for _, want := range []string{
		`aria-haspopup="dialog"`,
		`aria-expanded="false"`,
		`aria-controls="filtros-panel"`,
		`id="filtros-panel" role="dialog"`,
		`aria-label="Filtros"`,
		`style="display:none"`,
		`x-data="popover"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	// Non-modal: the page behind stays usable, so it must not claim otherwise.
	if strings.Contains(html, "aria-modal") {
		t.Errorf("a popover is not modal, got: %s", html)
	}
}

func TestPopoverPosition(t *testing.T) {
	html := renderWithChildren(t,
		components.Popover(components.PopoverProps{ID: "p", Label: "P", Position: "bottom-end"}), templ.Raw("x"))
	if !strings.Contains(html, "right-0") {
		t.Errorf("bottom-end should anchor to the right, got: %s", html)
	}
}

func TestPopoverTitleNamesTheDialog(t *testing.T) {
	html := renderWithChildren(t,
		components.Popover(components.PopoverProps{ID: "p", Label: "Abrir", Title: "Filtros da busca"}), templ.Raw("x"))
	if !strings.Contains(html, `aria-label="Filtros da busca"`) {
		t.Errorf("Title should name the dialog, got: %s", html)
	}
}
