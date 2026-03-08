<template>
  <div
    class="agent-card"
    :style="cardStyle"
  >
    <!-- Card Header with Avatar -->
    <div class="card-header">
      <AgentAvatar :agent="agent" size="large" />

      <div class="agent-info">
        <div class="agent-name-row">
          <h3 class="agent-name">{{ agent.name }}</h3>
          <span :class="['role-badge', `role-${agent.role}`]">
            {{ formatRole(agent.role) }}
          </span>
        </div>
        <p class="agent-description">{{ descriptionText }}</p>
      </div>

      <div class="card-actions">
        <button @click="handleEdit" class="btn-icon" title="Edit">
          ✏️
        </button>
        <button @click="handleCopyKey" class="btn-icon" title="Copy API Key">
          🔑
        </button>
        <button @click="handleDelete" class="btn-icon btn-danger" title="Delete">
          🗑️
        </button>
      </div>
    </div>

    <!-- Agent Stats -->
    <div class="agent-stats">
      <div class="stat">
        <span class="stat-icon">📋</span>
        <span class="stat-value">{{ taskCount }}</span>
        <span class="stat-label">Tasks</span>
      </div>

      <div class="stat">
        <span class="stat-icon">⏳</span>
        <span class="stat-value">{{ pendingCount }}</span>
        <span class="stat-label">Pending</span>
      </div>

      <div class="stat">
        <span class="stat-icon">✅</span>
        <span class="stat-value">{{ completedCount }}</span>
        <span class="stat-label">Done</span>
      </div>
    </div>

    <!-- API Key Section -->
    <div class="api-key-section">
      <div class="api-key-label">API Key</div>
      <div class="api-key-container">
        <code class="api-key">{{ maskedApiKey }}</code>
        <button @click="handleCopyKey" class="btn-copy">
          {{ copied ? '✓ Copied!' : 'Copy' }}
        </button>
      </div>
    </div>

    <!-- Task Progress -->
    <div v-if="taskCount > 0" class="progress-section">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
      </div>
      <div class="progress-text">
        {{ completedCount }} of {{ taskCount }} tasks completed
      </div>
    </div>

    <!-- View Tasks Button -->
    <div class="view-tasks-section">
      <router-link :to="`/tasks?agent_id=${agent.id}`" class="btn-view-tasks">
        <span class="btn-icon-left">📋</span>
        <span>View Tasks</span>
        <span class="btn-icon-right">→</span>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import AgentAvatar from './AgentAvatar.vue'
import { generateAgentColor } from '../utils/agentColors'

const props = defineProps({
  agent: {
    type: Object,
    required: true
  },
  tasks: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['edit', 'delete'])

const router = useRouter()
const copied = ref(false)

const descriptionText = computed(() => {
  return props.agent.description || 'No description provided'
})

const formatRole = (role) => {
  const roles = {
    regular: 'Regular',
    supervisor: 'Supervisor',
    admin: 'Admin'
  }
  return roles[role] || role
}

const maskedApiKey = computed(() => {
  return props.agent.api_key.slice(0, 12) + '...' + props.agent.api_key.slice(-4)
})

const taskCount = computed(() => props.tasks.length)
const pendingCount = computed(() => props.tasks.filter(t => t.status !== 'completed').length)
const completedCount = computed(() => props.tasks.filter(t => t.status === 'completed').length)

const progressPercent = computed(() => {
  if (taskCount.value === 0) return 0
  return Math.round((completedCount.value / taskCount.value) * 100)
})

const cardStyle = computed(() => {
  // Generate agent color
  const index = props.agent.id ? parseInt(props.agent.id.slice(-8), 16) % 10 : 0
  const colors = generateAgentColor(index)

  return {
    borderTop: `4px solid ${colors.primary}`,
    background: colors.gradient
  }
})

const handleEdit = () => {
  emit('edit', props.agent)
}

const handleDelete = () => {
  if (confirm(`Are you sure you want to delete ${props.agent.name}?`)) {
    emit('delete', props.agent.id)
  }
}

const handleCopyKey = async () => {
  try {
    await navigator.clipboard.writeText(props.agent.api_key)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}
</script>

<style scoped>
.agent-card {
  background: white;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  padding: 20px;
  transition: all 0.3s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.agent-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
}

/* Card Header */
.card-header {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.agent-info {
  flex: 1;
}

.agent-name-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.agent-name {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
  margin: 0;
}

.role-badge {
  font-size: 10px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.role-badge.role-regular {
  background: #DBEAFE;
  color: #1E40AF;
}

.role-badge.role-supervisor {
  background: #FEF3C7;
  color: #92400E;
}

.role-badge.role-admin {
  background: #FEE2E2;
  color: #991B1B;
}

.agent-description {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
  line-height: 1.4;
}

.card-actions {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.btn-icon {
  width: 32px;
  height: 32px;
  border: none;
  background: #F3F4F6;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 16px;
}

.btn-icon:hover {
  background: #E5E7EB;
  transform: scale(1.1);
}

.btn-icon.btn-danger:hover {
  background: #FEE2E2;
}

/* Agent Stats */
.agent-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 16px;
  padding: 12px;
  background: #F9FAFB;
  border-radius: 8px;
}

.stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.stat-icon {
  font-size: 20px;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
  color: #111827;
}

.stat-label {
  font-size: 11px;
  color: #6B7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* API Key Section */
.api-key-section {
  margin-bottom: 16px;
}

.api-key-label {
  font-size: 12px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 6px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.api-key-container {
  display: flex;
  gap: 8px;
  align-items: center;
}

.api-key {
  flex: 1;
  padding: 8px 12px;
  background: #F3F4F6;
  border: 1px solid #E5E7EB;
  border-radius: 6px;
  font-size: 12px;
  font-family: monospace;
  color: #374151;
  overflow: hidden;
  text-overflow: ellipsis;
}

.btn-copy {
  padding: 8px 14px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.btn-copy:hover {
  background: #2563EB;
}

/* Progress Section */
.progress-section {
  margin-top: 12px;
}

.progress-bar {
  height: 6px;
  background: #E5E7EB;
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 6px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #10B981 0%, #34D399 100%);
  border-radius: 3px;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 12px;
  color: #6B7280;
  text-align: center;
}

/* View Tasks Section */
.view-tasks-section {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #E5E7EB;
}

.btn-view-tasks {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 12px 16px;
  background: linear-gradient(135deg, #3B82F6 0%, #2563EB 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: none;
}

.btn-view-tasks:hover {
  background: linear-gradient(135deg, #2563EB 0%, #1D4ED8 100%);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4);
}

.btn-icon-left {
  font-size: 16px;
}

.btn-icon-right {
  font-size: 16px;
  opacity: 0;
  transform: translateX(-4px);
  transition: all 0.2s ease;
}

.btn-view-tasks:hover .btn-icon-right {
  opacity: 1;
  transform: translateX(0);
}
</style>
