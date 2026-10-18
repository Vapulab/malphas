package e2e_test

import (
	"strings"
	"testing"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
)

// Seventy demos share one document here. Any of them breaking another — an
// Alpine expression that throws, a script that assumes it owns the page —
// shows up as a console error.
func TestVitrineLoadsWithoutErrors(t *testing.T) {
	ctx := openPage(t, "/components/todos")
	var errs []string
	chromedp.ListenTarget(ctx, func(ev any) {
		switch e := ev.(type) {
		case *runtime.EventExceptionThrown:
			errs = append(errs, e.ExceptionDetails.Text)
		case *runtime.EventConsoleAPICalled:
			if e.Type == runtime.APITypeError || e.Type == runtime.APITypeWarning {
				errs = append(errs, string(e.Type))
			}
		}
	})
	navigate(t, ctx, "/components/todos")
	sleep(t, ctx, 800)
	if len(errs) != 0 {
		t.Errorf("the showcase logged %d error(s)/warning(s): %v", len(errs), errs)
	}
	if n := eval[int](t, ctx, `document.querySelectorAll('main section[id^="c-"]').length`); n < 70 {
		t.Errorf("expected every component on the page, found %d sections", n)
	}
}

// The server-driven demos find their regions by id; on the showcase they must
// still swap in place rather than navigate away.
func TestVitrineServerDrivenDemosWork(t *testing.T) {
	ctx := openPage(t, "/components/todos")

	click(t, ctx, `#dt-protocolos-sort-etapas`)
	sleep(t, ctx, 600)
	if loc := eval[string](t, ctx, `location.pathname`); loc != "/components/todos" && loc != "/components/data-table" {
		t.Fatalf("unexpected navigation to %q", loc)
	}
	if !eval[bool](t, ctx, `!!document.querySelector('#c-data-table th[aria-sort="ascending"] #dt-protocolos-sort-etapas')`) {
		t.Error("sorting the Data Table on the showcase did not swap its region in place")
	}

	click(t, ctx, `#cal-entrega-next`)
	sleep(t, ctx, 600)
	if !eval[bool](t, ctx, `!!document.querySelector('#c-calendar #cal-entrega-title')`) {
		t.Fatal("the Calendar left the showcase")
	}
}

func TestVitrineInteractiveDemosWork(t *testing.T) {
	ctx := openPage(t, "/components/todos")

	focus(t, ctx, "#cb-estado")
	eval[bool](t, ctx, `(document.querySelector('#cb-estado').select(), true)`)
	press(t, ctx, "zur")
	settle(t, ctx)
	press(t, ctx, kb.Enter)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.querySelector('input[name="city"]').value`); got != "zrh" {
		t.Errorf("the Combobox on the showcase picked %q, want zrh", got)
	}

	click(t, ctx, `[aria-controls="pop-filtros-panel"]`)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.activeElement.id`); got != "pop-autor" {
		t.Errorf("the Popover on the showcase did not take focus in: %q", got)
	}
	press(t, ctx, kb.Escape)
	settle(t, ctx)

	links := eval[string](t, ctx, `[...document.querySelectorAll('nav[aria-label="Showcase index"] a')].length + ''`)
	if !strings.HasPrefix(links, "7") {
		t.Errorf("the index should list every component, has %s links", links)
	}
}
