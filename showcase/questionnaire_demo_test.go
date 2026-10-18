package showcase_test

import (
	"net/url"
	"testing"

	"github.com/vapulab/malphas/showcase"
)

func form(pairs ...string) url.Values {
	v := url.Values{}
	for i := 0; i < len(pairs); i += 2 {
		v.Add(pairs[i], pairs[i+1])
	}
	return v
}

func TestQuestionnaireBranches(t *testing.T) {
	sim := showcase.QuestionnaireNext(form("atual", "experiencia", "nav", "next", "a.experiencia", "sim"))
	if sim.Question.Key != "areas" || sim.Step != 2 || sim.Total != 4 {
		t.Errorf("\"sim\" should lead to areas (2 of 4), got %s (%d of %d)", sim.Question.Key, sim.Step, sim.Total)
	}
	nao := showcase.QuestionnaireNext(form("atual", "experiencia", "nav", "next", "a.experiencia", "nao"))
	if nao.Question.Key != "objetivo" || nao.Step != 2 || nao.Total != 3 {
		t.Errorf("\"não\" should skip areas, landing on objetivo (2 of 3), got %s (%d of %d)", nao.Question.Key, nao.Step, nao.Total)
	}
}

func TestQuestionnaireValidatesOnlyForward(t *testing.T) {
	empty := showcase.QuestionnaireNext(form("atual", "areas", "nav", "next", "a.experiencia", "sim"))
	if empty.Question.Key != "areas" || empty.Error == "" {
		t.Errorf("a required question left empty should come back with an error, got %s %q", empty.Question.Key, empty.Error)
	}
	back := showcase.QuestionnaireNext(form("atual", "areas", "nav", "back", "a.experiencia", "sim"))
	if back.Question.Key != "experiencia" || back.Error != "" {
		t.Errorf("Back never validates: got %s with error %q", back.Question.Key, back.Error)
	}
	if got := back.Question.Value; len(got) != 1 || got[0] != "sim" {
		t.Errorf("going back should show the earlier answer, got %v", got)
	}
}

func TestQuestionnaireSkipOnlyWhenAllowed(t *testing.T) {
	skipped := showcase.QuestionnaireNext(form("atual", "objetivo", "nav", "skip", "a.experiencia", "nao", "a.objetivo", "rascunho"))
	if skipped.Question.Key != "prazo" {
		t.Errorf("skip on objetivo should advance to prazo, got %s", skipped.Question.Key)
	}
	if _, kept := skipped.Answers["objetivo"]; kept {
		t.Error("a skipped answer should not be kept")
	}
	// A forged skip on a required question is treated as next, and validated.
	forged := showcase.QuestionnaireNext(form("atual", "experiencia", "nav", "skip"))
	if forged.Question.Key != "experiencia" || forged.Error == "" {
		t.Errorf("experiencia is required: skip must not pass it, got %s %q", forged.Question.Key, forged.Error)
	}
}

func TestQuestionnaireFinishesAndResumes(t *testing.T) {
	done := showcase.QuestionnaireNext(form("atual", "prazo", "nav", "next", "a.experiencia", "nao", "a.prazo", "1m"))
	if !done.Done {
		t.Fatalf("answering the last question should finish, got %s", done.Question.Key)
	}
	resumed := showcase.QuestionnaireStart(true)
	if resumed.Question.Key != "objetivo" || resumed.Step != 3 || len(resumed.Answers["areas"]) == 0 {
		t.Errorf("resume should open at objetivo (3 of 4) with saved answers, got %s (%d) %v", resumed.Question.Key, resumed.Step, resumed.Answers)
	}
	fresh := showcase.QuestionnaireStart(false)
	if fresh.Question.Key != "experiencia" || fresh.Step != 1 {
		t.Errorf("a fresh start opens at the first question, got %s", fresh.Question.Key)
	}
}

// Unknown keys in the posted form are ignored rather than echoed back.
func TestQuestionnaireIgnoresUnknownAnswers(t *testing.T) {
	s := showcase.QuestionnaireNext(form("atual", "experiencia", "nav", "next", "a.experiencia", "sim", "a.injetado", "x"))
	if _, ok := s.Answers["injetado"]; ok {
		t.Error("an unknown answer key was kept")
	}
}

// Answers to questions that are no longer shown are dropped, not carried
// along to the end: switching "sim" to "nao" must not submit the areas.
func TestQuestionnaireDropsAnswersToHiddenQuestions(t *testing.T) {
	s := showcase.QuestionnaireNext(form("atual", "experiencia", "nav", "next", "a.experiencia", "nao", "a.areas", "genetica"))
	if _, kept := s.Answers["areas"]; kept {
		t.Errorf("areas is not shown after \"nao\" and must be dropped: %v", s.Answers)
	}
}
