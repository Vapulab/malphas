package components

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/vapulab/malphas"
)

func in(l malphas.Locale) context.Context { return malphas.WithLocale(context.Background(), l) }

func TestMonthAndDayLabels(t *testing.T) {
	d := time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC) // a Tuesday
	cases := map[malphas.Locale][2]string{
		malphas.EN: {"October 2026", "Tuesday, October 20, 2026"},
		malphas.PT: {"outubro de 2026", "terça-feira, 20 de outubro de 2026"},
		malphas.ES: {"octubre de 2026", "martes, 20 de octubre de 2026"},
		malphas.RU: {"Октябрь 2026", "вторник, 20 октября 2026 года"},
		malphas.ZH: {"2026年10月", "2026年10月20日 星期二"},
	}
	for l, want := range cases {
		if got := monthLabel(in(l), d); got != want[0] {
			t.Errorf("%s monthLabel = %q, want %q", l.Segment(), got, want[0])
		}
		if got := dayLabel(in(l), d); got != want[1] {
			t.Errorf("%s dayLabel = %q, want %q", l.Segment(), got, want[1])
		}
	}
}

func TestWeekdayHeadersFollowWeekStart(t *testing.T) {
	cases := map[malphas.Locale][]string{
		malphas.EN: {"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
		malphas.PT: {"Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"},
		malphas.RU: {"пн", "вт", "ср", "чт", "пт", "сб", "вс"},
		malphas.ZH: {"一", "二", "三", "四", "五", "六", "日"},
	}
	for l, want := range cases {
		got := make([]string, 0, 7)
		for _, h := range weekdayHeaders(in(l)) {
			got = append(got, h.short)
		}
		if !slices.Equal(got, want) {
			t.Errorf("%s headers = %v, want %v", l.Segment(), got, want)
		}
	}
}

// Monday-first grids: six weeks of seven, every day of the month exactly
// once, and the first filled cell sitting in its weekday's column.
func TestCalendarWeeksMondayStart(t *testing.T) {
	for _, month := range []time.Time{
		time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),  // starts on a Sunday: the last column
		time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC), // 31 days, starts on a Tuesday
		time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),  // starts on a Monday: the first column
	} {
		weeks := calendarWeeks(month, time.Monday)
		if len(weeks) != 6 {
			t.Fatalf("%s: %d weeks, want 6", month.Format("2006-01"), len(weeks))
		}
		seen := map[int]int{}
		for _, w := range weeks {
			if len(w) != 7 {
				t.Fatalf("%s: a week of %d days", month.Format("2006-01"), len(w))
			}
			for col, d := range w {
				if d.IsZero() {
					continue
				}
				seen[d.Day()]++
				if wantCol := (int(d.Weekday()) - int(time.Monday) + 7) % 7; col != wantCol {
					t.Errorf("%s: %s in column %d, want %d", month.Format("2006-01"), d.Format("Mon 02"), col, wantCol)
				}
			}
		}
		last := month.AddDate(0, 1, -1).Day()
		for day := 1; day <= last; day++ {
			if seen[day] != 1 {
				t.Errorf("%s: day %d appears %d times", month.Format("2006-01"), day, seen[day])
			}
		}
	}
}
