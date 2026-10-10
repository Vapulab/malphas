package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/darkanum/vapulab/malphas/components"
)

func attachment(state string) components.AttachmentProps {
	return components.AttachmentProps{ID: "a", Name: "protocolo.pdf", Meta: "PDF · 2,4 MB", State: state, Progress: 45, Href: "/arquivos/protocolo.pdf"}
}

func TestAttachmentStates(t *testing.T) {
	cases := map[string]struct {
		want    []string
		without []string
	}{
		"uploading":  {[]string{"Enviando 45%", `<progress`, `value="45"`}, []string{`href=`}},
		"processing": {[]string{"Processando", `<progress`}, []string{`value=`, `href=`}},
		"done":       {[]string{`href="/arquivos/protocolo.pdf"`, "PDF · 2,4 MB"}, []string{`<progress`}},
		"idle":       {[]string{"protocolo.pdf"}, []string{`<progress`, `href=`}},
	}
	for state, c := range cases {
		html := renderToString(t, components.Attachment(attachment(state)))
		for _, w := range c.want {
			if !strings.Contains(html, w) {
				t.Errorf("%s: missing %s in %s", state, w, html)
			}
		}
		for _, w := range c.without {
			if strings.Contains(html, w) {
				t.Errorf("%s: should not contain %s: %s", state, w, html)
			}
		}
	}
}

// The error is announced, and the actions (retry) stay reachable.
func TestAttachmentErrorIsAnAlert(t *testing.T) {
	p := attachment("error")
	p.Error = "Falha na conexão."
	p.Actions = templ.Raw(`<button>Tentar de novo</button>`)
	html := renderToString(t, components.Attachment(p))
	for _, want := range []string{`role="alert"`, "Falha na conexão.", "Tentar de novo"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}

func TestAttachmentSanitizesHref(t *testing.T) {
	p := attachment("done")
	p.Href = "javascript:alert(1)"
	if html := renderToString(t, components.Attachment(p)); strings.Contains(html, "javascript:") {
		t.Errorf("a javascript: href reached the markup: %s", html)
	}
}

func TestAttachmentGroupScrollsHorizontally(t *testing.T) {
	html := renderToString(t, components.AttachmentGroup(components.AttachmentGroupProps{
		Label: "Anexos", Items: []components.AttachmentProps{attachment("done"), attachment("uploading")},
	}))
	for _, want := range []string{`role="region"`, `aria-label="Anexos"`, "overflow-x-auto", "<ul", "<li"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	if strings.Count(html, "<li") != 2 {
		t.Errorf("expected two items: %s", html)
	}
}

// The live region exists in every state, so a change of state is announced
// as a change; a region born together with its text is often not read.
func TestAttachmentStatusRegionAlwaysPresent(t *testing.T) {
	for _, state := range []string{"idle", "uploading", "processing", "done"} {
		html := renderToString(t, components.Attachment(attachment(state)))
		if !strings.Contains(html, `role="status"`) {
			t.Errorf("%s: expected a status region: %s", state, html)
		}
	}
}
