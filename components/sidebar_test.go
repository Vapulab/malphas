package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
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

// The page layout owns <main>; a second one nested inside it is invalid HTML
// and a duplicate landmark.
func TestSidebarDoesNotClaimMain(t *testing.T) {
	html := renderWithChildren(t, components.Sidebar(sidebarFixture()), templ.Raw("<p>conteúdo</p>"))
	if strings.Contains(html, "<main") {
		t.Errorf("Sidebar must not render <main>: %s", html)
	}
}

// The remembered state is applied before first paint, by an inline script
// right after the aside, and the width only animates once that is done.
func TestSidebarRestoresBeforePaint(t *testing.T) {
	html := renderWithChildren(t, components.Sidebar(sidebarFixture()), templ.Raw(""))
	aside := strings.Index(html, `id="app-aside"`)
	script := strings.Index(html, `localStorage.getItem("malphas-sidebar:app")`)
	if script < aside {
		t.Fatalf("expected an inline restore script after the aside: %s", html)
	}
	if !strings.Contains(html, "data-[ready=true]:transition-[width]") || strings.Contains(html, " transition-[width] ") {
		t.Errorf("the width transition must wait for data-ready: %s", html)
	}
}

// The ID reaches an inline script: it must not be able to close it.
func TestSidebarScriptEscapesTheID(t *testing.T) {
	p := sidebarFixture()
	p.ID = "x</script><script>alert(1)"
	html := renderWithChildren(t, components.Sidebar(p), templ.Raw(""))
	if strings.Contains(html, "</script><script>alert(1)") {
		t.Errorf("the ID broke out of the inline script: %s", html)
	}
}

// Collapsed, a link without an Icon would be an empty focus stop; its
// label's first letter stands in, decorative.
func TestSidebarIconlessItemKeepsAVisibleMark(t *testing.T) {
	p := sidebarFixture()
	p.Groups = []components.SidebarGroup{{Items: []components.SidebarItem{{Label: "Equipe", Href: "/e"}}}}
	html := renderWithChildren(t, components.Sidebar(p), templ.Raw(""))
	if !strings.Contains(html, `>E</span>`) {
		t.Errorf("expected the initial as a collapsed fallback: %s", html)
	}
}

// One stable name plus aria-expanded; a label that also flips would announce
// the state twice ("Expandir menu, recolhido").
func TestSidebarToggleHasOneStableName(t *testing.T) {
	html := renderWithChildren(t, components.Sidebar(sidebarFixture()), templ.Raw(""))
	if strings.Contains(html, "x-bind:aria-label") || !strings.Contains(html, `aria-label="Sidebar"`) {
		t.Errorf("expected a static aria-label with aria-expanded: %s", html)
	}
}

// Following a link from the mobile drawer closes it.
func TestSidebarDrawerLinksClose(t *testing.T) {
	html := renderWithChildren(t, components.Sidebar(sidebarFixture()), templ.Raw(""))
	drawer := html[strings.Index(html, `id="app-drawer-title"`):]
	if !strings.Contains(drawer, `x-on:click="closeOverlay()"`) {
		t.Errorf("links in the drawer should close it: %s", drawer)
	}
	aside := html[strings.Index(html, `id="app-aside"`):strings.Index(html, `</aside>`)]
	if strings.Contains(aside, "closeOverlay") {
		t.Errorf("links in the aside have no drawer to close: %s", aside)
	}
}
