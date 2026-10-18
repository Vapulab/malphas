package showcase

// catalog is the source of truth for the docs. Adding a component here is all
// it takes to get an index entry and a /components/{slug} page — no new route.
func catalog() []Component {
	return []Component{
		{
			Slug:        "button",
			Name:        "Button",
			Category:    CategoryActions,
			Description: "Clickable button with color variants, sizes and support for arbitrary attributes (HTMX, Alpine).",
			Examples: []Example{
				{
					Title:  "Variants",
					Code:   `@components.Button(components.ButtonProps{Variant: "danger"}) { Danger }`,
					Render: demoButtonVariants(),
				},
				{
					Title:  "Sizes",
					Code:   `@components.Button(components.ButtonProps{Size: "lg"}) { Large }`,
					Render: demoButtonSizes(),
				},
				{
					Title:  "Disabled",
					Code:   `@components.Button(components.ButtonProps{Attrs: templ.Attributes{"disabled": true}}) { Unavailable }`,
					Render: demoButtonDisabled(),
				},
			},
		},
		{
			Slug:        "modal",
			Name:        "Modal",
			Category:    CategoryActions,
			Description: "Overlay dialog that is truly modal: focus is trapped inside, the page behind neither scrolls nor is reachable by a screen reader, and on close focus returns to whatever opened it.",
			Examples: []Example{
				{
					Title: "Open and close",
					Code: `<div x-data>
  @components.Button(components.ButtonProps{
    Attrs: templ.Attributes{"x-on:click": "$dispatch('open-modal', 'demo-modal')"},
  }) { Open modal }
</div>
@components.Modal(components.ModalProps{ID: "demo-modal", Title: "Confirm action"}) {
  <p>Modal content.</p>
}`,
					Render: demoModal(),
				},
			},
		},
		{
			Slug:        "alert-dialog",
			Name:        "Alert Dialog",
			Category:    CategoryActions,
			Description: "Modal that asks for confirmation of an action with consequences. role=alertdialog with the description announced, initial focus on Cancel, and clicking outside does not count as an answer. ConfirmAttrs carries hx-delete and the like.",
			Examples: []Example{
				{
					Title: "Delete with confirmation",
					Code: `@components.AlertDialog(components.AlertDialogProps{
  ID: "delete", Title: "Delete the protocol?", Description: "This cannot be undone.",
  ConfirmLabel: "Delete", ConfirmAttrs: templ.Attributes{"hx-delete": "/protocols/1"},
})`,
					Render: demoAlertDialog(),
				},
			},
		},
		{
			Slug:        "card",
			Name:        "Card",
			Category:    CategoryDataDisplay,
			Description: "Bordered, shadowed container that groups a block of content.",
			Examples: []Example{
				{
					Title:  "Simple",
					Code:   `@components.Card(components.CardProps{}) { <p>Content</p> }`,
					Render: demoCard(),
				},
				{
					Title:  "With an action",
					Code:   `@components.Card(components.CardProps{Class: "max-w-sm"}) { ... }`,
					Render: demoCardWithActions(),
				},
			},
		},
		{
			Slug:        "item",
			Name:        "Item",
			Category:    CategoryDataDisplay,
			Description: "Row with media, title, description and actions. With Href the whole row is clickable, while the actions stay outside the link (the stretched link pattern). ItemGroup builds the list.",
			Examples: []Example{
				{
					Title:  "Link with action, and a group",
					Code:   `@components.Item(components.ItemProps{Title: "Report", Href: "/reports/1", Actions: download()})`,
					Render: demoItem(),
				},
			},
		},
		{
			Slug:        "hover-card",
			Name:        "Hover Card",
			Category:    CategoryDataDisplay,
			Description: "Preview of what sits behind a link, on mouse hover or keyboard focus, with a delay so it does not flicker. It is only a complement: the trigger is a real link, and nothing may exist only in the card.",
			Examples: []Example{
				{
					Title: "Reviewer profile",
					Code: `@components.HoverCard(components.HoverCardProps{ID: "author", Label: "@lincoln", Href: "/u/lincoln"}) {
  <p>Researcher · 42 protocols</p>
}`,
					Render: demoHoverCard(),
				},
			},
		},
		{
			Slug:        "bubble",
			Name:        "Bubble",
			Category:    CategoryDataDisplay,
			Description: "Chat bubble with variants, aligned to the start or the end. Consecutive messages from the same author flatten the corners where they meet; reactions say what they mean and how many there are, not just the emoji.",
			Examples: []Example{
				{
					Title:  "Group, reaction, own message and quote",
					Code:   `@components.Bubble(components.BubbleProps{Variant: "primary", Align: "end"}) { Rerun it, and set aside a control aliquot. }`,
					Render: demoBubble(),
				},
			},
		},
		{
			Slug:        "marker",
			Name:        "Marker",
			Category:    CategoryDataDisplay,
			Description: "A conversation line nobody said: someone joined, a reply being generated, the date between days. Updates use role=status; the divider keeps its label readable for screen readers.",
			Examples: []Example{
				{
					Title:  "Divider, note and status",
					Code:   `@components.Marker(components.MarkerProps{Variant: "border", Status: true}) { Generating the summary… }`,
					Render: demoMarker(),
				},
			},
		},
		{
			Slug:        "message",
			Name:        "Message",
			Category:    CategoryDataDisplay,
			Description: "Message with avatar, name, time and a Bubble. In a run from the same author only the first one shows the header, but every message keeps author and time in its accessible name.",
			Examples: []Example{
				{
					Title: "Short conversation",
					Code: `@components.Message(components.MessageProps{
  Author: "Amara Okafor", Initials: "AO", Time: "14:02", Datetime: "2026-10-10T14:02", Group: "first",
}) { The curve for cycle 32 came out above the threshold. }`,
					Render: demoMessage(),
				},
			},
		},
		{
			Slug:        "attachment",
			Name:        "Attachment",
			Category:    CategoryDataDisplay,
			Description: "File or image in a composer, a message or an upload list, with its state in text (uploading, processing, error, done) as well as the bar. The name only becomes a link once the file is ready; the actions stay reachable, especially on error. AttachmentGroup scrolls horizontally.",
			Examples: []Example{
				{
					Title: "The four states",
					Code: `@components.AttachmentGroup(components.AttachmentGroupProps{Label: "Attachments", Items: []components.AttachmentProps{
  {Name: "protocol.pdf", Meta: "PDF · 2.4 MB", State: "done", Href: "/files/protocol.pdf"},
  {Name: "curve.png", State: "uploading", Progress: 45},
}})`,
					Render: demoAttachment(),
				},
			},
		},
		{
			Slug:        "message-scroller",
			Name:        "Message Scroller",
			Category:    CategoryDataDisplay,
			Description: "Conversation transcript that behaves like a chat: it follows new messages only if you are at the bottom, shows \"Jump to the end\" (with how many arrived) when you have scrolled up, and keeps your position when older messages are added at the top. It does not care how messages arrive: it watches its own content. role=log, polite.",
			Examples: []Example{
				{
					Title: "New via HTMX, older at the top, and streaming",
					Code: `@components.MessageScroller(components.MessageScrollerProps{ID: "chat", Label: "Conversation"}) {
  for _, m := range messages { @components.Message(...) { { m.Text } } }
}
// new ones: hx-get="/messages" hx-target="#chat-content" hx-swap="beforeend"`,
					Render: demoMessageScroller(),
				},
			},
		},
		{
			Slug:        "typography",
			Name:        "Typography",
			Category:    CategoryDataDisplay,
			Description: "Styling for long-form text whose HTML comes from elsewhere (Markdown, help pages, terms): headings, lists, quotes and code, using the Malphas tokens. Also available as the prose-malphas class.",
			Examples: []Example{
				{
					Title:  "Mixed content",
					Code:   `@components.Prose(components.ProseProps{}) { @templ.Raw(renderedMarkdown) }`,
					Render: demoProse(),
				},
			},
		},
		{
			Slug:        "input",
			Name:        "Text input",
			Category:    CategoryDataInput,
			Description: "Text field with a label associated by id. Accepts any native type.",
			Examples: []Example{
				{
					Title: "With label and placeholder",
					Code: `@components.Input(components.InputProps{
  ID: "email", Name: "email", Label: "Email",
  Type: "email", Placeholder: "you@example.com",
})`,
					Render: demoInput(),
				},
				{
					Title:  "With a value",
					Code:   `@components.Input(components.InputProps{ID: "name", Name: "name", Label: "Name", Value: "Lincoln"})`,
					Render: demoInputWithValue(),
				},
			},
		},
		{
			Slug:        "link",
			Name:        "Link",
			Category:    CategoryActions,
			Description: "Text link with color variants. Href goes through templ's URL sanitizer, and external links get rel=noopener.",
			Examples: []Example{
				{
					Title: "Variants",
					Code: `@components.Link(components.LinkProps{Href: "/all"}) { Internal link }
@components.Link(components.LinkProps{Href: "https://templ.guide", External: true}) { External }`,
					Render: demoLink(),
				},
			},
		},
		{
			Slug:        "badge",
			Name:        "Badge",
			Category:    CategoryDataDisplay,
			Description: "Compact status or count label, solid or outlined.",
			Examples: []Example{
				{
					Title:  "Solid",
					Code:   `@components.Badge(components.BadgeProps{Variant: "success"}) { Active }`,
					Render: demoBadgeSolid(),
				},
				{
					Title:  "Outline",
					Code:   `@components.Badge(components.BadgeProps{Variant: "primary", Outline: true}) { Primary }`,
					Render: demoBadgeOutline(),
				},
			},
		},
		{
			Slug:        "kbd",
			Name:        "Kbd",
			Category:    CategoryDataDisplay,
			Description: "Visual representation of a key, using the semantic <kbd> element.",
			Examples: []Example{
				{
					Title:  "Shortcut",
					Code:   `@components.Kbd(components.KbdProps{}) { Ctrl }`,
					Render: demoKbd(),
				},
			},
		},
		{
			Slug:        "divider",
			Name:        "Divider",
			Category:    CategoryLayout,
			Description: "Horizontal or vertical rule, with optional text in the middle. Without a label it is decorative and hidden from screen readers.",
			Examples: []Example{
				{
					Title: "Horizontal, labeled and vertical",
					Code: `@components.Divider(components.DividerProps{})
@components.Divider(components.DividerProps{Label: "or"})
@components.Divider(components.DividerProps{Vertical: true})`,
					Render: demoDivider(),
				},
			},
		},
		{
			Slug:        "aspect-ratio",
			Name:        "Aspect Ratio",
			Category:    CategoryLayout,
			Description: "Reserves a media element's aspect ratio before it loads, so the page does not jump. An unknown ratio falls back to 16/9 instead of collapsing the box.",
			Examples: []Example{
				{
					Title:  "16/9 and 1/1",
					Code:   `@components.AspectRatio(components.AspectRatioProps{Ratio: "16/9"}) { <img src="..." alt="..."> }`,
					Render: demoAspectRatio(),
				},
			},
		},
		{
			Slug:        "scroll-area",
			Name:        "Scroll Area",
			Category:    CategoryLayout,
			Description: "Native scrolling with a themed scrollbar. It takes focus (tabindex=0) so it can be scrolled from the keyboard, and therefore requires a Label.",
			Examples: []Example{
				{
					Title:  "Vertical and horizontal",
					Code:   `@components.ScrollArea(components.ScrollAreaProps{Label: "Terms of use", Class: "max-h-48"}) { ... }`,
					Render: demoScrollArea(),
				},
			},
		},
		{
			Slug:        "resizable",
			Name:        "Resizable",
			Category:    CategoryLayout,
			Description: "Two panels with a divider you drag or move from the keyboard (arrows, Shift for 1%, Home/End), following the APG window splitter pattern. The server renders the initial size; with an ID, the last size is remembered in this browser.",
			Examples: []Example{
				{
					Title: "Side by side and stacked",
					Code: `@components.Resizable(components.ResizableProps{
  ID: "editor", Label: "List width", DefaultSize: 35, Min: 20, Max: 70,
  First: list(), Second: editor(),
})`,
					Render: demoResizable(),
				},
			},
		},
		{
			Slug:        "footer",
			Name:        "Footer",
			Category:    CategoryLayout,
			Description: "Structural footer with link columns and a note line. It announces itself as a navigation landmark even when nested, and hands its links to Link to inherit href sanitizing.",
			Examples: []Example{
				{
					Title: "Columns and note",
					Code: `@components.Footer(components.FooterProps{
  Note: "© 2026 Vapulab",
  Columns: []components.FooterColumn{
    {Title: "Product", Links: []components.FooterLink{
      {Label: "Components", Href: "/components"},
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
			Description: "Round or square profile picture, falling back to initials, with an optional presence indicator.",
			Examples: []Example{
				{
					Title:  "Sizes and shapes",
					Code:   `@components.Avatar(components.AvatarProps{Initials: "LS", Size: "lg"})`,
					Render: demoAvatar(),
				},
				{
					Title:  "With status",
					Code:   `@components.Avatar(components.AvatarProps{Initials: "LS", Status: "online"})`,
					Render: demoAvatarStatus(),
				},
			},
		},
		{
			Slug:        "stat",
			Name:        "Stat",
			Category:    CategoryDataDisplay,
			Description: "Block that highlights a number with a title and context. It uses dt/dd, so it must sit inside a <dl>.",
			Examples: []Example{
				{
					Title:  "Metrics grid",
					Code:   `@components.Stat(components.StatProps{Title: "Revenue", Value: "$12,400", Description: "+12% this month"})`,
					Render: demoStat(),
				},
			},
		},
		{
			Slug:        "table",
			Name:        "Table",
			Category:    CategoryDataDisplay,
			Description: "Table for dense data, with its own horizontal scroll. Zebra striping and density are set once on the table and reach the cells through a descendant selector.",
			Examples: []Example{
				{
					Title: "Default",
					Code: `@components.Table(components.TableProps{}) {
  @components.THead() {
    @components.TR() {
      @components.TH(components.THProps{}) { Service }
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
					Title:  "Zebra and compact",
					Code:   `@components.Table(components.TableProps{Zebra: true, Compact: true}) { ... }`,
					Render: demoTableZebraCompact(),
				},
			},
		},
		{
			Slug:        "data-table",
			Name:        "Data Table",
			Category:    CategoryDataDisplay,
			Description: "Table with sorting, search and pagination done on the server. Everything is a real link or GET field (it works without JS); with JS, each action swaps only the table region and the URL keeps the state. The search box sits outside the swapped region, so typing is never interrupted.",
			Examples: []Example{
				{
					Title: "Protocols",
					Code: `@components.DataTable(components.DataTableProps{
  ID: "protocols", URL: "/protocols", Columns: columns,
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
			Description: "Highlight section for the top of a page. It handles the framing — spacing, centering and max width — and leaves the content to you.",
			Examples: []Example{
				{
					Title:  "Main call to action",
					Code:   `@components.Hero(components.HeroProps{Class: "bg-primary/5"}) { <h1>...</h1> }`,
					Render: demoHero(),
				},
			},
		},
		{
			Slug:        "drawer",
			Name:        "Drawer",
			Category:    CategoryActions,
			Description: "Modal panel that slides in from the side or from the bottom (Side: \"bottom\"). Same primitive as Modal — trapped focus, locked scroll, inert content behind — because a panel that lets Tab escape is worse than a dialog doing the same: the page looks like it is still usable.",
			Examples: []Example{
				{
					Title: "Left and right",
					Code: `@components.Drawer(components.DrawerProps{ID: "nav", Title: "Navigation"}) {
  @components.Menu(components.MenuProps{Entries: entries})
}
// trigger, inside an Alpine scope:
// x-on:click="$dispatch('open-drawer', 'nav')"`,
					Render: demoDrawer(),
				},
				{
					Title:  "From the bottom",
					Code:   `@components.Drawer(components.DrawerProps{ID: "share", Title: "Share", Side: "bottom"}) { ... }`,
					Render: demoDrawerBottom(),
				},
			},
		},
		{
			Slug:        "navbar",
			Name:        "Navbar",
			Category:    CategoryNavigation,
			Description: "Top navigation bar. It knows how to open a Drawer on mobile but does not require one: without a DrawerID, no menu button is rendered.",
			Examples: []Example{
				{
					Title: "With a mobile menu",
					Code: `@components.Navbar(components.NavbarProps{
  Brand: "Vapulab", DrawerID: "menu",
  Entries: []components.MenuEntry{
    {Label: "Components", Href: "/", Current: true},
  },
}) { ...right-hand content... }`,
					Render: demoNavbar(),
				},
			},
		},
		{
			Slug:        "menubar",
			Name:        "Menubar",
			Category:    CategoryNavigation,
			Description: "Application menu bar (File, Edit, View). One Tab stop for the whole bar; ← → move between menus and carry the open menu along, ↑ ↓ move inside it, Escape returns to the top level and Tab closes everything.",
			Examples: []Example{
				{
					Title: "Editor",
					Code: `@components.Menubar(components.MenubarProps{
  Label: "Editor",
  Menus: []components.MenubarMenu{{Label: "File", Entries: []components.MenuEntry{...}}, ...},
})`,
					Render: demoMenubar(),
				},
			},
		},
		{
			Slug:        "sidebar",
			Name:        "Sidebar",
			Category:    CategoryNavigation,
			Description: "Application frame: grouped side navigation with the content next to it. It collapses to icons only (the names remain for screen readers and in the title) and remembers the choice; on a phone it becomes a Menu button that opens the same navigation in a modal Drawer.",
			Examples: []Example{
				{
					Title: "App with two groups",
					Code: `@components.Sidebar(components.SidebarProps{
  ID: "app", Label: "Navigation", Collapsible: true,
  Groups: []components.SidebarGroup{{Label: "Work", Items: []components.SidebarItem{{Label: "Protocols", Href: "/protocols", Icon: "◆", Current: true}}}},
}) { page content }`,
					Render: demoSidebar(),
				},
			},
		},
		{
			Slug:        "command",
			Name:        "Command",
			Category:    CategoryNavigation,
			Description: "Command palette: a modal that searches grouped actions and pages, opened with Ctrl+K / ⌘K. Focus stays in the search box and the arrows move through the results (aria-activedescendant); the search starts empty every time it opens.",
			Examples: []Example{
				{
					Title: "Actions and pages",
					Code: `@components.Command(components.CommandProps{
  ID: "palette", Hotkey: true,
  Groups: []components.CommandGroup{{Heading: "Actions", Items: []components.CommandItem{{Label: "Toggle theme", OnClick: "..."}}}},
})`,
					Render: demoCommand(),
				},
			},
		},
		{
			Slug:        "dropdown",
			Name:        "Dropdown",
			Category:    CategoryActions,
			Description: "Menu button pattern, with role=menu and arrow-key navigation. It is the opposite of the Drawer: it does not trap focus — Escape returns to the trigger, but Tab leaves and closes it, as a menu should.",
			Examples: []Example{
				{
					Title: "Actions and account",
					Code: `@components.Dropdown(components.DropdownProps{
  ID: "actions", Label: "Actions",
  Entries: []components.MenuEntry{
    {Label: "Edit", OnClick: "$dispatch('edit')"},
    components.MenuDivider(),
    {Label: "Delete", OnClick: "$dispatch('del')", Disabled: true},
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
			Description: "Button that opens a non-modal panel with free-form content. Focus moves into the panel, Escape closes it and returns focus, and tabbing out closes it without trapping anything. Inside a Modal, Escape closes only the popover.",
			Examples: []Example{
				{
					Title: "Form in the panel, and inside a modal",
					Code: `@components.Popover(components.PopoverProps{ID: "filters", Label: "Filters"}) {
  @components.Input(components.InputProps{ID: "author", Name: "author", Label: "Author"})
}`,
					Render: demoPopover(),
				},
			},
		},
		{
			Slug:        "context-menu",
			Name:        "Context Menu",
			Category:    CategoryActions,
			Description: "Right-click action menu, opened at the click point and pushed back inside the viewport near the edges. The area takes focus and opens the same menu with Shift+F10 or the Menu key.",
			Examples: []Example{
				{
					Title: "File actions",
					Code: `@components.ContextMenu(components.ContextMenuProps{
  ID: "file", Label: "File actions",
  Entries: []components.MenuEntry{{Label: "Rename", OnClick: "..."}, ...},
}) { <div>report.pdf</div> }`,
					Render: demoContextMenu(),
				},
			},
		},
		{
			Slug:        "select",
			Name:        "Select",
			Category:    CategoryDataInput,
			Description: "Single choice on top of the native <select> — the second deliberate exception to Alpine, after Range. It gets the native picker on mobile, type-ahead and full keyboard support; in exchange, the open options cannot be styled.",
			Examples: []Example{
				{
					Title:  "With a value and with a placeholder",
					Code:   `@components.Select(components.SelectProps{ID: "plan", Name: "plan", Label: "Plan", Options: options, Value: "pro"})`,
					Render: demoSelect(),
				},
			},
		},
		{
			Slug:        "combobox",
			Name:        "Combobox",
			Category:    CategoryDataInput,
			Description: "Searchable select, for lists too long to scroll. Focus stays in the text field and the highlighted option is announced through aria-activedescendant; a hidden input carries the value, so free text is never submitted.",
			Examples: []Example{
				{
					Title: "Long list, and with a placeholder",
					Code: `@components.Combobox(components.ComboboxProps{
  ID: "city", Name: "city", Label: "City", Value: "ymq",
  Placeholder: "Search for a city…", Options: cities,
})`,
					Render: demoCombobox(),
				},
			},
		},
		{
			Slug:        "date-picker",
			Name:        "Date Picker",
			Category:    CategoryDataInput,
			Description: "The native <input type=\"date\">, styled — the same choice as Select: the platform picker (the wheel on a phone, the grid on desktop) is better at touch, keyboard and screen readers, and the value always arrives as YYYY-MM-DD. For a calendar that stays visible on the page, use Calendar.",
			Examples: []Example{
				{
					Title:  "With limits, inside a Field",
					Code:   `@components.DatePicker(components.DatePickerProps{ID: "delivery", Name: "delivery", Min: "2026-10-10", Max: "2026-12-31"})`,
					Render: demoDatePicker(),
				},
			},
		},
		{
			Slug:        "calendar",
			Name:        "Calendar",
			Category:    CategoryDataInput,
			Description: "A month grid for picking a day or a range, always visible on the page. The server builds the month; changing months fetches the neighbor and swaps only the calendar. APG keyboard pattern; month and weekday names, and the first day of the week, follow the locale.",
			Examples: []Example{
				{
					Title: "Single day with a minimum, and a range",
					Code: `@components.Calendar(components.CalendarProps{
  ID: "delivery", Month: month, Name: "delivery", Selected: chosen,
  Min: today, NavURL: "/schedule",
})`,
					Render: demoCalendar(),
				},
			},
		},
		{
			Slug:        "questionnaire",
			Name:        "Questionnaire",
			Category:    CategoryDataInput,
			Description: "Multi-step form, one question at a time (single choice, multiple choice, free text, skippable). Each step is a POST and the server decides what comes next: validation, conditional questions and resuming all live in Go. Keys 1–9 choose, Enter advances; Back never validates.",
			Examples: []Example{
				{
					Title: "With a conditional question and resuming",
					Code: `@components.Questionnaire(components.QuestionnaireProps{
  ID: "profile", Action: "/profile/answer",
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
			Description: "Field label with an optional required marker. The asterisk is aria-hidden: the field's own required attribute is what announces it as required.",
			Examples: []Example{
				{
					Title:  "Plain and required",
					Code:   `@components.Label(components.LabelProps{For: "email", Required: true}) { Email }`,
					Render: demoLabel(),
				},
			},
		},
		{
			Slug:        "field",
			Name:        "Field",
			Category:    CategoryDataInput,
			Description: "Brings together label, control, description and a server-side error. ControlAttrs() returns the aria-describedby and aria-invalid the control needs to carry.",
			Examples: []Example{
				{
					Title: "With a description and with an error",
					Code: `{{ f := components.FieldProps{ID: "site", Label: "Website", Error: "Invalid address."} }}
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
			Description: "Multi-line text field. With MaxLength it gets a counter that counts the way the browser does (an emoji counts as 2) and is already filled in by the server.",
			Examples: []Example{
				{
					Title:  "With a counter, and disabled",
					Code:   `@components.Textarea(components.TextareaProps{ID: "bio", Name: "bio", Label: "Bio", MaxLength: 160})`,
					Render: demoTextarea(),
				},
			},
		},
		{
			Slug:        "input-group",
			Name:        "Input Group",
			Category:    CategoryDataInput,
			Description: "Field with an attached prefix and suffix — currency, domain, a button. The focus ring wraps the whole group.",
			Examples: []Example{
				{
					Title: "Prefix, suffix and button",
					Code: `@components.InputGroup(components.InputGroupProps{
  ID: "price", Name: "price", Label: "Price",
  Prefix: components.InputGroupText("$"),
})`,
					Render: demoInputGroup(),
				},
			},
		},
		{
			Slug:        "toggle-button",
			Name:        "ToggleButton",
			Category:    CategoryDataInput,
			Description: "Two-state button with aria-pressed, like bold in an editor. Not the same as Toggle, which is a settings switch. With Name, it submits its value only while pressed.",
			Examples: []Example{
				{
					Title: "Icon, pressed and disabled",
					Code: `@components.ToggleButton(components.ToggleButtonProps{Name: "favorite", Label: "Favorite"}) {
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
			Description: "A set of ToggleButtons with a single Tab stop. Single presses at most one (and clicking the pressed one releases it); multiple toggles each one. The arrows move focus only.",
			Examples: []Example{
				{
					Title: "Single and multiple",
					Code: `@components.ToggleGroup(components.ToggleGroupProps{
  Name: "alignment", Label: "Alignment",
  Items: []components.ToggleGroupItem{{Label: "Left", Value: "left"}, ...},
  Value: []string{"center"},
})`,
					Render: demoToggleGroup(),
				},
			},
		},
		{
			Slug:        "otp",
			Name:        "OTP",
			Category:    CategoryDataInput,
			Description: "Verification code. The first field carries autocomplete=one-time-code, which is what makes iOS and Android offer the SMS code right on the keyboard.",
			Examples: []Example{
				{
					Title:  "6 and 4 digits",
					Code:   `@components.OTP(components.OTPProps{Name: "code", Label: "Verification code"})`,
					Render: demoOTP(),
				},
			},
		},
		{
			Slug:        "carousel",
			Name:        "Carousel",
			Category:    CategoryDataDisplay,
			Description: "Horizontal slides with scroll-snap. Fully navigable from the keyboard, with the position announced on each slide, and scrolling respects prefers-reduced-motion.",
			Examples: []Example{
				{
					Title:  "Three slides",
					Code:   `@components.Carousel(components.CarouselProps{ID: "shots", Label: "Screenshots", Slides: slides})`,
					Render: demoCarousel(),
				},
			},
		},
		{
			Slug:        "pagination",
			Name:        "Pagination",
			Category:    CategoryNavigation,
			Description: "Navigation between pages, collapsing long ranges. The current page is a span with aria-current, not a link — linking to where you already are is a dead control.",
			Examples: []Example{
				{
					Title:  "Interactive",
					Code:   `@components.Pagination(components.PaginationProps{Current: page, Total: 20, BaseHref: "/components/pagination?page="})`,
					Render: demoPagination(),
				},
				{
					Title:  "At the edges",
					Code:   `// Each example has its own parameter, and BaseHref carries the others along — without that, moving one resets the rest.`,
					Render: demoPaginationEdges(),
				},
			},
		},
		{
			Slug:        "menu",
			Name:        "Menu",
			Category:    CategoryNavigation,
			Description: "Vertical navigation. It deliberately does not use role=menu — that role is for application menus, and applying it to site links makes screen readers announce the page as a program menu. For the menu button pattern, use Dropdown.",
			Examples: []Example{
				{
					Title: "With a divider and an action",
					Code: `@components.Menu(components.MenuProps{Entries: []components.MenuEntry{
  {Label: "Home", Href: "/", Current: true},
  components.MenuDivider(),
  {Label: "Sign out", OnClick: "$dispatch('logout')"},
}})`,
					Render: demoMenu(),
				},
			},
		},
		{
			Slug:        "collapse",
			Name:        "Collapse",
			Category:    CategoryDataDisplay,
			Description: "Expandable area. The trigger is a real button — a clickable div does not take focus, does not respond to Enter or Space, and is announced as nothing.",
			Examples: []Example{
				{
					Title:  "Closed and open",
					Code:   `@components.Collapse(components.CollapseProps{ID: "c1", Title: "Details"}) { <p>body</p> }`,
					Render: demoCollapse(),
				},
			},
		},
		{
			Slug:        "accordion",
			Name:        "Accordion",
			Category:    CategoryDataDisplay,
			Description: "Group of expandable panels, with an optional exclusive mode where opening one closes the others.",
			Examples: []Example{
				{
					Title: "Independent",
					Code: `@components.Accordion(components.AccordionProps{
  ID: "faq",
  Items: []components.AccordionItem{
    {Title: "Does it need Node?", Content: answer()},
  },
})`,
					Render: demoAccordion(),
				},
				{
					Title:  "Exclusive",
					Code:   `@components.Accordion(components.AccordionProps{ID: "faq", Exclusive: true, ...})`,
					Render: demoAccordionExclusive(),
				},
			},
		},
		{
			Slug:        "tabs",
			Name:        "Tabs",
			Category:    CategoryNavigation,
			Description: "Tabs following the WAI-ARIA pattern: Tab reaches the strip once, then the arrows navigate, Home and End go to the ends, and disabled tabs are skipped.",
			Examples: []Example{
				{
					Title: "With a disabled tab",
					Code: `@components.Tabs(components.TabsProps{
  ID: "account", Label: "Account settings",
  Items: []components.Tab{
    {Label: "Profile", Content: profile()},
    {Label: "Billing", Content: billing(), Disabled: true},
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
			Description: "Chronological timeline of events, with a rail connecting the markers.",
			Examples: []Example{
				{
					Title: "Events of a deploy",
					Code: `@components.Timeline(components.TimelineProps{Items: []components.TimelineItem{
  {Time: "09:00", Title: "Deploy started", Description: "build #481"},
  {Time: "09:04", Title: "Deploy finished", Variant: "success"},
}})`,
					Render: demoTimeline(),
				},
			},
		},
		{
			Slug:        "countdown",
			Name:        "Countdown",
			Category:    CategoryDataDisplay,
			Description: "Countdown timer. The initial value comes from the server, so the clock is right from the first paint. aria-live is off on purpose: a number announced every second makes the page unusable with a screen reader.",
			Examples: []Example{
				{
					Title:  "Formats and sizes",
					Code:   `@components.Countdown(components.CountdownProps{Seconds: 90, Format: "mm:ss", Size: "lg"})`,
					Render: demoCountdown(),
				},
				{
					Title:  "Reacting to the end",
					Code:   `<div x-data="{ done: false }" x-on:countdown-end="done = true"> ... </div>`,
					Render: demoCountdownOnEnd(),
				},
			},
		},
		{
			Slug:        "breadcrumbs",
			Name:        "Breadcrumbs",
			Category:    CategoryNavigation,
			Description: "Navigation trail. The last item is never a link and gets aria-current=page.",
			Examples: []Example{
				{
					Title: "Trail",
					Code: `@components.Breadcrumbs(components.BreadcrumbsProps{Items: []components.BreadcrumbItem{
  {Label: "Home", Href: "/"},
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
			Description: "Step-by-step progress indicator. The current step gets aria-current=step, since the state is otherwise shown only by color.",
			Examples: []Example{
				{
					Title:  "Horizontal",
					Code:   `@components.Steps(components.StepsProps{Items: []components.Step{{Label: "Account", State: "done"}, {Label: "Payment", State: "current"}}})`,
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
			Description: "Pins a badge to the corner of any element — counters, status dots, \"new\" tags.",
			Examples: []Example{
				{
					Title:  "Counter and tag",
					Code:   `@components.Indicator(components.IndicatorProps{Label: "9"}) { ...button... }`,
					Render: demoIndicator(),
				},
			},
		},
		{
			Slug:        "join",
			Name:        "Join",
			Category:    CategoryLayout,
			Description: "Groups buttons or inputs into a single visual block, squaring the inner corners and overlapping adjacent borders.",
			Examples: []Example{
				{
					Title: "Horizontal and vertical",
					Code: `@components.Join(components.JoinProps{}) {
  @components.Button(components.ButtonProps{Variant: "secondary"}) { Day }
  @components.Button(components.ButtonProps{Variant: "secondary"}) { Week }
}`,
					Render: demoJoin(),
				},
			},
		},
		{
			Slug:        "skeleton",
			Name:        "Skeleton",
			Category:    CategoryFeedback,
			Description: "Animated placeholder shown while data loads. It is aria-hidden: it stands for content that does not exist yet.",
			Examples: []Example{
				{
					Title:  "Loading card",
					Code:   `@components.Skeleton(components.SkeletonProps{Class: "h-4 w-2/3"})`,
					Render: demoSkeleton(),
				},
			},
		},
		{
			Slug:        "loading",
			Name:        "Loading",
			Category:    CategoryFeedback,
			Description: "Loading indicator as a spinner, dots or bars. Announced with role=status and a text label.",
			Examples: []Example{
				{
					Title:  "Kinds",
					Code:   `@components.Loading(components.LoadingProps{Kind: "dots"})`,
					Render: demoLoadingKinds(),
				},
				{
					Title:  "Sizes",
					Code:   `@components.Loading(components.LoadingProps{Size: "lg", Variant: "muted"})`,
					Render: demoLoadingSizes(),
				},
			},
		},
		{
			Slug:        "progress",
			Name:        "Progress",
			Category:    CategoryFeedback,
			Description: "Linear progress bar on top of the native <progress> element, which comes with accessibility built in.",
			Examples: []Example{
				{
					Title:  "With a value",
					Code:   `@components.Progress(components.ProgressProps{Value: 60, Variant: "success"})`,
					Render: demoProgress(),
				},
				{
					Title:  "Indeterminate",
					Code:   `@components.Progress(components.ProgressProps{Indeterminate: true})`,
					Render: demoProgressIndeterminate(),
				},
			},
		},
		{
			Slug:        "alert",
			Name:        "Alert",
			Category:    CategoryFeedback,
			Description: "Message banner for info, success, warning or error. Always announced with role=alert.",
			Examples: []Example{
				{
					Title:  "Variants",
					Code:   `@components.Alert(components.AlertProps{Variant: "warning"}) { <span>Your token expires in 3 days.</span> }`,
					Render: demoAlertVariants(),
				},
			},
		},
		{
			Slug:        "toggle",
			Name:        "Toggle",
			Category:    CategoryDataInput,
			Description: "Sliding switch. The visual is Alpine; underneath sits a real sr-only checkbox, so it submits with a form as usual.",
			Examples: []Example{
				{
					Title:  "States",
					Code:   `@components.Toggle(components.ToggleProps{ID: "notif", Name: "notif", Label: "Notifications"})`,
					Render: demoToggle(),
				},
				{
					Title:  "Inside a form",
					Code:   `<form method="get">...</form> — submit it and see the values in the query string`,
					Render: demoFormSubmit(),
				},
			},
		},
		{
			Slug:        "checkbox",
			Name:        "Checkbox",
			Category:    CategoryDataInput,
			Description: "Checkbox with a third state. Indeterminate only exists as a JS property — there is no HTML attribute — so it is applied on init and announced with aria-checked=mixed.",
			Examples: []Example{
				{
					Title:  "States",
					Code:   `@components.Checkbox(components.CheckboxProps{ID: "terms", Name: "terms", Label: "I accept the terms"})`,
					Render: demoCheckbox(),
				},
				{
					Title:  "Indeterminate",
					Code:   `@components.Checkbox(components.CheckboxProps{ID: "all", Name: "all", Label: "Select all", Indeterminate: true})`,
					Render: demoCheckboxIndeterminate(),
				},
			},
		},
		{
			Slug:        "radio",
			Name:        "Radio",
			Category:    CategoryDataInput,
			Description: "Single choice. Navigate with the arrows: they move focus and selection together, as in a native radio group, and skip disabled options.",
			Examples: []Example{
				{
					Title: "Horizontal",
					Code: `@components.Radio(components.RadioProps{
  Name: "plan", Value: "pro",
  Options: []components.RadioOption{
    {Label: "Free", Value: "free"},
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
			Description: "Star rating, navigable with the arrows. In read-only mode it becomes output: no focus, nothing submitted, a single text label.",
			Examples: []Example{
				{
					Title:  "Interactive",
					Code:   `@components.Rating(components.RatingProps{Name: "score", Value: 3})`,
					Render: demoRating(),
				},
				{
					Title:  "Read-only",
					Code:   `@components.Rating(components.RatingProps{Value: 4, ReadOnly: true})`,
					Render: demoRatingReadOnly(),
				},
			},
		},
		{
			Slug:        "range",
			Name:        "Range",
			Category:    CategoryDataInput,
			Description: "Numeric slider. It is the only control in the set that keeps the native input visible — reimplementing dragging, touch and pointer capture would cost a lot for a worse result. Alpine only mirrors the value into the label.",
			Examples: []Example{
				{
					Title:  "With a visible value",
					Code:   `@components.Range(components.RangeProps{ID: "vol", Name: "vol", Label: "Volume", Value: 40, ShowValue: true})`,
					Render: demoRange(),
				},
			},
		},
		{
			Slug:        "swap",
			Name:        "Swap",
			Category:    CategoryActions,
			Description: "Switches between two visuals. It is a real button with aria-pressed — a toggle whose state assistive technology cannot see is just a drawing.",
			Examples: []Example{
				{
					Title: "Play/pause and favorite",
					Code: `@components.Swap(components.SwapProps{
  Label: "Play or pause",
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
			Description: "Contextual hint. It opens on hover and also on focus, closes with Escape, and delivers its text through aria-describedby — styling alone delivers it to no one.",
			Examples: []Example{
				{
					Title:  "Placements",
					Code:   `@components.Tooltip(components.TooltipProps{ID: "tip", Text: "Shows above"}) { ...trigger... }`,
					Render: demoTooltip(),
				},
				{
					Title:  "From the keyboard",
					Code:   `Tabbing to the trigger opens the hint; Escape closes it.`,
					Render: demoTooltipKeyboard(),
				},
			},
		},
		{
			Slug:        "toast",
			Name:        "Toast",
			Category:    CategoryFeedback,
			Description: "Temporary floating notification, pinned to the bottom-right corner. It goes away on its own after 3 seconds.",
			Examples: []Example{
				{
					Title:  "Success",
					Code:   `@components.Toast(components.ToastProps{Message: "Changes saved", Variant: "success"})`,
					Render: demoToast(),
				},
				{
					Title:  "Error",
					Code:   `@components.Toast(components.ToastProps{Message: "Could not save", Variant: "error"})`,
					Render: demoToastError(),
				},
			},
		},
		{
			Slug:        "empty",
			Name:        "Empty",
			Category:    CategoryFeedback,
			Description: "Empty state with icon, title, description and action. The title is a paragraph: the heading level depends on the page, which the component does not know.",
			Examples: []Example{
				{
					Title: "With an action",
					Code: `@components.Empty(components.EmptyProps{Title: "No projects yet", Description: "..."}) {
  @components.Button(components.ButtonProps{}) { New project }
}`,
					Render: demoEmpty(),
				},
			},
		},
	}
}
