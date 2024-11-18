const colors = require('tailwindcss/colors')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./view/*.html"],
  theme: {
    colors : {
      Navy : "#4A628A",
      Bone: "#DFF2EB",
    },
    extend: {
           colors: {
                transparent: 'transparent',
                current: 'currentColor',
                black: colors.black,
             green : colors.green,
                white: colors.white,
                emerald: colors.emerald,
                indigo: colors.indigo,
                yellow: colors.yellow,
                stone: colors.warmGray,
                sky: colors.lightBlue,
                neutral: colors.trueGray,
                gray: colors.coolGray,
                slate: colors.blueGray,
                lime: colors.lime,
                rose: colors.rose,
            },  
    },
  },
  plugins: [
    require('daisyui'),
  ],
  daisyui: {
    themes: ["cupcake"],
  },
}

