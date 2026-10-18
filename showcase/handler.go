package showcase

import (
	"context"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas"
)

// Frame wraps a catalog page in the host's chrome: layout.Base plus whatever
// navigation the host wants. wide asks for room for a side index, which only
// the showcase uses.
type Frame func(title string, wide bool, body templ.Component) templ.Component

// Handler serves the catalog, always in English: it is the technical
// material, read by people of every language, even inside a host that serves
// its own pages in another one. The demos link to /components/..., and two of
// them post to or fetch from endpoints below it, so the handler expects to see
// the full path: mount it at /components and /components/ without stripping
// the prefix.
func Handler(frame Frame) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /components", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, frame("Components", false, IndexPage(ByCategory())))
	})
	mux.HandleFunc("GET /components/todos", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, frame("Showcase", true, AllPage(ByCategory())))
	})
	mux.HandleFunc("GET /components/{slug}", func(w http.ResponseWriter, r *http.Request) {
		c, ok := BySlug(r.PathValue("slug"))
		if !ok {
			http.NotFound(w, r)
			return
		}
		render(w, r, frame(c.Name, false, ComponentPage(c)))
	})
	mux.HandleFunc("POST /components/questionnaire/responder", questionnaireAnswer(frame))
	mux.HandleFunc("GET /components/message-scroller/mensagem", messageDemo)
	return mux
}

// render passes the request query through the context. Demos are built once,
// at registry time, and have no request of their own; the few that follow
// the URL (Pagination, Data Table, Calendar) read it from there.
func render(w http.ResponseWriter, r *http.Request, page templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = page.Render(WithQuery(english(r.Context()), r), w)
}

func english(ctx context.Context) context.Context { return malphas.WithLocale(ctx, malphas.EN) }

// questionnaireAnswer answers one step of the Questionnaire demo. With HTMX
// it returns just the next step, swapped in place; a plain POST (no
// JavaScript) gets the whole component page showing that step. Stateless:
// it stores nothing.
func questionnaireAnswer(frame Frame) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		state := QuestionnaireNext(r.PostForm)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.Header.Get("HX-Request") == "true" {
			_ = QuestionnaireView(state, true).Render(english(r.Context()), w)
			return
		}
		c, _ := BySlug("questionnaire")
		ctx := WithQuestionnaire(WithQuery(english(r.Context()), r), state)
		_ = frame(c.Name, false, ComponentPage(c)).Render(ctx, w)
	}
}

// messageDemo feeds the Message Scroller demo: one new message, or with
// antes=1 a page of three older ones for the top. Stateless; the message is
// picked from the clock, which is all a demo needs.
func messageDemo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	n := int(time.Now().UnixNano() / int64(time.Millisecond))
	count := 1
	if r.URL.Query().Get("antes") == "1" {
		count = 3
	}
	for i := 0; i < count; i++ {
		_ = ChatMessageView(DemoChatMessage(n+i)).Render(english(r.Context()), w)
	}
}
