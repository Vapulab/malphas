package components_test

import (
	"strings"
	"testing"
	"time"

	"github.com/vapulab/malphas/components"
)

func day(y int, m time.Month, d int) time.Time { return time.Date(y, m, d, 0, 0, 0, 0, time.UTC) }

func calendarFixture() components.CalendarProps {
	return components.CalendarProps{
		ID: "cal", Month: day(2026, 10, 1), Name: "entrega",
		Selected: day(2026, 10, 20), Min: day(2026, 10, 5), Today: day(2026, 10, 10),
		NavURL: "/components/calendar",
	}
}

// October 2026 starts on a Thursday: with weeks starting on Sunday the grid
// opens with four blank cells (Sun to Wed), and every day is a button.
func TestCalendarGridStartsOnSunday(t *testing.T) {
	html := renderToString(t, components.Calendar(calendarFixture()))
	if n := strings.Count(html, `role="gridcell"`); n != 42 {
		t.Errorf("expected a 6x7 grid, got %d cells", n)
	}
	if n := strings.Count(html, `data-date="2026-10-`); n != 31 {
		t.Errorf("expected 31 day buttons, got %d", n)
	}
	first := strings.Index(html, `role="gridcell"`)
	oct1 := strings.Index(html, `data-date="2026-10-01"`)
	if n := strings.Count(html[first:oct1], `role="gridcell"`); n != 5 {
		t.Errorf("1 Oct 2026 is a Thursday: it should be the 5th cell, is the %dth", n)
	}
	for _, wd := range []string{"Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"} {
		if !strings.Contains(html, ">"+wd+"<") {
			t.Errorf("missing weekday %s", wd)
		}
	}
}

func TestCalendarGridARIA(t *testing.T) {
	html := renderToString(t, components.Calendar(calendarFixture()))
	for _, want := range []string{
		`role="grid"`,
		`aria-labelledby="cal-title"`,
		`id="cal-title"`,
		`outubro de 2026`,
		`aria-label="terça-feira, 20 de outubro de 2026"`,
		`<input type="hidden" name="entrega" value="2026-10-20"`,
		`aria-current="date"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	// Only the selected day is selected; days before Min are disabled but
	// stay focusable, as the APG date grid does.
	if n := strings.Count(html, `aria-selected="true"`); n != 1 {
		t.Errorf("expected one selected cell, got %d", n)
	}
	d4 := html[strings.Index(html, `data-date="2026-10-04"`)-400 : strings.Index(html, `data-date="2026-10-04"`)+300]
	if !strings.Contains(d4, `aria-disabled="true"`) {
		t.Errorf("4 Oct is before Min and should be aria-disabled: %s", d4)
	}
	// One tab stop: the selected day.
	if n := strings.Count(html, `tabindex="0"`); n != 1 {
		t.Errorf("expected one tab stop, got %d", n)
	}
}

func TestCalendarNavLinksCarryTheMonth(t *testing.T) {
	html := renderToString(t, components.Calendar(calendarFixture()))
	for _, want := range []string{
		`href="/components/calendar?mes=2026-09"`,
		`href="/components/calendar?mes=2026-11"`,
		`hx-select="#cal"`,
		`hx-target="#cal"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}

func TestCalendarRangeMarksTheSpan(t *testing.T) {
	props := components.CalendarProps{
		ID: "p", Month: day(2026, 10, 1), Mode: "range", StartName: "inicio", EndName: "fim",
		Start: day(2026, 10, 10), End: day(2026, 10, 14), NavURL: "/x?a=1",
	}
	html := renderToString(t, components.Calendar(props))
	if n := strings.Count(html, `aria-selected="true"`); n != 5 {
		t.Errorf("10 to 14 inclusive is 5 selected days, got %d", n)
	}
	for _, want := range []string{`name="inicio" value="2026-10-10"`, `name="fim" value="2026-10-14"`, `href="/x?a=1&amp;mes=2026-09"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}

// An inverted pair from a caller or the URL is reordered by the server too,
// not only by the client: the hidden inputs must never submit end < start.
func TestCalendarNormalisesAnInvertedRange(t *testing.T) {
	props := components.CalendarProps{
		ID: "p", Month: day(2026, 10, 1), Mode: "range", StartName: "inicio", EndName: "fim",
		Start: day(2026, 10, 20), End: day(2026, 10, 10), NavURL: "/x",
	}
	html := renderToString(t, components.Calendar(props))
	if !strings.Contains(html, `name="inicio" value="2026-10-10"`) || !strings.Contains(html, `name="fim" value="2026-10-20"`) {
		t.Errorf("the period should be reordered: %s", html)
	}
	if n := strings.Count(html, `aria-selected="true"`); n != 11 {
		t.Errorf("10 to 20 is 11 selected days, got %d", n)
	}
}

// A selection the user could not have made (before Min) is not shown as one.
func TestCalendarIgnoresADisabledSelection(t *testing.T) {
	props := calendarFixture()
	props.Selected = day(2026, 10, 1) // before Min (5 Oct)
	html := renderToString(t, components.Calendar(props))
	if strings.Contains(html, `aria-selected="true"`) || !strings.Contains(html, `name="entrega" value=""`) {
		t.Errorf("a disabled day must not render as selected: %s", html)
	}
}

// data-autofocus belongs to the requested day only; when that day is not in
// the rendered month, nothing takes focus on load.
func TestCalendarAutofocusOnlyOnTheRequestedDay(t *testing.T) {
	props := calendarFixture()
	props.Month = day(2026, 11, 1)
	props.Focus = day(2026, 10, 31)
	html := renderToString(t, components.Calendar(props))
	if strings.Contains(html, "data-autofocus") {
		t.Errorf("focus target is outside November; nothing should autofocus: %s", html)
	}
}

// Single mode keeps its selection in a hidden input even without a Name, so
// a month change does not lose it.
func TestCalendarSingleWithoutNameKeepsItsValue(t *testing.T) {
	props := calendarFixture()
	props.Name = ""
	html := renderToString(t, components.Calendar(props))
	if !strings.Contains(html, `data-cal="value"`) || strings.Contains(html, `name=""`) {
		t.Errorf("expected an unnamed hidden input carrying the value: %s", html)
	}
}
