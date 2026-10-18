package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/vapulab/malphas/components"
)

func tabsFixture() components.TabsProps {
	return components.TabsProps{
		ID:    "conta",
		Label: "Configurações da conta",
		Items: []components.Tab{
			{Label: "Perfil", Content: templ.Raw("<p>perfil</p>")},
			{Label: "Segurança", Content: templ.Raw("<p>segurança</p>")},
			{Label: "Faturamento", Content: templ.Raw("<p>faturamento</p>"), Disabled: true},
		},
	}
}

func TestTabsUsesTheTablistRoles(t *testing.T) {
	html := renderToString(t, components.Tabs(tabsFixture()))
	if !strings.Contains(html, `role="tablist"`) {
		t.Errorf("expected a tablist, got: %s", html)
	}
	if strings.Count(html, `role="tab"`) != 3 {
		t.Errorf("expected 3 tabs, got: %s", html)
	}
	if strings.Count(html, `role="tabpanel"`) != 3 {
		t.Errorf("expected 3 panels, got: %s", html)
	}
}

// Each tab must point at its panel and each panel back at its tab, or a screen
// reader cannot tell which content belongs to the selected tab.
func TestTabsPairsTabsAndPanels(t *testing.T) {
	html := renderToString(t, components.Tabs(tabsFixture()))
	for i, ids := range [][2]string{
		{"conta-tab-0", "conta-panel-0"},
		{"conta-tab-1", "conta-panel-1"},
		{"conta-tab-2", "conta-panel-2"},
	} {
		if !strings.Contains(html, `id="`+ids[0]+`"`) || !strings.Contains(html, `id="`+ids[1]+`"`) {
			t.Errorf("tab %d: missing ids %v, got: %s", i, ids, html)
		}
		if !strings.Contains(html, `aria-controls="`+ids[1]+`"`) {
			t.Errorf("tab %d: not wired to its panel, got: %s", i, html)
		}
		if !strings.Contains(html, `aria-labelledby="`+ids[0]+`"`) {
			t.Errorf("panel %d: not labelled by its tab, got: %s", i, html)
		}
	}
}

// tablistOf isolates the tab strip. Panels are focusable too (tabindex="0" is
// correct on a tabpanel), so counting across the whole component would conflate
// the two.
func tablistOf(t *testing.T, html string) string {
	t.Helper()
	start := strings.Index(html, `role="tablist"`)
	end := strings.Index(html, `role="tabpanel"`)
	if start < 0 || end < 0 || end < start {
		t.Fatalf("could not isolate the tablist, got: %s", html)
	}
	return html[start:end]
}

// Roving tabindex: Tab reaches the tablist once, then arrows move between tabs.
func TestTabsHasExactlyOneTabbableTab(t *testing.T) {
	tablist := tablistOf(t, renderToString(t, components.Tabs(tabsFixture())))
	if got := strings.Count(tablist, `tabindex="0"`); got != 1 {
		t.Errorf("exactly one tab may be tabbable, got %d in: %s", got, tablist)
	}
	if got := strings.Count(tablist, `tabindex="-1"`); got != 2 {
		t.Errorf("the other tabs must leave the tab order, got %d in: %s", got, tablist)
	}
}

func TestTabsMarksSelectedAndHidesOtherPanels(t *testing.T) {
	html := renderToString(t, components.Tabs(tabsFixture()))
	if strings.Count(html, `aria-selected="true"`) != 1 {
		t.Errorf("exactly one tab may be selected, got: %s", html)
	}
	// The hidden attribute, not display:none — an inactive tabpanel should
	// leave the accessibility tree, not merely stop being painted. Counting
	// with a leading space and no `=` skips the x-bind: counterpart.
	if got := strings.Count(html, " hidden "); got != 2 {
		t.Errorf("expected the two inactive panels hidden, got %d in: %s", got, html)
	}
}

func TestTabsDisabledTabOptsOutOfRoving(t *testing.T) {
	html := renderToString(t, components.Tabs(tabsFixture()))
	if !strings.Contains(html, "data-roving-disabled") {
		t.Errorf("a disabled tab must be skipped by the arrows, got: %s", html)
	}
}
