// Package showcase is the living documentation of the Malphas design system: a
// registry of components, each with runnable examples.
//
// It lives beside components/ and layout/, not inside them, on purpose. Demo
// pages use utility classes no real page needs; a project that points its
// Tailwind content at components/ and layout/, as the README says, never
// compiles them. A host that serves the catalog adds showcase/ to its content.
package showcase

import (
	"sync"

	"github.com/a-h/templ"
)

// Categories in the order they should be presented.
const (
	CategoryActions     = "Actions"
	CategoryDataDisplay = "Data display"
	CategoryDataInput   = "Data input"
	CategoryNavigation  = "Navigation"
	CategoryFeedback    = "Feedback"
	CategoryLayout      = "Layout & utilities"
)

var categoryOrder = []string{
	CategoryActions,
	CategoryDataDisplay,
	CategoryDataInput,
	CategoryNavigation,
	CategoryFeedback,
	CategoryLayout,
}

// Example is one runnable illustration of a component.
//
// Code is written by hand and is not derived from Render, so the two can drift.
// Keeping them in sync automatically would need embedded snippets or codegen —
// not worth it at this size.
type Example struct {
	Title  string
	Code   string
	Render templ.Component
}

// Component is a design-system component as presented in the docs.
type Component struct {
	Slug        string // URL segment, e.g. "button" for /components/button
	Name        string
	Category    string
	Description string
	Examples    []Example
}

// Group is a category with its components, for the index page.
type Group struct {
	Category   string
	Components []Component
}

// All returns every registered component. The catalog is built once; templ
// components are immutable and safe to share across requests.
var All = sync.OnceValue(func() []Component { return catalog() })

var bySlug = sync.OnceValue(func() map[string]Component {
	m := make(map[string]Component, len(All()))
	for _, c := range All() {
		m[c.Slug] = c
	}
	return m
})

// BySlug looks up a component for the /components/{slug} route.
func BySlug(slug string) (Component, bool) {
	c, ok := bySlug()[slug]
	return c, ok
}

// ByCategory groups components for the index, in categoryOrder. Categories with
// no components are omitted; a component in an unknown category would be
// dropped, which TestByCategoryCoversEveryComponent guards against.
var ByCategory = sync.OnceValue(func() []Group {
	byCat := map[string][]Component{}
	for _, c := range All() {
		byCat[c.Category] = append(byCat[c.Category], c)
	}
	groups := make([]Group, 0, len(categoryOrder))
	for _, cat := range categoryOrder {
		if cs := byCat[cat]; len(cs) > 0 {
			groups = append(groups, Group{Category: cat, Components: cs})
		}
	}
	return groups
})
