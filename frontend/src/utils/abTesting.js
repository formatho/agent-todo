// A/B Testing utilities for Agent Todo Platform
import { getTestVariant, trackImpression, TEST_CONFIGS } from '../services/abTestingService.js'

// Vue composable for A/B testing
export function useABTesting() {
  const getLayoutVariant = (userId) => {
    return getTestVariant('layout_test', userId, TEST_CONFIGS.LAYOUT_VARIANTS)
  }
  
  const getColorVariant = (userId) => {
    return getTestVariant('color_test', userId, TEST_CONFIGS.COLOR_VARIANTS)
  }
  
  const trackPageImpression = (userId, testNames = []) => {
    testNames.forEach(testName => {
      trackImpression(testName, getTestVariant(testName, userId, TEST_CONFIGS.LAYOUT_VARIANTS), userId)
    })
  }
  
  return {
    getLayoutVariant,
    getColorVariant,
    trackPageImpression
  }
}

// Generate CSS custom properties for color variants
export function generateColorCSS(variant) {
  const colorConfig = TEST_CONFIGS.COLOR_VARIANTS[variant]
  if (!colorConfig) return ''
  
  return `
    :root {
      --primary-color: ${colorConfig.primary};
      --secondary-color: ${colorConfig.secondary};
      --primary-bg: ${colorConfig.primary}10;
      --primary-border: ${colorConfig.primary}30;
      --primary-text: ${colorConfig.primary};
    }
  `
}

// A/B Testing components loader
export function loadABTestComponent(testName, variant) {
  // This would dynamically import components based on variant
  // For now, return placeholder component names
  const componentMap = {
    'layout_test': {
      'control': () => import('../components/DashboardControl.vue'),
      'variant-1': () => import('../components/DashboardVariant1.vue'),
      'variant-2': () => import('../components/DashboardVariant2.vue')
    }
  }
  
  const testComponents = componentMap[testName]
  if (!testComponents || !testComponents[variant]) {
    return null
  }
  
  return testComponents[variant]()
}

// Helper to create random user ID for testing
export function generateTestId() {
  return Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
}

// Debug function to test A/B testing assignments
export function debugABTesting(userId) {
  console.log('=== A/B Testing Debug Info ===')
  console.log('User ID:', userId)
  
  const layoutVariant = getLayoutVariant(userId)
  const colorVariant = getColorVariant(userId)
  
  console.log('Layout Variant:', layoutVariant, TEST_CONFIGS.LAYOUT_VARIANTS[layoutVariant]?.name)
  console.log('Color Variant:', colorVariant, TEST_CONFIGS.COLOR_VARIANTS[colorVariant]?.name)
  
  return {
    layoutVariant,
    colorVariant,
    layoutConfig: TEST_CONFIGS.LAYOUT_VARIANTS[layoutVariant],
    colorConfig: TEST_CONFIGS.COLOR_VARIANTS[colorVariant]
  }
}