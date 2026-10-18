package components

import "time"

var messagesRU = messages{
	Close: "Закрыть", Cancel: "Отмена", Confirm: "Подтвердить",

	CarouselRole: "карусель", CarouselPrevious: "Назад", CarouselNext: "Вперёд",
	Pagination: "Постраничная навигация", PagePrevious: "Предыдущая", PageNext: "Следующая",
	PrevMonth: "Предыдущий месяц", NextMonth: "Следующий месяц",

	Breadcrumb: "Навигационная цепочка", Search: "Поиск", NoResults: "Ничего не найдено", ShowOptions: "Показать варианты",
	CommandPalette: "Палитра команд", Commands: "Команды",
	Reactions: "Реакции", Attachments: "Вложения", SendFailed: "Не удалось отправить.",
	ContextMenuHint: "Щёлкните правой кнопкой мыши или нажмите Shift+F10, чтобы увидеть действия.",
	Loading:         "Загрузка", JumpToEnd: "В конец",
	OpenMenu: "Открыть меню", SideMenu: "Боковое меню", Menu: "Меню", Rating: "Оценка",

	QuestionnaireProgress: "Прогресс опроса", QuestionnaireNext: "Далее", QuestionnaireBack: "Назад",
	QuestionnaireFinish: "Завершить", QuestionnaireSkip: "Пропустить",

	StatusOnline: "В сети", StatusOffline: "Не в сети", StatusBusy: "Занят",

	Results: func(n int) string {
		return itoa(n) + " " + ruPlural(n, "результат", "результата", "результатов")
	},
	Page:       func(n int) string { return "Страница " + itoa(n) },
	QuestionOf: func(n, total int) string { return "Вопрос " + itoa(n) + " из " + itoa(total) },
	OfTotal:    func(n, total int) string { return itoa(n) + " из " + itoa(total) },
	Digit:      func(n, total int) string { return "Цифра " + itoa(n) + " из " + itoa(total) },

	UploadOf:     func(name string) string { return "Загрузка " + name },
	ProcessingOf: func(name string) string { return "Обработка " + name },

	WeekdayShort: [7]string{"вс", "пн", "вт", "ср", "чт", "пт", "сб"},
	WeekdayLong:  ruWeekdays,
	// The heading names the month on its own (nominative); inside a date the
	// month is genitive: "октябрь 2026", but "20 октября 2026 г.".
	MonthLabel: func(t time.Time) string { return ruMonths[t.Month()-1] + " " + itoa(t.Year()) },
	DayLabel: func(t time.Time) string {
		return ruWeekdays[t.Weekday()] + ", " + itoa(t.Day()) + " " + ruMonthsGenitive[t.Month()-1] + " " + itoa(t.Year()) + " г."
	},

	Unread: map[string]string{"one": "(%d новое)", "few": "(%d новых)", "many": "(%d новых)", "other": "(%d новых)"},
}

// ruPlural picks the Russian form for n: 1, 21, 101 take one; 2–4, 22–24
// take few; everything else, including 11–14, takes many.
func ruPlural(n int, one, few, many string) string {
	n10, n100 := n%10, n%100
	switch {
	case n10 == 1 && n100 != 11:
		return one
	case n10 >= 2 && n10 <= 4 && (n100 < 12 || n100 > 14):
		return few
	default:
		return many
	}
}

var (
	ruMonths         = [12]string{"январь", "февраль", "март", "апрель", "май", "июнь", "июль", "август", "сентябрь", "октябрь", "ноябрь", "декабрь"}
	ruMonthsGenitive = [12]string{"января", "февраля", "марта", "апреля", "мая", "июня", "июля", "августа", "сентября", "октября", "ноября", "декабря"}
)

var ruWeekdays = [7]string{"воскресенье", "понедельник", "вторник", "среда", "четверг", "пятница", "суббота"}
