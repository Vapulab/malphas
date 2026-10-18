package e2e_test

import (
	"testing"

	"github.com/chromedp/chromedp/kb"
)

// Week starts on Monday in Russian and on Sunday in English. Home and End go to
// the ends of the week as the grid shows it: from Sunday 18 Oct 2026, Home in
// Russian goes back to Monday 12, and End from Monday 19 goes to Sunday 25.
func TestCalendarHomeEndFollowWeekStart(t *testing.T) {
	cases := []struct {
		locale, from, key, want string
	}{
		{"ru", "2026-10-18", kb.Home, "2026-10-12"},
		{"ru", "2026-10-19", kb.End, "2026-10-25"},
		{"en", "2026-10-21", kb.Home, "2026-10-18"},
		{"en", "2026-10-18", kb.End, "2026-10-24"},
	}
	for _, c := range cases {
		ctx := openPage(t, "/fixtures/calendar?locale="+c.locale+"&month=2026-10")
		focus(t, ctx, `button[data-date="`+c.from+`"]`)
		press(t, ctx, c.key)
		settle(t, ctx)
		if got := eval[string](t, ctx, `document.activeElement.dataset.date || ''`); got != c.want {
			t.Errorf("%s: %s from %s went to %q, want %s", c.locale, c.key, c.from, got, c.want)
		}
	}
}

// The page's own words follow the language: the Russian grid is headed in
// Russian and its first column is Monday.
func TestFixtureRendersInRussian(t *testing.T) {
	ctx := openPage(t, "/fixtures/calendar?locale=ru&month=2026-10")
	if got := eval[string](t, ctx, `document.querySelector('#fx-cal-title').textContent.trim()`); got != "октябрь 2026" {
		t.Errorf("heading = %q, want октябрь 2026", got)
	}
	if got := eval[string](t, ctx, `document.querySelector('#fx-cal th').textContent.trim()`); got != "пн" {
		t.Errorf("first column = %q, want пн", got)
	}
}
