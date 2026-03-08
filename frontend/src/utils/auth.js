/**
 * Authentication mode utilities
 * Supports switching between User (JWT) and Agent (API Key) modes
 */

export const AuthMode = {
  USER: 'user',
  AGENT: 'agent'
}

/**
 * Get current authentication mode
 */
export function getAuthMode() {
  const apiKey = localStorage.getItem('agent_api_key')
  const token = localStorage.getItem('token')

  if (apiKey) return AuthMode.AGENT
  if (token) return AuthMode.USER
  return null
}

/**
 * Set agent API key and clear user token
 */
export function setAgentMode(apiKey) {
  localStorage.setItem('agent_api_key', apiKey)
  localStorage.removeItem('token')
  localStorage.removeItem('user')
}

/**
 * Set user token and clear agent API key
 */
export function setUserMode(token, user = null) {
  localStorage.setItem('token', token)
  if (user) {
    localStorage.setItem('user', JSON.stringify(user))
  }
  localStorage.removeItem('agent_api_key')
}

/**
 * Clear all authentication
 */
export function clearAuth() {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  localStorage.removeItem('agent_api_key')
}

/**
 * Check if in agent mode
 */
export function isAgentMode() {
  return !!localStorage.getItem('agent_api_key')
}

/**
 * Check if authenticated (either mode)
 */
export function isAuthenticated() {
  return !!(localStorage.getItem('token') || localStorage.getItem('agent_api_key'))
}

/**
 * Get auth info for display
 */
export function getAuthInfo() {
  const apiKey = localStorage.getItem('agent_api_key')
  const token = localStorage.getItem('token')
  const userStr = localStorage.getItem('user')

  if (apiKey) {
    return {
      mode: AuthMode.AGENT,
      identifier: `Agent (${apiKey.slice(0, 8)}...)`
    }
  }

  if (token && userStr) {
    try {
      const user = JSON.parse(userStr)
      return {
        mode: AuthMode.USER,
        identifier: user.email || 'User'
      }
    } catch (e) {
      return {
        mode: AuthMode.USER,
        identifier: 'User'
      }
    }
  }

  return null
}
