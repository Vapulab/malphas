package showcase

// catalog is the source of truth for the docs. Adding a component here is all
// it takes to get an index entry and a /components/{slug} page — no new route.
func catalog() []Component {
	return []Component{
		{
			Slug:        "button",
			Name:        "Button",
			Category:    CategoryActions,
			Description: "Botão clicável com variantes de cor, tamanhos e suporte a atributos arbitrários (HTMX, Alpine).",
			Examples: []Example{
				{
					Title:  "Variantes",
					Code:   `@components.Button(components.ButtonProps{Variant: "danger"}) { Danger }`,
					Render: demoButtonVariants(),
				},
				{
					Title:  "Tamanhos",
					Code:   `@components.Button(components.ButtonProps{Size: "lg"}) { Large }`,
					Render: demoButtonSizes(),
				},
				{
					Title:  "Desabilitado",
					Code:   `@components.Button(components.ButtonProps{Attrs: templ.Attributes{"disabled": true}}) { Indisponível }`,
					Render: demoButtonDisabled(),
				},
			},
		},
		{
			Slug:        "modal",
			Name:        "Modal",
			Category:    CategoryActions,
			Description: "Janela de diálogo sobreposta e realmente modal: o foco fica preso dentro, a página atrás não rola nem é alcançável por leitor de tela, e ao fechar o foco volta para quem abriu.",
			Examples: []Example{
				{
					Title: "Abrir e fechar",
					Code: `<div x-data>
  @components.Button(components.ButtonProps{
    Attrs: templ.Attributes{"x-on:click": "$dispatch('open-modal', 'demo-modal')"},
  }) { Abrir modal }
</div>
@components.Modal(components.ModalProps{ID: "demo-modal", Title: "Confirmar ação"}) {
  <p>Conteúdo do modal.</p>
}`,
					Render: demoModal(),
				},
			},
		},
		{
			Slug:        "alert-dialog",
			Name:        "Alert Dialog",
			Category:    CategoryActions,
			Description: "Modal que pede confirmação de uma ação com consequência. role=alertdialog com a descrição anunciada, foco inicial em Cancelar, e clicar fora não conta como resposta. ConfirmAttrs leva hx-delete e afins.",
			Examples: []Example{
				{
					Title: "Excluir com confirmação",
					Code: `@components.AlertDialog(components.AlertDialogProps{
  ID: "excluir", Title: "Excluir o protocolo?", Description: "Não dá para desfazer.",
  ConfirmLabel: "Excluir", ConfirmAttrs: templ.Attributes{"hx-delete": "/protocolos/1"},
})`,
					Render: demoAlertDialog(),
				},
			},
		},
		{
			Slug:        "card",
			Name:        "Card",
			Category:    CategoryDataDisplay,
			Description: "Contêiner com borda e sombra para agrupar um bloco de conteúdo.",
			Examples: []Example{
				{
					Title:  "Simples",
					Code:   `@components.Card(components.CardProps{}) { <p>Conteúdo</p> }`,
					Render: demoCard(),
				},
				{
					Title:  "Com ação",
					Code:   `@components.Card(components.CardProps{Class: "max-w-sm"}) { ... }`,
					Render: demoCardWithActions(),
				},
			},
		},
		{
			Slug:        "item",
			Name:        "Item",
			Category:    CategoryDataDisplay,
			Description: "Linha com mídia, título, descrição e ações. Com Href a linha inteira é clicável, mas as ações continuam fora do link (padrão stretched link). ItemGroup monta a lista.",
			Examples: []Example{
				{
					Title:  "Link com ação e grupo",
					Code:   `@components.Item(components.ItemProps{Title: "Relatório", Href: "/relatorios/1", Actions: baixar()})`,
					Render: demoItem(),
				},
			},
		},
		{
			Slug:        "hover-card",
			Name:        "Hover Card",
			Category:    CategoryDataDisplay,
			Description: "Prévia do que está atrás de um link, ao passar o mouse ou focar com o teclado, com atraso para não piscar. É só complemento: o gatilho é um link de verdade e nada pode existir apenas no card.",
			Examples: []Example{
				{
					Title: "Perfil de quem revisou",
					Code: `@components.HoverCard(components.HoverCardProps{ID: "autor", Label: "@lincoln", Href: "/u/lincoln"}) {
  <p>Pesquisador · 42 protocolos</p>
}`,
					Render: demoHoverCard(),
				},
			},
		},
		{
			Slug:        "bubble",
			Name:        "Bubble",
			Category:    CategoryDataDisplay,
			Description: "Balão de conversa com variantes, alinhado ao início ou ao fim. Mensagens seguidas do mesmo autor achatam os cantos onde se encontram; as reações dizem o que significam e quantas são, em vez de só o emoji.",
			Examples: []Example{
				{
					Title:  "Grupo, reação, própria e citação",
					Code:   `@components.Bubble(components.BubbleProps{Variant: "primary", Align: "end"}) { Repete, e separa uma alíquota. }`,
					Render: demoBubble(),
				},
			},
		},
		{
			Slug:        "marker",
			Name:        "Marker",
			Category:    CategoryDataDisplay,
			Description: "Linha da conversa que ninguém disse: alguém entrou, uma resposta sendo gerada, a data entre dias. Atualizações usam role=status; o divisor mantém o rótulo legível para leitor de tela.",
			Examples: []Example{
				{
					Title:  "Divisor, nota e status",
					Code:   `@components.Marker(components.MarkerProps{Variant: "border", Status: true}) { Gerando o resumo… }`,
					Render: demoMarker(),
				},
			},
		},
		{
			Slug:        "message",
			Name:        "Message",
			Category:    CategoryDataDisplay,
			Description: "Mensagem com avatar, nome, hora e Bubble. Em sequência do mesmo autor, só a primeira mostra o cabeçalho, mas toda mensagem mantém autor e hora no nome acessível.",
			Examples: []Example{
				{
					Title: "Conversa curta",
					Code: `@components.Message(components.MessageProps{
  Author: "Ana Lima", Initials: "AL", Time: "14:02", Datetime: "2026-10-10T14:02", Group: "first",
}) { A curva saiu acima do limiar. }`,
					Render: demoMessage(),
				},
			},
		},
		{
			Slug:        "attachment",
			Name:        "Attachment",
			Category:    CategoryDataDisplay,
			Description: "Arquivo ou imagem num composer, numa mensagem ou numa lista de envio, com o estado em texto (enviando, processando, erro, pronto) além da barra. O nome só vira link quando o arquivo está pronto; as ações continuam alcançáveis, sobretudo no erro. AttachmentGroup rola na horizontal.",
			Examples: []Example{
				{
					Title: "Os quatro estados",
					Code: `@components.AttachmentGroup(components.AttachmentGroupProps{Label: "Anexos", Items: []components.AttachmentProps{
  {Name: "protocolo.pdf", Meta: "PDF · 2,4 MB", State: "done", Href: "/arquivos/protocolo.pdf"},
  {Name: "curva.png", State: "uploading", Progress: 45},
}})`,
					Render: demoAttachment(),
				},
			},
		},
		{
			Slug:        "message-scroller",
			Name:        "Message Scroller",
			Category:    CategoryDataDisplay,
			Description: "Transcrição de conversa que se comporta como chat: acompanha mensagens novas só se você está no fim, mostra \"Ir para o fim\" (com quantas chegaram) quando você subiu, e mantém a posição quando mensagens antigas entram no topo. Não importa como as mensagens chegam: observa o próprio conteúdo. role=log, educado.",
			Examples: []Example{
				{
					Title: "Novas por HTMX, antigas no topo e streaming",
					Code: `@components.MessageScroller(components.MessageScrollerProps{ID: "chat", Label: "Conversa"}) {
  for _, m := range mensagens { @components.Message(...) { { m.Texto } } }
}
// novas: hx-get="/mensagens" hx-target="#chat-content" hx-swap="beforeend"`,
					Render: demoMessageScroller(),
				},
			},
		},
		{
			Slug:        "typography",
			Name:        "Typography",
			Category:    CategoryDataDisplay,
			Description: "Estilo para texto longo cujo HTML vem de fora (Markdown, ajuda, termos): headings, listas, citações e código, com os tokens do Malphas. Também disponível como a classe prose-malphas.",
			Examples: []Example{
				{
					Title:  "Conteúdo misto",
					Code:   `@components.Prose(components.ProseProps{}) { @templ.Raw(markdownRenderizado) }`,
					Render: demoProse(),
				},
			},
		},
		{
			Slug:        "input",
			Name:        "Text input",
			Category:    CategoryDataInput,
			Description: "Campo de texto com label associado por id. Aceita qualquer type nativo.",
			Examples: []Example{
				{
					Title: "Com label e placeholder",
					Code: `@components.Input(components.InputProps{
  ID: "email", Name: "email", Label: "Email",
  Type: "email", Placeholder: "voce@exemplo.com",
})`,
					Render: demoInput(),
				},
				{
					Title:  "Com valor",
					Code:   `@components.Input(components.InputProps{ID: "nome", Name: "nome", Label: "Nome", Value: "Lincoln"})`,
					Render: demoInputWithValue(),
				},
			},
		},
		{
			Slug:        "link",
			Name:        "Link",
			Category:    CategoryActions,
			Description: "Link de texto com variantes de cor. Href passa pelo sanitizador de URL do templ, e links externos ganham rel=noopener.",
			Examples: []Example{
				{
					Title: "Variantes",
					Code: `@components.Link(components.LinkProps{Href: "/todos"}) { Link interno }
@components.Link(components.LinkProps{Href: "https://templ.guide", External: true}) { Externo }`,
					Render: demoLink(),
				},
			},
		},
		{
			Slug:        "badge",
			Name:        "Badge",
			Category:    CategoryDataDisplay,
			Description: "Etiqueta compacta de status ou contagem, sólida ou apenas com contorno.",
			Examples: []Example{
				{
					Title:  "Sólido",
					Code:   `@components.Badge(components.BadgeProps{Variant: "success"}) { Ativo }`,
					Render: demoBadgeSolid(),
				},
				{
					Title:  "Contorno",
					Code:   `@components.Badge(components.BadgeProps{Variant: "primary", Outline: true}) { Primary }`,
					Render: demoBadgeOutline(),
				},
			},
		},
		{
			Slug:        "kbd",
			Name:        "Kbd",
			Category:    CategoryDataDisplay,
			Description: "Representação visual de uma tecla, usando o elemento semântico <kbd>.",
			Examples: []Example{
				{
					Title:  "Atalho",
					Code:   `@components.Kbd(components.KbdProps{}) { Ctrl }`,
					Render: demoKbd(),
				},
			},
		},
		{
			Slug:        "divider",
			Name:        "Divider",
			Category:    CategoryLayout,
			Description: "Linha divisória horizontal ou vertical, com texto opcional no meio. Sem rótulo é decorativa e fica oculta para leitores de tela.",
			Examples: []Example{
				{
					Title: "Horizontal, com rótulo e vertical",
					Code: `@components.Divider(components.DividerProps{})
@components.Divider(components.DividerProps{Label: "ou"})
@components.Divider(components.DividerProps{Vertical: true})`,
					Render: demoDivider(),
				},
			},
		},
		{
			Slug:        "aspect-ratio",
			Name:        "Aspect Ratio",
			Category:    CategoryLayout,
			Description: "Reserva a proporção de uma mídia antes de ela carregar, para a página não pular. Proporção desconhecida cai no 16/9 em vez de colapsar a caixa.",
			Examples: []Example{
				{
					Title:  "16/9 e 1/1",
					Code:   `@components.AspectRatio(components.AspectRatioProps{Ratio: "16/9"}) { <img src="..." alt="..."> }`,
					Render: demoAspectRatio(),
				},
			},
		},
		{
			Slug:        "scroll-area",
			Name:        "Scroll Area",
			Category:    CategoryLayout,
			Description: "Rolagem nativa com a barra no tema. Recebe foco (tabindex=0) para rolar pelo teclado, e por isso exige um Label.",
			Examples: []Example{
				{
					Title:  "Vertical e horizontal",
					Code:   `@components.ScrollArea(components.ScrollAreaProps{Label: "Termos de uso", Class: "max-h-48"}) { ... }`,
					Render: demoScrollArea(),
				},
			},
		},
		{
			Slug:        "resizable",
			Name:        "Resizable",
			Category:    CategoryLayout,
			Description: "Dois painéis com uma divisória que se arrasta ou se move pelo teclado (setas, Shift para 1%, Home/End), no padrão APG de window splitter. O servidor renderiza o tamanho inicial; com ID, o último tamanho fica guardado neste navegador.",
			Examples: []Example{
				{
					Title: "Lado a lado e empilhado",
					Code: `@components.Resizable(components.ResizableProps{
  ID: "editor", Label: "Largura da lista", DefaultSize: 35, Min: 20, Max: 70,
  First: lista(), Second: editor(),
})`,
					Render: demoResizable(),
				},
			},
		},
		{
			Slug:        "footer",
			Name:        "Footer",
			Category:    CategoryLayout,
			Description: "Rodapé estrutural com colunas de links e linha de nota. Anuncia-se como marco de navegação mesmo aninhado, e delega os links ao Link para herdar o saneamento de href.",
			Examples: []Example{
				{
					Title: "Colunas e nota",
					Code: `@components.Footer(components.FooterProps{
  Note: "© 2026 Vapulab",
  Columns: []components.FooterColumn{
    {Title: "Produto", Links: []components.FooterLink{
      {Label: "Componentes", Href: "/components"},
    }},
  },
}) {
  <span class="font-semibold text-fg">vapulab</span>
}`,
					Render: demoFooter(),
				},
			},
		},
		{
			Slug:        "avatar",
			Name:        "Avatar",
			Category:    CategoryDataDisplay,
			Description: "Imagem de perfil redonda ou quadrada, com fallback para iniciais e indicador de presença opcional.",
			Examples: []Example{
				{
					Title:  "Tamanhos e formatos",
					Code:   `@components.Avatar(components.AvatarProps{Initials: "LS", Size: "lg"})`,
					Render: demoAvatar(),
				},
				{
					Title:  "Com status",
					Code:   `@components.Avatar(components.AvatarProps{Initials: "LS", Status: "online"})`,
					Render: demoAvatarStatus(),
				},
			},
		},
		{
			Slug:        "stat",
			Name:        "Stat",
			Category:    CategoryDataDisplay,
			Description: "Bloco para destacar um número com título e contexto. Usa dt/dd, então deve ficar dentro de um <dl>.",
			Examples: []Example{
				{
					Title:  "Grade de métricas",
					Code:   `@components.Stat(components.StatProps{Title: "Receita", Value: "R$ 12.400", Description: "+12% no mês"})`,
					Render: demoStat(),
				},
			},
		},
		{
			Slug:        "table",
			Name:        "Table",
			Category:    CategoryDataDisplay,
			Description: "Tabela para dados densos, com scroll horizontal próprio. Zebra e densidade são definidos uma vez na tabela e alcançam as células por seletor descendente.",
			Examples: []Example{
				{
					Title: "Padrão",
					Code: `@components.Table(components.TableProps{}) {
  @components.THead() {
    @components.TR() {
      @components.TH(components.THProps{}) { Serviço }
    }
  }
  @components.TBody() {
    @components.TR() {
      @components.TD(components.TDProps{}) { vapulab-api }
    }
  }
}`,
					Render: demoTable(),
				},
				{
					Title:  "Zebrada e compacta",
					Code:   `@components.Table(components.TableProps{Zebra: true, Compact: true}) { ... }`,
					Render: demoTableZebraCompact(),
				},
			},
		},
		{
			Slug:        "data-table",
			Name:        "Data Table",
			Category:    CategoryDataDisplay,
			Description: "Tabela com ordenação, busca e paginação feitas no servidor. Tudo é link ou campo GET de verdade (funciona sem JS); com JS, cada ação troca só a região da tabela e a URL guarda o estado. A busca fica fora da região trocada, então digitar nunca é interrompido.",
			Examples: []Example{
				{
					Title: "Protocolos",
					Code: `@components.DataTable(components.DataTableProps{
  ID: "protocolos", URL: "/protocolos", Columns: colunas,
  Sort: r.Sort, Dir: r.Dir, Query: q, Page: r.Page, TotalPages: r.TotalPages, Total: r.Total,
  Searchable: true,
}) { for _, p := range r.Rows { <tr>...</tr> } }`,
					Render: demoDataTable(),
				},
			},
		},
		{
			Slug:        "hero",
			Name:        "Hero",
			Category:    CategoryLayout,
			Description: "Seção de destaque para o topo de uma página. Cuida do enquadramento — espaçamento, centralização e largura máxima — e deixa o conteúdo por sua conta.",
			Examples: []Example{
				{
					Title:  "Chamada principal",
					Code:   `@components.Hero(components.HeroProps{Class: "bg-primary/5"}) { <h1>...</h1> }`,
					Render: demoHero(),
				},
			},
		},
		{
			Slug:        "drawer",
			Name:        "Drawer",
			Category:    CategoryActions,
			Description: "Painel modal que abre pela lateral ou por baixo (Side: \"bottom\"). Mesma primitiva do Modal — foco preso, scroll travado, conteúdo atrás inerte — porque um painel que deixa o Tab escapar é pior que um diálogo fazendo o mesmo: a página parece continuar utilizável.",
			Examples: []Example{
				{
					Title: "Esquerda e direita",
					Code: `@components.Drawer(components.DrawerProps{ID: "nav", Title: "Navegação"}) {
  @components.Menu(components.MenuProps{Entries: entradas})
}
// gatilho, dentro de um escopo Alpine:
// x-on:click="$dispatch('open-drawer', 'nav')"`,
					Render: demoDrawer(),
				},
				{
					Title:  "Por baixo",
					Code:   `@components.Drawer(components.DrawerProps{ID: "compartilhar", Title: "Compartilhar", Side: "bottom"}) { ... }`,
					Render: demoDrawerBottom(),
				},
			},
		},
		{
			Slug:        "navbar",
			Name:        "Navbar",
			Category:    CategoryNavigation,
			Description: "Barra de navegação superior. Sabe abrir um Drawer no mobile, mas não exige um: sem DrawerID, nenhum botão de menu é renderizado.",
			Examples: []Example{
				{
					Title: "Com menu mobile",
					Code: `@components.Navbar(components.NavbarProps{
  Brand: "Vapulab", DrawerID: "menu",
  Entries: []components.MenuEntry{
    {Label: "Componentes", Href: "/", Current: true},
  },
}) { ...conteúdo à direita... }`,
					Render: demoNavbar(),
				},
			},
		},
		{
			Slug:        "menubar",
			Name:        "Menubar",
			Category:    CategoryNavigation,
			Description: "Barra de menus de aplicação (Arquivo, Editar, Exibir). Uma parada de Tab para a barra inteira; ← → andam entre menus e levam o menu aberto junto, ↑ ↓ andam dentro dele, Escape volta ao topo e Tab fecha tudo.",
			Examples: []Example{
				{
					Title: "Editor",
					Code: `@components.Menubar(components.MenubarProps{
  Label: "Editor",
  Menus: []components.MenubarMenu{{Label: "Arquivo", Entries: []components.MenuEntry{...}}, ...},
})`,
					Render: demoMenubar(),
				},
			},
		},
		{
			Slug:        "sidebar",
			Name:        "Sidebar",
			Category:    CategoryNavigation,
			Description: "Moldura de aplicação: navegação lateral em grupos e o conteúdo ao lado. Recolhe para só ícones (os nomes continuam para leitor de tela e no title) e lembra a escolha; no celular vira um botão Menu que abre a mesma navegação num Drawer modal.",
			Examples: []Example{
				{
					Title: "Aplicação com dois grupos",
					Code: `@components.Sidebar(components.SidebarProps{
  ID: "app", Label: "Navegação", Collapsible: true,
  Groups: []components.SidebarGroup{{Label: "Trabalho", Items: []components.SidebarItem{{Label: "Protocolos", Href: "/protocolos", Icon: "◆", Current: true}}}},
}) { conteúdo da página }`,
					Render: demoSidebar(),
				},
			},
		},
		{
			Slug:        "command",
			Name:        "Command",
			Category:    CategoryNavigation,
			Description: "Paleta de comandos: modal com busca sobre ações e páginas agrupadas, aberta por Ctrl+K / ⌘K. O foco fica na busca e as setas andam pelos resultados (aria-activedescendant); a busca começa limpa a cada abertura.",
			Examples: []Example{
				{
					Title: "Ações e páginas",
					Code: `@components.Command(components.CommandProps{
  ID: "paleta", Hotkey: true,
  Groups: []components.CommandGroup{{Heading: "Ações", Items: []components.CommandItem{{Label: "Alternar tema", OnClick: "..."}}}},
})`,
					Render: demoCommand(),
				},
			},
		},
		{
			Slug:        "dropdown",
			Name:        "Dropdown",
			Category:    CategoryActions,
			Description: "Padrão menu button, com role=menu e navegação por setas. É o oposto do Drawer: não prende o foco — Escape devolve ao gatilho, mas Tab sai e fecha, como um menu deve fazer.",
			Examples: []Example{
				{
					Title: "Ações e conta",
					Code: `@components.Dropdown(components.DropdownProps{
  ID: "acoes", Label: "Ações",
  Entries: []components.MenuEntry{
    {Label: "Editar", OnClick: "$dispatch('edit')"},
    components.MenuDivider(),
    {Label: "Excluir", OnClick: "$dispatch('del')", Disabled: true},
  },
})`,
					Render: demoDropdown(),
				},
			},
		},
		{
			Slug:        "popover",
			Name:        "Popover",
			Category:    CategoryActions,
			Description: "Botão que abre um painel não-modal com conteúdo livre. O foco entra no painel, Escape fecha e devolve o foco, e sair com Tab fecha sem prender nada. Dentro de um Modal, Escape fecha só o popover.",
			Examples: []Example{
				{
					Title: "Formulário no painel e dentro de um modal",
					Code: `@components.Popover(components.PopoverProps{ID: "filtros", Label: "Filtros"}) {
  @components.Input(components.InputProps{ID: "autor", Name: "autor", Label: "Autor"})
}`,
					Render: demoPopover(),
				},
			},
		},
		{
			Slug:        "context-menu",
			Name:        "Context Menu",
			Category:    CategoryActions,
			Description: "Menu de ações no botão direito, aberto no ponto do clique e empurrado para dentro da tela perto das bordas. A área recebe foco e abre o mesmo menu com Shift+F10 ou a tecla Menu.",
			Examples: []Example{
				{
					Title: "Ações de um arquivo",
					Code: `@components.ContextMenu(components.ContextMenuProps{
  ID: "arquivo", Label: "Ações do arquivo",
  Entries: []components.MenuEntry{{Label: "Renomear", OnClick: "..."}, ...},
}) { <div>relatorio.pdf</div> }`,
					Render: demoContextMenu(),
				},
			},
		},
		{
			Slug:        "select",
			Name:        "Select",
			Category:    CategoryDataInput,
			Description: "Seleção única sobre o <select> nativo — segunda exceção deliberada ao Alpine, depois do Range. Ganha seletor nativo no mobile, busca por digitação e teclado completo; em troca, não dá para estilizar as opções abertas.",
			Examples: []Example{
				{
					Title:  "Com valor e com placeholder",
					Code:   `@components.Select(components.SelectProps{ID: "plano", Name: "plano", Label: "Plano", Options: opcoes, Value: "pro"})`,
					Render: demoSelect(),
				},
			},
		},
		{
			Slug:        "combobox",
			Name:        "Combobox",
			Category:    CategoryDataInput,
			Description: "Select com busca, para listas longas demais para rolar. O foco fica no campo de texto e a opção destacada é anunciada por aria-activedescendant; um input oculto carrega o valor, então texto livre nunca é enviado.",
			Examples: []Example{
				{
					Title: "Lista longa e com placeholder",
					Code: `@components.Combobox(components.ComboboxProps{
  ID: "estado", Name: "estado", Label: "Estado", Value: "sp",
  Placeholder: "Busque um estado…", Options: estados,
})`,
					Render: demoCombobox(),
				},
			},
		},
		{
			Slug:        "date-picker",
			Name:        "Date Picker",
			Category:    CategoryDataInput,
			Description: "O <input type=\"date\"> nativo, estilizado — mesma escolha do Select: o seletor da plataforma (a roda do celular, a grade do desktop) é melhor em toque, teclado e leitor de tela, e o valor chega sempre como AAAA-MM-DD. Para um calendário visível na página, use o Calendar.",
			Examples: []Example{
				{
					Title:  "Com limites, dentro de um Field",
					Code:   `@components.DatePicker(components.DatePickerProps{ID: "entrega", Name: "entrega", Min: "2026-10-10", Max: "2026-12-31"})`,
					Render: demoDatePicker(),
				},
			},
		},
		{
			Slug:        "calendar",
			Name:        "Calendar",
			Category:    CategoryDataInput,
			Description: "Grade de um mês para escolher um dia ou um período, sempre visível na página. O servidor monta o mês; trocar de mês busca o vizinho e troca só o calendário. Teclado do padrão APG, nomes em pt-BR e semana começando no domingo.",
			Examples: []Example{
				{
					Title: "Dia único com mínimo, e período",
					Code: `@components.Calendar(components.CalendarProps{
  ID: "entrega", Month: mes, Name: "entrega", Selected: escolhido,
  Min: hoje, NavURL: "/agenda",
})`,
					Render: demoCalendar(),
				},
			},
		},
		{
			Slug:        "questionnaire",
			Name:        "Questionnaire",
			Category:    CategoryDataInput,
			Description: "Formulário de várias etapas, uma pergunta por vez (única, múltipla, texto livre, pulável). Cada etapa é um POST e o servidor decide o que vem depois: validação, perguntas condicionais e retomada ficam no Go. Teclas 1–9 escolhem, Enter avança; Voltar nunca valida.",
			Examples: []Example{
				{
					Title: "Com pergunta condicional e retomada",
					Code: `@components.Questionnaire(components.QuestionnaireProps{
  ID: "perfil", Action: "/perfil/responder",
  Step: s.Step, Total: s.Total, Question: s.Question, Answers: s.Answers, Error: s.Error,
})`,
					Render: demoQuestionnaire(),
				},
			},
		},
		{
			Slug:        "label",
			Name:        "Label",
			Category:    CategoryDataInput,
			Description: "Rótulo de campo com marcador opcional de obrigatório. O asterisco é aria-hidden: quem anuncia a obrigatoriedade é o required do próprio campo.",
			Examples: []Example{
				{
					Title:  "Simples e obrigatório",
					Code:   `@components.Label(components.LabelProps{For: "email", Required: true}) { Email }`,
					Render: demoLabel(),
				},
			},
		},
		{
			Slug:        "field",
			Name:        "Field",
			Category:    CategoryDataInput,
			Description: "Junta label, campo, descrição e erro vindo do servidor. ControlAttrs() entrega o aria-describedby e o aria-invalid que o campo precisa carregar.",
			Examples: []Example{
				{
					Title: "Com descrição e com erro",
					Code: `{{ f := components.FieldProps{ID: "site", Label: "Site", Error: "Endereço inválido."} }}
@components.Field(f) {
  @components.Input(components.InputProps{ID: f.ID, Name: "site", Attrs: f.ControlAttrs()})
}`,
					Render: demoField(),
				},
			},
		},
		{
			Slug:        "textarea",
			Name:        "Textarea",
			Category:    CategoryDataInput,
			Description: "Campo de texto de várias linhas. Com MaxLength, ganha um contador que conta como o navegador (emoji vale 2) e já vem pintado do servidor.",
			Examples: []Example{
				{
					Title:  "Com contador e desabilitado",
					Code:   `@components.Textarea(components.TextareaProps{ID: "bio", Name: "bio", Label: "Bio", MaxLength: 160})`,
					Render: demoTextarea(),
				},
			},
		},
		{
			Slug:        "input-group",
			Name:        "Input Group",
			Category:    CategoryDataInput,
			Description: "Campo com prefixo e sufixo colados — moeda, domínio, um botão. O anel de foco envolve o grupo inteiro.",
			Examples: []Example{
				{
					Title: "Prefixo, sufixo e botão",
					Code: `@components.InputGroup(components.InputGroupProps{
  ID: "preco", Name: "preco", Label: "Preço",
  Prefix: components.InputGroupText("R$"),
})`,
					Render: demoInputGroup(),
				},
			},
		},
		{
			Slug:        "toggle-button",
			Name:        "ToggleButton",
			Category:    CategoryDataInput,
			Description: "Botão de dois estados com aria-pressed, como o negrito de um editor. Diferente do Toggle, que é um switch de configuração. Com Name, envia o valor só enquanto pressionado.",
			Examples: []Example{
				{
					Title: "Ícone, pressionado e desabilitado",
					Code: `@components.ToggleButton(components.ToggleButtonProps{Name: "favorito", Label: "Favorito"}) {
  <span aria-hidden="true">★</span>
}`,
					Render: demoToggleButton(),
				},
			},
		},
		{
			Slug:        "toggle-group",
			Name:        "ToggleGroup",
			Category:    CategoryDataInput,
			Description: "Conjunto de ToggleButtons com uma única parada de Tab. Single pressiona no máximo um (e clicar no pressionado solta); multiple alterna cada um. As setas movem só o foco.",
			Examples: []Example{
				{
					Title: "Single e multiple",
					Code: `@components.ToggleGroup(components.ToggleGroupProps{
  Name: "alinhamento", Label: "Alinhamento",
  Items: []components.ToggleGroupItem{{Label: "Esquerda", Value: "esq"}, ...},
  Value: []string{"centro"},
})`,
					Render: demoToggleGroup(),
				},
			},
		},
		{
			Slug:        "otp",
			Name:        "OTP",
			Category:    CategoryDataInput,
			Description: "Código de verificação. O primeiro campo leva autocomplete=one-time-code, que é o que faz iOS e Android oferecerem o código do SMS direto no teclado.",
			Examples: []Example{
				{
					Title:  "6 e 4 dígitos",
					Code:   `@components.OTP(components.OTPProps{Name: "codigo", Label: "Código de verificação"})`,
					Render: demoOTP(),
				},
			},
		},
		{
			Slug:        "carousel",
			Name:        "Carousel",
			Category:    CategoryDataDisplay,
			Description: "Slides horizontais com scroll-snap. Navegável só pelo teclado, com posição anunciada em cada slide, e o rolar respeita prefers-reduced-motion.",
			Examples: []Example{
				{
					Title:  "Três slides",
					Code:   `@components.Carousel(components.CarouselProps{ID: "shots", Label: "Capturas", Slides: slides})`,
					Render: demoCarousel(),
				},
			},
		},
		{
			Slug:        "pagination",
			Name:        "Pagination",
			Category:    CategoryNavigation,
			Description: "Navegação entre páginas, colapsando intervalos longos. A página atual é um span com aria-current, não um link — linkar para onde já se está é um controle morto.",
			Examples: []Example{
				{
					Title:  "Interativa",
					Code:   `@components.Pagination(components.PaginationProps{Current: page, Total: 20, BaseHref: "/components/pagination?page="})`,
					Render: demoPagination(),
				},
				{
					Title:  "Nas pontas",
					Code:   `// Cada exemplo tem seu próprio parâmetro, e o BaseHref carrega o dos outros — sem isso, mover um zera os demais.`,
					Render: demoPaginationEdges(),
				},
			},
		},
		{
			Slug:        "menu",
			Name:        "Menu",
			Category:    CategoryNavigation,
			Description: "Navegação vertical. Deliberadamente não usa role=menu — esse papel é para menus de aplicação, e aplicá-lo a links de site faz o leitor de tela anunciar a página como um menu de programa. Para o padrão menu button, use o Dropdown.",
			Examples: []Example{
				{
					Title: "Com divisor e ação",
					Code: `@components.Menu(components.MenuProps{Entries: []components.MenuEntry{
  {Label: "Início", Href: "/", Current: true},
  components.MenuDivider(),
  {Label: "Sair", OnClick: "$dispatch('logout')"},
}})`,
					Render: demoMenu(),
				},
			},
		},
		{
			Slug:        "collapse",
			Name:        "Collapse",
			Category:    CategoryDataDisplay,
			Description: "Área expansível. O gatilho é um botão de verdade — uma div clicável não recebe foco, não responde a Enter ou Espaço, e é anunciada como nada.",
			Examples: []Example{
				{
					Title:  "Fechado e aberto",
					Code:   `@components.Collapse(components.CollapseProps{ID: "c1", Title: "Detalhes"}) { <p>corpo</p> }`,
					Render: demoCollapse(),
				},
			},
		},
		{
			Slug:        "accordion",
			Name:        "Accordion",
			Category:    CategoryDataDisplay,
			Description: "Grupo de painéis expansíveis, com modo exclusivo opcional em que abrir um fecha os demais.",
			Examples: []Example{
				{
					Title: "Independentes",
					Code: `@components.Accordion(components.AccordionProps{
  ID: "faq",
  Items: []components.AccordionItem{
    {Title: "Precisa de Node?", Content: resposta()},
  },
})`,
					Render: demoAccordion(),
				},
				{
					Title:  "Exclusivo",
					Code:   `@components.Accordion(components.AccordionProps{ID: "faq", Exclusive: true, ...})`,
					Render: demoAccordionExclusive(),
				},
			},
		},
		{
			Slug:        "tabs",
			Name:        "Tabs",
			Category:    CategoryNavigation,
			Description: "Abas seguindo o padrão WAI-ARIA: Tab alcança a faixa uma vez, depois as setas navegam, Home e End vão às pontas e abas desabilitadas são puladas.",
			Examples: []Example{
				{
					Title: "Com aba desabilitada",
					Code: `@components.Tabs(components.TabsProps{
  ID: "conta", Label: "Configurações da conta",
  Items: []components.Tab{
    {Label: "Perfil", Content: perfil()},
    {Label: "Faturamento", Content: faturamento(), Disabled: true},
  },
})`,
					Render: demoTabs(),
				},
			},
		},
		{
			Slug:        "timeline",
			Name:        "Timeline",
			Category:    CategoryDataDisplay,
			Description: "Linha do tempo cronológica de eventos, com trilho conectando os marcadores.",
			Examples: []Example{
				{
					Title: "Eventos de um deploy",
					Code: `@components.Timeline(components.TimelineProps{Items: []components.TimelineItem{
  {Time: "09:00", Title: "Deploy iniciado", Description: "build #481"},
  {Time: "09:04", Title: "Deploy concluído", Variant: "success"},
}})`,
					Render: demoTimeline(),
				},
			},
		},
		{
			Slug:        "countdown",
			Name:        "Countdown",
			Category:    CategoryDataDisplay,
			Description: "Temporizador regressivo. O valor inicial vem do servidor, então o relógio já nasce correto. aria-live fica desligado de propósito: um número anunciado a cada segundo torna a página inutilizável com leitor de tela.",
			Examples: []Example{
				{
					Title:  "Formatos e tamanhos",
					Code:   `@components.Countdown(components.CountdownProps{Seconds: 90, Format: "mm:ss", Size: "lg"})`,
					Render: demoCountdown(),
				},
				{
					Title:  "Reagindo ao fim",
					Code:   `<div x-data="{ done: false }" x-on:countdown-end="done = true"> ... </div>`,
					Render: demoCountdownOnEnd(),
				},
			},
		},
		{
			Slug:        "breadcrumbs",
			Name:        "Breadcrumbs",
			Category:    CategoryNavigation,
			Description: "Trilha de navegação. O último item nunca é link e recebe aria-current=page.",
			Examples: []Example{
				{
					Title: "Trilha",
					Code: `@components.Breadcrumbs(components.BreadcrumbsProps{Items: []components.BreadcrumbItem{
  {Label: "Início", Href: "/"},
  {Label: "Breadcrumbs"},
}})`,
					Render: demoBreadcrumbs(),
				},
			},
		},
		{
			Slug:        "steps",
			Name:        "Steps",
			Category:    CategoryNavigation,
			Description: "Indicador de progresso por etapas. O passo atual recebe aria-current=step, já que o estado é mostrado por cor.",
			Examples: []Example{
				{
					Title:  "Horizontal",
					Code:   `@components.Steps(components.StepsProps{Items: []components.Step{{Label: "Conta", State: "done"}, {Label: "Pagamento", State: "current"}}})`,
					Render: demoSteps(),
				},
				{
					Title:  "Vertical",
					Code:   `@components.Steps(components.StepsProps{Vertical: true, Items: ...})`,
					Render: demoStepsVertical(),
				},
			},
		},
		{
			Slug:        "indicator",
			Name:        "Indicator",
			Category:    CategoryLayout,
			Description: "Fixa um badge no canto de qualquer elemento — contadores, pontos de status, selos de novidade.",
			Examples: []Example{
				{
					Title:  "Contador e selo",
					Code:   `@components.Indicator(components.IndicatorProps{Label: "9"}) { ...botão... }`,
					Render: demoIndicator(),
				},
			},
		},
		{
			Slug:        "join",
			Name:        "Join",
			Category:    CategoryLayout,
			Description: "Agrupa botões ou inputs num único bloco visual, quadrando os cantos internos e sobrepondo as bordas adjacentes.",
			Examples: []Example{
				{
					Title: "Horizontal e vertical",
					Code: `@components.Join(components.JoinProps{}) {
  @components.Button(components.ButtonProps{Variant: "secondary"}) { Dia }
  @components.Button(components.ButtonProps{Variant: "secondary"}) { Semana }
}`,
					Render: demoJoin(),
				},
			},
		},
		{
			Slug:        "skeleton",
			Name:        "Skeleton",
			Category:    CategoryFeedback,
			Description: "Placeholder animado exibido enquanto os dados carregam. É aria-hidden: representa conteúdo que ainda não existe.",
			Examples: []Example{
				{
					Title:  "Card carregando",
					Code:   `@components.Skeleton(components.SkeletonProps{Class: "h-4 w-2/3"})`,
					Render: demoSkeleton(),
				},
			},
		},
		{
			Slug:        "loading",
			Name:        "Loading",
			Category:    CategoryFeedback,
			Description: "Indicador de carregamento em spinner, pontos ou barras. Anunciado com role=status e um rótulo textual.",
			Examples: []Example{
				{
					Title:  "Tipos",
					Code:   `@components.Loading(components.LoadingProps{Kind: "dots"})`,
					Render: demoLoadingKinds(),
				},
				{
					Title:  "Tamanhos",
					Code:   `@components.Loading(components.LoadingProps{Size: "lg", Variant: "muted"})`,
					Render: demoLoadingSizes(),
				},
			},
		},
		{
			Slug:        "progress",
			Name:        "Progress",
			Category:    CategoryFeedback,
			Description: "Barra de progresso linear sobre o elemento nativo <progress>, que já traz a acessibilidade pronta.",
			Examples: []Example{
				{
					Title:  "Com valor",
					Code:   `@components.Progress(components.ProgressProps{Value: 60, Variant: "success"})`,
					Render: demoProgress(),
				},
				{
					Title:  "Indeterminado",
					Code:   `@components.Progress(components.ProgressProps{Indeterminate: true})`,
					Render: demoProgressIndeterminate(),
				},
			},
		},
		{
			Slug:        "alert",
			Name:        "Alert",
			Category:    CategoryFeedback,
			Description: "Banner de mensagem informativa, de sucesso, aviso ou erro. Sempre anunciado com role=alert.",
			Examples: []Example{
				{
					Title:  "Variantes",
					Code:   `@components.Alert(components.AlertProps{Variant: "warning"}) { <span>Seu token expira em 3 dias.</span> }`,
					Render: demoAlertVariants(),
				},
			},
		},
		{
			Slug:        "toggle",
			Name:        "Toggle",
			Category:    CategoryDataInput,
			Description: "Interruptor deslizante. O visual é Alpine; por baixo vai um checkbox real com sr-only, então envia em formulário normalmente.",
			Examples: []Example{
				{
					Title:  "Estados",
					Code:   `@components.Toggle(components.ToggleProps{ID: "notif", Name: "notif", Label: "Notificações"})`,
					Render: demoToggle(),
				},
				{
					Title:  "Dentro de um formulário",
					Code:   `<form method="get">...</form> — envie e veja os valores na query string`,
					Render: demoFormSubmit(),
				},
			},
		},
		{
			Slug:        "checkbox",
			Name:        "Checkbox",
			Category:    CategoryDataInput,
			Description: "Caixa de seleção com terceiro estado. Indeterminate só existe como propriedade JS — não há atributo HTML — então é aplicado na inicialização e anunciado com aria-checked=mixed.",
			Examples: []Example{
				{
					Title:  "Estados",
					Code:   `@components.Checkbox(components.CheckboxProps{ID: "termos", Name: "termos", Label: "Aceito os termos"})`,
					Render: demoCheckbox(),
				},
				{
					Title:  "Indeterminado",
					Code:   `@components.Checkbox(components.CheckboxProps{ID: "todos", Name: "todos", Label: "Selecionar todos", Indeterminate: true})`,
					Render: demoCheckboxIndeterminate(),
				},
			},
		},
		{
			Slug:        "radio",
			Name:        "Radio",
			Category:    CategoryDataInput,
			Description: "Seleção única. Navegue com as setas: elas movem foco e seleção juntos, como num radiogroup nativo, e pulam opções desabilitadas.",
			Examples: []Example{
				{
					Title: "Horizontal",
					Code: `@components.Radio(components.RadioProps{
  Name: "plano", Value: "pro",
  Options: []components.RadioOption{
    {Label: "Grátis", Value: "free"},
    {Label: "Pro", Value: "pro"},
  },
})`,
					Render: demoRadio(),
				},
				{
					Title:  "Vertical",
					Code:   `@components.Radio(components.RadioProps{Vertical: true, ...})`,
					Render: demoRadioVertical(),
				},
			},
		},
		{
			Slug:        "rating",
			Name:        "Rating",
			Category:    CategoryDataInput,
			Description: "Avaliação por estrelas, navegável pelas setas. Em modo somente leitura vira saída: sem foco, sem envio, com um rótulo textual único.",
			Examples: []Example{
				{
					Title:  "Interativo",
					Code:   `@components.Rating(components.RatingProps{Name: "nota", Value: 3})`,
					Render: demoRating(),
				},
				{
					Title:  "Somente leitura",
					Code:   `@components.Rating(components.RatingProps{Value: 4, ReadOnly: true})`,
					Render: demoRatingReadOnly(),
				},
			},
		},
		{
			Slug:        "range",
			Name:        "Range",
			Category:    CategoryDataInput,
			Description: "Slider numérico. É o único controle da fatia que mantém o input nativo visível — reimplementar arrasto, toque e captura de ponteiro custaria muito para um resultado pior. Alpine só espelha o valor no rótulo.",
			Examples: []Example{
				{
					Title:  "Com valor visível",
					Code:   `@components.Range(components.RangeProps{ID: "vol", Name: "vol", Label: "Volume", Value: 40, ShowValue: true})`,
					Render: demoRange(),
				},
			},
		},
		{
			Slug:        "swap",
			Name:        "Swap",
			Category:    CategoryActions,
			Description: "Alterna entre dois visuais. É um botão de verdade com aria-pressed — um toggle cujo estado a tecnologia assistiva não enxerga é só um desenho.",
			Examples: []Example{
				{
					Title: "Play/pause e favoritar",
					Code: `@components.Swap(components.SwapProps{
  Label: "Reproduzir ou pausar",
  On:    pauseIcon(), Off: playIcon(),
})`,
					Render: demoSwap(),
				},
			},
		},
		{
			Slug:        "tooltip",
			Name:        "Tooltip",
			Category:    CategoryLayout,
			Description: "Dica contextual. Abre no hover e também no foco, fecha com Escape, e entrega o texto por aria-describedby — o estilo sozinho não entrega para ninguém.",
			Examples: []Example{
				{
					Title:  "Posições",
					Code:   `@components.Tooltip(components.TooltipProps{ID: "tip", Text: "Aparece acima"}) { ...gatilho... }`,
					Render: demoTooltip(),
				},
				{
					Title:  "Pelo teclado",
					Code:   `Tab até o gatilho abre a dica; Escape fecha.`,
					Render: demoTooltipKeyboard(),
				},
			},
		},
		{
			Slug:        "toast",
			Name:        "Toast",
			Category:    CategoryFeedback,
			Description: "Notificação flutuante temporária, fixada no canto inferior direito. Some sozinha após 3 segundos.",
			Examples: []Example{
				{
					Title:  "Sucesso",
					Code:   `@components.Toast(components.ToastProps{Message: "Alterações salvas", Variant: "success"})`,
					Render: demoToast(),
				},
				{
					Title:  "Erro",
					Code:   `@components.Toast(components.ToastProps{Message: "Não foi possível salvar", Variant: "error"})`,
					Render: demoToastError(),
				},
			},
		},
		{
			Slug:        "empty",
			Name:        "Empty",
			Category:    CategoryFeedback,
			Description: "Estado vazio com ícone, título, descrição e ação. O título é um parágrafo: o nível de heading depende da página, que o componente não conhece.",
			Examples: []Example{
				{
					Title: "Com ação",
					Code: `@components.Empty(components.EmptyProps{Title: "Nenhum projeto ainda", Description: "..."}) {
  @components.Button(components.ButtonProps{}) { Novo projeto }
}`,
					Render: demoEmpty(),
				},
			},
		},
	}
}
