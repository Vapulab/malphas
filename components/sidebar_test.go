package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/darkanum/vapulab/malphas/components"
)

func sidebarFixture() components.SidebarProps {
	return components.SidebarProps{
		ID: "app", Label: "Navegação principal", Collapsible: true,
		Groups: []components.SidebarGroup{
			{Label: "Trabalho", Items: []components.SidebarItem{
				{Label: "Protocolos", Href: "/protocolos", Icon: "◆", Current: true, Badge: "3"},
				{Label: "Resultados", Href: "/resultados", Icon: "▲"},
			}},
			{Label: "Conta", Items: []components.SidebarItem{
				{Label: "Mau", Href: "javascript:alert(1)", Icon: "!"},
			}},
		},
	}
}

func TestSidebarNavAndCurrent(t *testing.T) {
	html := renderWithChildren(t, components.Sidebar(sidebarFixture()), templ.Raw("<p>conteúdo</p>"))
	for _, want := range []string{
		`<nav aria-label="Navegação principal"`,
		`href="/protocolos"`,
		`aria-current="page"`,
		`title="Protocolos"`,
		`>Trabalho<`,
		`>Conta<`,
		`<main`,
		"conteúdo",
		`x-data="sidebar(&#34;app&#34;)"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	// The nav appears twice (aside + mobile drawer); aria-current once in each.
	if n := strings.Count(html, `aria-current="page"`); n != 2 {
		t.Errorf("expected aria-current on the current item in both navs, got %d", n)
	}
	if strings.Contains(html, `href="javascript:`) {
		t.Errorf("hrefs must go through templ's sanitizer: %s", html)
	}
}

func TestSidebarCollapsibleControl(t *testing.T) {
	html := renderWithChildren(t, components.Sidebar(sidebarFixture()), templ.Raw(""))
	if !strings.Contains(html, `aria-expanded="true"`) || !strings.Contains(html, `aria-controls="app-aside"`) {
		t.Errorf("expected a collapse control wired to the aside: %s", html)
	}
	p := sidebarFixture()
	p.Collapsible = false
	html = renderWithChildren(t, components.Sidebar(p), templ.Raw(""))
	if strings.Contains(html, `aria-controls="app-aside"`) {
		t.Errorf("not collapsible, no collapse control: %s", html)
	}
}

// On a phone the same navigation opens in a modal Drawer.
func TestSidebarHasAMobileDrawer(t *testing.T) {
	html := renderWithChildren(t, components.Sidebar(sidebarFixture()), templ.Raw(""))
	for _, want := range []string{`$dispatch(&#39;open-drawer&#39;, &#39;app-drawer&#39;)`, `x-data="drawer"`, `aria-labelledby="app-drawer-title"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}
