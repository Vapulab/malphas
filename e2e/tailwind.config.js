/** @type {import('tailwindcss').Config} */
// Runs from e2e/. Malphas is the tree above, so the paths are relative.
module.exports = {
  presets: [require("../theme/preset.js")],
  content: [
    "./*.templ",
    "../components/**/*.templ",
    "../layout/**/*.templ",
    "../showcase/**/*.templ",
  ],
};
