package showcase_test

import (
	"html"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"slices"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas"
	"github.com/vapulab/malphas/layout"
	"github.com/vapulab/malphas/showcase"
)

// frame is the least chrome a host can give: layout.Base around the body.
func frame(title string, _ bool, body templ.Component) templ.Component {
	return layout.Base(layout.BaseProps{Title: title, Body: body})
}

func serve(t *testing.T, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	showcase.Handler(frame).ServeHTTP(rec, req)
	return rec
}

func get(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()
	return serve(t, httptest.NewRequest(http.MethodGet, path, nil))
}
func TestIndexListsEveryComponent(t *testing.T) {
	rec := get(t, "/components")
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	body := rec.Body.String()
	for _, c := range showcase.All() {
		if !strings.Contains(body, "/components/"+c.Slug) {
			t.Errorf("index is missing a link to %q", c.Slug)
		}
	}
}

// Walking the registry means a new component is covered by this test the
// moment it is registered — there is no per-component test to forget.
func TestEveryComponentPageRenders(t *testing.T) {
	for _, c := range showcase.All() {
		t.Run(c.Slug, func(t *testing.T) {
			rec := get(t, "/components/"+c.Slug)
			if rec.Code != http.StatusOK {
				t.Fatalf("want 200, got %d", rec.Code)
			}
			body := rec.Body.String()
			if !strings.Contains(body, c.Name) {
				t.Errorf("page does not mention the component name %q", c.Name)
			}
			for _, ex := range c.Examples {
				if !strings.Contains(body, ex.Title) {
					t.Errorf("page is missing example %q", ex.Title)
				}
			}
		})
	}
}

// Demos are built once at catalog time, so without the query reaching them
// through the context, clicking a page changed the URL and the marker stayed
// put — the component looked broken while being correct.
func TestPaginationDemoFollowsTheQueryString(t *testing.T) {
	body := get(t, "/components/pagination?page=7").Body.String()
	if !strings.Contains(body, `aria-current="page"`) {
		t.Fatalf("expected a current page marker, got: %s", body)
	}
	// The marked page is a span; page 7 must not also appear as a link.
	if strings.Contains(body, `href="/components/pagination?page=7"`) {
		t.Error("page 7 was requested but is still rendered as a link, so the demo ignored the query")
	}
	if !strings.Contains(body, `href="/components/pagination?page=6"`) {
		t.Error("expected neighbouring pages linked around 7")
	}
}

// Four paginators share one page. With one query parameter between them,
// clicking inside an example left that example frozen and moved a different
// one further up the page — the interaction looked broken at the exact spot
// the user clicked.
func TestEachPaginationExampleMovesOnItsOwn(t *testing.T) {
	// Unescaped, because templ writes & as &amp; inside an attribute — correct
	// HTML, and url.ParseQuery would otherwise read "amp;page" as a parameter.
	body := html.UnescapeString(get(t, "/components/pagination?page=7&e1=4").Body.String())

	if !strings.Contains(body, `aria-current="page"`) {
		t.Fatal("no current-page marker anywhere on the page")
	}

	// Which page each example's links move it to. Parsed rather than matched as
	// a string: the owner parameter is appended last, so the order varies.
	moves := map[string]map[string]bool{}
	for _, m := range regexp.MustCompile(`href="/components/pagination\?([^"]*)"`).FindAllStringSubmatch(body, -1) {
		q, err := url.ParseQuery(m[1])
		if err != nil {
			t.Fatalf("unparseable pagination link %q: %v", m[1], err)
		}
		for param, values := range q {
			if moves[param] == nil {
				moves[param] = map[string]bool{}
			}
			moves[param][values[0]] = true
		}
	}

	for param, current := range map[string]string{"page": "7", "e1": "4"} {
		// The current page is a span, so no link may point back at it.
		for _, other := range []string{"page", "e1"} {
			if other != param && !moves[param][current] {
				t.Errorf("no %s link carries %s=%s: navigating %s resets %s", other, param, current, other, param)
			}
		}
		// A link that moves the example must move it somewhere else.
		if len(moves[param]) == 1 && moves[param][current] {
			t.Errorf("every %s link points at %s=%s, so the example cannot move", param, param, current)
		}
	}
	for _, want := range []struct{ param, page string }{{"page", "6"}, {"e1", "3"}} {
		if !moves[want.param][want.page] {
			t.Errorf("no link moves %s to %s", want.param, want.page)
		}
	}
}

func TestUnknownComponentIs404(t *testing.T) {
	if rec := get(t, "/components/nao-existe"); rec.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", rec.Code)
	}
}

func postAnswer(t *testing.T, form url.Values, htmx bool) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/components/questionnaire/responder", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if htmx {
		req.Header.Set("HX-Request", "true")
	}
	return serve(t, req)
}

// With HTMX the response is just the next step, to swap into place.
func TestQuestionnaireAnswerReturnsTheNextStepFragment(t *testing.T) {
	rec := postAnswer(t, url.Values{"atual": {"experiencia"}, "nav": {"next"}, "a.experiencia": {"sim"}}, true)
	body := rec.Body.String()
	if rec.Code != http.StatusOK || !strings.Contains(body, "In which areas?") {
		t.Fatalf("status %d, body: %s", rec.Code, body)
	}
	if strings.Contains(body, "<html") {
		t.Error("an HTMX response should be a fragment, not a page")
	}
}

// Without JavaScript the form is a plain POST: the whole page comes back,
// showing the next step.
func TestQuestionnaireAnswerWorksWithoutJavaScript(t *testing.T) {
	rec := postAnswer(t, url.Values{"atual": {"experiencia"}, "nav": {"next"}, "a.experiencia": {"nao"}}, false)
	body := rec.Body.String()
	if rec.Code != http.StatusOK || !strings.Contains(body, "<html") || !strings.Contains(body, "What is the goal of your next experiment?") {
		t.Fatalf("expected the full page at the objetivo step; status %d", rec.Code)
	}
}

// Any page number past the end clamps to the last page, however large.
func TestDataTableDemoClampsHugePages(t *testing.T) {
	for _, page := range []string{"999", "1001", "999999"} {
		rec := get(t, "/components/data-table?page="+page)
		if !strings.Contains(rec.Body.String(), `aria-current="page" class="inline-flex min-w-9 justify-center rounded border border-primary bg-primary px-3 py-1.5 text-sm font-medium text-primary-fg">5<`) {
			t.Errorf("page=%s should clamp to page 5", page)
		}
	}
}

// Every component, live, on one page.
func TestVitrineRendersEveryComponent(t *testing.T) {
	rec := get(t, "/components/all")
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
	body := rec.Body.String()
	for _, c := range showcase.All() {
		if !strings.Contains(body, `id="c-`+c.Slug+`"`) {
			t.Errorf("component %q has no section on the showcase", c.Slug)
		}
	}
	// Live demos, not cards or code: no example code blocks on this page.
	if strings.Contains(body, "<pre") {
		t.Error("the showcase shows components, not their code")
	}
}

// The showcase page moved from /components/todos; old links still land on it,
// query string included.
func TestOldVitrineURLRedirects(t *testing.T) {
	for path, want := range map[string]string{
		"/components/todos":          "/components/all",
		"/components/todos?page=2&x": "/components/all?page=2&x",
	} {
		rec := get(t, path)
		if rec.Code != http.StatusMovedPermanently {
			t.Errorf("%s: want 301, got %d", path, rec.Code)
		}
		if loc := rec.Header().Get("Location"); loc != want {
			t.Errorf("%s: want Location %q, got %q", path, want, loc)
		}
	}
}

// Demos written for their own page now share one document: a repeated id
// would make a modal, a calendar or an htmx target find the wrong element.
func TestVitrineHasUniqueIDs(t *testing.T) {
	body := get(t, "/components/all").Body.String()
	// Only real attributes: x-bind:id and friends are not ids.
	re := regexp.MustCompile(`[\s"]id="([^"]+)"`)
	seen := map[string]int{}
	for _, m := range re.FindAllStringSubmatch(body, -1) {
		seen[m[1]]++
	}
	for id, n := range seen {
		if n > 1 {
			t.Errorf("id %q appears %d times", id, n)
		}
	}
}

// The Message Scroller demo pulls messages from this endpoint with
// hx-swap="beforeend" (new) and "afterbegin" (older).
func TestMessageDemoFragment(t *testing.T) {
	rec := get(t, "/components/message-scroller/mensagem")
	if rec.Code != 200 || strings.Count(rec.Body.String(), "<article") != 1 || strings.Contains(rec.Body.String(), "<html") {
		t.Fatalf("expected one message fragment, status %d: %s", rec.Code, rec.Body.String())
	}
	older := get(t, "/components/message-scroller/mensagem?antes=1")
	if strings.Count(older.Body.String(), "<article") != 3 {
		t.Errorf("antes=1 should return a page of three older messages: %s", older.Body.String())
	}
}

// The catalog is the technical material: English even when the host serves
// its own pages in another language and puts that language in the context.
func TestCatalogIsAlwaysEnglish(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/components/modal", nil)
	req = req.WithContext(malphas.WithLocale(req.Context(), malphas.RU))
	html := serve(t, req).Body.String()
	if !strings.Contains(html, `aria-label="Close"`) || strings.Contains(html, "Закрыть") {
		t.Errorf("the catalog should render its components in English")
	}
	if !strings.Contains(html, `<html lang="en">`) {
		t.Errorf("the catalog page is English: %s", html[:120])
	}
}

func TestCategoriesAreEnglish(t *testing.T) {
	var got []string
	for _, g := range showcase.ByCategory() {
		got = append(got, g.Category)
	}
	want := []string{"Actions", "Data display", "Data input", "Navigation", "Feedback", "Layout & utilities"}
	if !slices.Equal(got, want) {
		t.Errorf("categories = %v, want %v", got, want)
	}
}
