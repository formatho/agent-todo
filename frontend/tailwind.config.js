/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class', // Enable dark mode via class selector
  theme: {
    extend: {
      // Performance optimized color palette
      colors: {
        primary: {
          50: '#EFF6FF',
          100: '#DBEAFE',
          200: '#BFDBFE',
          300: '#93C5FD',
          400: '#60A5FA',
          500: '#3B82F6',
          600: '#2563EB',
          700: '#1D4ED8',
          800: '#1E40AF',
          900: '#1E3A8A',
        },
        optimized: {
          // Optimized colors for better contrast and performance
          text: {
            primary: '#111827',
            secondary: '#374151',
            muted: '#6B7280',
            subtle: '#9CA3AF',
          },
          background: {
            primary: '#FFFFFF',
            secondary: '#F9FAFB',
            muted: '#F3F4F6',
            card: '#FFFFFF',
          }
        }
      },
      // Optimized spacing for better performance
      spacing: {
        '2xs': '0.25rem',
        'xs': '0.5rem',
        'sm': '0.625rem',
        'base': '1rem',
        'lg': '1.125rem',
        'xl': '1.25rem',
        '2xl': '1.5rem',
        '3xl': '2rem',
        '4xl': '2.25rem',
        '5xl': '3rem',
      },
      // Optimized transitions for better performance
      transition: {
        'all': 'all 0.2s ease-in-out',
        'color': 'color 0.2s ease-in-out',
        'background': 'background 0.2s ease-in-out',
        'border': 'border 0.2s ease-in-out',
        'opacity': 'opacity 0.2s ease-in-out',
        'transform': 'transform 0.2s ease-in-out',
        'shadow': 'box-shadow 0.2s ease-in-out',
        'none': 'none',
      },
      // Optimized animations for better performance
      animation: {
        'fade-in': 'fadeIn 0.3s ease-in-out',
        'slide-in': 'slideIn 0.3s ease-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'scale-in': 'scaleIn 0.2s ease-out',
        'bounce-in': 'bounceIn 0.5s ease-out',
      },
      // Custom animations optimized for performance
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideIn: {
          '0%': { transform: 'translateX(100%)' },
          '100%': { transform: 'translateX(0)' },
        },
        slideUp: {
          '0%': { transform: 'translateY(20px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        scaleIn: {
          '0%': { transform: 'scale(0.95)', opacity: '0' },
          '100%': { transform: 'scale(1)', opacity: '1' },
        },
        bounceIn: {
          '0%': { transform: 'scale(0.3)', opacity: '0' },
          '50%': { transform: 'scale(1.05)' },
          '70%': { transform: 'scale(0.9)' },
          '100%': { transform: 'scale(1)', opacity: '1' },
        },
      },
      // Optimized typography for better readability
      fontFamily: {
        sans: [
          'Inter',
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'Roboto',
          'Oxygen',
          'Ubuntu',
          'Cantarell',
          'Fira Sans',
          'Droid Sans',
          'Helvetica Neue',
          'sans-serif',
        ],
        mono: [
          'JetBrains Mono',
          'Fira Code',
          'Monaco',
          'Consolas',
          'Liberation Mono',
          'Courier New',
          'monospace',
        ],
      },
      // Optimized border radius
      borderRadius: {
        'none': '0px',
        'sm': '0.125rem',
        'base': '0.375rem',
        'md': '0.5rem',
        'lg': '0.625rem',
        'xl': '0.75rem',
        '2xl': '1rem',
        '3xl': '1.5rem',
        'full': '9999px',
      },
      // Optimized box shadows for better performance
      boxShadow: {
        'sm': '0 1px 2px 0 rgb(0 0 0 / 0.05)',
        'base': '0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)',
        'md': '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
        'lg': '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
        'xl': '0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1)',
        '2xl': '0 25px 50px -12px rgb(0 0 0 / 0.25)',
        'inner': 'inset 0 2px 4px 0 rgb(0 0 0 / 0.05)',
        'none': '0 0 #0000',
      },
    },
  },
  // Optimized plugins for performance
  plugins: [
    // Remove unused CSS for better performance
    function({ addUtilities, theme }) {
      const utilities = {
        '.scrollbar-hide': {
          '-ms-overflow-style': 'none',
          'scrollbar-width': 'none',
        },
        '.scrollbar-hide::-webkit-scrollbar': {
          display: 'none',
        },
        '.optimize-performance': {
          'will-change': 'transform, opacity',
          'backface-visibility': 'hidden',
          'transform': 'translateZ(0)',
        },
        '.optimize-motion': {
          'prefers-reduced-motion': 'reduce',
        },
      }
      addUtilities(utilities)
    },
  ],
  // Purge unused CSS for production
  purge: {
    enabled: process.env.NODE_ENV === 'production',
    content: [
      './index.html',
      './src/**/*.{vue,js,ts,jsx,tsx}',
    ],
    options: {
      safelist: {
        standard: [
          'dark',
          'bg-white',
          'bg-gray-50',
          'bg-blue-50',
          'text-gray-900',
          'text-gray-700',
          'text-gray-600',
          'border-gray-200',
          'border-gray-300',
        ]
      }
    }
  },
}