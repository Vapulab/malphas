package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

func TestStatRendersAllFields(t *testing.T) {
	html := renderToString(t, components.Stat(components.StatProps{
		Title:       "Receita",
		Value:       "R$ 12.400",
		Description: "+12% no mês",
	}))
	for _, want := range []string{"Receita", "R$ 12.400", "+12% no mês"} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %q, got: %s", want, html)
		}
	}
}

func TestStatOmitsEmptyDescription(t *testing.T) {
	html := renderToString(t, components.Stat(components.StatProps{Title: "Usuários", Value: "42"}))
	if strings.Count(html, "<p") > 1 {
		t.Errorf("expected no empty description paragraph, got: %s", html)
	}
}

func TestStatVariantColorsTheValue(t *testing.T) {
	html := renderToString(t, components.Stat(components.StatProps{
		Title: "Falhas", Value: "3", Variant: "danger",
	}))
	if !strings.Contains(html, "text-danger") {
		t.Errorf("expected the danger color on the value, got: %s", html)
	}
}
