package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

func TestAvatarRendersImageWithAlt(t *testing.T) {
	html := renderToString(t, components.Avatar(components.AvatarProps{
		Src: "/img/lincoln.png", Alt: "Lincoln",
	}))
	if !strings.Contains(html, `src="/img/lincoln.png"`) {
		t.Errorf("expected the image, got: %s", html)
	}
	if !strings.Contains(html, `alt="Lincoln"`) {
		t.Errorf("expected the alt text, got: %s", html)
	}
}

func TestAvatarFallsBackToInitials(t *testing.T) {
	html := renderToString(t, components.Avatar(components.AvatarProps{Initials: "LS"}))
	if strings.Contains(html, "<img") {
		t.Errorf("no Src should mean no <img>, got: %s", html)
	}
	if !strings.Contains(html, "LS") {
		t.Errorf("expected the initials, got: %s", html)
	}
}

func TestAvatarSquareShape(t *testing.T) {
	html := renderToString(t, components.Avatar(components.AvatarProps{Initials: "LS", Shape: "square"}))
	if strings.Contains(html, "rounded-full") {
		t.Errorf("square avatar must not be a circle, got: %s", html)
	}
}

// The status dot is decorative next to a name; without a label it would be
// announced as nothing useful.
func TestAvatarStatusDotIsLabelled(t *testing.T) {
	html := renderToString(t, components.Avatar(components.AvatarProps{Initials: "LS", Status: "online"}))
	if !strings.Contains(html, "bg-success") {
		t.Errorf("expected an online dot, got: %s", html)
	}
	if !strings.Contains(html, "sr-only") {
		t.Errorf("expected a text label for the status, got: %s", html)
	}
}
