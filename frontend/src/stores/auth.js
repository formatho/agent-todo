import { defineStore } from 'pinia'
import { authService } from '../services/authService'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: authService.getUser(),
    token: authService.getToken(),
    isAuthenticated: !!authService.getToken()
  }),

  actions: {
    async login(email, password) {
      const data = await authService.login(email, password)
      this.token = data.token
      this.user = data.user
      this.isAuthenticated = true
      authService.setAuth(data.token, data.user)
    },

    async register(email, password) {
      const data = await authService.register(email, password)
      this.token = data.token
      this.user = data.user
      this.isAuthenticated = true
      authService.setAuth(data.token, data.user)
    },

    logout() {
      authService.logout()
      this.user = null
      this.token = null
      this.isAuthenticated = false
    }
  }
})
