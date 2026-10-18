package e2e_test

import (
	"strings"
	"testing"

	"github.com/chromedp/chromedp/kb"
)

// --- Calendar ---------------------------------------------------------------

const calEntrega = `#cal-entrega`

func calDay(cal, iso string) string { return cal + ` button[data-date="` + iso + `"]` }

const focusedDate = `document.activeElement.dataset.date || ''`

func TestCalendarArrowsMoveFocusWithinTheMonth(t *testing.T) {
	ctx := openPage(t, "/components/calendar?mes=2026-10")

	focus(t, ctx, calDay(calEntrega, "2026-10-15"))
	for _, step := range []struct{ key, want string }{
		{kb.ArrowRight, "2026-10-16"},
		{kb.ArrowDown, "2026-10-23"},
		{kb.ArrowLeft, "2026-10-22"},
		{kb.ArrowUp, "2026-10-15"},
		{kb.Home, "2026-10-11"}, // Sunday of that week
		{kb.End, "2026-10-17"},  // Saturday
	} {
		press(t, ctx, step.key)
		settle(t, ctx)
		if got := eval[string](t, ctx, focusedDate); got != step.want {
			t.Fatalf("after %q focus is on %q, want %q", step.key, got, step.want)
		}
	}
	if n := eval[int](t, ctx, `document.querySelectorAll('#cal-entrega button[data-date][tabindex="0"]').length`); n != 1 {
		t.Errorf("roving tabindex: %d tab stops", n)
	}
}

// Walking off the month fetches the neighbouring one from the server, and the
// new grid puts focus on the day the key was heading for.
func TestCalendarArrowAcrossTheMonthLoadsItAndKeepsFocus(t *testing.T) {
	ctx := openPage(t, "/components/calendar?mes=2026-10")

	focus(t, ctx, calDay(calEntrega, "2026-10-01"))
	press(t, ctx, kb.ArrowLeft)
	sleep(t, ctx, 600)
	if got := eval[string](t, ctx, `document.querySelector('#cal-entrega-title').textContent.trim()`); got != "setembro de 2026" {
		t.Fatalf("the previous month did not load: title %q", got)
	}
	if got := eval[string](t, ctx, focusedDate); got != "2026-09-30" {
		t.Errorf("focus should land on 30 Sep, is on %q", got)
	}
}

func TestCalendarPageDownGoesToNextMonth(t *testing.T) {
	ctx := openPage(t, "/components/calendar?mes=2026-10")

	focus(t, ctx, calDay(calEntrega, "2026-10-31"))
	press(t, ctx, kb.PageDown)
	sleep(t, ctx, 600)
	// 31 Oct + 1 month clamps to 30 Nov.
	if got := eval[string](t, ctx, focusedDate); got != "2026-11-30" {
		t.Errorf("PageDown from 31 Oct should land on 30 Nov, is on %q", got)
	}
}

func TestCalendarSelectsAndSubmits(t *testing.T) {
	ctx := openPage(t, "/components/calendar?mes=2026-10")

	focus(t, ctx, calDay(calEntrega, "2026-10-20"))
	press(t, ctx, kb.Enter)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.querySelector('#cal-entrega input[name="entrega"]').value`); got != "2026-10-20" {
		t.Errorf("hidden input is %q, want 2026-10-20", got)
	}
	if got := eval[string](t, ctx, `document.querySelector('`+calDay(calEntrega, "2026-10-20")+`').parentElement.getAttribute('aria-selected')`); got != "true" {
		t.Errorf("the picked day's cell has aria-selected=%s", got)
	}

	// A day before Min is focusable but cannot be picked.
	focus(t, ctx, calDay(calEntrega, "2026-10-04"))
	press(t, ctx, kb.Enter)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.querySelector('#cal-entrega input[name="entrega"]').value`); got != "2026-10-20" {
		t.Errorf("a disabled day was picked: %q", got)
	}

	// The selection survives a month round trip through the server.
	click(t, ctx, "#cal-entrega-next")
	sleep(t, ctx, 500)
	click(t, ctx, "#cal-entrega-prev")
	sleep(t, ctx, 500)
	if got := eval[string](t, ctx, `document.querySelector('#cal-entrega input[name="entrega"]').value`); got != "2026-10-20" {
		t.Errorf("selection lost across months: %q", got)
	}
}

// The second pick before the first still makes a forward period.
func TestCalendarRangeOrdersItsEnds(t *testing.T) {
	ctx := openPage(t, "/components/calendar?mes=2026-10")

	click(t, ctx, calDay(`#cal-periodo`, "2026-10-20"))
	click(t, ctx, calDay(`#cal-periodo`, "2026-10-10"))
	settle(t, ctx)
	start := eval[string](t, ctx, `document.querySelector('#cal-periodo input[name="inicio"]').value`)
	end := eval[string](t, ctx, `document.querySelector('#cal-periodo input[name="fim"]').value`)
	if start != "2026-10-10" || end != "2026-10-20" {
		t.Errorf("period is %s → %s, want 2026-10-10 → 2026-10-20", start, end)
	}
	selected := eval[[]string](t, ctx, `[...document.querySelectorAll('#cal-periodo td[aria-selected="true"] button')].map(b => b.dataset.date)`)
	if len(selected) != 11 || !strings.HasSuffix(selected[0], "10-10") {
		t.Errorf("expected 11 selected days from the 10th, got %v", selected)
	}
}

// --- Data Table ---------------------------------------------------------------

const dtFirstRow = `document.querySelector('#dt-protocolos-region tbody tr td').textContent.trim()`
const dtRows = `[...document.querySelectorAll('#dt-protocolos-region tbody tr')].map(tr => tr.firstElementChild.textContent.trim())`

func TestDataTableSortsByHeader(t *testing.T) {
	ctx := openPage(t, "/components/data-table")

	click(t, ctx, `#dt-protocolos-sort-etapas`)
	sleep(t, ctx, 500)
	if got := eval[string](t, ctx, `document.querySelector('th[aria-sort="ascending"]').textContent.trim()`); !strings.HasPrefix(got, "Etapas") {
		t.Fatalf("Etapas should be ascending, aria-sort is on %q", got)
	}
	asc := eval[string](t, ctx, `[...document.querySelectorAll('#dt-protocolos-region tbody tr td:last-child')].map(td => td.textContent.trim()).join(',')`)

	click(t, ctx, `#dt-protocolos-sort-etapas`)
	sleep(t, ctx, 500)
	if !eval[bool](t, ctx, `!!document.querySelector('th[aria-sort="descending"]')`) {
		t.Fatal("a second click should sort descending")
	}
	desc := eval[string](t, ctx, `[...document.querySelectorAll('#dt-protocolos-region tbody tr td:last-child')].map(td => td.textContent.trim()).join(',')`)
	if asc == desc {
		t.Errorf("descending shows the same rows as ascending: %v", asc)
	}
	if loc := eval[string](t, ctx, `location.search`); !strings.Contains(loc, "sort=etapas") || !strings.Contains(loc, "dir=desc") {
		t.Errorf("the URL should keep the state, is %q", loc)
	}
	// The header link keeps focus across the swap (same id in the new region).
	if got := eval[string](t, ctx, `document.activeElement.id`); got != "dt-protocolos-sort-etapas" {
		t.Errorf("focus should stay on the header link, is on %q", got)
	}
}

// Typing while responses arrive: the field is outside the swapped region, so
// neither its focus nor its text is ever replaced.
func TestDataTableSearchFiltersWithoutLosingFocus(t *testing.T) {
	ctx := openPage(t, "/components/data-table")

	focus(t, ctx, `input[type="search"][aria-controls="dt-protocolos-region"]`)
	press(t, ctx, "pc")
	sleep(t, ctx, 450) // a response lands mid-typing
	press(t, ctx, "r")
	sleep(t, ctx, 700)

	if got := eval[string](t, ctx, `document.activeElement.value`); got != "pcr" {
		t.Errorf("the field should still hold \"pcr\" and have focus, active value is %q", got)
	}
	rows := eval[[]string](t, ctx, dtRows)
	if len(rows) == 0 {
		t.Fatal("no rows for pcr")
	}
	for _, r := range rows {
		if !strings.Contains(strings.ToLower(r), "pcr") {
			t.Errorf("row %q does not match the search", r)
		}
	}
	if loc := eval[string](t, ctx, `location.search`); !strings.Contains(loc, "q=pcr") {
		t.Errorf("the URL should carry the search, is %q", loc)
	}
}

func TestDataTablePaginatesAndKeepsKeyboardPlace(t *testing.T) {
	ctx := openPage(t, "/components/data-table")
	first := eval[string](t, ctx, dtFirstRow)

	click(t, ctx, `#dt-protocolos-region nav a[aria-label="Página 2"]`)
	sleep(t, ctx, 500)
	if got := eval[string](t, ctx, `document.querySelector('#dt-protocolos-region [aria-current="page"]').textContent.trim()`); got != "2" {
		t.Errorf("aria-current should be on page 2, is on %q", got)
	}
	if got := eval[string](t, ctx, dtFirstRow); got == first {
		t.Errorf("page 2 shows the same first row %q", got)
	}
	// The clicked link is gone with the old region; focus lands on the
	// results line, not on <body>.
	if got := eval[string](t, ctx, `document.activeElement.id`); got != "dt-protocolos-status" {
		t.Errorf("focus fell to %q; it should land on the results line", got)
	}
}

func TestDataTableSurvivesBadQuery(t *testing.T) {
	ctx := openPage(t, "/components/data-table?sort=xx&dir=up&page=999")

	if !eval[bool](t, ctx, `!!document.querySelector('#dt-protocolos-sort-nome').closest('th[aria-sort="ascending"]')`) {
		t.Error("an unknown sort should fall back to Protocolo ascending")
	}
	if got := eval[string](t, ctx, `document.querySelector('#dt-protocolos-region [aria-current="page"]').textContent.trim()`); got != "5" {
		t.Errorf("page 999 should clamp to the last page, aria-current is on %q", got)
	}
}

// --- Questionnaire ------------------------------------------------------------

const qnTitle = `(document.querySelector('#qn-demo legend, #qn-demo label[for]')?.textContent || '').trim()`

func TestQuestionnaireAdvancesAndBranches(t *testing.T) {
	ctx := openPage(t, "/components/questionnaire")

	focus(t, ctx, `#qn-demo input[data-option-number="1"]`)
	press(t, ctx, "1") // Sim
	pressEnterNative(t, ctx)
	sleep(t, ctx, 500)
	if got := eval[string](t, ctx, qnTitle); got != "Em quais áreas?" {
		t.Fatalf("\"Sim\" should lead to the areas question, shows %q", got)
	}
	// The new question takes focus, so it is read out.
	if got := eval[string](t, ctx, `document.activeElement.tagName`); got != "LEGEND" {
		t.Errorf("focus should be on the new question, is on %s", got)
	}

	click(t, ctx, `#qn-demo button[value="back"]`)
	sleep(t, ctx, 500)
	if !eval[bool](t, ctx, `document.querySelector('#qn-demo input[value="sim"]').checked`) {
		t.Error("going back should keep the earlier answer checked")
	}
	click(t, ctx, `#qn-demo input[value="nao"]`)
	click(t, ctx, `#qn-demo button[value="next"]`)
	sleep(t, ctx, 500)
	if got := eval[string](t, ctx, qnTitle); got != "Qual o objetivo do próximo experimento?" {
		t.Errorf("\"Não\" should skip the areas question, shows %q", got)
	}
	if got := eval[string](t, ctx, `document.querySelector('#qn-demo p.text-xs').textContent`); got != "Pergunta 2 de 3" {
		t.Errorf("progress should count the shorter path: %q", got)
	}
}

func TestQuestionnaireRequiredShowsErrorAndBackDoesNot(t *testing.T) {
	ctx := openPage(t, "/components/questionnaire")

	click(t, ctx, `#qn-demo button[value="next"]`)
	sleep(t, ctx, 500)
	if got := eval[string](t, ctx, `document.querySelector('#qn-demo [role="alert"]')?.textContent.trim() || ''`); got == "" {
		t.Fatal("submitting a required question empty should show an error")
	}

	click(t, ctx, `#qn-demo input[value="sim"]`)
	click(t, ctx, `#qn-demo button[value="next"]`)
	sleep(t, ctx, 500)
	// On "areas", empty, go back: no error.
	click(t, ctx, `#qn-demo button[value="back"]`)
	sleep(t, ctx, 500)
	if eval[bool](t, ctx, `!!document.querySelector('#qn-demo [role="alert"]')`) {
		t.Error("Back must never validate")
	}
}

// Typing a number in a free-text answer is text, not an option shortcut.
func TestQuestionnaireDigitInTextAnswerIsText(t *testing.T) {
	ctx := openPage(t, "/components/questionnaire?retomar=1")

	focus(t, ctx, `#qn-demo input[type="text"]`)
	press(t, ctx, "2 semanas")
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.querySelector('#qn-demo input[type="text"]').value`); got != "2 semanas" {
		t.Errorf("the answer field holds %q, want \"2 semanas\"", got)
	}
}

func TestQuestionnaireResumesAndFinishes(t *testing.T) {
	ctx := openPage(t, "/components/questionnaire?retomar=1")

	if got := eval[string](t, ctx, `document.querySelector('#qn-demo p.text-xs').textContent`); got != "Pergunta 3 de 4" {
		t.Fatalf("resume should open at question 3 of 4, shows %q", got)
	}
	click(t, ctx, `#qn-demo button[value="skip"]`)
	sleep(t, ctx, 500)
	click(t, ctx, `#qn-demo input[value="1m"]`)
	click(t, ctx, `#qn-demo button[value="next"]`)
	sleep(t, ctx, 500)
	summary := eval[string](t, ctx, `document.querySelector('#qn-demo dl')?.textContent || ''`)
	for _, want := range []string{"Sim", "Genética, Imunologia", "—", "Este mês"} {
		if !strings.Contains(summary, want) {
			t.Errorf("summary is missing %q: %s", want, summary)
		}
	}
}
