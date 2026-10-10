package malphas_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/darkanum/vapulab/malphas"
)

// Behavior names reach the markup two ways, and both have to be scanned:
//
//	x-data="tabs"                              a literal attribute
//	x-data={ "rating(" + intJS(v) + ")" }      a templ expression building one
//
// An object literal (x-data="{ ... }") or a bare x-data names no behavior, so
// neither matches.
var xDataName = regexp.MustCompile(
	`x-data="([A-Za-z_$][A-Za-z0-9_$]*)\s*(?:\(|")` + // literal
		`|x-data=\{\s*"([A-Za-z_$][A-Za-z0-9_$]*)`, // templ expression
)

// behaviorNamesIn returns every Alpine behavior the source references.
func behaviorNamesIn(src string) []string {
	var names []string
	for _, m := range xDataName.FindAllStringSubmatch(src, -1) {
		// Exactly one of the two groups matches per hit.
		for _, g := range m[1:] {
			if g != "" {
				names = append(names, g)
			}
		}
	}
	return names
}

func readUIJS(t *testing.T) string {
	t.Helper()
	data, err := fs.ReadFile(malphas.Assets(), "malphas.js")
	if err != nil {
		t.Fatalf("read ui.js: %v", err)
	}
	return string(data)
}

// Components reference Alpine behaviors by name, as a string. Nothing in the
// Go type system connects the two, so a typo or a forgotten registration is
// invisible until someone clicks the thing in a browser — and with no Node
// toolchain there is no JS test suite to catch it either. This test is the
// only automated link between the markup and the behaviors it depends on.
func TestEveryReferencedBehaviorIsRegistered(t *testing.T) {
	js := readUIJS(t)

	templFiles, err := filepath.Glob(filepath.Join("components", "*.templ"))
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	layoutFiles, err := filepath.Glob(filepath.Join("layout", "*.templ"))
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	files := append(templFiles, layoutFiles...)
	if len(files) == 0 {
		t.Fatal("found no .templ files to scan — has the layout moved?")
	}

	found := 0
	for _, path := range files {
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		for _, name := range behaviorNamesIn(string(src)) {
			found++
			if !strings.Contains(js, `Alpine.data("`+name+`"`) {
				t.Errorf("%s uses x-data=%q, but ui.js never registers it", path, name)
			}
		}
	}

	// A scan that matches nothing passes for the wrong reason. This test was
	// silently vacuous once already, when every component still used object
	// literals and the regex only understood literal attributes.
	if found == 0 {
		t.Fatal("scanned every .templ and found no named behavior — the regex has stopped matching the markup")
	}
	t.Logf("checked %d behavior references", found)
}

// rovingList is the single implementation of the arrow-key pattern behind Tabs,
// Rating and Radio. Losing it would silently strip keyboard navigation from all
// three, and no render test would notice.
func TestRovingListHelperExists(t *testing.T) {
	js := readUIJS(t)
	if !strings.Contains(js, "function rovingList(") {
		t.Fatal("rovingList() is gone from ui.js; Tabs, Rating and Radio depend on it")
	}
	for _, key := range []string{"ArrowLeft", "ArrowRight", "ArrowUp", "ArrowDown", "Home", "End"} {
		if !strings.Contains(js, key) {
			t.Errorf("rovingList no longer handles %s", key)
		}
	}
}

// listboxNav is the single implementation of aria-activedescendant navigation
// behind Combobox and Command. Two copies would drift, and the Combobox fixes
// (retyping, Escape propagation) would reach only one of them.
func TestListboxNavHelperExists(t *testing.T) {
	js := readUIJS(t)
	if !strings.Contains(js, "function listboxNav(") {
		t.Fatal("listboxNav() is missing from malphas.js")
	}
	if !strings.Contains(js, "...listboxNav(") {
		t.Error("no behavior spreads listboxNav(); Combobox should")
	}
}
