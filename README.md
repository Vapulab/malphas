# Malphas

A design system for Go: **templ** components styled with **Tailwind**, with
behavior in **Alpine** and ready for **HTMX**. No Node/npm, no CDN at runtime.

71 components, from Button and Modal to Data Table, Command and Sidebar. The
live catalog, with a page per component, is at
[malphas.vapulab.com](https://malphas.vapulab.com/components).

Built-in component text is in English by default, with Portuguese, Spanish,
Russian and Chinese available: see [Languages](#languages).

Malphas is the UI layer of [Vapulab](https://vapulab.com).

## Using it in a project

### 1. Install the module

```bash
go get github.com/vapulab/malphas@latest
```

### 2. Point Tailwind at the module

Malphas requires Tailwind **v3** (standalone CLI or npm). v4 dropped JS
`presets`, which Malphas relies on.

The module lives in Go's cache, at a path that changes per version and per
machine. `tailwind.config.js` asks Go for it:

```js
const { execSync } = require("child_process");

// split/join rather than a regex: the standalone Tailwind config parser
// trips on /\\/g. Forward slashes work for globs on Windows too.
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

The preset gives Tailwind the colors, radius and fonts **and** injects
`theme/tokens.css` into the base layer. Your input CSS only needs the usual
directives:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

To change the palette, redefine the `--color-*` variables in your CSS, after
`@tailwind base`. `theme/tokens.css` is the reference for what exists. If you
override dark colors, change both dark blocks, `@media (prefers-color-scheme: dark)`
and `:root[data-theme="dark"]`, the same way: CSS can't share one declaration
block between a media query and a selector.

### 3. Serve the assets

The components' JavaScript, the font and the favicon are embedded in the
binary:

```go
import "github.com/vapulab/malphas"

mux.Handle("/malphas-assets/", http.StripPrefix("/malphas-assets",
	http.FileServer(http.FS(malphas.Assets()))))
```

`layout.Base` also references the stylesheet Tailwind builds and the
project's vendored scripts, which are yours to serve:

| Path | What it is |
|---|---|
| `/static/app.css` | Tailwind's output |
| `/static/vendor/htmx.min.js` | HTMX |
| `/static/vendor/alpine.min.js` | Alpine.js 3 |

### 4. Use the components

```go
import (
	"github.com/vapulab/malphas/components"
	"github.com/vapulab/malphas/layout"
)
```

```templ
@layout.Base(layout.BaseProps{Title: "Invoices", Body: page()})

@components.Button(components.ButtonProps{}) {
	Issue invoice
}
```

The generated `*_templ.go` files are committed, so the module builds without
`templ generate`.

## Languages

Some components render built-in text: accessible names ("Close",
"Pagination"), empty states ("No results"), counters ("3 results",
"Question 2 of 4") and the Calendar's months and weekdays. That text ships in
five languages:

| `malphas.Locale` | Tag | Week starts on |
|---|---|---|
| `EN` (default) | `en` | Sunday |
| `PT` | `pt-BR` | Sunday |
| `ES` | `es` | Monday |
| `RU` | `ru` | Monday |
| `ZH` | `zh-Hans` | Monday |

The language travels in the request's `context.Context`, which templ already
passes to every component. Set it once, in a middleware, and every component
follows. `layout.Base` sets `<html lang>` from it too.

```go
func withLanguage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l, _ := malphas.ParseLocale(r.PathValue("lang")) // "en", "pt", "pt-BR"…; EN if unknown
		next.ServeHTTP(w, r.WithContext(malphas.WithLocale(r.Context(), l)))
	})
}
```

Malphas never decides a request's language: whether it comes from the URL, a
cookie or `Accept-Language` is up to your app.

A label passed as a prop always wins over the language's default:

```templ
@components.AlertDialog(components.AlertDialogProps{
	ID: "delete", Title: "Delete the project?",
	CancelLabel: "Keep it", ConfirmLabel: "Delete",
})
```

Plurals follow each language's rules (Russian has three forms, Chinese none).
The one counter built in the browser, the Message Scroller's "(3 new)", picks
its form with `Intl.PluralRules` for the language the server rendered it in.

`ParseLocale` takes a tag or a short code, and regional tags by their language
(`es-MX` is `ES`, `zh-CN` is `ZH`). Traditional Chinese (`zh-TW`, `zh-Hant`) is
not matched: the Chinese texts are Simplified.

## Upgrading from v0.2

v0.3.0 changes three things a v0.2 app sees without touching its code:

1. **Built-in text is English by default.** It used to be Portuguese. To keep
   it, set the language once per request:
   `ctx = malphas.WithLocale(ctx, malphas.PT)`.
2. **`<html lang>` follows that language** (`en` by default) instead of a fixed
   `pt-BR`.
3. **The Calendar's query parameters are `month` and `focus`**, no longer `mes`
   and `foco`. A handler that still reads `mes` gets the default month every
   time, so the arrows and PageUp/PageDown stop moving the calendar, with no
   error. Rename them where you build the Calendar's month and focus from the
   URL. The Calendar's `x-data` also takes a third argument, the week's first
   day, which the component fills in itself.

If you link to the catalog's all-in-one page, it moved from
`/components/todos` to `/components/all`; the old path redirects.

## How the components behave

### Real forms

Each form control is a real `<input>`, visually hidden with `sr-only`, under an
Alpine-driven visual. The input handles focus, submission and screen reader
output; Alpine only handles presentation. So **they work in a plain `<form>`
and with `hx-post` out of the box.**

**Range, Select and Date Picker are the deliberate exceptions**: they keep the
native control visible and styled. Hiding it and reimplementing dragging, touch,
pointer capture and the native pickers would be a lot of work to end up worse
than what the browser already does.

### Keyboard

Following the WAI-ARIA patterns:

| Component | Keys |
|---|---|
| Tabs | ← →, Home, End; disabled tabs are skipped |
| Radio, Rating | arrows move focus **and** selection |
| Collapse, Accordion | Enter, Space |
| Tooltip | opens on focus, closes on Escape |
| Toggle, Swap | Space |
| Calendar | arrows by day and week, Home/End to the start/end of the week (per locale), PageUp/PageDown by month |

The list components share a single `rovingList()` in `assets/malphas.js`: one
roving tabindex implementation, not one per component.

### Overlays: modal and non-modal

| | Components | Focus |
|---|---|---|
| **Modal** | Modal, Drawer, Alert Dialog, Command (palette) | trapped inside; Tab cycles, Escape closes and returns it |
| **Non-modal** | Dropdown | free; arrows navigate, Escape returns it, **Tab leaves and closes** |

Trapping focus in a menu is as wrong as not trapping it in a dialog.

Modals do the **three** things that make an overlay truly modal:

- **focus trap**: without it Tab wanders through the content behind
- **scroll lock**: without it the page scrolls under the overlay
- **inert**: without it the screen reader's virtual cursor reads the whole page

The third is the one people forget. Trapping only Tab gives a modal to
keyboard users and a non-modal to screen reader users.

### Light and dark

By default the theme follows the operating system (`prefers-color-scheme`).
An explicit choice, saved in `localStorage`, overrides the system **in either direction**:
you can force light on a dark OS.

`layout.Base` applies the theme with an **inline, synchronous** `<script>`
before the stylesheet, so the page paints in the right theme from the start.
The `localStorage` key is the same in `layout/base.templ` and
`assets/malphas.js`, and a test makes sure of it.

## Development

Go 1.25+ and [templ](https://templ.guide) at the version in `go.mod`.

```bash
go install github.com/a-h/templ/cmd/templ@v0.3.1020
templ generate   # always from the repo root
go vet ./...
go test ./...
```

These tests cover rendering, the accessibility of the markup, every locale's
strings (a locale missing a string doesn't compile, and a test checks that none
is empty), the catalog pages and the contract between the `.templ` files and
`malphas.js`: every `x-data="name"` needs a registered `Alpine.data()`.

### The catalog

`showcase/` is the living documentation: a registry with every component and
its examples, the pages that show them and an `http.Handler` that serves it
all under `/components`, always in English. It is what runs at
[malphas.vapulab.com/components](https://malphas.vapulab.com/components).
Adding a component to the catalog means adding a registry entry; the route,
the render tests and the showcase page pick it up automatically.

A host serving the catalog provides the page frame and adds `showcase/` to
Tailwind's `content`:

```go
mux.Handle("/components", showcase.Handler(frame))
mux.Handle("/components/", showcase.Handler(frame))
```

`showcase/` sits outside `components/` and `layout/` on purpose: the demos use
classes no real page needs, and a project that only uses the components never
compiles them.

### Browser tests

`e2e/` drives a real Chrome against the catalog, via
[chromedp](https://github.com/chromedp/chromedp), which speaks CDP straight
from Go, no Node. It covers what no render test reaches: arrows moving focus
and selection, focus trapped in modals and free in menus, Escape returning
focus, scroll lock, `inert`, the `FormData` the browser actually sends, the
demos that talk to the server over HTMX, and the Calendar's week in each
language.

It is a separate module (Go 1.26, because of chromedp), so none of it reaches
people who `go get` Malphas. It always tests the code in this tree.

```bash
cd e2e
tailwindcss -c tailwind.config.js -i input.css -o static/app.css --minify
go test ./...
# or, if Chrome isn't on PATH:
CHROME_PATH=/path/to/chrome go test ./...
```

Without Chrome or without `app.css`, the tests **skip** rather than fail, and
say why. Things worth knowing when you touch them:

- They use Chrome's **new headless**. In the old one `document.hasFocus()` is
  `false`, and the browser sends no focus events to an unfocused document:
  focus-driven components would look broken while being correct.
- Use `$root`, not `$el`, when an Alpine method looks for something "in this
  component". Inside a method called by a child's handler, `$el` is the child,
  and the lookup finds nothing without an error.
- **Wait for the condition, not for a time.** A `sleep` followed by a single
  check passes on a fast machine and fails on a slow runner with the component
  working: a timer, a transition or a frame takes what the CPU allows. Use
  `waitUntil` for what has to happen. A `sleep` only makes sense to prove
  something does *not* happen, and then it scales with `cpuThrottle`.

To find timing-sensitive tests before CI does, run with a slower Chrome:

```bash
E2E_CPU_THROTTLE=10 go test ./...   # the page's CPU 10x slower
```

### Working on Malphas from an app

Clone both side by side and create a `go.work`, outside version control, at
the app's root:

```
go 1.26

use (
	./app
	../malphas
)
```

The app then builds against the local clone, and the `go list -m` in
`tailwind.config.js` points at it too, so Tailwind's watch sees changes to the
`.templ` files. To release, merge into Malphas's `main`, tag `vX.Y.Z` and update the app with
`go get github.com/vapulab/malphas@vX.Y.Z`.

## License

[MIT](LICENSE). The Space Grotesk font (`assets/fonts/`) is distributed under
the [SIL Open Font License 1.1](assets/fonts/OFL.txt).
