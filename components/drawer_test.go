package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
)

func drawerHTML(t *testing.T, side string) string {
	t.Helper()
	return renderWithChildren(t,
		components.Drawer(components.DrawerProps{ID: "nav", Title: "Navegação", Side: side}),
		templ.Raw("<p>links</p>"))
}

func TestDrawerHasDialogSemantics(t *testing.T) {
	html := drawerHTML(t, "")
	for _, want := range []string{`role="dialog"`, `aria-modal="true"`, `aria-labelledby="nav-title"`, `id="nav-title"`} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %s, got: %s", want, html)
		}
	}
}

// A side panel that lets Tab wander behind it is worse than a dialog doing the
// same, because the page still looks usable.
func TestDrawerIsModal(t *testing.T) {
	html := drawerHTML(t, "")
	for _, want := range []string{`x-ref="panel"`, "onTab", "openOverlay()", "closeOverlay()"} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %s — the drawer must go through modalOverlay(), got: %s", want, html)
		}
	}
}

func TestDrawerSide(t *testing.T) {
	left := drawerHTML(t, "left")
	right := drawerHTML(t, "right")
	if !strings.Contains(left, "mr-auto") {
		t.Errorf("left drawer should hug the left edge, got: %s", left)
	}
	if !strings.Contains(right, "ml-auto") {
		t.Errorf("right drawer should hug the right edge, got: %s", right)
	}
	if left == right {
		t.Error("Side must change the rendered markup")
	}
}

func TestDrawerDefaultsToLeft(t *testing.T) {
	if drawerHTML(t, "") != drawerHTML(t, "left") {
		t.Error("an empty Side should render the same as left")
	}
}

// Side bottom is a sheet rising from the bottom edge: full width, capped
// height, no side-panel width.
func TestDrawerBottom(t *testing.T) {
	html := renderWithChildren(t,
		components.Drawer(components.DrawerProps{ID: "d", Title: "Filtros", Side: "bottom"}), templ.Raw("x"))
	for _, want := range []string{"mt-auto", "w-full", "max-h-[85vh]", "border-t", "rounded-t-lg"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	if strings.Contains(html, "w-80") || strings.Contains(html, "h-full") {
		t.Errorf("a bottom sheet must not keep the side panel's size, got: %s", html)
	}
}
