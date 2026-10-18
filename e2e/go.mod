module github.com/vapulab/malphas/e2e

go 1.26

// O e2e e um modulo a parte para que o chromedp e o Go 1.26 que ele pede nao
// cheguem a quem so faz go get do Malphas. Ele sempre testa o Malphas desta
// arvore, nunca uma versao publicada.
replace github.com/vapulab/malphas => ../

require (
	github.com/a-h/templ v0.3.1020
	github.com/chromedp/cdproto v0.0.0-20260714215040-dc233986426f
	github.com/chromedp/chromedp v0.16.0
	github.com/vapulab/malphas v0.0.0
)

require (
	github.com/chromedp/sysutil v1.1.0 // indirect
	github.com/go-json-experiment/json v0.0.0-20260623181947-01eb4420fa68 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)
