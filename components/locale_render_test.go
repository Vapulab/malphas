package components_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas"
	"github.com/vapulab/malphas/components"
)

func renderIn(t *testing.T, l malphas.Locale, c templ.Component) string {
	t.Helper()
	return renderInContext(t, malphas.WithLocale(context.Background(), l), c)
}

// The components write their own words in the language of the request: a
// host sets it once and every component follows.
func TestComponentsSpeakTheContextLocale(t *testing.T) {
	modal := renderInContext(t, templ.WithChildren(malphas.WithLocale(context.Background(), malphas.RU), templ.Raw("<p>x</p>")),
		components.Modal(components.ModalProps{ID: "m", Title: "T"}))
	if !strings.Contains(modal, `aria-label="Закрыть"`) {
		t.Errorf("the Russian modal should close with Закрыть: %s", modal)
	}

	pages := renderIn(t, malphas.RU, components.Pagination(components.PaginationProps{Current: 2, Total: 5, BaseHref: "/p?page="}))
	for _, want := range []string{`aria-label="Постраничная навигация"`, "Предыдущая", "Следующая", `aria-label="Страница 3"`} {
		if !strings.Contains(pages, want) {
			t.Errorf("Russian pagination is missing %s", want)
		}
	}

	table := renderIn(t, malphas.RU, components.DataTable(components.DataTableProps{
		ID: "t", Columns: []components.DataTableColumn{{Key: "name", Label: "Name"}}, URL: "/t", Page: 1, TotalPages: 1,
	}))
	if !strings.Contains(table, "Ничего не найдено") || !strings.Contains(table, "0 результатов") {
		t.Errorf("an empty Russian table should say so in Russian: %s", table)
	}
}

// A label passed as a prop is the host's word and beats the language.
func TestExplicitLabelBeatsLocale(t *testing.T) {
	html := renderInContext(t, templ.WithChildren(malphas.WithLocale(context.Background(), malphas.RU), templ.Raw("x")),
		components.AlertDialog(components.AlertDialogProps{ID: "d", Title: "T", CancelLabel: "Voltar", ConfirmLabel: "Apagar"}))
	if !strings.Contains(html, ">Voltar<") || strings.Contains(html, ">Отмена<") {
		t.Errorf("CancelLabel should win over the Russian default: %s", html)
	}
}

// The jump button's counter is built in the browser; it gets every plural
// form of the language from the server.
func TestScrollerCarriesTheUnreadForms(t *testing.T) {
	html := renderIn(t, malphas.RU, components.MessageScroller(components.MessageScrollerProps{ID: "s"}))
	for _, want := range []string{"В конец", "новое", "новых"} {
		if !strings.Contains(html, want) {
			t.Errorf("the Russian scroller is missing %q", want)
		}
	}
	if strings.Contains(html, "nova") {
		t.Error("the Portuguese counter is gone")
	}
}
