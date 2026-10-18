package components_test

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
)

func carouselFixture() components.CarouselProps {
	return components.CarouselProps{
		ID:    "shots",
		Label: "Capturas de tela",
		Slides: []templ.Component{
			templ.Raw("<p>um</p>"),
			templ.Raw("<p>dois</p>"),
			templ.Raw("<p>três</p>"),
		},
	}
}

// A carousel reachable only by mouse is decoration.
func TestCarouselIsKeyboardReachable(t *testing.T) {
	html := renderToString(t, components.Carousel(carouselFixture()))
	if !strings.Contains(html, `tabindex="0"`) {
		t.Errorf("the track must be focusable, got: %s", html)
	}
	if !strings.Contains(html, "arrow-left") || !strings.Contains(html, "arrow-right") {
		t.Errorf("expected arrow keys to move between slides, got: %s", html)
	}
}

// Position has to be announced, not just visible.
func TestCarouselSlidesAreLabelled(t *testing.T) {
	html := renderToString(t, components.Carousel(carouselFixture()))
	if !strings.Contains(html, `aria-roledescription="carrossel"`) {
		t.Errorf("expected the carousel role description, got: %s", html)
	}
	if strings.Count(html, `aria-roledescription="slide"`) != 3 {
		t.Errorf("expected 3 labelled slides, got: %s", html)
	}
	for _, want := range []string{"1 de 3", "2 de 3", "3 de 3"} {
		if !strings.Contains(html, want) {
			t.Errorf("expected position %q announced, got: %s", want, html)
		}
	}
}

// Glyph-only buttons read as punctuation.
func TestCarouselButtonsHaveText(t *testing.T) {
	html := renderToString(t, components.Carousel(carouselFixture()))
	if !strings.Contains(html, "Anterior") || !strings.Contains(html, "Próximo") {
		t.Errorf("expected text on the controls, got: %s", html)
	}
}

func TestCarouselRendersSlideContent(t *testing.T) {
	html := renderToString(t, components.Carousel(carouselFixture()))
	for _, want := range []string{"<p>um</p>", "<p>dois</p>", "<p>três</p>"} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %s, got: %s", want, html)
		}
	}
}
