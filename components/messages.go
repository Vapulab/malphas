package components

import (
	"context"
	"strconv"

	"github.com/vapulab/malphas"
)

// messages is every piece of text a component writes on its own: accessible
// names, empty states, counters. Each language fills the whole struct, so a
// language missing a text does not compile. Props that already name a label
// (CancelLabel, Empty, Label…) still win over these.
type messages struct {
	Close, Cancel, Confirm string

	CarouselRole, CarouselPrevious, CarouselNext string
	Pagination, PagePrevious, PageNext           string
	PrevMonth, NextMonth                         string

	Breadcrumb, Search, NoResults, ShowOptions string
	CommandPalette, Commands                   string
	Reactions, Attachments, SendFailed         string
	ContextMenuHint, Loading, JumpToEnd        string
	OpenMenu, SideMenu, Menu, Rating           string

	QuestionnaireProgress, QuestionnaireNext, QuestionnaireBack string

	StatusOnline, StatusOffline, StatusBusy string

	Results    func(n int) string
	Page       func(n int) string
	QuestionOf func(n, total int) string
	OfTotal    func(n, total int) string
	Digit      func(n, total int) string

	UploadOf     func(name string) string
	ProcessingOf func(name string) string

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
