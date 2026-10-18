package components

import "time"

var messagesEN = messages{
	Close: "Close", Cancel: "Cancel", Confirm: "Confirm",

	CarouselRole: "carousel", SlideRole: "slide", CarouselPrevious: "Previous", CarouselNext: "Next",
	Pagination: "Pagination", PagePrevious: "Previous", PageNext: "Next",
	PrevMonth: "Previous month", NextMonth: "Next month",

	Breadcrumb: "Breadcrumb", Search: "Search", NoResults: "No results", ShowOptions: "Show options",
	CommandPalette: "Command palette", Commands: "Commands",
	Reactions: "Reactions", Attachments: "Attachments", SendFailed: "Upload failed.",
	ContextMenuHint: "Right-click or press Shift+F10 for actions.",
	Loading:         "Loading", JumpToEnd: "Jump to latest",
	OpenMenu: "Open menu", SideMenu: "Sidebar", Menu: "Menu", Rating: "Rating",

	QuestionnaireProgress: "Questionnaire progress", QuestionnaireNext: "Next", QuestionnaireBack: "Back",
	QuestionnaireFinish: "Finish", QuestionnaireSkip: "Skip",

	StatusOnline: "Online", StatusOffline: "Offline", StatusBusy: "Busy",

	Results: func(n int) string {
		if n == 1 {
			return "1 result"
		}
		return itoa(n) + " results"
	},
	Page:       func(n int) string { return "Page " + itoa(n) },
	QuestionOf: func(n, total int) string { return "Question " + itoa(n) + " of " + itoa(total) },
	SlideOf:    func(n, total int) string { return itoa(n) + " of " + itoa(total) },
	RatingOf:   func(n, total int) string { return itoa(n) + " of " + itoa(total) },
	Digit:      func(n, total int) string { return "Digit " + itoa(n) + " of " + itoa(total) },

	UploadOf:     func(name string) string { return "Uploading " + name },
	ProcessingOf: func(name string) string { return "Processing " + name },

	UploadingStatus:  func(percent int) string { return "Uploading " + itoa(percent) + "%" },
	ProcessingStatus: "Processing…",

	WeekdayShort: [7]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
	WeekdayLong:  [7]string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
	MonthLabel:   func(t time.Time) string { return t.Month().String() + " " + itoa(t.Year()) },
	DayLabel: func(t time.Time) string {
		return t.Weekday().String() + ", " + t.Month().String() + " " + itoa(t.Day()) + ", " + itoa(t.Year())
	},

	Unread: map[string]string{"one": "(%d new)", "other": "(%d new)"},
}
