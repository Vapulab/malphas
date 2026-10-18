package components

var messagesEN = messages{
	Close: "Close", Cancel: "Cancel", Confirm: "Confirm",

	CarouselRole: "carousel", CarouselPrevious: "Previous", CarouselNext: "Next",
	Pagination: "Pagination", PagePrevious: "Previous", PageNext: "Next",
	PrevMonth: "Previous month", NextMonth: "Next month",

	Breadcrumb: "Breadcrumb", Search: "Search", NoResults: "No results", ShowOptions: "Show options",
	CommandPalette: "Command palette", Commands: "Commands",
	Reactions: "Reactions", Attachments: "Attachments", SendFailed: "Couldn't upload.",
	ContextMenuHint: "Right-click or press Shift+F10 to see the actions.",
	Loading:         "Loading", JumpToEnd: "Jump to the end",
	OpenMenu: "Open menu", SideMenu: "Sidebar", Menu: "Menu", Rating: "Rating",

	QuestionnaireProgress: "Questionnaire progress", QuestionnaireNext: "Next", QuestionnaireBack: "Back",

	StatusOnline: "Online", StatusOffline: "Offline", StatusBusy: "Busy",

	Results: func(n int) string {
		if n == 1 {
			return "1 result"
		}
		return itoa(n) + " results"
	},
	Page:       func(n int) string { return "Page " + itoa(n) },
	QuestionOf: func(n, total int) string { return "Question " + itoa(n) + " of " + itoa(total) },
	OfTotal:    func(n, total int) string { return itoa(n) + " of " + itoa(total) },
	Digit:      func(n, total int) string { return "Digit " + itoa(n) + " of " + itoa(total) },

	UploadOf:     func(name string) string { return "Uploading " + name },
	ProcessingOf: func(name string) string { return "Processing " + name },

	Unread: map[string]string{"one": "(%d new)", "other": "(%d new)"},
}
