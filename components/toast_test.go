package components_test

import (
	"strings"
	"testing"

	"github.com/darkanum/vapulab/malphas/components"
)

// Each fill pairs with its own foreground token. A hardcoded text-white was
// fine until the dark theme lightened primary and danger: white on red-400 is
// 2.77:1. The -fg tokens move with their fill in both themes.
func TestToastTextFollowsItsFill(t *testing.T) {
	for variant, want := range map[string]string{
		"success": "bg-primary text-primary-fg",
		"error":   "bg-danger text-danger-fg",
	} {
		html := renderToString(t, components.Toast(components.ToastProps{Message: "x", Variant: variant}))
		if strings.Contains(html, "text-white") {
			t.Errorf("%s: hardcoded text-white ignores the theme, got: %s", variant, html)
		}
		for _, class := range strings.Fields(want) {
			if !strings.Contains(html, class) {
				t.Errorf("%s: missing %s, got: %s", variant, class, html)
			}
		}
	}
}
