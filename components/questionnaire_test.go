package components_test

import (
	"strings"
	"testing"

	"github.com/darkanum/vapulab/malphas/components"
)

func questionnaireFixture() components.QuestionnaireProps {
	return components.QuestionnaireProps{
		ID: "q", Action: "/responder", Step: 2, Total: 4,
		Question: components.QuestionnaireQuestion{
			Key: "areas", Title: "Quais áreas?", Kind: "multiple", Required: true,
			Options: []components.QuestionnaireOption{{Label: "Genética", Value: "gen"}, {Label: "Imunologia", Value: "imu"}},
			Value:   []string{"imu"},
		},
		Answers: map[string][]string{"experiencia": {"sim"}, "areas": {"imu"}},
	}
}

func TestQuestionnaireMarkup(t *testing.T) {
	html := renderToString(t, components.Questionnaire(questionnaireFixture()))
	for _, want := range []string{
		`hx-post="/responder"`,
		`method="post"`,
		`action="/responder"`,
		`<fieldset`,
		`<legend`,
		`Quais áreas?`,
		`type="checkbox" name="a.areas" value="gen"`,
		`type="checkbox" name="a.areas" value="imu" checked`,
		// Earlier answers ride along; the current one comes from its inputs.
		`<input type="hidden" name="a.experiencia" value="sim">`,
		`<input type="hidden" name="atual" value="areas">`,
		`Pergunta 2 de 4`,
		`formnovalidate`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	if strings.Contains(html, `type="hidden" name="a.areas"`) {
		t.Errorf("the current answer must not also be a hidden field: %s", html)
	}
	// Required: no Skip button.
	if strings.Contains(html, `value="skip"`) {
		t.Errorf("a required question cannot be skipped: %s", html)
	}
	// Enter submits with the first submit button in the form: it must be Next,
	// never Back.
	next := strings.Index(html, `value="next"`)
	back := strings.Index(html, `value="back"`)
	if next < 0 || back < 0 || next > back {
		t.Errorf("Next must come before Back in the DOM (Enter uses the first submit): %s", html)
	}
}

func TestQuestionnaireFirstStepHasNoBack(t *testing.T) {
	p := questionnaireFixture()
	p.Step = 1
	html := renderToString(t, components.Questionnaire(p))
	if strings.Contains(html, `value="back"`) {
		t.Errorf("step 1 has nowhere to go back to: %s", html)
	}
}

func TestQuestionnaireErrorIsAnnouncedAndLinked(t *testing.T) {
	p := questionnaireFixture()
	p.Error = "Escolha ao menos uma área."
	html := renderToString(t, components.Questionnaire(p))
	for _, want := range []string{`role="alert"`, `id="q-error"`, `aria-describedby="q-error"`, "Escolha ao menos uma área."} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}

func TestQuestionnaireSkippableText(t *testing.T) {
	p := questionnaireFixture()
	p.Question = components.QuestionnaireQuestion{Key: "objetivo", Title: "Qual o objetivo?", Kind: "text", Value: []string{"medir"}}
	p.Step, p.Total = 4, 4
	html := renderToString(t, components.Questionnaire(p))
	for _, want := range []string{`type="text"`, `name="a.objetivo"`, `value="medir"`, `value="skip"`, `Concluir`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
}
