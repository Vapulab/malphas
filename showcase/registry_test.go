package showcase_test

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/vapulab/malphas/showcase"
)

func TestAllComponentsAreWellFormed(t *testing.T) {
	all := showcase.All()
	if len(all) == 0 {
		t.Fatal("registry is empty")
	}

	seen := map[string]bool{}
	for _, c := range all {
		if c.Slug == "" || c.Name == "" || c.Category == "" || c.Description == "" {
			t.Errorf("incomplete entry: %+v", c)
		}
		if seen[c.Slug] {
			t.Errorf("duplicate slug %q", c.Slug)
		}
		seen[c.Slug] = true
		if strings.ContainsAny(c.Slug, " /_") || c.Slug != strings.ToLower(c.Slug) {
			t.Errorf("slug %q must be lowercase and URL-safe", c.Slug)
		}
		if len(c.Examples) == 0 {
			t.Errorf("component %q has no examples", c.Slug)
		}
		for _, ex := range c.Examples {
			if ex.Title == "" {
				t.Errorf("component %q has an untitled example", c.Slug)
			}
			if ex.Render == nil {
				t.Errorf("component %q example %q has no Render", c.Slug, ex.Title)
			}
		}
	}
}

// Rendering every example catches a broken demo without needing a browser.
func TestEveryExampleRenders(t *testing.T) {
	for _, c := range showcase.All() {
		for _, ex := range c.Examples {
			if ex.Render == nil {
				continue // already reported by TestAllComponentsAreWellFormed
			}
			var sb strings.Builder
			if err := ex.Render.Render(context.Background(), &sb); err != nil {
				t.Errorf("%s / %s: render failed: %v", c.Slug, ex.Title, err)
				continue
			}
			if sb.Len() == 0 {
				t.Errorf("%s / %s: rendered nothing", c.Slug, ex.Title)
			}
		}
	}
}

// A demo link that goes nowhere makes the component look broken. The
// Pagination edge examples shipped pointing at "#p", so clicking a page did
// nothing at all and read as a bug in Pagination itself.
func TestNoExampleRendersADeadLink(t *testing.T) {
	dead := regexp.MustCompile(`href="(#[a-z]*\d*|javascript:|)"`)

	for _, c := range showcase.All() {
		for _, ex := range c.Examples {
			if ex.Render == nil {
				continue
			}
			var sb strings.Builder
			if err := ex.Render.Render(context.Background(), &sb); err != nil {
				continue // reported by TestEveryExampleRenders
			}
			if m := dead.FindString(sb.String()); m != "" {
				t.Errorf("%s / %s renders a link that goes nowhere: %s", c.Slug, ex.Title, m)
			}
		}
	}
}

func TestBySlug(t *testing.T) {
	all := showcase.All()
	want := all[0]
	got, ok := showcase.BySlug(want.Slug)
	if !ok {
		t.Fatalf("BySlug(%q) not found", want.Slug)
	}
	if got.Name != want.Name {
		t.Fatalf("BySlug(%q) = %q, want %q", want.Slug, got.Name, want.Name)
	}
	if _, ok := showcase.BySlug("nao-existe"); ok {
		t.Fatal("BySlug returned ok for an unknown slug")
	}
}

// Every registered component must appear in exactly one group, or the index
// page would silently omit it.
func TestByCategoryCoversEveryComponent(t *testing.T) {
	grouped := 0
	seen := map[string]bool{}
	for _, g := range showcase.ByCategory() {
		if g.Category == "" {
			t.Error("group with empty category")
		}
		if len(g.Components) == 0 {
			t.Errorf("category %q has no components", g.Category)
		}
		for _, c := range g.Components {
			if seen[c.Slug] {
				t.Errorf("component %q appears in more than one group", c.Slug)
			}
			seen[c.Slug] = true
			grouped++
		}
	}
	if grouped != len(showcase.All()) {
		t.Fatalf("ByCategory covers %d components, registry has %d", grouped, len(showcase.All()))
	}
}
