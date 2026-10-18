package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

func toggleGroupFixture() components.ToggleGroupProps {
	return components.ToggleGroupProps{
		Name: "alinhamento", Label: "Alignment",
		Items: []components.ToggleGroupItem{
			{Label: "Left", Value: "esq"},
			{Label: "Center", Value: "centro"},
			{Label: "Right", Value: "dir", Disabled: true},
		},
		Value: []string{"centro"},
	}
}

func TestToggleGroupIsALabelledGroupOfPressableButtons(t *testing.T) {
	html := renderToString(t, components.ToggleGroup(toggleGroupFixture()))
	for _, want := range []string{
		`role="group"`,
		`aria-label="Alignment"`,
		`x-data="toggleGroup([false,true,false], false, 1)"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in: %s", want, html)
		}
	}
	if strings.Count(html, `aria-pressed="true"`) != 1 {
		t.Errorf("expected one pressed item, got: %s", html)
	}
}

// Roving tabindex: exactly one item is tabbable, and it is the pressed one.
func TestToggleGroupTabStopIsThePressedItem(t *testing.T) {
	html := renderToString(t, components.ToggleGroup(toggleGroupFixture()))
	if strings.Count(html, `tabindex="0"`) != 1 {
		t.Fatalf("expected exactly one tab stop, got: %s", html)
	}
	stop := strings.Index(html, `tabindex="0"`)
	centro := strings.Index(html, `>Center<`)
	esquerda := strings.Index(html, `>Left<`)
	if !(esquerda < stop && stop < centro) {
		t.Errorf("the tab stop should be on Center, got: %s", html)
	}
}

// With nothing pressed and the first item disabled, the stop moves to the
// first enabled item — otherwise the group is unreachable by keyboard.
func TestToggleGroupTabStopSkipsDisabledWhenNothingPressed(t *testing.T) {
	props := toggleGroupFixture()
	props.Value = nil
	props.Items[0].Disabled = true
	html := renderToString(t, components.ToggleGroup(props))
	if !strings.Contains(html, `x-data="toggleGroup([false,false,false], false, 1)"`) {
		t.Errorf("expected the active index on the first enabled item, got: %s", html)
	}
}

// One hidden input per item; only pressed, enabled ones are submitted.
func TestToggleGroupHiddenInputs(t *testing.T) {
	html := renderToString(t, components.ToggleGroup(toggleGroupFixture()))
	if !strings.Contains(html, `<input type="hidden" name="alinhamento" value="esq" disabled`) {
		t.Errorf("unpressed item must not submit, got: %s", html)
	}
	if strings.Contains(html, `value="centro" disabled`) {
		t.Errorf("pressed item must submit, got: %s", html)
	}
}
