/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "/templates/**/*.{gohtml,html}"
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms')
  ],
}

