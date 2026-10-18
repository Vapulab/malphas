package showcase

import (
	"net/url"
	"slices"
	"strings"

	"github.com/vapulab/malphas/components"
)

// The Questionnaire demo's server side: the questions, the conditional one,
// validation and resuming. It is stateless — every answer so far comes back
// with each POST — so it needs no session and no storage.

var demoQuestions = []components.QuestionnaireQuestion{
	{
		Key: "experiencia", Kind: "single", Required: true,
		Title: "Have you ever run an experiment from a written protocol?",
		Options: []components.QuestionnaireOption{
			{Label: "Yes", Value: "sim"},
			{Label: "No", Value: "nao"},
		},
	},
	{
		Key: "areas", Kind: "multiple", Required: true,
		Title:       "In which areas?",
		Description: "Check all that apply.",
		Options: []components.QuestionnaireOption{
			{Label: "Genetics", Value: "genetica"},
			{Label: "Immunology", Value: "imunologia"},
			{Label: "Biochemistry", Value: "bioquimica"},
			{Label: "Neuroscience", Value: "neurociencia"},
		},
	},
	{
		Key: "objetivo", Kind: "text",
		Title:       "What is the goal of your next experiment?",
		Description: "One sentence is enough. You can skip this one.",
	},
	{
		Key: "prazo", Kind: "single", Required: true,
		Title: "When do you need the results?",
		Options: []components.QuestionnaireOption{
			{Label: "This week", Value: "1s"},
			{Label: "This month", Value: "1m"},
			{Label: "No rush", Value: "livre"},
		},
	},
}

// QuestionnaireState is what the page renders: a question with its position
// and answers, or Done with the answers to summarise.
type QuestionnaireState struct {
	Question components.QuestionnaireQuestion
	Step     int
	Total    int
	Answers  map[string][]string
	Error    string
	Done     bool
}

// visibleQuestions applies the condition: "areas" only follows a "sim".
func visibleQuestions(answers map[string][]string) []components.QuestionnaireQuestion {
	var out []components.QuestionnaireQuestion
	for _, q := range demoQuestions {
		if q.Key == "areas" && !slices.Contains(answers["experiencia"], "sim") {
			continue
		}
		out = append(out, q)
	}
	return out
}

func questionByKey(key string) (components.QuestionnaireQuestion, bool) {
	for _, q := range demoQuestions {
		if q.Key == key {
			return q, true
		}
	}
	return components.QuestionnaireQuestion{}, false
}

// cleanAnswers keeps only known questions, and for choice questions only the
// values they offer: the form is user input and must not echo arbitrary
// fields back into hidden inputs.
func cleanAnswers(form url.Values) map[string][]string {
	answers := map[string][]string{}
	for name, values := range form {
		key, ok := strings.CutPrefix(name, "a.")
		if !ok {
			continue
		}
		q, known := questionByKey(key)
		if !known {
			continue
		}
		var kept []string
		for _, v := range values {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			if q.Kind == "text" || slices.ContainsFunc(q.Options, func(o components.QuestionnaireOption) bool { return o.Value == v }) {
				kept = append(kept, v)
			}
		}
		if q.Kind == "single" && len(kept) > 1 {
			kept = kept[:1]
		}
		if len(kept) > 0 {
			answers[key] = kept
		}
	}
	return answers
}

func stateAt(visible []components.QuestionnaireQuestion, i int, answers map[string][]string, errMsg string) QuestionnaireState {
	q := visible[i]
	q.Value = answers[q.Key]
	return QuestionnaireState{Question: q, Step: i + 1, Total: len(visible), Answers: answers, Error: errMsg}
}

// pruneHidden drops answers to questions the current answers no longer show:
// switching "sim" to "nao" must not carry the areas through to the end.
func pruneHidden(answers map[string][]string, visible []components.QuestionnaireQuestion) {
	for key := range answers {
		if !slices.ContainsFunc(visible, func(q components.QuestionnaireQuestion) bool { return q.Key == key }) {
			delete(answers, key)
		}
	}
}

// QuestionnaireStart opens the questionnaire, or resumes it at "objetivo"
// with saved answers — what a returning visitor would get from storage.
func QuestionnaireStart(resume bool) QuestionnaireState {
	if !resume {
		return stateAt(visibleQuestions(nil), 0, map[string][]string{}, "")
	}
	answers := map[string][]string{"experiencia": {"sim"}, "areas": {"genetica", "imunologia"}}
	visible := visibleQuestions(answers)
	return stateAt(visible, slices.IndexFunc(visible, func(q components.QuestionnaireQuestion) bool { return q.Key == "objetivo" }), answers, "")
}

// QuestionnaireNext answers one POST: nav is "back", "skip" or "next".
// Back never validates; skip is honoured only on optional questions (a forged
// skip on a required one is handled as next); next validates required ones.
func QuestionnaireNext(form url.Values) QuestionnaireState {
	answers := cleanAnswers(form)
	visible := visibleQuestions(answers)
	pruneHidden(answers, visible)
	i := slices.IndexFunc(visible, func(q components.QuestionnaireQuestion) bool { return q.Key == form.Get("atual") })
	if i == -1 {
		return stateAt(visible, 0, answers, "")
	}
	q := visible[i]

	nav := form.Get("nav")
	if nav == "skip" && q.Required {
		nav = "next"
	}
	switch nav {
	case "back":
		return stateAt(visible, max(i-1, 0), answers, "")
	case "skip":
		delete(answers, q.Key)
		visible = visibleQuestions(answers)
	default:
		if q.Required && len(answers[q.Key]) == 0 {
			msg := "Choose an option to continue."
			if q.Kind == "multiple" {
				msg = "Choose at least one option to continue."
			}
			return stateAt(visible, i, answers, msg)
		}
	}
	// The answer just given may have changed which questions follow.
	visible = visibleQuestions(answers)
	i = slices.IndexFunc(visible, func(v components.QuestionnaireQuestion) bool { return v.Key == q.Key })
	if i+1 >= len(visible) {
		return QuestionnaireState{Answers: answers, Done: true, Total: len(visible), Step: len(visible)}
	}
	return stateAt(visible, i+1, answers, "")
}

// QuestionnaireSummary turns the answers into label pairs for the final screen.
func QuestionnaireSummary(answers map[string][]string) [][2]string {
	var out [][2]string
	for _, q := range visibleQuestions(answers) {
		vals := answers[q.Key]
		if len(vals) == 0 {
			out = append(out, [2]string{q.Title, "—"})
			continue
		}
		var labels []string
		for _, v := range vals {
			label := v
			for _, o := range q.Options {
				if o.Value == v {
					label = o.Label
				}
			}
			labels = append(labels, label)
		}
		out = append(out, [2]string{q.Title, strings.Join(labels, ", ")})
	}
	return out
}
