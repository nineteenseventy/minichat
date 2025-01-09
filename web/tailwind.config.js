const { join } = require('path');

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    join(__dirname, './index.html'),
    join(__dirname, 'src/**/*.{vue,js,ts,jsx,tsx}'),
  ],
  theme: {
    extend: {},
  },
  plugins: [import('tailwindcss-primeui')]
}

