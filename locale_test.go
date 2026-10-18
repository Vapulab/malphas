package malphas_test

import (
	"context"
	"slices"
	"testing"

	"github.com/vapulab/malphas"
)

func TestLocaleCodes(t *testing.T) {
	want := map[malphas.Locale][4]string{
		malphas.EN: {"en", "en", "English", "Sunday"},
		malphas.PT: {"pt-BR", "pt", "Português", "Sunday"},
		malphas.ES: {"es", "es", "Español", "Monday"},
		malphas.RU: {"ru", "ru", "Русский", "Monday"},
		malphas.ZH: {"zh-Hans", "zh", "中文", "Monday"},
	}
	for l, w := range want {
		got := [4]string{l.Tag(), l.Segment(), l.NativeName(), l.WeekStart().String()}
		if got != w {
			t.Errorf("locale %d: got %v, want %v", l, got, w)
		}
	}
}

func TestParseLocale(t *testing.T) {
	ok := map[string]malphas.Locale{
		"en": malphas.EN, "pt-BR": malphas.PT, "pt": malphas.PT, "PT-br": malphas.PT,
		"es": malphas.ES, "ru": malphas.RU, "zh-Hans": malphas.ZH, "zh": malphas.ZH,
	}
	for s, want := range ok {
		if got, found := malphas.ParseLocale(s); !found || got != want {
			t.Errorf("ParseLocale(%q) = %v, %v; want %v, true", s, got, found, want)
		}
	}
	for _, s := range []string{"de", "", "pt-PT-x", "english"} {
		if _, found := malphas.ParseLocale(s); found {
			t.Errorf("ParseLocale(%q) should not be found", s)
		}
	}
}

func TestLocaleFromDefaultsToEnglish(t *testing.T) {
	if got := malphas.LocaleFrom(context.Background()); got != malphas.EN {
		t.Errorf("got %v, want EN", got)
	}
}

func TestLocaleFromRoundTrip(t *testing.T) {
	ctx := malphas.WithLocale(context.Background(), malphas.RU)
	if got := malphas.LocaleFrom(ctx); got != malphas.RU {
		t.Errorf("got %v, want RU", got)
	}
}

// A host bug that puts a value out of range in the context must not index
// past the tables: everything behaves as English.
func TestInvalidLocaleFallsBackToEnglish(t *testing.T) {
	ctx := malphas.WithLocale(context.Background(), malphas.Locale(99))
	if got := malphas.LocaleFrom(ctx); got != malphas.EN {
		t.Errorf("LocaleFrom = %v, want EN", got)
	}
	if got := malphas.Locale(99).Tag(); got != "en" {
		t.Errorf("Locale(99).Tag() = %q, want en", got)
	}
	if got := malphas.Locale(-1).Segment(); got != "en" {
		t.Errorf("Locale(-1).Segment() = %q, want en", got)
	}
}

func TestLocalesOrder(t *testing.T) {
	want := []malphas.Locale{malphas.EN, malphas.PT, malphas.ES, malphas.RU, malphas.ZH}
	if got := malphas.Locales(); !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
