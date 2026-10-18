package components

import "time"

var messagesES = messages{
	Close: "Cerrar", Cancel: "Cancelar", Confirm: "Confirmar",

	CarouselRole: "carrusel", CarouselPrevious: "Anterior", CarouselNext: "Siguiente",
	Pagination: "Paginación", PagePrevious: "Anterior", PageNext: "Siguiente",
	PrevMonth: "Mes anterior", NextMonth: "Mes siguiente",

	Breadcrumb: "Ruta de navegación", Search: "Buscar", NoResults: "Sin resultados", ShowOptions: "Mostrar opciones",
	CommandPalette: "Paleta de comandos", Commands: "Comandos",
	Reactions: "Reacciones", Attachments: "Archivos adjuntos", SendFailed: "No se pudo enviar.",
	ContextMenuHint: "Haz clic derecho o pulsa Shift+F10 para ver las acciones.",
	Loading:         "Cargando", JumpToEnd: "Ir al final",
	OpenMenu: "Abrir menú", SideMenu: "Menú lateral", Menu: "Menú", Rating: "Valoración",

	QuestionnaireProgress: "Progreso del cuestionario", QuestionnaireNext: "Siguiente", QuestionnaireBack: "Atrás",

	StatusOnline: "En línea", StatusOffline: "Desconectado", StatusBusy: "Ocupado",

	Results: func(n int) string {
		if n == 1 {
			return "1 resultado"
		}
		return itoa(n) + " resultados"
	},
	Page:       func(n int) string { return "Página " + itoa(n) },
	QuestionOf: func(n, total int) string { return "Pregunta " + itoa(n) + " de " + itoa(total) },
	OfTotal:    func(n, total int) string { return itoa(n) + " de " + itoa(total) },
	Digit:      func(n, total int) string { return "Dígito " + itoa(n) + " de " + itoa(total) },

	UploadOf:     func(name string) string { return "Subiendo " + name },
	ProcessingOf: func(name string) string { return "Procesando " + name },

	WeekdayShort: [7]string{"dom", "lun", "mar", "mié", "jue", "vie", "sáb"},
	WeekdayLong:  esWeekdays,
	MonthLabel:   func(t time.Time) string { return esMonths[t.Month()-1] + " de " + itoa(t.Year()) },
	DayLabel: func(t time.Time) string {
		return esWeekdays[t.Weekday()] + ", " + itoa(t.Day()) + " de " + esMonths[t.Month()-1] + " de " + itoa(t.Year())
	},

	Unread: map[string]string{"one": "(%d nuevo)", "other": "(%d nuevos)"},
}

var esMonths = [12]string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}

var esWeekdays = [7]string{"domingo", "lunes", "martes", "miércoles", "jueves", "viernes", "sábado"}
