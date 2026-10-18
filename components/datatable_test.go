package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
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

// Without JavaScript, Enter in the search field must still search: the field
// sits in a GET form to URL, carrying sort and dir along. With JavaScript the
// form never submits (the field's own hx-get does the work).
func TestDataTableSearchWorksAsAPlainForm(t *testing.T) {
	html := renderWithChildren(t, components.DataTable(dataTableFixture()), templ.Raw("<tr><td>x</td></tr>"))
	form := strings.Index(html, `<form role="search" method="get" action="/protocolos"`)
	field := strings.Index(html, `type="search"`)
	end := strings.Index(html, `</form>`)
	if form < 0 || !(form < field && field < end) {
		t.Fatalf("the search field must sit in a GET form to the table URL: %s", html)
	}
	state := strings.Index(html, `name="sort"`)
	if !(form < state && state < end) {
		t.Errorf("sort and dir must travel with the form: %s", html)
	}
	if !strings.Contains(html, `x-on:submit.prevent`) {
		t.Errorf("with JavaScript the form must not also navigate: %s", html)
	}
}

// Caller URLs go through templ's sanitizer like every link in Malphas
// (see link.templ): a javascript: URL must never reach an action.
func TestDataTableSanitizesItsURL(t *testing.T) {
	p := dataTableFixture()
	p.URL = "javascript:alert(1)"
	html := renderWithChildren(t, components.DataTable(p), templ.Raw(""))
	if strings.Contains(html, `action="javascript:`) || strings.Contains(html, `href="javascript:`) {
		t.Errorf("a javascript: URL reached the markup: %s", html)
	}
}

// Each debounced keystroke replaces the history entry instead of pushing a
// new one, so Back does not walk through "p", "pc", "pcr".
func TestDataTableSearchReplacesHistory(t *testing.T) {
	html := renderWithChildren(t, components.DataTable(dataTableFixture()), templ.Raw(""))
	field := html[strings.Index(html, `type="search"`):]
	field = field[:strings.Index(field, ">")]
	if !strings.Contains(field, `hx-replace-url="true"`) || strings.Contains(field, "hx-push-url") {
		t.Errorf("search should replace, not push, the URL: %s", field)
	}
}

// Without a search there is no state block, and the region must not point
// hx-select-oob at an element that does not exist.
func TestDataTableWithoutSearchHasNoStateOOB(t *testing.T) {
	p := dataTableFixture()
	p.Searchable = false
	html := renderWithChildren(t, components.DataTable(p), templ.Raw(""))
	if strings.Contains(html, "hx-select-oob") {
		t.Errorf("no search, no state to refresh: %s", html)
	}
}
