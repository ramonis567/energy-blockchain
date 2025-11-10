/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        blueColor: "#002F5F",
        greenColor: "#78BE20",
        highlightColor: "#00ADEF",
        surfaceColor: "#F5F7FA",
        textDark: "#222222",
      },
    },
  },
  plugins: [],
};
