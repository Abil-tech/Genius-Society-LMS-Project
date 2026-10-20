/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        brand: {
          orange: 'rgb(var(--gs-orange) / <alpha-value>)',
          'orange-dark': 'rgb(var(--gs-orange-dark) / <alpha-value>)',
          peach: 'rgb(var(--gs-peach) / <alpha-value>)',
          navy: 'rgb(var(--gs-navy) / <alpha-value>)',
        },
        surface: 'rgb(var(--gs-surface) / <alpha-value>)',
        canvas: 'rgb(var(--gs-canvas-landing) / <alpha-value>)',
        muted: 'rgb(var(--gs-text-muted) / <alpha-value>)',
        line: 'rgb(var(--gs-border) / <alpha-value>)',
      },
      fontFamily: {
        sans: ['"Plus Jakarta Sans Variable"', 'Inter', 'system-ui', 'sans-serif'],
      },
      borderRadius: {
        card: '16px',
        pill: '9999px',
      },
      boxShadow: {
        card: '0 1px 2px rgb(11 36 71 / 0.04), 0 4px 16px rgb(11 36 71 / 0.04)',
      },
    },
  },
  plugins: [],
}
