package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
)

func modalHTML(t *testing.T) string {
	t.Helper()
	return renderWithChildren(t,
		components.Modal(components.ModalProps{ID: "confirm", Title: "Confirmar"}),
		templ.Raw("<p>corpo</p>"))
}

// A styled div is not a dialog. Without these three attributes a screen reader
// announces the panel as ordinary page content.
func TestModalHasDialogSemantics(t *testing.T) {
	html := modalHTML(t)
	for _, want := range []string{`role="dialog"`, `aria-modal="true"`, `aria-labelledby="confirm-title"`} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %s, got: %s", want, html)
		}
	}
	if !strings.Contains(html, `id="confirm-title"`) {
		t.Errorf("aria-labelledby points at nothing, got: %s", html)
	}
}

// The panel must be focusable so focus can land on it, and must intercept Tab
// so focus cannot leave.
func TestModalPanelTrapsFocus(t *testing.T) {
	html := modalHTML(t)
	if !strings.Contains(html, `x-ref="panel"`) {
		t.Errorf("modalOverlay() locates the panel by ref; without it there is no trap: %s", html)
	}
	if !strings.Contains(html, `tabindex="-1"`) {
		t.Errorf("the panel must be programmatically focusable, got: %s", html)
	}
	if !strings.Contains(html, "onTab") {
		t.Errorf("expected Tab to be intercepted, got: %s", html)
	}
}

// openOverlay/closeOverlay carry the scroll lock, the inert siblings and the
// focus restore. Going back to a bare `open = true` would silently drop all of
// it while still looking like it works.
func TestModalUsesTheOverlayPrimitive(t *testing.T) {
	html := modalHTML(t)
	if !strings.Contains(html, "openOverlay()") || !strings.Contains(html, "closeOverlay()") {
		t.Errorf("expected the modal to open and close through modalOverlay(), got: %s", html)
	}
	if strings.Contains(html, "open = true") {
		t.Errorf("setting open directly skips the scroll lock, the inert siblings and the focus restore: %s", html)
	}
}

func TestModalCloseButtonIsLabelled(t *testing.T) {
	html := modalHTML(t)
	if !strings.Contains(html, `aria-label="Close"`) {
		t.Errorf("a × with no label reads as \"times\", got: %s", html)
	}
}

func TestModalRendersTitleAndChildren(t *testing.T) {
	html := modalHTML(t)
	if !strings.Contains(html, "Confirmar") || !strings.Contains(html, "<p>corpo</p>") {
		t.Errorf("expected title and children, got: %s", html)
	}
}
