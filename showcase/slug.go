package showcase

import "strings"

// categorySlug turns "Entrada de dados" into "entrada-de-dados" for anchors.
func categorySlug(category string) string {
	r := strings.NewReplacer(" ", "-", "ç", "c", "õ", "o", "ã", "a", "é", "e")
	return r.Replace(strings.ToLower(category))
}
