/** @type {import('tailwindcss').Config} */
// Roda de e2e/. O Malphas e a arvore acima, entao os caminhos sao relativos.
module.exports = {
  presets: [require("../theme/preset.js")],
  content: [
    "./*.templ",
    "../components/**/*.templ",
    "../layout/**/*.templ",
    "../showcase/**/*.templ",
  ],
};
