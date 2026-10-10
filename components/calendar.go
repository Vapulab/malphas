package components

import (
	"strconv"
	"strings"
	"time"
)

// Plain Go for the Calendar's date arithmetic, kept out of the template so it
// can be read (and tested) as dates rather than markup. Every date here is a
// day at UTC midnight: a calendar picks days, not instants, and mixing zones
// is how a picked date turns into the day before.

var monthNames = [...]string{
	"janeiro", "fevereiro", "março", "abril", "maio", "junho",
	"julho", "agosto", "setembro", "outubro", "novembro", "dezembro",
}

var weekdayShort = [...]string{"Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"}

var weekdayLong = [...]string{
	"domingo", "segunda-feira", "terça-feira", "quarta-feira",
	"quinta-feira", "sexta-feira", "sábado",
}

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

func monthLabel(t time.Time) string {
	return monthNames[t.Month()-1] + " de " + strconv.Itoa(t.Year())
}

// dayLabel is the button's accessible name: "terça-feira, 20 de outubro de
// 2026". The visible text is only the number, which on its own does not say
// which month or weekday a screen reader user has landed on.
func dayLabel(t time.Time) string {
	return weekdayLong[t.Weekday()] + ", " + strconv.Itoa(t.Day()) + " de " + monthLabel(t)
}

// calendarWeeks returns six weeks of seven cells starting on Sunday; a zero
// time is a cell outside the month, rendered blank. Six rows always, so the
// grid does not change height between months.
func calendarWeeks(month time.Time) [][]time.Time {
	first := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	cursor := first.AddDate(0, 0, -int(first.Weekday()))
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
