package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

func TestSkeletonAnimatesAndIsHiddenFromScreenReaders(t *testing.T) {
	html := renderToString(t, components.Skeleton(components.SkeletonProps{}))
	if !strings.Contains(html, "animate-pulse") {
		t.Errorf("expected a pulse animation, got: %s", html)
	}
	// A placeholder has no content worth announcing.
	if !strings.Contains(html, "aria-hidden") {
		t.Errorf("expected aria-hidden, got: %s", html)
	}
}

func TestSkeletonAcceptsSizingClasses(t *testing.T) {
	html := renderToString(t, components.Skeleton(components.SkeletonProps{Class: "h-4 w-32"}))
	if !strings.Contains(html, "h-4 w-32") {
		t.Errorf("expected the caller's sizing classes, got: %s", html)
	}
}
