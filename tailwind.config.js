/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./view/*.html"],
  theme: {
    colors : {
      Navy : "#4A628A",
      Bone: "#DFF2EB",
    },
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],
}

