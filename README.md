# Malphas

Design system em Go: componentes **templ** estilizados com **Tailwind**, com
comportamento em **Alpine** e prontos para **HTMX**. Sem Node/npm e sem CDN em
runtime.

São 71 componentes, de Button e Modal a Data Table, Command e Sidebar. O
catálogo vivo, com uma página por componente, está em
[malphas.vapulab.com](https://malphas.vapulab.com).

O Malphas é a camada de UI da [Vapulab](https://vapulab.com).

## Usando em um projeto

### 1. Instale o módulo

```bash
go get github.com/vapulab/malphas@latest
```

### 2. Aponte o Tailwind para o módulo

Tailwind **v3**, CLI standalone ou via npm. A v4 remove os `presets` em JS,
dos quais o Malphas depende.

O módulo fica no cache do Go, num caminho que muda de versão para versão e de
máquina para máquina. O `tailwind.config.js` pergunta esse caminho ao próprio
Go:

```js
const { execSync } = require("child_process");

// split/join em vez de regex: o parser de config do Tailwind standalone
// tropeça em /\\/g. As barras normais servem aos globs no Windows.
const malphasDir = execSync("go list -m -f {{.Dir}} github.com/vapulab/malphas")
  .toString()
  .trim()
  .split("\\")
  .join("/");

module.exports = {
  presets: [require(malphasDir + "/theme/preset.js")],
  content: [
    "./views/**/*.templ",
    malphasDir + "/components/**/*.templ",
    malphasDir + "/layout/**/*.templ",
  ],
};
```

O preset entrega cores, raio e fontes ao Tailwind **e** injeta
`theme/tokens.css` na camada base. O CSS de entrada só precisa das diretivas de
sempre:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

Para trocar a paleta, redefina as variáveis `--color-*` no seu CSS, depois do
`@tailwind base`. `theme/tokens.css` é a referência do que existe. Os dois
blocos escuros precisam andar juntos: o CSS não permite compartilhar um bloco
de declarações entre uma media query e um seletor.

### 3. Sirva os assets

O JavaScript dos componentes, a fonte e o favicon vão embutidos no binário:

```go
import "github.com/vapulab/malphas"

mux.Handle("/malphas-assets/", http.StripPrefix("/malphas-assets",
	http.FileServer(http.FS(malphas.Assets()))))
```

`layout.Base` também referencia o CSS gerado pelo Tailwind e os vendors do
projeto, que ficam a cargo de quem consome:

| Caminho | O que é |
|---|---|
| `/static/app.css` | a saída do Tailwind |
| `/static/vendor/htmx.min.js` | HTMX |
| `/static/vendor/alpine.min.js` | Alpine.js 3 |

### 4. Use os componentes

```go
import (
	"github.com/vapulab/malphas/components"
	"github.com/vapulab/malphas/layout"
)
```

```templ
@layout.Base(layout.BaseProps{Title: "Notas", Body: page()})

@components.Button(components.ButtonProps{}) {
	Emitir nota
}
```

Os `*_templ.go` gerados são commitados, então o módulo compila sem
`templ generate`.

## Como os componentes se comportam

### Formulários de verdade

Os controles de formulário são visual em Alpine sobre um `<input>` real com
`sr-only`. O input carrega foco, envio e leitura de tela; o Alpine cuida só da
aparência. Na prática, **eles funcionam dentro de `<form>` e de `hx-post` sem
nada extra.**

**Range é a exceção deliberada**: mantém o input nativo visível e estilizado.
Esconder e reimplementar arrasto por ponteiro, toque e `aria-valuenow` custaria
muito para um resultado pior que o do browser.

### Teclado

Seguindo os padrões WAI-ARIA:

| Componente | Teclas |
|---|---|
| Tabs | ← →, Home, End; abas desabilitadas são puladas |
| Radio, Rating | setas movem foco **e** seleção |
| Collapse, Accordion | Enter, Espaço |
| Tooltip | abre no foco, fecha com Escape |
| Toggle, Swap | Espaço |

Os componentes de lista compartilham um único `rovingList()` em
`assets/malphas.js`: uma implementação do roving tabindex, não uma por
componente.

### Sobreposições: modal e não-modal

| | Componentes | Foco |
|---|---|---|
| **Modal** | Modal, Drawer, Alert Dialog, Command (paleta) | preso dentro; Tab circula, Escape fecha e devolve |
| **Não-modal** | Dropdown | livre; setas navegam, Escape devolve, **Tab sai e fecha** |

Prender o foco num menu é tão errado quanto não prender num diálogo.

Os modais fazem as **três** coisas que tornam uma sobreposição realmente modal:

- **focus trap**: sem isso o Tab passeia pelo conteúdo atrás
- **scroll lock**: sem isso a página rola sob a sobreposição
- **inert**: sem isso o cursor virtual do leitor de tela lê a página inteira

O terceiro é o mais esquecido. Prender só o Tab entrega um modal para quem usa
teclado e um não-modal para quem usa leitor de tela.

### Tema claro e escuro

Por padrão o tema segue o sistema operacional (`prefers-color-scheme`). Uma
escolha explícita, gravada em `localStorage`, vence o sistema **nos dois
sentidos**: dá para forçar claro num SO escuro.

`layout.Base` aplica o tema com um `<script>` **inline e síncrono** antes do
stylesheet, para que a página já pinte no tema certo. A chave de
`localStorage` é a mesma em `layout/base.templ` e `assets/malphas.js`, e um
teste garante isso.

## Desenvolvimento

Go 1.25+ e [templ](https://templ.guide) na versão do `go.mod`.

```bash
go install github.com/a-h/templ/cmd/templ@v0.3.1020
templ generate
go vet ./...
go test ./...
```

Os testes daqui cobrem renderização, acessibilidade da marcação e o contrato
entre os `.templ` e o `malphas.js`: cada `x-data="nome"` precisa de um
`Alpine.data()` registrado. O comportamento no browser (teclado, foco,
`FormData`) é exercitado pela suíte chromedp da aplicação que hospeda o
catálogo.

### Mexendo no Malphas a partir de um app

Clone os dois lado a lado e crie um `go.work`, fora do controle de versão, na
raiz do app:

```
go 1.26

use (
	./app
	../malphas
)
```

O app passa a compilar contra o clone local, e o `go list -m` do
`tailwind.config.js` também aponta para ele. Assim o watch do Tailwind enxerga
as mudanças nos `.templ`. Para publicar, faça o merge aqui, crie a tag
`vX.Y.Z` e atualize o app com `go get github.com/vapulab/malphas@vX.Y.Z`.

## Licença

[MIT](LICENSE). A fonte Space Grotesk (`assets/fonts/`) é distribuída sob a
[SIL Open Font License 1.1](assets/fonts/OFL.txt).
