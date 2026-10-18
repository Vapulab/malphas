package components

import "time"

var messagesPT = messages{
	Close: "Fechar", Cancel: "Cancelar", Confirm: "Confirmar",

	CarouselRole: "carrossel", CarouselPrevious: "Anterior", CarouselNext: "Próximo",
	Pagination: "Paginação", PagePrevious: "Anterior", PageNext: "Próxima",
	PrevMonth: "Mês anterior", NextMonth: "Próximo mês",

	Breadcrumb: "Trilha de navegação", Search: "Buscar", NoResults: "Nenhum resultado", ShowOptions: "Mostrar opções",
	CommandPalette: "Paleta de comandos", Commands: "Comandos",
	Reactions: "Reações", Attachments: "Anexos", SendFailed: "Não foi possível enviar.",
	ContextMenuHint: "Clique com o botão direito ou use Shift+F10 para ver as ações.",
	Loading:         "Carregando", JumpToEnd: "Ir para o fim",
	OpenMenu: "Abrir menu", SideMenu: "Menu lateral", Menu: "Menu", Rating: "Avaliação",

	QuestionnaireProgress: "Progresso do questionário", QuestionnaireNext: "Próxima", QuestionnaireBack: "Voltar",

	StatusOnline: "Online", StatusOffline: "Offline", StatusBusy: "Ocupado",

	Results: func(n int) string {
		if n == 1 {
			return "1 resultado"
		}
		return itoa(n) + " resultados"
	},
	Page:       func(n int) string { return "Página " + itoa(n) },
	QuestionOf: func(n, total int) string { return "Pergunta " + itoa(n) + " de " + itoa(total) },
	OfTotal:    func(n, total int) string { return itoa(n) + " de " + itoa(total) },
	Digit:      func(n, total int) string { return "Dígito " + itoa(n) + " de " + itoa(total) },

	UploadOf:     func(name string) string { return "Envio de " + name },
	ProcessingOf: func(name string) string { return "Processamento de " + name },

	WeekdayShort: [7]string{"Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"},
	WeekdayLong:  ptWeekdays,
	MonthLabel:   func(t time.Time) string { return ptMonths[t.Month()-1] + " de " + itoa(t.Year()) },
	DayLabel: func(t time.Time) string {
		return ptWeekdays[t.Weekday()] + ", " + itoa(t.Day()) + " de " + ptMonths[t.Month()-1] + " de " + itoa(t.Year())
	},

	Unread: map[string]string{"one": "(%d nova)", "other": "(%d novas)"},
}

var ptMonths = [12]string{"janeiro", "fevereiro", "março", "abril", "maio", "junho", "julho", "agosto", "setembro", "outubro", "novembro", "dezembro"}

var ptWeekdays = [7]string{"domingo", "segunda-feira", "terça-feira", "quarta-feira", "quinta-feira", "sexta-feira", "sábado"}
