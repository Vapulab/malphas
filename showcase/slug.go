package showcase

import "strings"

// categorySlug turns "Layout & utilities" into "layout-utilities" for
// anchors: letters and digits, words joined by single hyphens.
func categorySlug(category string) string {
	return strings.Join(strings.FieldsFunc(strings.ToLower(category), func(r rune) bool {
		return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9')
	}), "-")
}
