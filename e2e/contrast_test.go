package e2e_test

import (
	"context"
	"math"
	"strconv"
	"strings"
	"testing"
)

// relativeLuminance implementa a fórmula de luminância relativa da WCAG 2.x:
// https://www.w3.org/TR/WCAG21/#dfn-relative-luminance
func relativeLuminance(r, g, b int) float64 {
	channel := func(c int) float64 {
		cs := float64(c) / 255
		if cs <= 0.03928 {
			return cs / 12.92
		}
		return math.Pow((cs+0.055)/1.055, 2.4)
	}
	return 0.2126*channel(r) + 0.7152*channel(g) + 0.0722*channel(b)
}

// contrastRatio segue a mesma norma: (L1+0.05)/(L2+0.05), L1 sendo a cor mais
// clara. É simétrica de propósito — o par (texto, fundo) dá a mesma razão que
// o par (fundo, texto), então não importa qual das duas cores é qual.
func contrastRatio(a, b [3]int) float64 {
	la, lb := relativeLuminance(a[0], a[1], a[2]), relativeLuminance(b[0], b[1], b[2])
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

// parseRGB lê o formato que getComputedStyle devolve: "rgb(r, g, b)" ou
// "rgba(r, g, b, a)". O canal alfa, quando presente, é ignorado de propósito —
// os elementos que este arquivo mede (texto, preenchimento sólido) são opacos;
// mediar alfa exigiria compor contra o fundo real, que é mais do que o teste
// precisa provar.
func parseRGB(t *testing.T, s string) [3]int {
	t.Helper()
	inner := strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(s, "rgba("), "rgb("), ")")
	parts := strings.Split(inner, ",")
	if len(parts) < 3 {
		t.Fatalf("não foi possível interpretar a cor %q", s)
	}
	var out [3]int
	for i := 0; i < 3; i++ {
		n, err := strconv.Atoi(strings.TrimSpace(parts[i]))
		if err != nil {
			t.Fatalf("não foi possível interpretar a cor %q: %v", s, err)
		}
		out[i] = n
	}
	return out
}

// openPageForcedDark abre path com o tema explicitamente escuro, gravado em
// localStorage antes da navegação — não com prefers-color-scheme emulado.
// Um dado explícito em localStorage vence o SO nos dois sentidos (ver
// TestExplicitThemeBeatsTheOperatingSystem em theme_test.go), então isso
// funciona não importa o que outro teste tenha deixado no mesmo perfil de
// Chrome, ou qual seja o esquema de cores da máquina que roda a suíte.
// t.Cleanup remove a chave depois: sem isso, todo teste que rodasse depois
// herdaria tema escuro silenciosamente.
func openPageForcedDark(t *testing.T, path string) context.Context {
	t.Helper()
	ctx := openPage(t, "/components")
	t.Cleanup(func() {
		eval[bool](t, ctx, `(() => { localStorage.removeItem('malphas-theme'); return true; })()`)
	})
	eval[bool](t, ctx, `(() => { localStorage.setItem('malphas-theme', 'dark'); return true; })()`)
	navigate(t, ctx, path)
	settle(t, ctx)
	if got := eval[string](t, ctx, `document.documentElement.dataset.theme || ''`); got != "dark" {
		t.Fatalf("não foi possível forçar o tema escuro: data-theme=%q", got)
	}
	return ctx
}

// minAA é o piso da WCAG 2.x para texto normal (nível AA). Textos grandes e
// limites de UI têm piso 3:1, mas medir pelo mais exigente cobre os dois.
const minAA = 4.5

// TestDarkThemePrimaryTextMeetsAA guarda o contraste que rendering_test.go
// nunca cobriu: nunca em tema escuro, e nunca como razão de verdade — só
// cores comparadas entre si (mixed !== plain), o que não pega uma cor que
// contrasta tecnicamente diferente mas ainda ilegível. --color-primary no
// bloco escuro de tokens.css não tinha override antes desta correção: o valor
// claro (indigo-700) dá ~2.2:1 sobre a superfície escura, bem abaixo do piso
// AA. O índice do catálogo tem um link text-primary logo no topo, então ele
// já basta.
func TestDarkThemePrimaryTextMeetsAA(t *testing.T) {
	ctx := openPageForcedDark(t, "/components")

	textColor := eval[string](t, ctx, `getComputedStyle(document.querySelector('.text-primary')).color`)
	bgColor := eval[string](t, ctx, `getComputedStyle(document.body).backgroundColor`)

	ratio := contrastRatio(parseRGB(t, textColor), parseRGB(t, bgColor))
	if ratio < minAA {
		t.Errorf("text-primary sobre a superfície escura: %.2f:1, abaixo do piso AA de %.1f:1 (texto %s, fundo %s)",
			ratio, minAA, textColor, bgColor)
	}
}

// TestDarkThemePrimaryFillMeetsAA cobre o outro lado da mesma armadilha: subir
// --color-primary sozinho conserta text-primary e quebra todo preenchimento
// bg-primary text-primary-fg (Button, Badge sólido, Pagination, Steps done,
// Checkbox, Toggle) — o branco original cairia para ~2.9:1 sobre o azul mais
// claro. --color-primary-fg tem de andar junto, nos mesmos blocos de
// tokens.css. /components/button renderiza um Button de variante "primary"
// (bg-primary text-primary-fg) na primeira demo da página.
func TestDarkThemePrimaryFillMeetsAA(t *testing.T) {
	ctx := openPageForcedDark(t, "/components/button")

	textColor := eval[string](t, ctx, `getComputedStyle(document.querySelector('.text-primary-fg')).color`)
	bgColor := eval[string](t, ctx, `getComputedStyle(document.querySelector('.text-primary-fg')).backgroundColor`)

	ratio := contrastRatio(parseRGB(t, textColor), parseRGB(t, bgColor))
	if ratio < minAA {
		t.Errorf("text-primary-fg sobre bg-primary: %.2f:1, abaixo do piso AA de %.1f:1 (texto %s, fundo %s)",
			ratio, minAA, textColor, bgColor)
	}
}

// O mesmo buraco do primário, agora no danger: o bloco escuro não tinha
// override, e o vermelho claro (red-600) dá ~3.9:1 sobre a superfície escura.
// O erro do Field é texto que a pessoa precisa ler para corrigir o campo.
func TestDarkThemeDangerTextMeetsAA(t *testing.T) {
	ctx := openPageForcedDark(t, "/components/field")

	textColor := eval[string](t, ctx, `getComputedStyle(document.querySelector('#f-site-error')).color`)
	bgColor := eval[string](t, ctx, `getComputedStyle(document.body).backgroundColor`)

	ratio := contrastRatio(parseRGB(t, textColor), parseRGB(t, bgColor))
	if ratio < minAA {
		t.Errorf("text-danger sobre a superfície escura: %.2f:1, abaixo do piso AA de %.1f:1 (texto %s, fundo %s)",
			ratio, minAA, textColor, bgColor)
	}
}

// E o outro lado, como no primário: clarear --color-danger sozinho quebraria o
// branco do Button danger. /components/button tem um Button "danger".
func TestDarkThemeDangerFillMeetsAA(t *testing.T) {
	ctx := openPageForcedDark(t, "/components/button")

	textColor := eval[string](t, ctx, `getComputedStyle(document.querySelector('.text-danger-fg')).color`)
	bgColor := eval[string](t, ctx, `getComputedStyle(document.querySelector('.text-danger-fg')).backgroundColor`)

	ratio := contrastRatio(parseRGB(t, textColor), parseRGB(t, bgColor))
	if ratio < minAA {
		t.Errorf("text-danger-fg sobre bg-danger: %.2f:1, abaixo do piso AA de %.1f:1 (texto %s, fundo %s)",
			ratio, minAA, textColor, bgColor)
	}
}
