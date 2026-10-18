package showcase

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Examples are built once, at catalog time, so a demo has no request to read.
// That is fine for almost all of them — but not for Pagination, whose entire
// job is to reflect which page you are on. Clicking a page changed the URL and
// the demo went on showing its hardcoded value, which made a working component
// look broken.
//
// The handler puts the request's query into the context; demos that need it
// read it back at render time.

type queryKey struct{}

// WithQuery stores the request's query string for the demos to read.
func WithQuery(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, queryKey{}, r.URL.Query())
}

// QueryInt reads a positive integer from the request query, falling back to
// def when absent, unparseable or out of range.
func QueryInt(ctx context.Context, name string, def, min, max int) int {
	q, ok := ctx.Value(queryKey{}).(url.Values)
	if !ok {
		return def
	}
	n, err := strconv.Atoi(q.Get(name))
	if err != nil || n < min || n > max {
		return def
	}
	return n
}

func query(ctx context.Context) (url.Values, bool) {
	q, ok := ctx.Value(queryKey{}).(url.Values)
	return q, ok
}

// QueryString reads a query parameter, falling back to def when absent.
func QueryString(ctx context.Context, name, def string) string {
	q, ok := query(ctx)
	if !ok || !q.Has(name) || q.Get(name) == "" {
		return def
	}
	return q.Get(name)
}

// QueryDate reads a date as YYYY-MM-DD or a month as YYYY-MM (its first day),
// in UTC, falling back to def when absent or malformed.
func QueryDate(ctx context.Context, name string, def time.Time) time.Time {
	raw := QueryString(ctx, name, "")
	for _, layout := range []string{"2006-01-02", "2006-01"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return t
		}
	}
	return def
}

type questionnaireKey struct{}

// WithQuestionnaire carries a posted step into the page render, for the
// no-JavaScript path: the questionnaire handler answers a plain POST with the
// whole component page, and the demo shows this state instead of a fresh one.
func WithQuestionnaire(ctx context.Context, s QuestionnaireState) context.Context {
	return context.WithValue(ctx, questionnaireKey{}, s)
}

func questionnaireFrom(ctx context.Context) QuestionnaireState {
	if s, ok := ctx.Value(questionnaireKey{}).(QuestionnaireState); ok {
		return s
	}
	return QuestionnaireStart(QueryString(ctx, "retomar", "") == "1")
}
