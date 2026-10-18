package showcase_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/vapulab/malphas/showcase"
)

func TestDemoProtocolsSortsFiltersAndPages(t *testing.T) {
	r := showcase.DemoProtocols("etapas", "desc", "", 1, 8)
	if r.Total != 40 || len(r.Rows) != 8 || r.TotalPages != 5 {
		t.Fatalf("unexpected shape: total=%d rows=%d pages=%d", r.Total, len(r.Rows), r.TotalPages)
	}
	if !sort.SliceIsSorted(r.Rows, func(i, j int) bool { return r.Rows[i].Etapas > r.Rows[j].Etapas }) {
		t.Errorf("rows not sorted by etapas desc: %+v", r.Rows)
	}

	// Search is case- and accent-insensitive, like the Combobox.
	r = showcase.DemoProtocols("nome", "asc", "PCR", 1, 8)
	if r.Total == 0 || r.Total == 40 {
		t.Fatalf("search for PCR should match some but not all, got %d", r.Total)
	}
	for _, row := range r.Rows {
		if !strings.Contains(strings.ToLower(row.Nome+" "+row.Area), "pcr") {
			t.Errorf("row %q does not match pcr", row.Nome)
		}
	}
	if showcase.DemoProtocols("nome", "asc", "chloe", 1, 8).Total == 0 {
		t.Error("\"chloe\" should find \"Chloé\" (accents ignored)")
	}
}

// The query string is user input on a public page: anything unexpected falls
// back to a valid state rather than an error.
func TestDemoProtocolsClampsBadInput(t *testing.T) {
	r := showcase.DemoProtocols("xx", "sideways", "", 999, 8)
	if r.Sort != "nome" || r.Dir != "asc" {
		t.Errorf("bad sort/dir should fall back to nome/asc, got %s/%s", r.Sort, r.Dir)
	}
	if r.Page != r.TotalPages {
		t.Errorf("page 999 should clamp to the last page %d, got %d", r.TotalPages, r.Page)
	}
	if showcase.DemoProtocols("nome", "asc", "", -3, 8).Page != 1 {
		t.Error("a negative page should clamp to 1")
	}
	empty := showcase.DemoProtocols("nome", "asc", "nothing-matches-this", 1, 8)
	if empty.Total != 0 || empty.Page != 1 || len(empty.Rows) != 0 {
		t.Errorf("no match: %+v", empty)
	}
}
