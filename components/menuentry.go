package components

// MenuEntry is one line in a Menu, Navbar or Dropdown.
//
// A divider is a kind of entry rather than a boolean on an item: it has no
// label, takes no focus and is not clickable, so folding it into the item
// struct would force every consumer to check a flag before reading any other
// field. Build one with MenuDivider().
type MenuEntry struct {
	Kind MenuEntryKind // empty means a link

	Label string
	Href  string
	// OnClick is a JS expression for entries that act rather than navigate,
	// e.g. "$dispatch('logout')". Ignored when Href is set.
	OnClick  string
	Current  bool // marks the page the user is already on
	Disabled bool
}

type MenuEntryKind string

const (
	MenuEntryLink    MenuEntryKind = ""
	MenuEntryDivider MenuEntryKind = "divider"
)

// MenuDivider returns a separator entry.
func MenuDivider() MenuEntry {
	return MenuEntry{Kind: MenuEntryDivider}
}

func (e MenuEntry) isDivider() bool {
	return e.Kind == MenuEntryDivider
}

func dropdownTriggerClasses(variant string) string {
	return map[string]string{
		"primary":   "bg-primary text-primary-fg hover:bg-primary/90",
		"secondary": "border border-border bg-surface text-fg hover:bg-muted/10",
		"ghost":     "text-fg hover:bg-muted/10",
	}[variant]
}

func dropdownItemClasses(e MenuEntry) string {
	if e.Disabled {
		return "cursor-not-allowed text-muted/50"
	}
	return "text-fg hover:bg-muted/10"
}

// navLinkClasses styles a link or action entry. Current and disabled states are
// carried by attributes too (aria-current, disabled); these are only the paint.
func navLinkClasses(e MenuEntry) string {
	switch {
	case e.Disabled:
		return "cursor-not-allowed text-muted/50"
	case e.Current:
		return "font-medium text-fg"
	default:
		return "text-muted transition-colors hover:text-fg"
	}
}
