package showcase

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/vapulab/malphas/components"
)

func intToString(n int) string { return strconv.Itoa(n) }

// The Pagination page hosts four independent paginators, so each owns a query
// parameter. Sharing one would make clicking inside an example move a
// different example — which is exactly how the page read before, since the
// three edge examples had a fixed page and only the interactive one moved.
type paginationDemo struct {
	param string
	def   int // page shown before anything is clicked
	total int
}

var paginationDemos = []paginationDemo{
	{param: "page", def: 3, total: 20},
	{param: "e1", def: 1, total: 5},
	{param: "e2", def: 1, total: 20},
	{param: "e3", def: 20, total: 20},
}

func paginationTotal(param string) int {
	for _, d := range paginationDemos {
		if d.param == param {
			return d.total
		}
	}
	return 0
}

func paginationPage(ctx context.Context, param string) int {
	for _, d := range paginationDemos {
		if d.param == param {
			return QueryInt(ctx, d.param, d.def, 1, d.total)
		}
	}
	return 1
}

// paginationHref carries the other examples' pages along, so navigating one
// does not silently reset the rest. Only non-default values ride along, which
// keeps the URL readable.
func paginationHref(ctx context.Context, param string) string {
	q := url.Values{}
	for _, d := range paginationDemos {
		if d.param == param {
			continue
		}
		if page := paginationPage(ctx, d.param); page != d.def {
			q.Set(d.param, strconv.Itoa(page))
		}
	}
	base := "/components/pagination?"
	if len(q) > 0 {
		base += q.Encode() + "&"
	}
	return base + param + "="
}

// carouselSlideClass paints the demo slides. Lives in Go rather than inline in
// the .templ so the class strings stay in one place for Tailwind to scan.
func carouselSlideClass(variant string) string {
	return map[string]string{
		"primary": "bg-primary/10 text-primary",
		"success": "bg-success/10 text-success",
		"warning": "bg-warning/10 text-warning",
	}[variant]
}

// demoEstados is long enough that scrolling it would be a chore — the case
// Combobox exists for — and full of accents to show the filter ignoring them.
func demoEstados() []components.SelectOption {
	names := [][2]string{
		{"AC", "Acre"}, {"AL", "Alagoas"}, {"AP", "Amapá"}, {"AM", "Amazonas"},
		{"BA", "Bahia"}, {"CE", "Ceará"}, {"DF", "Distrito Federal"},
		{"ES", "Espírito Santo"}, {"GO", "Goiás"}, {"MA", "Maranhão"},
		{"MT", "Mato Grosso"}, {"MS", "Mato Grosso do Sul"}, {"MG", "Minas Gerais"},
		{"PA", "Pará"}, {"PB", "Paraíba"}, {"PR", "Paraná"}, {"PE", "Pernambuco"},
		{"PI", "Piauí"}, {"RJ", "Rio de Janeiro"}, {"RN", "Rio Grande do Norte"},
		{"RS", "Rio Grande do Sul"}, {"RO", "Rondônia"}, {"RR", "Roraima"},
		{"SC", "Santa Catarina"}, {"SP", "São Paulo"}, {"SE", "Sergipe"},
		{"TO", "Tocantins"},
	}
	opts := make([]components.SelectOption, len(names))
	for i, n := range names {
		opts[i] = components.SelectOption{Label: n[1], Value: strings.ToLower(n[0])}
	}
	return opts
}
