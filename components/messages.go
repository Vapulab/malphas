package components

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/vapulab/malphas"
)

// messages is every piece of text a component writes on its own: accessible
// names, empty states, counters. Each language fills the whole struct, so a
// language missing a text does not compile. Props that already name a label
// (CancelLabel, Empty, Label…) still win over these.
type messages struct {
	Close, Cancel, Confirm string

	CarouselRole, SlideRole, CarouselPrevious, CarouselNext string
	Pagination, PagePrevious, PageNext                      string
	PrevMonth, NextMonth                                    string

	Breadcrumb, Search, NoResults, ShowOptions string
	CommandPalette, Commands                   string
	Reactions, Attachments, SendFailed         string
	ContextMenuHint, Loading, JumpToEnd        string
	OpenMenu, SideMenu, Menu, Rating           string

	QuestionnaireProgress, QuestionnaireNext, QuestionnaireBack string
	QuestionnaireFinish, QuestionnaireSkip                      string

	StatusOnline, StatusOffline, StatusBusy string

	Results    func(n int) string
	Page       func(n int) string
	QuestionOf func(n, total int) string
	SlideOf    func(n, total int) string // a carousel slide's position
	RatingOf   func(n, total int) string // a rating's value, in stars
	Digit      func(n, total int) string

	UploadOf     func(name string) string
	ProcessingOf func(name string) string

	UploadingStatus  func(percent int) string // the visible status line
	ProcessingStatus string

	// Calendar. Weekday arrays are indexed by time.Weekday (Sunday = 0); the
	// grid reorders them from the language's first day of the week.
	WeekdayShort, WeekdayLong [7]string
	MonthLabel                func(t time.Time) string // the grid's heading: "October 2026"
	DayLabel                  func(t time.Time) string // a day's accessible name

	// Unread is the jump button's counter, built in the browser. Keys are CLDR
	// plural categories (one, few, many, other); %d is the count. The script
	// picks the form with Intl.PluralRules for the page's lang.
	Unread map[string]string
}

var allMessages = map[malphas.Locale]messages{
	malphas.EN: messagesEN,
	malphas.PT: messagesPT,
	malphas.ES: messagesES,
	malphas.RU: messagesRU,
	malphas.ZH: messagesZH,
}

// msg is the text for the language in ctx.
func msg(ctx context.Context) messages { return allMessages[malphas.LocaleFrom(ctx)] }

func itoa(n int) string { return strconv.Itoa(n) }

// unreadFormsJSON hands the jump button's plural forms to the script.
func unreadFormsJSON(ctx context.Context) string {
	b, _ := json.Marshal(msg(ctx).Unread)
	return string(b)
}
