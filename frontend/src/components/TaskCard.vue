<template>
  <div
    class="task-card"
    :class="[
      `status-${task.status}`,
      `priority-${task.priority}`,
      { 'has-agent': task.assigned_agent }
    ]"
    :style="cardStyle"
    @click="handleCardClick"
  >
    <!-- Critical Priority Banner -->
    <div v-if="task.priority === 'critical'" class="priority-banner">
      <span class="priority-icon">{{ PRIORITY_COLORS.critical.icon }}</span>
      <span class="priority-text">CRITICAL</span>
    </div>

    <!-- Card Header -->
    <div class="card-header">
      <h3 class="task-title">{{ task.title }}</h3>

      <!-- Status Badge -->
      <span class="status-badge" :style="statusBadgeStyle">
        {{ STATUS_COLORS[task.status].icon }}
        {{ formatStatus(task.status) }}
      </span>

      <!-- Priority Badge (if not critical) -->
      <span v-if="task.priority !== 'critical'" class="priority-badge" :style="priorityBadgeStyle">
        {{ PRIORITY_COLORS[task.priority].icon }}
        {{ PRIORITY_COLORS[task.priority].label }}
      </span>
    </div>

    <!-- Project Badge -->
    <div v-if="task.project" class="project-indicator">
      <span class="project-icon">📁</span>
      <span class="project-name">{{ task.project.name }}</span>
    </div>

    <!-- Agent Indicator -->
    <div v-if="task.assigned_agent" class="agent-indicator">
      <AgentAvatar :agent="task.assigned_agent" size="small" />
      <span class="agent-name">{{ task.assigned_agent.name }}</span>
    </div>
    <div v-else class="agent-indicator unassigned">
      <span class="no-agent">Unassigned</span>
    </div>

    <!-- Card Body -->
    <div class="card-body">
      <p class="task-description">{{ truncatedDescription }}</p>

      <!-- Task Metadata -->
      <div class="task-meta">
        <span v-if="task.due_date" class="meta-item">
          <span class="meta-icon">📅</span>
          <span>{{ formatDate(task.due_date) }}</span>
        </span>
        <span class="meta-item">
          <span class="meta-icon">👤</span>
          <span>{{ task.created_by?.email || 'Unknown' }}</span>
        </span>
      </div>

      <!-- Progress Bar (for in_progress) -->
      <div v-if="task.status === 'in_progress'" class="progress-container">
        <div class="progress-bar">
          <div class="progress-fill"></div>
        </div>
        <span class="progress-text">In Progress</span>
      </div>
    </div>

    <!-- Card Footer -->
    <div class="card-footer">
      <div class="action-buttons">
        <button @click.stop="handleView" class="btn-view">View</button>
        <button @click.stop="handleEdit" class="btn-edit">Edit</button>
      </div>

      <div class="task-stats">
        <span class="stat" title="Comments">
          <span>💬</span>
          <span>{{ commentCount }}</span>
        </span>
        <span class="stat" title="Events">
          <span>📜</span>
          <span>{{ eventCount }}</span>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import AgentAvatar from './AgentAvatar.vue'
import { STATUS_COLORS, PRIORITY_COLORS, generateAgentColor } from '../utils/agentColors'

const props = defineProps({
  task: {
    type: Object,
    required: true
  }
})

const router = useRouter()

const truncatedDescription = computed(() => {
  if (!props.task.description) return 'No description'
  return props.task.description.length > 100
    ? props.task.description.substring(0, 100) + '...'
    : props.task.description
})

const commentCount = computed(() => {
  if (!Array.isArray(props.task.comments)) return 0
  return props.task.comments.length
})

const eventCount = computed(() => {
  if (!Array.isArray(props.task.events)) return 0
  return props.task.events.length
})

const cardStyle = computed(() => {
  if (!props.task.assigned_agent) {
    return {}
  }

  // Generate agent color for border
  const index = props.task.assigned_agent.id
    ? parseInt(props.task.assigned_agent.id.slice(-8), 16) % 10
    : 0
  const colors = generateAgentColor(index)

  return {
    borderLeftColor: colors.primary,
    background: colors.gradient
  }
})

const statusBadgeStyle = computed(() => {
  const colors = STATUS_COLORS[props.task.status]
  return {
    backgroundColor: colors.bg,
    color: colors.text,
    borderColor: colors.border
  }
})

const priorityBadgeStyle = computed(() => {
  const colors = PRIORITY_COLORS[props.task.priority]
  return {
    backgroundColor: colors.bg,
    color: colors.text,
    borderColor: colors.border
  }
})

const formatStatus = (status) => {
  return status.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diffTime = date - now
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

  if (diffDays === 0) return 'Today'
  if (diffDays === 1) return 'Tomorrow'
  if (diffDays === -1) return 'Yesterday'
  if (diffDays < -1) return `${Math.abs(diffDays)} days overdue`
  if (diffDays <= 7) return `In ${diffDays} days`

  return date.toLocaleDateString()
}

const handleCardClick = () => {
  router.push(`/tasks/${props.task.id}`)
}

const handleView = () => {
  router.push(`/tasks/${props.task.id}`)
}

const handleEdit = () => {
  // Emit edit event or open edit modal
  console.log('Edit task:', props.task.id)
}
</script>

<style scoped>
.task-card {
  background: white;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  border-left: 4px solid #9CA3AF;
  padding: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  position: relative;
  overflow: hidden;
}

.task-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
}

/* Priority Banner */
.priority-banner {
  background: linear-gradient(90deg, #EF4444 0%, #F87171 100%);
  color: white;
  padding: 4px 12px;
  margin: -16px -16px 12px -16px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.5px;
  display: flex;
  align-items: center;
  gap: 6px;
}

/* Card Header */
.card-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.task-title {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  margin: 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.status-badge,
.priority-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
  border: 1px solid;
  align-self: flex-start;
}

/* Project Indicator */
.project-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  background: #E0E7FF;
  border-radius: 6px;
  margin-bottom: 8px;
  width: fit-content;
}

.project-icon {
  font-size: 12px;
}

.project-name {
  font-size: 12px;
  font-weight: 500;
  color: #3730A3;
}

/* Agent Indicator */
.agent-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background: #F9FAFB;
  border-radius: 8px;
  margin-bottom: 12px;
}

.agent-indicator.unassigned {
  opacity: 0.6;
}

.agent-name {
  font-size: 13px;
  font-weight: 500;
  color: #374151;
}

.no-agent {
  font-size: 12px;
  color: #9CA3AF;
  font-style: italic;
}

/* Card Body */
.card-body {
  margin-bottom: 12px;
}

.task-description {
  font-size: 14px;
  color: #6B7280;
  line-height: 1.5;
  margin: 0 0 12px 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.task-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #6B7280;
}

.meta-icon {
  font-size: 14px;
}

/* Progress Bar */
.progress-container {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 12px;
}

.progress-bar {
  flex: 1;
  height: 6px;
  background: #E5E7EB;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #3B82F6 0%, #60A5FA 100%);
  border-radius: 3px;
  width: 50%;
  animation: progressPulse 2s ease-in-out infinite;
}

@keyframes progressPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

.progress-text {
  font-size: 11px;
  color: #3B82F6;
  font-weight: 600;
}

/* Card Footer */
.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #F3F4F6;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.btn-view,
.btn-edit {
  padding: 6px 14px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.btn-view {
  background: #F3F4F6;
  color: #374151;
}

.btn-view:hover {
  background: #E5E7EB;
}

.btn-edit {
  background: #DBEAFE;
  color: #1E40AF;
}

.btn-edit:hover {
  background: #BFDBFE;
}

.task-stats {
  display: flex;
  gap: 12px;
}

.stat {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #6B7280;
}

/* Status-specific styles */
.task-card.status-pending {
  border-left-color: #F59E0B;
}

.task-card.status-in_progress {
  border-left-color: #3B82F6;
}

.task-card.status-completed {
  border-left-color: #10B981;
  opacity: 0.85;
}

.task-card.status-failed {
  border-left-color: #EF4444;
}

/* Animations */
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.task-card {
  animation: slideIn 0.3s ease;
}
</style>
