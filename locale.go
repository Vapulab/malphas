package malphas

import (
	"context"
	"strings"
	"time"
)

// Locale is one of the languages the components speak. The zero value is
// English, the default. The package never decides a request's language: the
// host reads the URL, a cookie or a header and calls WithLocale.
type Locale int

const (
	EN Locale = iota
	PT
	ES
	RU
	ZH
)

type localeInfo struct {
	tag, segment, native string
	weekStart            time.Weekday
}

// Indexed by Locale. Any value out of range behaves as EN.
var locales = [...]localeInfo{
	EN: {"en", "en", "English", time.Sunday},
	PT: {"pt-BR", "pt", "Português", time.Sunday},
	ES: {"es", "es", "Español", time.Monday},
	RU: {"ru", "ru", "Русский", time.Monday},
	ZH: {"zh-Hans", "zh", "中文", time.Monday},
}

// Locales lists every supported language, in a stable order.
func Locales() []Locale { return []Locale{EN, PT, ES, RU, ZH} }

func (l Locale) valid() Locale {
	if l < 0 || int(l) >= len(locales) {
		return EN
	}
	return l
}

// Tag is the BCP 47 tag, for the lang attribute and hreflang.
func (l Locale) Tag() string { return locales[l.valid()].tag }

// Segment is the short code for URLs: en, pt, es, ru, zh.
func (l Locale) Segment() string { return locales[l.valid()].segment }

// NativeName is the language's name in itself, for language pickers.
func (l Locale) NativeName() string { return locales[l.valid()].native }

// WeekStart is the first day of the week in calendars.
func (l Locale) WeekStart() time.Weekday { return locales[l.valid()].weekStart }

// ParseLocale accepts a tag or a segment, case-insensitively, and regional
// tags by their language: "es-MX" is ES, "zh-CN" is ZH. Traditional Chinese
// (zh-Hant, zh-TW, zh-HK, zh-MO) is not: the ZH texts are Simplified.
func ParseLocale(s string) (Locale, bool) {
	for _, l := range Locales() {
		if strings.EqualFold(s, l.Tag()) || strings.EqualFold(s, l.Segment()) {
			return l, true
		}
	}
	parts := strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return r == '-' || r == '_' })
	if len(parts) < 2 {
		return EN, false
	}
	if parts[0] == "zh" {
		for _, sub := range parts[1:] {
			if sub == "hant" || sub == "tw" || sub == "hk" || sub == "mo" {
				return EN, false
			}
		}
	}
	for _, l := range Locales() {
		if parts[0] == l.Segment() {
			return l, true
		}
	}
	return EN, false
}

type localeKey struct{}

// WithLocale returns ctx carrying the language every component under it
// renders in.
func WithLocale(ctx context.Context, l Locale) context.Context {
	return context.WithValue(ctx, localeKey{}, l)
}

// LocaleFrom is the language in ctx, EN when none was set.
func LocaleFrom(ctx context.Context) Locale {
	l, _ := ctx.Value(localeKey{}).(Locale)
	return l.valid()
}
