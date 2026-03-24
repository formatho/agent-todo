import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme', {
  state: () => ({
    theme: 'system', // 'light', 'dark', or 'system'
    isDark: false
  }),

  actions: {
    initTheme() {
      // Get saved theme preference or default to 'system'
      const savedTheme = localStorage.getItem('theme-preference')
      this.theme = savedTheme || 'system'
      this.applyTheme()
    },

    setTheme(theme) {
      this.theme = theme
      localStorage.setItem('theme-preference', theme)
      this.applyTheme()
    },

    toggleTheme() {
      if (this.isDark) {
        this.setTheme('light')
      } else {
        this.setTheme('dark')
      }
    },

    applyTheme() {
      let isDark = false

      if (this.theme === 'dark') {
        isDark = true
      } else if (this.theme === 'light') {
        isDark = false
      } else if (this.theme === 'system') {
        // Check system preference
        isDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      }

      this.isDark = isDark

      // Apply dark mode class to HTML element
      if (isDark) {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }

      // Listen for system theme changes if using system preference
      if (this.theme === 'system') {
        this.setupSystemThemeListener()
      } else {
        this.removeSystemThemeListener()
      }
    },

    setupSystemThemeListener() {
      // Set up listener for system theme changes
      this.removeSystemThemeListener()

      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
      mediaQuery.addEventListener('change', this.handleSystemThemeChange)
      this._systemThemeListener = mediaQuery
    },

    removeSystemThemeListener() {
      if (this._systemThemeListener) {
        this._systemThemeListener.removeEventListener('change', this.handleSystemThemeChange)
        this._systemThemeListener = null
      }
    },

    handleSystemThemeChange(e) {
      // Only apply if using system preference
      if (this.theme === 'system') {
        this.isDark = e.matches
        if (e.matches) {
          document.documentElement.classList.add('dark')
        } else {
          document.documentElement.classList.remove('dark')
        }
      }
    }
  }
})
