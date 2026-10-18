package components_test

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/vapulab/malphas/components"
)

func renderToString(t *testing.T, c templ.Component) string {
	t.Helper()
	return renderInContext(t, context.Background(), c)
}

// renderWithChildren renders a component that declares { children... }. templ
// passes children through the context, not as a function argument.
func renderWithChildren(t *testing.T, c templ.Component, children templ.Component) string {
	t.Helper()
	return renderInContext(t, templ.WithChildren(context.Background(), children), c)
}

func renderInContext(t *testing.T, ctx context.Context, c templ.Component) string {
	t.Helper()
	var sb strings.Builder
	if err := c.Render(ctx, &sb); err != nil {
		t.Fatalf("render: %v", err)
	}
	return sb.String()
}

func TestButtonPrimary(t *testing.T) {
	html := renderToString(t, components.Button(components.ButtonProps{Variant: "primary", Size: "md"}))
	if !strings.Contains(html, "<button") {
		t.Fatalf("expected a button element, got: %s", html)
	}
	if !strings.Contains(html, "bg-primary") {
		t.Fatalf("expected primary variant class, got: %s", html)
	}
}

func TestInputRendersLabelAndField(t *testing.T) {
	html := renderToString(t, components.Input(components.InputProps{ID: "email", Name: "email", Label: "Email"}))
	if !strings.Contains(html, `for="email"`) || !strings.Contains(html, `name="email"`) {
		t.Fatalf("expected label bound to input, got: %s", html)
	}
}

func TestCardRendersChildren(t *testing.T) {
	html := renderWithChildren(t, components.Card(components.CardProps{}), templ.Raw("<p>hi</p>"))
	if !strings.Contains(html, "<p>hi</p>") {
		t.Fatalf("expected children rendered, got: %s", html)
	}
}

func TestModalHasAlpineState(t *testing.T) {
	html := renderWithChildren(t, components.Modal(components.ModalProps{ID: "edit", Title: "Editar"}), templ.Raw("body"))
	if !strings.Contains(html, "x-data") || !strings.Contains(html, "Editar") {
		t.Fatalf("expected alpine modal with title, got: %s", html)
	}
}

func TestToastSuccess(t *testing.T) {
	html := renderToString(t, components.Toast(components.ToastProps{Message: "Salvo", Variant: "success"}))
	if !strings.Contains(html, "Salvo") {
		t.Fatalf("expected toast message, got: %s", html)
	}
}
