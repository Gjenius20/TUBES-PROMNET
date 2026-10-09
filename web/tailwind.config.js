/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,jsx}'],
  theme: {
    extend: {
      colors: {
        paper: '#EDF1F5',
        ink: '#13203A',
        'ink-soft': '#4A5873',
        signal: '#2F5BFF',
        rule: '#CBD5E1',
        pass: '#1E8A5A',
        fail: '#C93A3A',
        warn: '#B7791F',
      },
      fontFamily: {
        display: ['"Bricolage Grotesque"', 'system-ui', 'sans-serif'],
        sans: ['"Instrument Sans"', 'system-ui', 'sans-serif'],
        mono: ['ui-monospace', 'SFMono-Regular', 'Menlo', 'Consolas', 'monospace'],
      },
    },
  },
  plugins: [],
}
