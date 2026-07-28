/** Malphas design-system Tailwind preset. Import from a project's tailwind.config.js. */
const rgb = (v) => `rgb(var(${v}) / <alpha-value>)`;

// Semantic colors pair a background with a readable foreground: bg-success +
// text-success-fg. Neutrals (surface, fg, muted, border) stay flat.
const semantic = (name) => ({
  DEFAULT: rgb(`--color-${name}`),
  fg: rgb(`--color-${name}-fg`),
});

module.exports = {
  // Mirrors the two dark blocks in tokens.css, so `dark:` utilities agree with
  // the tokens: an explicit data-theme wins over the OS in both directions.
  // Keep these three in sync — CSS cannot share the condition between them.
  darkMode: [
    "variant",
    [
      '&:is([data-theme="dark"] *)',
      '@media (prefers-color-scheme: dark) { &:not([data-theme="light"] *) }',
    ],
  ],

  theme: {
    extend: {
      colors: {
        primary: semantic("primary"),
        danger: semantic("danger"),
        success: semantic("success"),
        warning: semantic("warning"),
        info: semantic("info"),
        surface: rgb("--color-surface"),
        fg: rgb("--color-fg"),
        muted: rgb("--color-muted"),
        border: rgb("--color-border"),
      },
      borderRadius: { DEFAULT: "var(--radius)" },
      fontFamily: {
        sans: "var(--font-sans)",
        display: "var(--font-display)",
      },
    },
  },
};
