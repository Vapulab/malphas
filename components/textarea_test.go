package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

func TestTextareaBasics(t *testing.T) {
	html := renderToString(t, components.Textarea(components.TextareaProps{
		ID: "bio", Name: "bio", Label: "Bio", Required: true,
	}))
	for _, want := range []string{`for="bio"`, `name="bio"`, `rows="4"`, `required`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	if strings.Contains(html, "charCounter") {
		t.Errorf("no MaxLength, no counter: %s", html)
	}
}

// maxlength and the JS counter both count UTF-16 code units, so the
// server-rendered count must too: "São 😀" is 6, not 5.
func TestTextareaCounterCountsLikeTheBrowser(t *testing.T) {
	html := renderToString(t, components.Textarea(components.TextareaProps{
		ID: "bio", Name: "bio", Value: "São 😀", MaxLength: 10,
	}))
	if !strings.Contains(html, `maxlength="10"`) {
		t.Errorf("expected maxlength, got: %s", html)
	}
	if !strings.Contains(html, `x-data="charCounter(6)"`) {
		t.Errorf("expected the counter to start at 6, got: %s", html)
	}
	if !strings.Contains(html, `<span x-text="count">6</span>/10`) {
		t.Errorf("expected 6/10 painted before Alpine, got: %s", html)
	}
}

// The HTML parser drops one newline right after <textarea>. Without a guard
// newline, a saved value that starts with a blank line loses it on every
// round trip.
func TestTextareaKeepsALeadingNewline(t *testing.T) {
	html := renderToString(t, components.Textarea(components.TextareaProps{
		ID: "bio", Name: "bio", Value: "\nsecond line",
	}))
	if !strings.Contains(html, ">\n\nsecond line</textarea>") {
		t.Errorf("expected a guard newline before the value, got: %q", html)
	}
}

// Browsers submit line breaks as CRLF, and the HTML parser turns them back
// into LF, so the value the server echoes after a validation error has two
// units per break where the browser counts one.
func TestTextareaCounterIgnoresCRLF(t *testing.T) {
	html := renderToString(t, components.Textarea(components.TextareaProps{
		ID: "bio", Name: "bio", Value: "a\r\nb\rc", MaxLength: 10,
	}))
	if !strings.Contains(html, `x-data="charCounter(5)"`) {
		t.Errorf("\"a\r\nb\rc\" is 5 units in the browser, got: %s", html)
	}
}
