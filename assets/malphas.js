// Malphas design-system Alpine behaviors.
//
// Loaded with `defer` BEFORE alpine.min.js: Alpine dispatches alpine:init the
// moment it loads, so a listener attached afterwards never runs.

// Must match the key used by the inline anti-FOUC script in layout.Base.
const MALPHAS_THEME_KEY = "malphas-theme";

// A note on $root vs $el, because getting it wrong fails silently:
//
// $el is *contextual*. Inside a method invoked from a handler bound to a child
// element, $el is that child — not the component. So this.$el.querySelectorAll
// in a method called from an <input>'s x-on:input searches inside the input and
// finds nothing, and the component quietly does nothing at all.
//
// Any lookup that means "everything in this component" uses $root. $el is only
// correct when the current element really is the subject, as when dispatching
// an event from it.

// rovingList implements the roving tabindex pattern shared by Tabs, Rating and
// Radio: exactly one item is tabbable, and the arrow keys move focus and
// selection between items.
//
// Implemented once on purpose. These three are the components that lost the
// browser's built-in keyboard handling when we chose to drive state from
// Alpine, and three separate implementations would be three chances to get it
// wrong. If arrow keys break anywhere, they break here.
//
// Items are located by the data-roving-item attribute (whose value is the item
// index) inside the component root, so the markup stays declarative and no
// per-item x-ref bookkeeping is needed. data-roving-disabled opts an item out.
function rovingList({ orientation = "horizontal", loop = true, active = 0 } = {}) {
  const keys = {
    horizontal: { prev: ["ArrowLeft"], next: ["ArrowRight"] },
    vertical: { prev: ["ArrowUp"], next: ["ArrowDown"] },
    both: { prev: ["ArrowLeft", "ArrowUp"], next: ["ArrowRight", "ArrowDown"] },
  }[orientation];

  return {
    active,

    // Disabled items are skipped rather than focused, so a disabled tab never
    // becomes a dead end for keyboard users.
    _rovingItems() {
      return Array.from(this.$root.querySelectorAll("[data-roving-item]")).filter(
        (el) => !el.hasAttribute("data-roving-disabled"),
      );
    },

    isActive(i) {
      return this.active === i;
    },

    // Exactly one item is reachable with Tab; the arrows do the rest.
    tabindex(i) {
      return this.active === i ? 0 : -1;
    },

    // Components override onRovingSelect to mirror the change elsewhere, e.g.
    // into a hidden input.
    select(i) {
      this.active = i;
      if (typeof this.onRovingSelect === "function") this.onRovingSelect(i);
    },

    _focusAt(position) {
      const items = this._rovingItems();
      const el = items[position];
      if (!el) return;
      this.select(Number(el.dataset.rovingItem));
      el.focus();
    },

    _move(delta) {
      const items = this._rovingItems();
      if (!items.length) return;
      const current = items.findIndex(
        (el) => Number(el.dataset.rovingItem) === this.active,
      );
      let next = current + delta;
      if (next < 0) next = loop ? items.length - 1 : 0;
      if (next >= items.length) next = loop ? 0 : items.length - 1;
      this._focusAt(next);
    },

    onKeydown(e) {
      if (keys.prev.includes(e.key)) {
        e.preventDefault();
        this._move(-1);
      } else if (keys.next.includes(e.key)) {
        e.preventDefault();
        this._move(1);
      } else if (e.key === "Home") {
        e.preventDefault();
        this._focusAt(0);
      } else if (e.key === "End") {
        e.preventDefault();
        this._focusAt(this._rovingItems().length - 1);
      }
    },
  };
}

// foldText makes filtering case- and accent-insensitive, so "sao" finds
// "São Paulo" — in Portuguese, users routinely skip the accents when typing.
function foldText(s) {
  return s.normalize("NFD").replace(/[̀-ͯ]/g, "").toLowerCase();
}

// listboxNav is the aria-activedescendant pattern shared by Combobox and
// Command: focus stays in a text field while the arrows move a highlight over
// [role="option"] elements, each carrying data-option-index and data-label.
// Implemented once for the same reason as rovingList — and because the
// Combobox already paid for its edge cases (retyping, Unicode folding).
//
// The spreading behavior owns `query`. It may override _filterActive() to
// show the whole list until the user starts typing.
function listboxNav() {
  return {
    active: -1,
    // Folding normalises Unicode, and matches() runs for every option on every
    // keystroke. Labels never change, so each is folded once; the query is
    // folded once per distinct value.
    _folded: new Map(),
    _foldedQuery: { raw: null, folded: "" },

    _fold(text) {
      let folded = this._folded.get(text);
      if (folded === undefined) {
        folded = foldText(text);
        this._folded.set(text, folded);
      }
      return folded;
    },

    _query() {
      if (this._foldedQuery.raw !== this.query) {
        this._foldedQuery = { raw: this.query, folded: foldText(this.query.trim()) };
      }
      return this._foldedQuery.folded;
    },

    _filterActive() {
      return true;
    },

    matches(optionLabel) {
      if (!this._filterActive()) return true;
      const q = this._query();
      return !q || this._fold(optionLabel).includes(q);
    },

    hasMatches() {
      return this._options().some((el) => this.matches(el.dataset.label));
    },

    _options() {
      return Array.from(this.$root.querySelectorAll('[role="option"]'));
    },

    // Options the arrows can land on: visible under the filter and enabled.
    _navigable() {
      return this._options().filter(
        (el) =>
          el.getAttribute("aria-disabled") !== "true" && this.matches(el.dataset.label),
      );
    },

    activeId() {
      // null, not "": Alpine drops the attribute for null, and an empty
      // aria-activedescendant is an invalid reference.
      return this.active >= 0 ? this._options()[this.active]?.id ?? null : null;
    },

    _highlight(el) {
      this.active = el ? Number(el.dataset.optionIndex) : -1;
      // The list may have just been shown; scrolling a display:none element
      // does nothing.
      if (el) this.$nextTick(() => el.scrollIntoView({ block: "nearest" }));
    },

    _move(delta) {
      const items = this._navigable();
      if (!items.length) return;
      const current = items.findIndex(
        (el) => Number(el.dataset.optionIndex) === this.active,
      );
      const next =
        current === -1
          ? delta > 0 ? 0 : items.length - 1
          : (current + delta + items.length) % items.length;
      this._highlight(items[next]);
    },
  };
}

// --- Non-modal floating panels ----------------------------------------------
//
// floatingPanel is the non-modal half of the overlay pair, for a trigger that
// opens a panel of free content (Popover). Opening moves focus in, Escape
// closes and returns focus to the trigger, and focus leaving the component
// closes it without pulling focus back. No trap, no scroll lock, no inert:
// the page behind stays usable, which is the whole point of not being modal.
//
// Escape calls stopPropagation because Modal and Drawer close on Escape at
// window: a popover inside a dialog must close alone (the Combobox lesson).
// Expects x-ref="trigger" and x-ref="panel" inside the component root.
function floatingPanel() {
  return {
    open: false,

    show() {
      if (this.open) return;
      this.open = true;
      this.$nextTick(() => {
        const panel = this.$refs.panel;
        if (!panel) return;
        (panel.querySelector(FOCUSABLE) || panel).focus();
      });
    },

    hide({ restoreFocus = true } = {}) {
      if (!this.open) return;
      this.open = false;
      if (restoreFocus) this.$refs.trigger?.focus();
    },

    toggle() {
      this.open ? this.hide() : this.show();
    },

    onPanelKeydown(e) {
      if (e.key !== "Escape" || !this.open) return;
      e.preventDefault();
      e.stopPropagation();
      this.hide();
    },

    // relatedTarget is where focus is going; null means it left the document
    // (another window, or a click on nothing focusable), which click.outside
    // already handles.
    onFocusOut(e) {
      if (e.relatedTarget && !this.$root.contains(e.relatedTarget)) {
        this.hide({ restoreFocus: false });
      }
    },
  };
}

// --- Modal overlays ---------------------------------------------------------
//
// Making an overlay modal takes three things, and the last two are the ones
// usually missing:
//
//   focus trap   without it, Tab walks the content behind — unusable by keyboard
//   scroll lock  without it, the page scrolls under the overlay
//   inert        without it, a screen reader's virtual cursor reads the whole
//                page as if no dialog were open
//
// Trapping Tab alone gets you a modal for keyboard users and a non-modal for
// screen reader users. All three live here so Modal and Drawer cannot
// implement two thirds of it each.

const FOCUSABLE = [
  "a[href]",
  "button:not([disabled])",
  "input:not([disabled]):not([type=hidden])",
  "select:not([disabled])",
  "textarea:not([disabled])",
  '[tabindex]:not([tabindex="-1"])',
].join(",");

// Nested overlays share one lock, so closing the inner one does not hand the
// page back while the outer is still open.
let scrollLockDepth = 0;
let savedScrollBehavior = "";

function lockBodyScroll() {
  if (scrollLockDepth++ > 0) return;
  const root = document.documentElement;
  // Compensating for the scrollbar keeps the page from jumping sideways as it
  // disappears. Note: overflow:hidden does not stop body scroll on iOS Safari;
  // the position:fixed workaround costs a scroll-restore dance and a repaint,
  // and is not worth it here.
  const gutter = window.innerWidth - root.clientWidth;
  savedScrollBehavior = root.style.scrollBehavior;
  root.style.overflow = "hidden";
  if (gutter > 0) root.style.paddingRight = `${gutter}px`;
}

function unlockBodyScroll() {
  if (--scrollLockDepth > 0) return;
  scrollLockDepth = 0;
  const root = document.documentElement;
  root.style.overflow = "";
  root.style.paddingRight = "";
  root.style.scrollBehavior = savedScrollBehavior;
}

// modalOverlay is spread into a component that owns an `open` flag and a panel
// marked with x-ref="panel".
function modalOverlay() {
  return {
    open: false,
    _returnFocusTo: null,
    _inerted: [],

    openOverlay() {
      if (this.open) return;
      // Remember who opened us before the DOM changes under our feet.
      this._returnFocusTo = document.activeElement;
      this.open = true;
      // The panel lives inside x-if, so it does not exist until Alpine has
      // rendered this tick.
      this.$nextTick(() => {
        const panel = this.$refs.panel;
        if (!panel) return;
        lockBodyScroll();
        this._applyInert(panel);
        (panel.querySelector(FOCUSABLE) || panel).focus();
      });
    },

    closeOverlay() {
      if (!this.open) return;
      this.open = false;
      unlockBodyScroll();
      this._releaseInert();
      // Returning focus is what keeps a keyboard user from being dumped at the
      // top of the document every time they dismiss something.
      if (this._returnFocusTo && document.contains(this._returnFocusTo)) {
        this._returnFocusTo.focus();
      }
      this._returnFocusTo = null;
    },

    // Everything that is not an ancestor of the panel becomes unreachable —
    // to Tab, to the pointer, and to assistive tech.
    _applyInert(panel) {
      this._inerted = [];
      for (const el of Array.from(document.body.children)) {
        if (el.contains(panel) || el.hasAttribute("inert")) continue;
        el.setAttribute("inert", "");
        this._inerted.push(el);
      }
    },

    _releaseInert() {
      for (const el of this._inerted) el.removeAttribute("inert");
      this._inerted = [];
    },

    // Bound to keydown.tab on the panel: wrap at both ends so focus cannot
    // leave. inert already blocks most escapes, but not the browser chrome
    // round-trip, and not browsers without inert support.
    onTab(e) {
      const panel = this.$refs.panel;
      if (!panel) return;
      const items = Array.from(panel.querySelectorAll(FOCUSABLE)).filter(
        (el) => el.offsetParent !== null,
      );
      if (!items.length) return;
      const first = items[0];
      const last = items[items.length - 1];
      if (e.shiftKey && document.activeElement === first) {
        e.preventDefault();
        last.focus();
      } else if (!e.shiftKey && document.activeElement === last) {
        e.preventDefault();
        first.focus();
      }
    },
  };
}

document.addEventListener("alpine:init", () => {
  Alpine.data("modal", () => ({ ...modalOverlay() }));
  Alpine.data("drawer", () => ({ ...modalOverlay() }));

  Alpine.data("carousel", () => ({
    // Respecting the preference matters most in a component whose entire job
    // is to move. matchMedia rather than CSS because the scroll is scripted.
    get _behavior() {
      return window.matchMedia("(prefers-reduced-motion: reduce)").matches
        ? "auto"
        : "smooth";
    },

    _scrollBy(direction) {
      const track = this.$refs.track;
      if (!track) return;
      const slide = track.firstElementChild;
      const step = slide ? slide.getBoundingClientRect().width + 16 : track.clientWidth;
      track.scrollBy({ left: direction * step, behavior: this._behavior });
    },

    prev() { this._scrollBy(-1); },
    next() { this._scrollBy(1); },
  }));

  Alpine.data("otp", (length = 6) => ({
    length,
    value: "",

    _boxes() {
      return Array.from(this.$root.querySelectorAll("[data-otp-index]"));
    },

    _sync() {
      this.value = this._boxes().map((b) => b.value).join("");
    },

    onInput(e, index) {
      // Keep one character per box: a fast typist or an IME can deliver more.
      const typed = e.target.value.replace(/\s/g, "");
      e.target.value = typed.slice(-1);
      this._sync();
      if (e.target.value && index < this.length - 1) {
        this._boxes()[index + 1].focus();
      }
    },

    onKeydown(e, index) {
      const boxes = this._boxes();
      if (e.key === "Backspace" && !e.target.value && index > 0) {
        // Stepping back on an empty box is what makes correcting a typo feel
        // like one field instead of six.
        e.preventDefault();
        boxes[index - 1].focus();
        boxes[index - 1].value = "";
        this._sync();
      } else if (e.key === "ArrowLeft" && index > 0) {
        e.preventDefault();
        boxes[index - 1].focus();
      } else if (e.key === "ArrowRight" && index < this.length - 1) {
        e.preventDefault();
        boxes[index + 1].focus();
      }
    },

    // Pasting the whole code is how most people enter it; without this it
    // would land entirely in one box.
    onPaste(e) {
      const digits = (e.clipboardData || window.clipboardData)
        .getData("text")
        .replace(/\D/g, "")
        .slice(0, this.length);
      const boxes = this._boxes();
      digits.split("").forEach((d, i) => (boxes[i].value = d));
      this._sync();
      boxes[Math.min(digits.length, this.length - 1)].focus();
    },
  }));

  Alpine.data("popover", () => ({ ...floatingPanel() }));

  // Timers rather than CSS :hover, because the card must survive the gap
  // between link and card, and open on keyboard focus too. One pending timer
  // at a time: entering cancels a pending close and vice versa.
  Alpine.data("hoverCard", (openDelay = 400, closeDelay = 200) => ({
    open: false,
    _timer: null,

    _schedule(fn, ms) {
      clearTimeout(this._timer);
      this._timer = setTimeout(fn, ms);
    },

    enter() {
      if (this.open) clearTimeout(this._timer);
      else this._schedule(() => (this.open = true), openDelay);
    },

    leave() {
      this._schedule(() => (this.open = false), closeDelay);
    },

    onKeydown(e) {
      if (e.key !== "Escape" || !this.open) return;
      e.stopPropagation();
      clearTimeout(this._timer);
      this.open = false;
    },
  }));

  // The menu-button pattern, and the mirror image of modalOverlay: this
  // overlay is NOT modal. Arrows move between items and Escape returns focus,
  // but Tab must leave and close — trapping focus here would be as wrong as
  // failing to trap it in a dialog.
  Alpine.data("dropdownMenu", () => ({
    ...rovingList({ orientation: "vertical", active: 0 }),
    open: false,

    toggle() {
      this.open ? this.close() : this.show();
    },

    show() {
      this.open = true;
      this.active = 0;
      // $root, not $el: show() is called from the trigger's click handler, so
      // $el is the button — which contains no menu items. See the note on
      // $root above.
      //
      // The menu lives inside x-if, so it does not exist until Alpine has
      // rendered it. One tick is not always enough for a nested template, so
      // this waits for the element rather than assuming.
      this.$nextTick(() => {
        const focusFirst = () => {
          const first = this.$root.querySelector(
            "[data-roving-item]:not([data-roving-disabled])",
          );
          if (first) first.focus();
          else requestAnimationFrame(focusFirst);
        };
        focusFirst();
      });
    },

    // restoreFocus is false when Tab closes the menu: the browser is already
    // moving focus onward, and yanking it back to the trigger would fight it.
    close({ restoreFocus = true } = {}) {
      if (!this.open) return;
      this.open = false;
      if (restoreFocus) this.$refs.trigger?.focus();
    },

    onMenuKeydown(e) {
      if (e.key === "Tab") {
        this.close({ restoreFocus: false });
        return;
      }
      if (e.key === "Escape") {
        e.preventDefault();
        this.close();
        return;
      }
      this.onKeydown(e);
    },
  }));

  // The combobox pattern with a list-autocomplete listbox. Unlike the dropdown,
  // focus never moves into the list: the text field keeps it, and the
  // highlighted option is only pointed at via aria-activedescendant. That is
  // what lets the user keep typing while arrowing.
  //
  // value/label are the committed selection; query is what the field shows.
  // They only diverge while the user is typing, and commit() reconciles them.
  //
  // filtering is set by typing and cleared by any commit. It, not a comparison
  // of query against label, decides whether the list is filtered: retyping the
  // current label is still a search, and must not reset the list and move the
  // highlight to the first option.
  Alpine.data("combobox", (value = "", label = "") => ({
    ...listboxNav(),
    open: false,
    filtering: false,
    value,
    label,
    query: label,

    // Until the user types, the whole list is offered — opening a filled
    // combobox should not show a single option.
    _filterActive() {
      return this.filtering;
    },

    isSelected(v) {
      return this.value !== "" && this.value === v;
    },

    // Opening lands on the current selection, so ↓ from a filled field
    // continues from where the user already is.
    _openAt(fallback) {
      this.open = true;
      const items = this._navigable();
      const selected = items.find((el) => el.dataset.value === this.value);
      this._highlight(selected || (fallback === "last" ? items[items.length - 1] : items[0]));
    },

    close() {
      this.open = false;
      this.active = -1;
    },

    toggle() {
      if (this.open) this.close();
      else this._openAt("first");
      this.$root.querySelector('[role="combobox"]').focus();
    },

    // Typing highlights the first match, so Enter picks the top result.
    onInput() {
      this.filtering = true;
      this.open = true;
      this._highlight(this._navigable()[0]);
    },

    // Written to the element directly as well as to query: x-model applies
    // its update a tick later, and an Enter-triggered submit validates
    // `required` against the field before then.
    _show(text) {
      this.query = text;
      this.filtering = false;
      this.$root.querySelector('[role="combobox"]').value = text;
    },

    // The hidden input is what the form submits, and what hx-trigger="change"
    // listens to — the text field never fires a change worth acting on.
    _set(value, label) {
      const changed = value !== this.value;
      this.value = value;
      this.label = label;
      this._show(label);
      this.$refs.hidden.value = value;
      if (changed) {
        this.$refs.hidden.dispatchEvent(new Event("change", { bubbles: true }));
      }
    },

    choose(i) {
      const el = this._options()[i];
      if (!el || el.getAttribute("aria-disabled") === "true") return;
      this._set(el.dataset.value, el.dataset.label);
      this.close();
    },

    // On blur. Free text is never kept: an exact label match is accepted
    // (typing the whole name should not require picking it), an empty field
    // clears, and anything else reverts to the last real selection.
    commit() {
      this.close();
      const q = this._query();
      if (!q) {
        this._set("", "");
        return;
      }
      const exact = this._options().find(
        (el) =>
          el.getAttribute("aria-disabled") !== "true" && this._fold(el.dataset.label) === q,
      );
      if (exact) this._set(exact.dataset.value, exact.dataset.label);
      else this._show(this.label);
    },

    onKeydown(e) {
      switch (e.key) {
        case "ArrowDown":
          e.preventDefault();
          // Alt+↓ only opens, per the pattern; plain ↓ opens and highlights.
          if (!this.open) {
            if (e.altKey) this.open = true;
            else this._openAt("first");
          } else this._move(1);
          break;
        case "ArrowUp":
          e.preventDefault();
          if (!this.open) this._openAt("last");
          else this._move(-1);
          break;
        case "Enter":
          if (this.open && this.active >= 0) {
            e.preventDefault();
            this.choose(this.active);
          } else {
            // Without a highlight, Enter submits the form — and it does so
            // before blur, so the text is reconciled here first. Otherwise the
            // field would show one thing while the form sent another.
            this.commit();
          }
          break;
        case "Tab":
          // Leaving after arrowing to an option takes it, as Enter would.
          // No preventDefault: Tab still moves on.
          if (this.open && this.active >= 0) this.choose(this.active);
          break;
        case "Escape":
          // First Escape closes; a second one clears, as the pattern asks.
          // stopPropagation because Modal and Drawer close on Escape at
          // window: closing the list must not close the dialog around it.
          if (this.open) {
            e.preventDefault();
            e.stopPropagation();
            this.close();
          } else if (this.query) {
            e.preventDefault();
            e.stopPropagation();
            this._set("", "");
          }
          break;
      }
    },
  }));

  // Reads the theme the inline script already applied, so the toggle's initial
  // state matches what is on screen. "" means "follow the OS".
  Alpine.data("themeToggle", () => ({
    theme: document.documentElement.dataset.theme || "",

    get isDark() {
      if (this.theme) return this.theme === "dark";
      return window.matchMedia("(prefers-color-scheme: dark)").matches;
    },

    set(theme) {
      this.theme = theme;
      if (theme) {
        document.documentElement.dataset.theme = theme;
        localStorage.setItem(MALPHAS_THEME_KEY, theme);
      } else {
        delete document.documentElement.dataset.theme;
        localStorage.removeItem(MALPHAS_THEME_KEY);
      }
    },

    toggle() { this.set(this.isDark ? "light" : "dark"); },

    // Drop the override and go back to following the OS.
    useSystem() { this.set(""); },
  }));

  // --- Form controls -------------------------------------------------------
  // Each of these drives a visual built from spans while a real, sr-only input
  // carries focus, form submission and screen-reader interaction.

  Alpine.data("toggle", (checked = false) => ({
    on: checked,
  }));

  // The server paints the starting count; this only keeps it current. Typing
  // fires input, but a form reset and a bfcache or autofill restore change the
  // value silently, so those are synced here too.
  Alpine.data("charCounter", (count = 0) => ({
    count,

    init() {
      const field = this.$root.querySelector("textarea");
      if (!field) return;
      const sync = () => (this.count = field.value.length);
      sync();
      window.addEventListener("pageshow", sync);
      // reset fires before the browser restores the defaults; read after.
      field.form?.addEventListener("reset", () => setTimeout(sync));
    },
  }));

  Alpine.data("toggleButton", (pressed = false) => ({ pressed }));

  // Arrows move focus only, not the pressed state: unlike a radio group, a
  // toolbar of toggles must not change formatting just because focus passed
  // over it. So onRovingSelect is deliberately not defined.
  Alpine.data("toggleGroup", (pressed = [], multiple = false, active = 0) => ({
    ...rovingList({ orientation: "both", active }),
    pressed,

    isPressed(i) {
      return !!this.pressed[i];
    },

    // Single mode can still end with nothing pressed: pressing the pressed
    // item releases it.
    press(i) {
      const next = !this.pressed[i];
      if (!multiple) this.pressed = this.pressed.map(() => false);
      this.pressed[i] = next;
      this.active = i;
    },
  }));

  Alpine.data("checkbox", (checked = false, indeterminate = false) => ({
    checked,
    indeterminate,

    init() {
      // indeterminate has no HTML attribute — it exists only as a property, so
      // it has to be set here rather than rendered by the server.
      this.$refs.input.indeterminate = this.indeterminate;

      // Any deliberate interaction resolves the third state; leaving it on
      // would make the box lie about what will be submitted.
      this.$watch("checked", () => {
        if (!this.indeterminate) return;
        this.indeterminate = false;
        this.$refs.input.indeterminate = false;
      });
    },
  }));

  Alpine.data("radioGroup", (active = 0) => ({
    ...rovingList({ orientation: "both", active }),

    // Arrow keys move selection, not just focus, which is what a radio group
    // does natively. Checking the real input keeps the form value in step.
    onRovingSelect(i) {
      const input = this.$root.querySelectorAll('input[type="radio"]')[i];
      if (input && !input.disabled) input.checked = true;
    },
  }));

  // --- Time ----------------------------------------------------------------

  // Formatting is duplicated in formatDuration() in components/helpers.go so
  // the server-rendered clock and the first tick agree. Keep the two in step.
  Alpine.data("countdown", (seconds = 0, format = "hh:mm:ss", onEnd = "") => ({
    remaining: Math.max(0, seconds),
    _timer: null,

    init() {
      if (this.remaining <= 0) return;
      // setInterval, then cleared on destroy: an HTMX swap can remove this
      // element while the clock is still running, and an orphaned interval
      // would keep firing against a detached node.
      this._timer = setInterval(() => this.tick(), 1000);
    },

    destroy() {
      if (this._timer) clearInterval(this._timer);
    },

    tick() {
      this.remaining = Math.max(0, this.remaining - 1);
      if (this.remaining > 0) return;
      clearInterval(this._timer);
      this._timer = null;
      if (onEnd) this.$nextTick(() => this.$el.dispatchEvent(new CustomEvent("countdown-end", { bubbles: true })));
    },

    get display() {
      const pad = (n) => String(n).padStart(2, "0");
      const t = this.remaining;
      if (format === "mm:ss") return `${pad(Math.floor(t / 60))}:${pad(t % 60)}`;
      return `${pad(Math.floor(t / 3600))}:${pad(Math.floor((t % 3600) / 60))}:${pad(t % 60)}`;
    },
  }));

  // --- Disclosure and navigation -------------------------------------------

  Alpine.data("collapse", (open = false) => ({
    open,
    toggle() { this.open = !this.open; },
  }));

  Alpine.data("accordion", (open = [], exclusive = false) => ({
    openPanels: [...open],

    isOpen(i) { return !!this.openPanels[i]; },

    toggle(i) {
      const next = !this.openPanels[i];
      // Rebuilding the array rather than mutating in place keeps Alpine's
      // reactivity simple and makes exclusive mode a single expression.
      this.openPanels = this.openPanels.map((was, j) =>
        j === i ? next : exclusive ? false : was,
      );
    },
  }));

  Alpine.data("tabs", (active = 0) => ({
    ...rovingList({ orientation: "horizontal", active }),
  }));

  // Range keeps its native input; this only mirrors the value into the label.
  Alpine.data("rangeValue", (value = 0) => ({ value }));

  Alpine.data("swap", (active = false) => ({
    active,
    toggle() { this.active = !this.active; },
  }));

  Alpine.data("tooltip", () => ({
    open: false,
    show() { this.open = true; },
    hide() { this.open = false; },
  }));

  Alpine.data("rating", (value = 0, max = 5) => ({
    ...rovingList({ orientation: "horizontal", active: Math.max(0, value - 1) }),
    value,
    max,

    // Stars are 1-based; rovingList indexes from 0.
    onRovingSelect(i) {
      this.value = i + 1;
    },
  }));
});
