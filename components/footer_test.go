package components_test

import (
	"strings"
	"testing"

	"github.com/darkanum/vapulab/malphas/components"
)

// O rodape e um marco de navegacao: <footer> sozinho ja mapeia para
// contentinfo, mas so quando e filho direto de body. Aqui ele e aninhado no
// shell da pagina, entao o papel precisa ser explicito — sem isso o leitor de
// tela nao oferece o atalho para ele.
func TestFooterIsALandmark(t *testing.T) {
	html := renderToString(t, components.Footer(components.FooterProps{
		Note: "Vapulab 2026",
		Columns: []components.FooterColumn{{
			Title: "Produto",
			Links: []components.FooterLink{{Label: "Componentes", Href: "/components"}},
		}},
	}))

	if !strings.Contains(html, "<footer") {
		t.Error("Footer does not render a <footer> element")
	}
	if !strings.Contains(html, `role="contentinfo"`) {
		t.Error("nested <footer> without role=contentinfo is not exposed as a landmark")
	}
	if !strings.Contains(html, "Produto") || !strings.Contains(html, "/components") {
		t.Error("column title or link is missing from the output")
	}
	if !strings.Contains(html, "Vapulab 2026") {
		t.Error("the note line is missing")
	}
}

// Link externo sem rel=noopener deixa a pagina aberta alcancar a origem via
// window.opener. O Link do sistema ja resolve isso, mas so checar rel/target
// nao prova que o Footer de fato delega a ele: um <a> escrito na mao com os
// mesmos dois atributos passaria igual. A prova que so a delegacao satisfaz e
// o saneamento de href: templ.URL rejeita esquemas fora de
// http/https/mailto/tel/ftp/ftps e devolve FailedSanitizationURL. Um href cru
// nao passaria por isso.
func TestFooterExternalLinksAreSafe(t *testing.T) {
	html := renderToString(t, components.Footer(components.FooterProps{
		Columns: []components.FooterColumn{{
			Title: "Código",
			Links: []components.FooterLink{
				{Label: "GitHub", Href: "https://github.com/darkanum/vapulab", External: true},
				{Label: "Malicioso", Href: "javascript:alert(1)"},
			},
		}},
	}))

	if !strings.Contains(html, `rel="noopener noreferrer"`) {
		t.Error("external footer link is missing rel=noopener noreferrer")
	}
	if !strings.Contains(html, `target="_blank"`) {
		t.Error("external footer link does not open in a new tab")
	}
	if strings.Contains(html, "javascript:alert(1)") {
		t.Error("javascript: href leaked into the output — link was not sanitised, possible XSS")
	}
	if !strings.Contains(html, "about:invalid#TemplFailedSanitizationURL") {
		t.Error("expected templ's sanitised placeholder for an unsafe scheme — Footer is not delegating to Link")
	}
}
