package showcase

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
)

// Frame wraps a catalog page in the host's chrome: layout.Base plus whatever
// navigation the host wants. wide asks for room for a side index, which only
// the vitrine uses.
type Frame func(title string, wide bool, body templ.Component) templ.Component

// Handler serves the catalog. The demos link to /components/..., and two of
// them post to or fetch from endpoints below it, so the handler expects to see
// the full path: mount it at /components and /components/ without stripping
// the prefix.
func Handler(frame Frame) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /components", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, frame("Componentes", false, IndexPage(ByCategory())))
	})
	mux.HandleFunc("GET /components/todos", func(w http.ResponseWriter, r *http.Request) {
		render(w, r, frame("Vitrine", true, AllPage(ByCategory())))
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
	_ = page.Render(WithQuery(r.Context(), r), w)
}

// questionnaireAnswer answers one step of the Questionnaire demo. With HTMX
// it returns just the next step, swapped in place; a plain POST (no
// JavaScript) gets the whole component page showing that step. Stateless:
// it stores nothing.
func questionnaireAnswer(frame Frame) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "formulário inválido", http.StatusBadRequest)
			return
		}
		state := QuestionnaireNext(r.PostForm)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if r.Header.Get("HX-Request") == "true" {
			_ = QuestionnaireView(state, true).Render(r.Context(), w)
			return
		}
		c, _ := BySlug("questionnaire")
		ctx := WithQuestionnaire(WithQuery(r.Context(), r), state)
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
		_ = ChatMessageView(DemoChatMessage(n+i)).Render(r.Context(), w)
	}
}
