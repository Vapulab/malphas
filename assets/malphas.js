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
