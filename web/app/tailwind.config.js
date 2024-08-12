// noinspection JSUnresolvedVariable

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../components/*.templ"],
  daisyui: {
    themes: true,
  },
  darkMode: ['class', [
    '[data-theme="dark"] &',
    '[data-theme="night"] &',
    '[data-theme="synthwave"] &',
    '[data-theme="halloween"] &',
    '[data-theme="forest"] &',
    '[data-theme="black"] &',
    '[data-theme="luxury"] &',
    '[data-theme="dracula"] &',
    '[data-theme="business"] &',
    '[data-theme="coffee"] &',
    '[data-theme="dim"] &',
    '[data-theme="sunset"] &',
  ]],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
}
