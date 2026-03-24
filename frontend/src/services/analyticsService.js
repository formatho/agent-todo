import api from './api'

// Generate or get session ID for anonymous tracking
const getSessionId = () => {
  let sessionId = localStorage.getItem('analytics_session_id')
  if (!sessionId) {
    sessionId = 'sess_' + Math.random().toString(36).substring(2) + Date.now().toString(36)
    localStorage.setItem('analytics_session_id', sessionId)
  }
  return sessionId
}

/**
 * Track an analytics event
 * @param {string} eventType - Type of event: 'page_view', 'checkout_start', 'checkout_complete'
 * @param {string} page - Page name: 'pricing', 'checkout'
 * @param {string} plan - Plan name (optional): 'starter', 'pro', 'enterprise'
 * @param {object} metadata - Additional metadata (optional)
 */
const trackEvent = async (eventType, page, plan = null, metadata = {}) => {
  try {
    await api.post('/analytics/track', {
      event_type: eventType,
      page,
      plan,
      session_id: getSessionId(),
      metadata
    })
  } catch (error) {
    // Silently fail - analytics shouldn't break the user experience
    console.debug('Analytics tracking failed:', error)
  }
}

/**
 * Track pricing page view
 */
const trackPricingPageView = () => {
  return trackEvent('page_view', 'pricing')
}

/**
 * Track checkout button click
 * @param {string} plan - The plan being selected
 */
const trackCheckoutStart = (plan) => {
  return trackEvent('checkout_start', 'checkout', plan)
}

/**
 * Track completed checkout
 * @param {string} plan - The plan that was purchased
 * @param {object} metadata - Additional checkout details
 */
const trackCheckoutComplete = (plan, metadata = {}) => {
  return trackEvent('checkout_complete', 'checkout', plan, metadata)
}

/**
 * Get funnel statistics (requires auth)
 * @param {number} days - Number of days to look back
 */
const getFunnelStats = async (days = 30) => {
  const response = await api.get(`/analytics/funnel?days=${days}`)
  return response.data
}

/**
 * Get recent events (requires auth)
 * @param {number} limit - Number of events to return
 * @param {string} eventType - Filter by event type (optional)
 */
const getRecentEvents = async (limit = 100, eventType = null) => {
  let url = `/analytics/events?limit=${limit}`
  if (eventType) {
    url += `&event_type=${eventType}`
  }
  const response = await api.get(url)
  return response.data
}

export default {
  trackEvent,
  trackPricingPageView,
  trackCheckoutStart,
  trackCheckoutComplete,
  getFunnelStats,
  getRecentEvents
}
