package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/vapulab/malphas/components"
)

func TestTimelineRendersItems(t *testing.T) {
	html := renderToString(t, components.Timeline(components.TimelineProps{
		Items: []components.TimelineItem{
			{Time: "09:00", Title: "Deploy iniciado", Description: "build #481"},
			{Time: "09:04", Title: "Deploy concluído", Variant: "success"},
		},
	}))
	for _, want := range []string{"09:00", "Deploy iniciado", "build #481", "Deploy concluído"} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %q, got: %s", want, html)
		}
	}
	if !strings.Contains(html, "<ol") {
		t.Errorf("a timeline is an ordered list, got: %s", html)
	}
}

func TestTimelineOmitsEmptyDescription(t *testing.T) {
	html := renderToString(t, components.Timeline(components.TimelineProps{
		Items: []components.TimelineItem{{Title: "Só título"}},
	}))
	if strings.Contains(html, "<p") {
		t.Errorf("expected no description paragraph, got: %s", html)
	}
}

// Times are machine-readable data, not decoration.
func TestTimelineUsesTimeElement(t *testing.T) {
	html := renderToString(t, components.Timeline(components.TimelineProps{
		Items: []components.TimelineItem{{Time: "09:00", Title: "x"}},
	}))
	if !strings.Contains(html, "<time") {
		t.Errorf("expected a <time> element, got: %s", html)
	}
}
