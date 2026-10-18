package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

// The footer is a navigation landmark: <footer> alone already maps to
// contentinfo, but only when it is a direct child of body. Here it is nested
// in the page shell, so the role has to be explicit — without it the screen
// reader does not offer the shortcut to it.
func TestFooterIsALandmark(t *testing.T) {
	html := renderToString(t, components.Footer(components.FooterProps{
		Note: "Vapulab 2026",
		Columns: []components.FooterColumn{{
			Title: "Product",
			Links: []components.FooterLink{{Label: "Components", Href: "/components"}},
		}},
	}))

	if !strings.Contains(html, "<footer") {
		t.Error("Footer does not render a <footer> element")
	}
	if !strings.Contains(html, `role="contentinfo"`) {
		t.Error("nested <footer> without role=contentinfo is not exposed as a landmark")
	}
	if !strings.Contains(html, "Product") || !strings.Contains(html, "/components") {
		t.Error("column title or link is missing from the output")
	}
	if !strings.Contains(html, "Vapulab 2026") {
		t.Error("the note line is missing")
	}
}

// An external link without rel=noopener lets the opened page reach the origin
// through window.opener. The system's Link already handles that, but checking
// only rel/target does not prove Footer actually delegates to it: a hand-written
// <a> with the same two attributes would pass just the same. The proof only
// delegation satisfies is href sanitisation: templ.URL rejects schemes outside
// http/https/mailto/tel/ftp/ftps and returns FailedSanitizationURL. A raw href
// would not go through that.
func TestFooterExternalLinksAreSafe(t *testing.T) {
	html := renderToString(t, components.Footer(components.FooterProps{
		Columns: []components.FooterColumn{{
			Title: "Code",
			Links: []components.FooterLink{
				{Label: "GitHub", Href: "https://github.com/vapulab/vapulab", External: true},
				{Label: "Malicious", Href: "javascript:alert(1)"},
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
