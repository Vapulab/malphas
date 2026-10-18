package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

func stepsFixture() components.StepsProps {
	return components.StepsProps{Items: []components.Step{
		{Label: "Conta", State: "done"},
		{Label: "Pagamento", State: "current"},
		{Label: "Confirmação"},
	}}
}

func TestStepsRendersEveryLabel(t *testing.T) {
	html := renderToString(t, components.Steps(stepsFixture()))
	for _, want := range []string{"Conta", "Pagamento", "Confirmação"} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %q, got: %s", want, html)
		}
	}
	if !strings.Contains(html, "<ol") {
		t.Errorf("steps are an ordered list, got: %s", html)
	}
}

// Progress is shown with colour and a checkmark, so the current step also needs
// aria-current for anyone not seeing it.
func TestStepsMarksCurrentStep(t *testing.T) {
	html := renderToString(t, components.Steps(stepsFixture()))
	if !strings.Contains(html, `aria-current="step"`) {
		t.Errorf("expected aria-current=step, got: %s", html)
	}
	if strings.Count(html, `aria-current="step"`) != 1 {
		t.Errorf("exactly one step may be current, got: %s", html)
	}
}

func TestStepsDoneStepUsesPrimary(t *testing.T) {
	html := renderToString(t, components.Steps(stepsFixture()))
	if !strings.Contains(html, "bg-primary") {
		t.Errorf("expected a completed step to be filled, got: %s", html)
	}
}
