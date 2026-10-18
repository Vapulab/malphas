package components

import (
	"context"
	"strings"
	"time"

	"github.com/vapulab/malphas"
)

// Plain Go for the Calendar's date arithmetic, kept out of the template so it
// can be read (and tested) as dates rather than markup. Every date here is a
// day at UTC midnight: a calendar picks days, not instants, and mixing zones
// is how a picked date turns into the day before.

func dayOf(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func isoDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func isoMonth(t time.Time) string { return t.Format("2006-01") }

// monthLabel is the grid's heading in the request's language.
func monthLabel(ctx context.Context, t time.Time) string { return msg(ctx).MonthLabel(t) }

// dayLabel is the button's accessible name: "Tuesday, October 20, 2026". The
// visible text is only the number, which on its own does not say which month
// or weekday a screen reader user has landed on.
func dayLabel(ctx context.Context, t time.Time) string { return msg(ctx).DayLabel(t) }

type weekdayHeader struct{ short, long string }

// weekdayHeaders are the grid's column heads, starting on the language's
// first day of the week.
func weekdayHeaders(ctx context.Context) []weekdayHeader {
	m, start := msg(ctx), malphas.LocaleFrom(ctx).WeekStart()
	out := make([]weekdayHeader, 7)
	for i := range out {
		wd := (int(start) + i) % 7
		out[i] = weekdayHeader{m.WeekdayShort[wd], m.WeekdayLong[wd]}
	}
	return out
}

// calendarWeeks returns six weeks of seven cells starting on start; a zero
// time is a cell outside the month, rendered blank. Six rows always, so the
// grid does not change height between months.
func calendarWeeks(month time.Time, start time.Weekday) [][]time.Time {
	first := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	cursor := first.AddDate(0, 0, -((int(first.Weekday()) - int(start) + 7) % 7))
	weeks := make([][]time.Time, 6)
	for w := range weeks {
		weeks[w] = make([]time.Time, 7)
		for d := range weeks[w] {
			if cursor.Month() == first.Month() {
				weeks[w][d] = cursor
			}
			cursor = cursor.AddDate(0, 0, 1)
		}
	}
	return weeks
}

// withQuery appends a query parameter to a URL that may already have some.
func withQuery(base, params string) string {
	if strings.Contains(base, "?") {
		return base + "&" + params
	}
	return base + "?" + params
}
