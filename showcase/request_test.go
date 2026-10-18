package showcase_test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vapulab/malphas/showcase"
)

func ctxWith(t *testing.T, rawQuery string) context.Context {
	t.Helper()
	r := httptest.NewRequest("GET", "/components/pagination?"+rawQuery, nil)
	return showcase.WithQuery(context.Background(), r)
}

func TestQueryIntReadsTheRequest(t *testing.T) {
	if got := showcase.QueryInt(ctxWith(t, "page=7"), "page", 1, 1, 20); got != 7 {
		t.Errorf("QueryInt = %d, want 7", got)
	}
}

// Every rejection path must fall back rather than surface a broken value: the
// query string is user input on a public page.
func TestQueryIntFallsBackOnBadInput(t *testing.T) {
	cases := map[string]string{
		"absent":        "",
		"non-numeric":   "page=abc",
		"zero":          "page=0",
		"negative":      "page=-3",
		"above the max": "page=99",
		"empty":         "page=",
	}
	for name, query := range cases {
		if got := showcase.QueryInt(ctxWith(t, query), "page", 3, 1, 20); got != 3 {
			t.Errorf("%s: QueryInt = %d, want the fallback 3", name, got)
		}
	}
}

// A context with no query at all — as when a demo is rendered from a test —
// must not panic.
func TestQueryIntWithoutARequest(t *testing.T) {
	if got := showcase.QueryInt(context.Background(), "page", 5, 1, 20); got != 5 {
		t.Errorf("QueryInt = %d, want the fallback 5", got)
	}
}

func TestQueryStringReadsAndFallsBack(t *testing.T) {
	if got := showcase.QueryString(ctxWith(t, "q=quantitative%20pcr"), "q", "x"); got != "quantitative pcr" {
		t.Errorf("QueryString = %q", got)
	}
	if got := showcase.QueryString(ctxWith(t, ""), "q", "default"); got != "default" {
		t.Errorf("absent should fall back, got %q", got)
	}
	if got := showcase.QueryString(context.Background(), "q", "default"); got != "default" {
		t.Errorf("no request in context should fall back, got %q", got)
	}
}

func TestQueryDateAcceptsDayAndMonth(t *testing.T) {
	def := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cases := map[string]time.Time{
		"month=2026-10":    time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
		"month=2026-10-20": time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC),
		"month=":           def,
		"month=yesterday":  def,
		"month=2026-13":    def,
		"":                 def,
	}
	for query, want := range cases {
		if got := showcase.QueryDate(ctxWith(t, query), "month", def); !got.Equal(want) {
			t.Errorf("%q: QueryDate = %v, want %v", query, got, want)
		}
	}
}
