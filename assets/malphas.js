// Malphas design-system Alpine behaviors.
//
// Loaded with `defer` BEFORE alpine.min.js: Alpine dispatches alpine:init the
// moment it loads, so a listener attached afterwards never runs.

// Must match the key used by the inline anti-FOUC script in layout.Base.
const MALPHAS_THEME_KEY = "malphas-theme";

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
      return Array.from(this.$el.querySelectorAll("[data-roving-item]")).filter(
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

document.addEventListener("alpine:init", () => {
  Alpine.data("dropdown", () => ({
    open: false,
    toggle() { this.open = !this.open; },
    close() { this.open = false; },
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
      const input = this.$el.querySelectorAll('input[type="radio"]')[i];
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
