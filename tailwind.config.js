/** @type {import('tailwindcss').Config} */
module.exports = {
  mode: 'jit',
  content: [
    "./internal/templates/**/*.templ",
    "./internal/templates/*.templ",
    "./internal/templates/*.go",
    "./internal/templates/**/*.go",
    "./internal/templates/*/*.go",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
} 