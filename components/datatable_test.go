package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/darkanum/vapulab/malphas/components"
)

func dataTableFixture() components.DataTableProps {
	return components.DataTableProps{
		ID: "dt", URL: "/protocolos",
		Columns: []components.DataTableColumn{
			{Key: "nome", Label: "Nome", Sortable: true},
			{Key: "area", Label: "Área", Sortable: true},
			{Key: "notas", Label: "Notas"},
		},
		Sort: "nome", Dir: "asc", Query: "pcr", Page: 2, TotalPages: 3, Total: 25,
		Searchable: true, SearchLabel: "Buscar protocolos",
	}
}

func TestDataTableAriaSortAndLinks(t *testing.T) {
	html := renderWithChildren(t, components.DataTable(dataTableFixture()), templ.Raw("<tr><td>x</td></tr>"))
	for _, want := range []string{
		`aria-sort="ascending"`,
		// The active column flips direction; the search survives; page resets.
		`href="/protocolos?dir=desc&amp;q=pcr&amp;sort=nome"`,
		`href="/protocolos?dir=asc&amp;q=pcr&amp;sort=area"`,
		`hx-boost="true"`,
		`hx-select="#dt-region"`,
		`hx-target="#dt-region"`,
		`id="dt-region"`,
		`type="search"`,
		`aria-label="Buscar protocolos"`,
		`value="pcr"`,
		`hx-select-oob="#dt-state"`,
		`25 resultados`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	if n := strings.Count(html, `aria-sort="none"`); n != 1 {
		t.Errorf("the other sortable column is aria-sort=none; got %d", n)
	}
	// A column that cannot be sorted has no link and no aria-sort.
	notas := html[strings.Index(html, ">Notas<")-200 : strings.Index(html, ">Notas<")]
	if strings.Contains(notas, "<a") || strings.Contains(notas, "aria-sort") {
		t.Errorf("Notas is not sortable: %s", notas)
	}
	// Pagination keeps sort, dir and q, with the page last.
	if !strings.Contains(html, `href="/protocolos?dir=asc&amp;q=pcr&amp;sort=nome&amp;page=3"`) {
		t.Errorf("pagination should carry the state: %s", html)
	}
}

func TestDataTableHrefEncodes(t *testing.T) {
	p := dataTableFixture()
	got := p.Href("nome", "asc", "ação rápida", 1)
	if got != "/protocolos?dir=asc&q=a%C3%A7%C3%A3o+r%C3%A1pida&sort=nome" {
		t.Errorf("Href = %q", got)
	}
}

func TestDataTableEmpty(t *testing.T) {
	p := dataTableFixture()
	p.Total, p.TotalPages, p.Empty = 0, 0, "Nenhum protocolo encontrado."
	html := renderWithChildren(t, components.DataTable(p), templ.Raw(""))
	if !strings.Contains(html, "Nenhum protocolo encontrado.") || !strings.Contains(html, `colspan="3"`) {
		t.Errorf("expected the empty message across all columns: %s", html)
	}
	if strings.Contains(html, `aria-label="Paginação"`) {
		t.Errorf("no pages, no pagination: %s", html)
	}
}
