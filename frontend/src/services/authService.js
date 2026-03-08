import api from './api'

export const authService = {
  async register(email, password) {
    const response = await api.post('/auth/register', { email, password })
    return response.data
  },

  async login(email, password) {
    const response = await api.post('/auth/login', { email, password })
    return response.data
  },

  async getCurrentUser() {
    const response = await api.get('/auth/me')
    return response.data
  },

  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  },

  getToken() {
    return localStorage.getItem('token')
  },

  getUser() {
    return JSON.parse(localStorage.getItem('user') || 'null')
  },

  setAuth(token, user) {
    localStorage.setItem('token', token)
    localStorage.setItem('user', JSON.stringify(user))
  }
}
