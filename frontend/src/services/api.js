import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor to add auth token (JWT or API Key)
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    const apiKey = localStorage.getItem('agent_api_key')

    // Prefer API key if set (agent mode), otherwise use JWT (user mode)
    if (apiKey) {
      config.headers['X-API-KEY'] = apiKey
    } else if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const apiKey = localStorage.getItem('agent_api_key')
      const token = localStorage.getItem('token')

      // Clear the appropriate auth based on mode
      if (apiKey) {
        localStorage.removeItem('agent_api_key')
      } else if (token) {
        localStorage.removeItem('token')
      }

      // Only redirect to login if not in agent mode
      if (!apiKey) {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default api
