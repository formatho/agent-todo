<template>
  <div class="kanban-board">
    <!-- Kanban Columns -->
    <div class="kanban-columns">
      <div
        v-for="column in columns"
        :key="column.id"
        class="kanban-column"
        :class="`column-${column.id}`"
      >
        <!-- Column Header -->
        <div class="column-header">
          <div class="column-title">
            <span class="column-icon">{{ column.icon }}</span>
            <span class="column-name">{{ column.name }}</span>
            <span class="column-count">{{ column.tasks.length }}</span>
          </div>
        </div>

        <!-- Column Content (Draggable) -->
        <draggable
          :list="column.tasks"
          :animation="200"
          group="tasks"
          ghost-class="ghost-card"
          drag-class="drag-card"
          item-key="id"
          @change="onStatusChange($event, column.id)"
          class="column-content"
        >
          <template #item="{ element: task }">
            <div class="kanban-task-card">
              <!-- Priority Banner for Critical -->
              <div v-if="task.priority === 'critical'" class="priority-banner">
                <span>⚠️ CRITICAL</span>
              </div>

              <!-- Task Title -->
              <h4 class="task-title">{{ task.title }}</h4>

              <!-- Task Description -->
              <p class="task-description">{{ truncatedDescription(task) }}</p>

              <!-- Badges -->
              <div class="task-badges">
                <span class="priority-badge" :class="task.priority">
                  {{ getPriorityIcon(task.priority) }}
                  {{ formatPriority(task.priority) }}
                </span>
              </div>

              <!-- Agent Assignment -->
              <div v-if="task.assigned_agent" class="task-agent">
                <AgentAvatar :agent="task.assigned_agent" size="small" />
                <span class="agent-name">{{ task.assigned_agent.name }}</span>
              </div>
              <div v-else class="task-agent unassigned">
                <span class="unassigned-text">Unassigned</span>
              </div>

              <!-- Task Footer -->
              <div class="task-footer">
                <div class="task-meta">
                  <span v-if="task.due_date" class="meta-item">
                    <span>📅</span>
                    <span>{{ formatDate(task.due_date) }}</span>
                  </span>
                </div>

                <div class="task-actions">
                  <router-link :to="`/tasks/${task.id}`" class="btn-view">
                    View
                  </router-link>
                </div>
              </div>
            </div>
          </template>
        </draggable>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed, watch, onMounted } from 'vue'
import draggable from 'vuedraggable'
import AgentAvatar from './AgentAvatar.vue'
import { useTaskStore } from '../stores/tasks'

const taskStore = useTaskStore()

const columns = reactive([
  {
    id: 'pending',
    name: 'To Do',
    icon: '📋',
    color: '#F59E0B',
    tasks: []
  },
  {
    id: 'in_progress',
    name: 'In Progress',
    icon: '🔄',
    color: '#3B82F6',
    tasks: []
  },
  {
    id: 'completed',
    name: 'Done',
    icon: '✅',
    color: '#10B981',
    tasks: []
  },
  {
    id: 'failed',
    name: 'Failed',
    icon: '❌',
    color: '#EF4444',
    tasks: []
  }
])

// Priority order mapping (higher number = higher priority)
const priorityOrder = {
  critical: 4,
  high: 3,
  medium: 2,
  low: 1
}

// Sort tasks by priority (highest first)
const sortByPriority = (tasks) => {
  return [...tasks].sort((a, b) => {
    const priorityA = priorityOrder[a.priority] || 0
    const priorityB = priorityOrder[b.priority] || 0
    return priorityB - priorityA
  })
}

// Load tasks into columns
const loadTasks = () => {
  // Ensure tasks is an array
  const tasks = Array.isArray(taskStore.tasks) ? taskStore.tasks : []
  columns[0].tasks = sortByPriority(tasks.filter(task => task.status === 'pending'))
  columns[1].tasks = sortByPriority(tasks.filter(task => task.status === 'in_progress'))
  columns[2].tasks = sortByPriority(tasks.filter(task => task.status === 'completed'))
  columns[3].tasks = sortByPriority(tasks.filter(task => task.status === 'failed'))
}

// Watch for store changes and refresh columns
watch(() => taskStore.tasks, () => {
  loadTasks()
}, { deep: true })

onMounted(async () => {
  // Fetch tasks if not already loaded
  if (!taskStore.tasks || taskStore.tasks.length === 0) {
    await taskStore.fetchTasks()
  }
  loadTasks()
})

const truncatedDescription = (task) => {
  if (!task.description) return 'No description'
  return task.description.length > 80
    ? task.description.substring(0, 80) + '...'
    : task.description
}

const formatPriority = (priority) => {
  return priority.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const getPriorityIcon = (priority) => {
  const icons = {
    critical: '🔴',
    high: '🟠',
    medium: '🟡',
    low: '🟢'
  }
  return icons[priority] || '⚪'
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
  if (diffDays < -1) return `${Math.abs(diffDays)}d overdue`
  if (diffDays <= 7) return `${diffDays}d`

  return date.toLocaleDateString()
}

const onStatusChange = async (event, newStatus) => {
  if (event.added) {
    const task = event.added.element
    await updateTaskStatus(task, newStatus)
  }
}

const updateTaskStatus = async (task, newStatus) => {
  try {
    await taskStore.updateTask(task.id, { status: newStatus })
    // Refresh all tasks from the server to get latest state
    await taskStore.fetchTasks()
    console.log(`Task "${task.title}" moved to ${newStatus}`)
  } catch (error) {
    console.error('Failed to update task status:', error)
    // Reload tasks to revert the UI change
    await taskStore.fetchTasks()
  }
}
</script>

<style scoped>
.kanban-board {
  width: 100%;
  height: calc(100vh - 200px);
  overflow-x: auto;
  overflow-y: hidden;
}

.kanban-columns {
  display: flex;
  gap: 20px;
  height: 100%;
  min-width: min-content;
  padding: 4px;
}

.kanban-column {
  flex: 1;
  min-width: 280px;
  max-width: 350px;
  background: #F9FAFB;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  border: 2px solid #E5E7EB;
  height: 100%;
}

.column-pending { border-top: 4px solid #F59E0B; }
.column-in_progress { border-top: 4px solid #3B82F6; }
.column-completed { border-top: 4px solid #10B981; }
.column-failed { border-top: 4px solid #EF4444; }

/* Column Header */
.column-header {
  padding: 16px;
  border-bottom: 1px solid #E5E7EB;
  background: white;
  border-radius: 10px 10px 0 0;
  position: sticky;
  top: 0;
  z-index: 10;
}

.column-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: #111827;
  font-size: 14px;
}

.column-icon {
  font-size: 18px;
}

.column-name {
  flex: 1;
}

.column-count {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  height: 28px;
  background: #F3F4F6;
  color: #374151;
  border-radius: 14px;
  font-size: 13px;
  font-weight: 700;
  padding: 0 8px;
}

/* Column Content */
.column-content {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 100px;
}

.column-content::-webkit-scrollbar {
  width: 6px;
}

.column-content::-webkit-scrollbar-track {
  background: #F3F4F6;
  border-radius: 3px;
}

.column-content::-webkit-scrollbar-thumb {
  background: #D1D5DB;
  border-radius: 3px;
}

.column-content::-webkit-scrollbar-thumb:hover {
  background: #9CA3AF;
}

/* Task Card */
.kanban-task-card {
  background: white;
  border-radius: 8px;
  padding: 12px;
  cursor: grab;
  border: 1px solid #E5E7EB;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;
  position: relative;
}

.kanban-task-card:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
  border-color: #D1D5DB;
}

.kanban-task-card:active {
  cursor: grabbing;
}

/* Priority Banner */
.priority-banner {
  background: linear-gradient(90deg, #EF4444 0%, #F87171 100%);
  color: white;
  padding: 4px 8px;
  margin: -12px -12px 8px -12px;
  border-radius: 8px 8px 0 0;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.5px;
}

/* Task Title */
.task-title {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 6px 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Task Description */
.task-description {
  font-size: 12px;
  color: #6B7280;
  line-height: 1.5;
  margin: 0 0 8px 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Badges */
.task-badges {
  display: flex;
  gap: 6px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}

.priority-badge {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  padding: 3px 8px;
  border-radius: 10px;
  font-size: 10px;
  font-weight: 600;
}

.priority-badge.critical {
  background: #FEE2E2;
  color: #991B1B;
}

.priority-badge.high {
  background: #FFEDD5;
  color: #9A3412;
}

.priority-badge.medium {
  background: #DBEAFE;
  color: #1E40AF;
}

.priority-badge.low {
  background: #F3F4F6;
  color: #374151;
}

/* Agent */
.task-agent {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px;
  background: #F9FAFB;
  border-radius: 6px;
  margin-bottom: 8px;
}

.task-agent.unassigned {
  opacity: 0.6;
}

.agent-name {
  font-size: 11px;
  font-weight: 500;
  color: #374151;
}

.unassigned-text {
  font-size: 11px;
  color: #9CA3AF;
  font-style: italic;
}

/* Task Footer */
.task-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 8px;
  border-top: 1px solid #F3F4F6;
}

.task-meta {
  display: flex;
  gap: 8px;
  font-size: 11px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 3px;
  color: #6B7280;
}

.task-actions {
  display: flex;
  gap: 6px;
}

.btn-view {
  padding: 4px 10px;
  background: #F3F4F6;
  color: #374151;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.2s ease;
}

.btn-view:hover {
  background: #E5E7EB;
}

/* Drag & Drop States */
.ghost-card {
  opacity: 0.5;
  background: #F3F4F6;
  border: 2px dashed #D1D5DB;
}

.drag-card {
  transform: rotate(3deg);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.2);
}

/* Empty State */
.column-content:empty::after {
  content: 'Drop tasks here';
  display: flex;
  align-items: center;
  justify-content: center;
  height: 60px;
  color: #9CA3AF;
  font-size: 13px;
  font-style: italic;
  border: 2px dashed #E5E7EB;
  border-radius: 8px;
}
</style>
