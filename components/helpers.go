package components

import "strconv"

// Plain Go, not a .templ file: these are helpers for the templates, with no
// markup of their own. templ would emit an unused import for a file with no
// components in it.

// boolJS and intJS render Go values into the JS expressions that x-data
// initialisers take. Writing them inline invites quoting mistakes, and templ
// gives no help inside an attribute string.

func boolJS(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func intJS(v int) string {
	return strconv.Itoa(v)
}

func strJS(v string) string {
	// strconv.Quote escapes quotes and control characters, which is exactly
	// what a JS string literal needs.
	return strconv.Quote(v)
}

// clampSeconds keeps a negative duration from rendering as -1:-1:-1.
func clampSeconds(s int) int {
	if s < 0 {
		return 0
	}
	return s
}

// formatDuration mirrors the formatting in countdown() so the server-rendered
// clock and the first tick agree. Keep the two in step.
func formatDuration(total int, format string) string {
	pad := func(n int) string {
		if n < 10 {
			return "0" + strconv.Itoa(n)
		}
		return strconv.Itoa(n)
	}
	if format == "mm:ss" {
		return pad(total/60) + ":" + pad(total%60)
	}
	return pad(total/3600) + ":" + pad((total%3600)/60) + ":" + pad(total%60)
}

// ariaBool and ariaChecked render the *initial* state into the server HTML.
// The x-bind: counterparts only evaluate once Alpine has loaded, so without
// these the markup would ship with no state at all and a screen reader hitting
// the page early would read the control wrong.

func ariaBool(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func ariaChecked(checked, indeterminate bool) string {
	if indeterminate {
		return "mixed"
	}
	return ariaBool(checked)
}

func ratingStarClass(filled bool) string {
	if filled {
		return "text-warning"
	}
	return "text-muted/40"
}

// ratingActiveIndex keeps the group reachable when nothing is rated yet: with
// value 0 every star would otherwise sit at tabindex=-1.
func ratingActiveIndex(value int) int {
	if value < 1 {
		return 0
	}
	return value - 1
}

// rovingTabindex renders the initial tab order for a roving list, matching what
// rovingList() will compute once Alpine takes over: exactly one item reachable
// with Tab, the rest driven by the arrow keys.
func rovingTabindex(index, active int) string {
	if index == active {
		return "0"
	}
	return "-1"
}

// floatingPosition places a panel anchored to its trigger: start|end
// horizontally, top|bottom vertically, combined as "bottom-start" (default).
// No viewport measurement — near an edge, the page author picks the side.
// Shared by every non-modal overlay that hangs off a trigger.
func floatingPosition(position string) string {
	return map[string]string{
		"bottom-start": "top-full left-0 mt-1",
		"bottom-end":   "top-full right-0 mt-1",
		"top-start":    "bottom-full left-0 mb-1",
		"top-end":      "bottom-full right-0 mb-1",
	}[or(position, "bottom-start")]
}
