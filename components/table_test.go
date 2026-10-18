package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/vapulab/malphas/components"
)

func TestTableWrapsInAHorizontalScroller(t *testing.T) {
	html := renderWithChildren(t, components.Table(components.TableProps{}), templ.Raw(""))
	// A dense table must not force the whole page to scroll sideways.
	if !strings.Contains(html, "overflow-x-auto") {
		t.Errorf("expected the table to scroll inside its own container, got: %s", html)
	}
	if !strings.Contains(html, "<table") {
		t.Errorf("expected a table element, got: %s", html)
	}
}

// Zebra and density are set once on the table and reach the cells through
// descendant selectors, so TH and TD stay free of layout knobs.
func TestTableZebraTargetsBodyRowsOnly(t *testing.T) {
	html := renderWithChildren(t, components.Table(components.TableProps{Zebra: true}), templ.Raw(""))
	if !strings.Contains(html, "tbody") {
		t.Errorf("zebra striping must be scoped to tbody rows, got: %s", html)
	}
}

func TestTableCompactChangesCellPadding(t *testing.T) {
	normal := renderWithChildren(t, components.Table(components.TableProps{}), templ.Raw(""))
	compact := renderWithChildren(t, components.Table(components.TableProps{Compact: true}), templ.Raw(""))
	if normal == compact {
		t.Error("compact must change the rendered classes")
	}
}

func TestTableHeaderCellDefaultsToColumnScope(t *testing.T) {
	html := renderWithChildren(t, components.TH(components.THProps{}), templ.Raw("Nome"))
	if !strings.Contains(html, `scope="col"`) {
		t.Errorf("a header cell needs a scope to be useful, got: %s", html)
	}
	if !strings.Contains(html, "Nome") {
		t.Errorf("expected children rendered, got: %s", html)
	}
}

func TestTableRowCellsRender(t *testing.T) {
	row := renderWithChildren(t, components.TR(), templ.Raw("<td>x</td>"))
	if !strings.Contains(row, "<tr") {
		t.Errorf("expected a tr, got: %s", row)
	}
	cell := renderWithChildren(t, components.TD(components.TDProps{}), templ.Raw("42"))
	if !strings.Contains(cell, "<td") || !strings.Contains(cell, "42") {
		t.Errorf("expected a td with content, got: %s", cell)
	}
}
