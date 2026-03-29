// A/B Testing Service for Agent Todo Platform
// Provides functionality to run A/B tests on different page variations

// Test configurations
const TEST_CONFIGS = {
  // Landing page layout tests
  LAYOUT_VARIANTS: {
    'control': {
      name: 'Current Layout',
      weight: 50,
      component: 'DashboardControl'
    },
    'variant-1': {
      name: 'Enhanced Activity Feed',
      weight: 25,
      component: 'DashboardVariant1'
    },
    'variant-2': {
      name: 'Compact Task Grid',
      weight: 25,
      component: 'DashboardVariant2'
    }
  },
  
  // Color scheme tests
  COLOR_VARIANTS: {
    'default': {
      name: 'Default Blue',
      weight: 60,
      primary: '#3B82F6',
      secondary: '#60A5FA'
    },
    'green': {
      name: 'Professional Green',
      weight: 20,
      primary: '#10B981',
      secondary: '#34D399'
    },
    'purple': {
      name: 'Modern Purple',
      weight: 20,
      primary: '#8B5CF6',
      secondary: '#A78BFA'
    }
  }
}

// User assignment storage
class ABTestingStorage {
  constructor() {
    this.storageKey = 'agent_todo_ab_tests'
  }
  
  getUserAssignment(testName) {
    if (typeof window === 'undefined') return null
    
    const assignments = JSON.parse(localStorage.getItem(this.storageKey) || '{}')
    return assignments[testName]
  }
  
  setUserAssignment(testName, variant) {
    if (typeof window === 'undefined') return
    
    const assignments = JSON.parse(localStorage.getItem(this.storageKey) || '{}')
    assignments[testName] = {
      variant,
      assignedAt: new Date().toISOString()
    }
    localStorage.setItem(this.storageKey, JSON.stringify(assignments))
  }
  
  clearAssignments() {
    if (typeof window === 'undefined') return
    localStorage.removeItem(this.storageKey)
  }
}

const storage = new ABTestingStorage()

// Hash-based user assignment (consistent per user)
function hashUserAssignment(userId, testKey, variants) {
  const hash = String(userId).split('').reduce((a, b) => {
    a = ((a << 5) - a) + b.charCodeAt(0)
    return a & a
  }, 0)
  
  const totalWeight = Object.values(variants).reduce((sum, v) => sum + v.weight, 0)
  let currentWeight = 0
  
  for (const [variant, config] of Object.entries(variants)) {
    currentWeight += config.weight
    if (Math.abs(hash) % totalWeight < currentWeight {
      return variant
    }
  }
  
  // Fallback to first variant
  return Object.keys(variants)[0]
}

// Get user variant for a test
function getTestVariant(testName, userId, variants) {
  // Check if user already has an assignment
  const existingAssignment = storage.getUserAssignment(testName)
  if (existingAssignment) {
    return existingAssignment.variant
  }
  
  // Assign new variant based on user ID
  const variant = hashUserAssignment(userId, testName, variants)
  storage.setUserAssignment(testName, variant)
  return variant
}

// Track test impression
function trackImpression(testName, variant, userId) {
  // In a real implementation, this would send data to analytics
  console.log(`A/B Test Impression: ${testName}=${variant} for user=${userId}`)
  
  // Store in local analytics for now
  if (typeof window !== 'undefined') {
    const impressions = JSON.parse(localStorage.getItem('ab_test_impressions') || '{}')
    if (!impressions[testName]) impressions[testName] = {}
    if (!impressions[testName][variant]) impressions[testName][variant] = 0
    impressions[testName][variant] += 1
    
    localStorage.setItem('ab_test_impressions', JSON.stringify(impressions))
  }
}

// Track conversion event
function trackConversion(testName, variant, userId, conversionType) {
  // In a real implementation, this would send data to analytics
  console.log(`A/B Test Conversion: ${testName}=${variant}, ${conversionType} for user=${userId}`)
  
  // Store in local analytics for now
  if (typeof window !== 'undefined') {
    const conversions = JSON.parse(localStorage.getItem('ab_test_conversions') || '{}')
    if (!conversions[testName]) conversions[testName] = {}
    if (!conversions[testName][variant]) conversions[testName][variant] = []
    conversions[testName][variant].push({
      type: conversionType,
      timestamp: new Date().toISOString()
    })
    
    localStorage.setItem('ab_test_conversions', JSON.stringify(conversions))
  }
}

// Get test results for analysis
function getTestResults(testName) {
  const impressions = JSON.parse(localStorage.getItem('ab_test_impressions') || '{}')
  const conversions = JSON.parse(localStorage.getItem('ab_test_conversions') || '{}')
  
  return {
    testName,
    impressions: impressions[testName] || {},
    conversions: conversions[testName] || {},
    generatedAt: new Date().toISOString()
  }
}

// Reset all test data (admin/development use)
function resetTestData() {
  storage.clearAssignments()
  if (typeof window !== 'undefined') {
    localStorage.removeItem('ab_test_impressions')
    localStorage.removeItem('ab_test_conversions')
  }
}

export {
  TEST_CONFIGS,
  getTestVariant,
  trackImpression,
  trackConversion,
  getTestResults,
  resetTestData
}