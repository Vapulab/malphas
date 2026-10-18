package showcase

import (
	"sort"
	"strings"
)

// DemoProtocol is one row of the Data Table demo. The data is fixed and in
// memory: the demo exists to exercise the component's server round trips,
// not a database.
type DemoProtocol struct {
	Nome        string
	Area        string
	Responsavel string
	Etapas      int
}

// DemoQuery is the outcome of a Data Table request, with every input already
// normalised: these are the values the page must render, not the raw query.
type DemoQuery struct {
	Rows       []DemoProtocol
	Total      int
	Sort, Dir  string
	Page       int
	TotalPages int
}

var demoProtocols = func() []DemoProtocol {
	tecnicas := []string{"PCR quantitativo", "Cultura de célula", "Western blot", "ELISA", "Sequenciamento", "Microscopia confocal", "Citometria de fluxo", "Cromatografia"}
	areas := []string{"Genética", "Imunologia", "Bioquímica", "Microbiologia", "Neurociência"}
	pessoas := []string{"Ana Lima", "Bruno Reis", "Carla Souza", "Diego Alves", "Elisa Prado"}
	var out []DemoProtocol
	for i := 0; i < 40; i++ {
		out = append(out, DemoProtocol{
			Nome:        tecnicas[i%len(tecnicas)] + " " + string(rune('A'+i/len(tecnicas))),
			Area:        areas[(i*3)%len(areas)],
			Responsavel: pessoas[(i*7)%len(pessoas)],
			Etapas:      3 + (i*11)%14,
		})
	}
	return out
}()

var accentFolder = strings.NewReplacer(
	"á", "a", "à", "a", "â", "a", "ã", "a", "é", "e", "ê", "e", "í", "i",
	"ó", "o", "ô", "o", "õ", "o", "ú", "u", "ç", "c",
)

// fold matches the Combobox's search: case- and accent-insensitive.
func fold(s string) string { return accentFolder.Replace(strings.ToLower(s)) }

var demoSortKeys = map[string]func(a, b DemoProtocol) bool{
	"nome":        func(a, b DemoProtocol) bool { return a.Nome < b.Nome },
	"area":        func(a, b DemoProtocol) bool { return a.Area < b.Area },
	"responsavel": func(a, b DemoProtocol) bool { return a.Responsavel < b.Responsavel },
	"etapas":      func(a, b DemoProtocol) bool { return a.Etapas < b.Etapas },
}

// DemoProtocols filters, sorts and pages the demo data. Every input comes
// from the query string of a public page, so anything unexpected falls back
// to a valid state instead of an error: an unknown sort to "nome", an unknown
// direction to "asc", a page out of range to the nearest real page.
func DemoProtocols(sortKey, dir, q string, page, perPage int) DemoQuery {
	less, ok := demoSortKeys[sortKey]
	if !ok {
		sortKey, less = "nome", demoSortKeys["nome"]
	}
	if dir != "desc" {
		dir = "asc"
	}

	needle := fold(strings.TrimSpace(q))
	var rows []DemoProtocol
	for _, p := range demoProtocols {
		if needle == "" || strings.Contains(fold(p.Nome+" "+p.Area+" "+p.Responsavel), needle) {
			rows = append(rows, p)
		}
	}
	sort.SliceStable(rows, func(i, j int) bool {
		if dir == "desc" {
			return less(rows[j], rows[i])
		}
		return less(rows[i], rows[j])
	})

	total := len(rows)
	pages := (total + perPage - 1) / perPage
	if page > pages {
		page = pages
	}
	if page < 1 {
		page = 1
	}
	from := (page - 1) * perPage
	to := min(from+perPage, total)
	if from > total {
		from = total
	}
	return DemoQuery{Rows: rows[from:to], Total: total, Sort: sortKey, Dir: dir, Page: page, TotalPages: pages}
}
