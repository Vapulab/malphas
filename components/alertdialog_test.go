package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/darkanum/vapulab/malphas/components"
)

func alertDialogFixture() components.AlertDialogProps {
	return components.AlertDialogProps{
		ID: "excluir", Title: "Excluir protocolo?", Description: "Não dá para desfazer.",
		ConfirmLabel: "Excluir", ConfirmAttrs: templ.Attributes{"hx-delete": "/protocolos/1"},
	}
}

func TestAlertDialogRolesAndOrder(t *testing.T) {
	html := renderToString(t, components.AlertDialog(alertDialogFixture()))
	for _, want := range []string{
		`role="alertdialog"`,
		`aria-modal="true"`,
		`aria-labelledby="excluir-title"`,
		`aria-describedby="excluir-description"`,
		`x-data="alertDialog"`,
		`hx-delete="/protocolos/1"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	// Cancel comes first, so modalOverlay's "focus the first focusable" lands
	// on the safe choice, never on the destructive one.
	if c, x := strings.Index(html, ">Cancelar<"), strings.Index(html, ">Excluir<"); c < 0 || x < 0 || c > x {
		t.Errorf("Cancelar must come before the confirm button, got: %s", html)
	}
	// No close button: the only ways out are the two decisions (and Escape).
	if strings.Contains(html, `aria-label="Fechar"`) {
		t.Errorf("an alert dialog offers no × button, got: %s", html)
	}
}
