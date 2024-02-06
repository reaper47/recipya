// noinspection JSUnresolvedVariable

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../**/*.gohtml"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
}
