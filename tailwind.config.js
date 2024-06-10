/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./frontend/**/*.templ"],
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["bumblebee"],
  },
};
