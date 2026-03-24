// Agent color generation utility

const AGENT_HUES = [
  220, // Blue
  160, // Green
  30,  // Orange
  330, // Pink
  180, // Teal
  0,   // Red
  270, // Purple
  45,  // Yellow
  200, // Sky
  280  // Violet
]

/**
 * Generate color palette for an agent based on index
 * @param {number} index - Agent index
 * @returns {Object} Color palette with primary, light, dark, and gradient
 */
export function generateAgentColor(index) {
  const hue = AGENT_HUES[index % AGENT_HUES.length]

  return {
    primary: `hsl(${hue}, 70%, 50%)`,
    light: `hsl(${hue}, 70%, 95%)`,
    dark: `hsl(${hue}, 70%, 30%)`,
    hex: hslToHex(hue, 70, 50),
    lightHex: hslToHex(hue, 70, 95),
    darkHex: hslToHex(hue, 70, 30),
    gradient: `linear-gradient(135deg, hsl(${hue}, 70%, 95%) 0%, white 100%)`
  }
}

/**
 * Get agent initials for avatar
 * @param {string} name - Agent name
 * @returns {string} Initials (max 2 characters)
 */
export function getAgentInitials(name) {
  if (!name) return '??'

  const words = name.trim().split(/\s+/)
  if (words.length === 1) {
    return words[0].substring(0, 2).toUpperCase()
  }

  return words
    .slice(0, 2)
    .map(word => word[0])
    .join('')
    .toUpperCase()
}

/**
 * HSL to Hex conversion
 */
function hslToHex(h, s, l) {
  l /= 100
  const a = s * Math.min(l, 1 - l) / 100
  const f = n => {
    const k = (n + h / 30) % 12
    const color = l - a * Math.max(Math.min(k - 3, 9 - k, 1), -1)
    return Math.round(255 * color).toString(16).padStart(2, '0')
  }
  return `#${f(0)}${f(8)}${f(4)}`
}

/**
 * Default agent colors
 */
export const DEFAULT_AGENT_COLORS = {
  unassigned: {
    primary: '#9CA3AF',
    light: '#F3F4F6',
    dark: '#4B5563',
    hex: '#9CA3AF',
    lightHex: '#F3F4F6',
    darkHex: '#4B5563',
    gradient: 'linear-gradient(135deg, #F3F4F6 0%, white 100%)'
  }
}

/**
 * Status colors
 */
export const STATUS_COLORS = {
  pending: {
    bg: '#FEF3C7',
    text: '#92400E',
    border: '#F59E0B',
    gradient: 'linear-gradient(to right, #FFFBEB, white)',
    icon: '⏳'
  },
  in_progress: {
    bg: '#DBEAFE',
    text: '#1E40AF',
    border: '#3B82F6',
    gradient: 'linear-gradient(to right, #DBEAFE, white)',
    icon: '🔄'
  },
  blocked: {
    bg: '#EDE9FE',
    text: '#5B21B6',
    border: '#8B5CF6',
    gradient: 'linear-gradient(to right, #EDE9FE, white)',
    icon: '🚧'
  },
  completed: {
    bg: '#D1FAE5',
    text: '#065F46',
    border: '#10B981',
    gradient: 'linear-gradient(to right, #D1FAE5, white)',
    icon: '✅'
  },
  failed: {
    bg: '#FEE2E2',
    text: '#991B1B',
    border: '#EF4444',
    gradient: 'linear-gradient(to right, #FEE2E2, white)',
    icon: '❌'
  }
}

/**
 * Priority colors
 */
export const PRIORITY_COLORS = {
  low: {
    bg: '#F3F4F6',
    text: '#374151',
    border: '#9CA3AF',
    icon: '⬇️',
    label: 'LOW'
  },
  medium: {
    bg: '#DBEAFE',
    text: '#1E40AF',
    border: '#3B82F6',
    icon: '➡️',
    label: 'MEDIUM'
  },
  high: {
    bg: '#FED7AA',
    text: '#9A3412',
    border: '#F59E0B',
    icon: '⬆️',
    label: 'HIGH'
  },
  critical: {
    bg: '#FECACA',
    text: '#991B1B',
    border: '#EF4444',
    icon: '🔥',
    label: 'CRITICAL'
  }
}
