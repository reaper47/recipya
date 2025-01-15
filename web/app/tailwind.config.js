// noinspection JSUnresolvedVariable

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../components/*.templ"],
  daisyui: {
    darkTheme: "forest",
    themes: true,
  },
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
}
