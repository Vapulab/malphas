package components_test

import (
	"strings"
	"testing"

	"github.com/darkanum/vapulab/malphas/components"
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
		ID: "bio", Name: "bio", Value: "\nsegunda linha",
	}))
	if !strings.Contains(html, ">\n\nsegunda linha</textarea>") {
		t.Errorf("expected a guard newline before the value, got: %q", html)
	}
}
