<template>
  <div class="task-grid-enhanced">
    <div class="grid-header">
      <div class="header-left">
        <h3 class="text-lg font-semibold text-gray-900">Task Overview</h3>
        <div class="quick-stats">
          <div class="quick-stat">
            <span class="stat-number">{{ totalTasks }}</span>
            <span class="stat-label">Total</span>
          </div>
          <div class="quick-stat">
            <span class="stat-number">{{ activeTasks }}</span>
            <span class="stat-label">Active</span>
          </div>
          <div class="quick-stat">
            <span class="stat-number">{{ completedTasks }}</span>
            <span class="stat-label">Done</span>
          </div>
        </div>
      </div>
      <div class="header-right">
        <div class="view-controls">
          <button 
            @click="setView('grid')" 
            :class="['view-btn', view === 'grid' ? 'active' : '']"
            title="Grid View"
          >
            <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
            </svg>
          </button>
          <button 
            @click="setView('list')" 
            :class="['view-btn', view === 'list' ? 'active' : '']"
            title="List View"
          >
            <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          <button 
            @click="setView('kanban')" 
            :class="['view-btn', view === 'kanban' ? 'active' : '']"
            title="Kanban View"
          >
            <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </button>
        </div>
        <button @click="showAddTask = true" class="btn-add-task">
          <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Add Task
        </button>
      </div>
    </div>

    <!-- Add Task Modal -->
    <div v-if="showAddTask" class="add-task-modal">
      <div class="modal-overlay" @click="showAddTask = false"></div>
      <div class="modal-content">
        <div class="modal-header">
          <h4>Create New Task</h4>
          <button @click="showAddTask = false" class="btn-close">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">Task Title</label>
            <input 
              v-model="newTask.title" 
              type="text" 
              class="form-input"
              placeholder="Enter task title"
            />
          </div>
          <div class="form-group">
            <label class="form-label">Description</label>
            <textarea 
              v-model="newTask.description" 
              class="form-textarea"
              placeholder="Enter task description"
              rows="3"
            ></textarea>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">Priority</label>
              <select v-model="newTask.priority" class="form-select">
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
                <option value="critical">Critical</option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label">Assign to</label>
              <select v-model="newTask.assignee" class="form-select">
                <option value="">Unassigned</option>
                <option v-for="agent in availableAgents" :key="agent.id" :value="agent.id">
                  {{ agent.name }}
                </option>
              </select>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showAddTask = false" class="btn-secondary">Cancel</button>
          <button @click="addTask" class="btn-primary" :disabled="!newTask.title.trim()">
            Create Task
          </button>
        </div>
      </div>
    </div>

    <!-- Task Filters -->
    <div class="task-filters">
      <div class="filter-group">
        <label class="filter-label">Status:</label>
        <button 
          v-for="status in taskStatuses" 
          :key="status"
          @click="toggleStatusFilter(status)"
          :class="['filter-btn', statusFilters.includes(status) ? 'active' : '']"
        >
          {{ status }}
        </button>
      </div>
      <div class="filter-group">
        <label class="filter-label">Priority:</label>
        <button 
          v-for="priority in taskPriorities" 
          :key="priority"
          @click="togglePriorityFilter(priority)"
          :class="['filter-btn', priorityFilters.includes(priority) ? 'active' : '']"
        >
          {{ priority }}
        </button>
      </div>
      <div class="filter-group">
        <label class="filter-label">Assignee:</label>
        <button 
          v-for="agent in availableAgents" 
          :key="agent.id"
          @click="toggleAgentFilter(agent.id)"
          :class="['filter-btn', agentFilters.includes(agent.id) ? 'active' : '']"
        >
          {{ agent.name }}
        </button>
      </div>
    </div>

    <!-- Task Grid/List/Kanban Views -->
    <div class="task-content">
      <!-- Grid View -->
      <div v-if="view === 'grid'" class="grid-view">
        <div 
          v-for="task in filteredTasks" 
          :key="task.id"
          class="task-card"
          :class="getTaskClass(task)"
        >
          <div class="task-header">
            <div class="task-priority" :class="task.priority"></div>
            <div class="task-status" :class="task.status"></div>
          </div>
          <div class="task-title">{{ task.title }}</div>
          <div class="task-description">{{ task.description || 'No description' }}</div>
          <div class="task-meta">
            <span v-if="task.assignee" class="meta-item">
              🤖 {{ getAgentName(task.assignee) }}
            </span>
            <span class="meta-item">
              📅 {{ formatDate(task.createdAt) }}
            </span>
          </div>
          <div class="task-actions">
            <button @click="editTask(task)" class="btn-action">Edit</button>
            <button @click="toggleTaskStatus(task)" class="btn-action">
              {{ task.status === 'completed' ? 'Reopen' : 'Complete' }}
            </button>
          </div>
        </div>
      </div>

      <!-- List View -->
      <div v-if="view === 'list'" class="list-view">
        <div 
          v-for="task in filteredTasks" 
          :key="task.id"
          class="task-row"
          :class="getTaskClass(task)"
        >
          <div class="task-checkbox">
            <input 
              type="checkbox" 
              :checked="task.status === 'completed'"
              @change="toggleTaskStatus(task)"
            />
          </div>
          <div class="task-content">
            <div class="task-title-row">
              <span class="task-title">{{ task.title }}</span>
              <span class="task-status" :class="task.status">{{ task.status }}</span>
            </div>
            <div class="task-description">{{ task.description || 'No description' }}</div>
            <div class="task-meta">
              <span v-if="task.assignee" class="meta-item">
                🤖 {{ getAgentName(task.assignee) }}
              </span>
              <span class="meta-item">
                📅 {{ formatDate(task.createdAt) }}
              </span>
              <span class="meta-item">
                ⚡ {{ task.priority }}
              </span>
            </div>
          </div>
          <div class="task-actions">
            <button @click="editTask(task)" class="btn-action">Edit</button>
            <button @click="deleteTask(task)" class="btn-action danger">Delete</button>
          </div>
        </div>
      </div>

      <!-- Kanban View -->
      <div v-if="view === 'kanban'" class="kanban-view">
        <div v-for="status in kanbanColumns" :key="status" class="kanban-column">
          <div class="column-header">
            <h4>{{ status }}</h4>
            <span class="column-count">{{ tasksByStatus[status].length }}</span>
          </div>
          <div class="column-content">
            <div 
              v-for="task in tasksByStatus[status]" 
              :key="task.id"
              class="task-card-kanban"
              :class="getTaskClass(task)"
            >
              <div class="task-priority" :class="task.priority"></div>
              <div class="task-title">{{ task.title }}</div>
              <div class="task-meta">
                <span v-if="task.assignee" class="meta-item">
                  🤖 {{ getAgentName(task.assignee) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="filteredTasks.length === 0" class="empty-tasks">
      <div class="empty-icon">📋</div>
      <h4>No tasks found</h4>
      <p>{{ getEmptyMessage() }}</p>
      <button @click="clearFilters" class="btn-secondary">Clear Filters</button>
      <button @click="showAddTask = true" class="btn-primary">Create Task</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// Mock task data
const tasks = ref([
  {
    id: 1,
    title: 'Setup database schema',
    description: 'Design and implement the database schema for the task management system',
    status: 'in-progress',
    priority: 'high',
    assignee: 1,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 2), // 2 hours ago
    completedAt: null
  },
  {
    id: 2,
    title: 'Create API endpoints',
    description: 'Implement REST API endpoints for task CRUD operations',
    status: 'pending',
    priority: 'high',
    assignee: 1,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 4), // 4 hours ago
    completedAt: null
  },
  {
    id: 3,
    title: 'Design user interface',
    description: 'Create wireframes and mockups for the dashboard',
    status: 'completed',
    priority: 'medium',
    assignee: 2,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24), // 1 day ago
    completedAt: new Date(Date.now() - 1000 * 60 * 60 * 6) // 6 hours ago
  },
  {
    id: 4,
    title: 'Write documentation',
    description: 'Create comprehensive documentation for API usage',
    status: 'pending',
    priority: 'low',
    assignee: null,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 6), // 6 hours ago
    completedAt: null
  }
])

// Mock agent data
const availableAgents = ref([
  { id: 1, name: 'code-assistant', type: 'developer' },
  { id: 2, name: 'doc-writer', type: 'general' },
  { id: 3, name: 'design-bot', type: 'designer' }
])

// View state
const view = ref('grid')
const showAddTask = ref(false)

// New task form
const newTask = ref({
  title: '',
  description: '',
  priority: 'medium',
  assignee: null
})

// Filter states
const statusFilters = ref(['pending', 'in-progress', 'completed'])
const priorityFilters = ref(['low', 'medium', 'high', 'critical'])
const agentFilters = ref([])

// Task options
const taskStatuses = ['pending', 'in-progress', 'completed']
const taskPriorities = ['low', 'medium', 'high', 'critical']

// Computed properties
const totalTasks = computed(() => tasks.value.length)
const activeTasks = computed(() => tasks.value.filter(task => task.status !== 'completed').length)
const completedTasks = computed(() => tasks.value.filter(task => task.status === 'completed').length)

const filteredTasks = computed(() => {
  let filtered = tasks.value

  // Apply status filters
  if (statusFilters.value.length > 0) {
    filtered = filtered.filter(task => statusFilters.value.includes(task.status))
  }

  // Apply priority filters
  if (priorityFilters.value.length > 0) {
    filtered = filtered.filter(task => priorityFilters.value.includes(task.priority))
  }

  // Apply agent filters
  if (agentFilters.value.length > 0) {
    filtered = filtered.filter(task => agentFilters.value.includes(task.assignee))
  }

  return filtered
})

const tasksByStatus = computed(() => {
  const grouped = {}
  taskStatuses.forEach(status => {
    grouped[status] = filteredTasks.value.filter(task => task.status === status)
  })
  return grouped
})

// Methods
const setView = (newView) => {
  view.value = newView
}

const toggleStatusFilter = (status) => {
  const index = statusFilters.value.indexOf(status)
  if (index > -1) {
    statusFilters.value.splice(index, 1)
  } else {
    statusFilters.value.push(status)
  }
}

const togglePriorityFilter = (priority) => {
  const index = priorityFilters.value.indexOf(priority)
  if (index > -1) {
    priorityFilters.value.splice(index, 1)
  } else {
    priorityFilters.value.push(priority)
  }
}

const toggleAgentFilter = (agentId) => {
  const index = agentFilters.value.indexOf(agentId)
  if (index > -1) {
    agentFilters.value.splice(index, 1)
  } else {
    agentFilters.value.push(agentId)
  }
}

const clearFilters = () => {
  statusFilters.value = ['pending', 'in-progress', 'completed']
  priorityFilters.value = ['low', 'medium', 'high', 'critical']
  agentFilters.value = []
}

const addTask = () => {
  if (!newTask.value.title.trim()) return

  const newTaskObj = {
    id: tasks.value.length + 1,
    title: newTask.value.title.trim(),
    description: newTask.value.description,
    status: 'pending',
    priority: newTask.value.priority,
    assignee: newTask.value.assignee,
    createdAt: new Date(),
    completedAt: null
  }

  tasks.value.push(newTaskObj)

  // Reset form
  newTask.value = {
    title: '',
    description: '',
    priority: 'medium',
    assignee: null
  }
  showAddTask.value = false
}

const toggleTaskStatus = (task) => {
  if (task.status === 'completed') {
    task.status = 'pending'
    task.completedAt = null
  } else {
    task.status = 'completed'
    task.completedAt = new Date()
  }
}

const editTask = (task) => {
  // In a real app, this would open an edit modal
  console.log('Editing task:', task)
}

const deleteTask = (task) => {
  const index = tasks.value.findIndex(t => t.id === task.id)
  if (index > -1) {
    tasks.value.splice(index, 1)
  }
}

const getTaskClass = (task) => {
  return {
    [`task-${task.status}`]: true,
    [`task-${task.priority}`]: true,
    'task-assigned': task.assignee !== null
  }
}

const getAgentName = (agentId) => {
  const agent = availableAgents.value.find(a => a.id === agentId)
  return agent ? agent.name : 'Unassigned'
}

const formatDate = (date) => {
  const now = new Date()
  const diff = now - new Date(date)
  
  if (diff < 60000) return 'Just now'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`
  return `${Math.floor(diff / 86400000)}d ago`
}

const getEmptyMessage = () => {
  if (statusFilters.value.length === 1 && statusFilters.value[0] === 'completed') {
    return 'All tasks are completed! Great job!'
  }
  if (filteredTasks.value.length === 0) {
    return 'Create your first task to get started'
  }
  return 'No tasks match your current filters'
}

// Kanban columns
const kanbanColumns = ['pending', 'in-progress', 'completed']
</script>

<style scoped>
.task-grid-enhanced {
  background: white;
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  overflow: hidden;
}

.grid-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.header-left {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.grid-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.quick-stats {
  display: flex;
  gap: 16px;
}

.quick-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-number {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
}

.stat-label {
  font-size: 11px;
  color: #6B7280;
  font-weight: 500;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.view-controls {
  display: flex;
  gap: 4px;
  padding: 4px;
  background: #E5E7EB;
  border-radius: 6px;
}

.view-btn {
  padding: 6px 8px;
  border: none;
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #6B7280;
}

.view-btn:hover {
  background: #F3F4F6;
}

.view-btn.active {
  background: white;
  color: #111827;
}

.btn-add-task {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-add-task:hover {
  background: #2563EB;
}

.add-task-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
}

.modal-content {
  position: relative;
  background: white;
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #E5E7EB;
}

.modal-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.btn-close {
  background: none;
  border: none;
  font-size: 20px;
  color: #6B7280;
  cursor: pointer;
  padding: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.btn-close:hover {
  background: #F3F4F6;
  color: #374151;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 4px;
}

.form-input,
.form-textarea,
.form-select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #D1D5DB;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.2s ease;
}

.form-input:focus,
.form-textarea:focus,
.form-select:focus {
  outline: none;
  border-color: #3B82F6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 20px;
  border-top: 1px solid #E5E7EB;
}

.task-filters {
  padding: 16px 20px;
  border-bottom: 1px solid #E5E7EB;
  background: #F9FAFB;
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-label {
  font-size: 12px;
  font-weight: 500;
  color: #374151;
  min-width: 50px;
}

.filter-btn {
  padding: 4px 8px;
  border: 1px solid #D1D5DB;
  background: white;
  border-radius: 4px;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.filter-btn:hover {
  background: #F3F4F6;
  border-color: #9CA3AF;
}

.filter-btn.active {
  background: #3B82F6;
  color: white;
  border-color: #3B82F6;
}

.task-content {
  padding: 20px;
  min-height: 400px;
}

/* Grid View */
.grid-view {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.task-card {
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  padding: 16px;
  transition: all 0.2s ease;
  position: relative;
}

.task-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: #D1D5DB;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.task-priority {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.task-priority.high {
  background: #EF4444;
}

.task-priority.medium {
  background: #F59E0B;
}

.task-priority.low {
  background: #10B981;
}

.task-priority.critical {
  background: #8B5CF6;
}

.task-status {
  font-size: 10px;
  font-weight: 500;
  padding: 2px 6px;
  border-radius: 4px;
  text-transform: uppercase;
}

.task-status.pending {
  background: #FEF3C7;
  color: #92400E;
}

.task-status.in-progress {
  background: #DBEAFE;
  color: #1E40AF;
}

.task-status.completed {
  background: #D1FAE5;
  color: #065F46;
}

.task-title {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 8px;
  line-height: 1.4;
}

.task-description {
  font-size: 12px;
  color: #6B7280;
  line-height: 1.4;
  margin-bottom: 12px;
}

.task-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  font-size: 11px;
  color: #6B7280;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.task-actions {
  display: flex;
  gap: 6px;
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #F3F4F6;
}

.btn-action {
  padding: 4px 8px;
  background: #F3F4F6;
  color: #374151;
  border: 1px solid #D1D5DB;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-action:hover {
  background: #E5E7EB;
  border-color: #9CA3AF;
}

.btn-action.danger {
  color: #DC2626;
}

.btn-action.danger:hover {
  background: #FEE2E2;
  border-color: #FCA5A5;
}

/* List View */
.list-view {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-row {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.task-row:hover {
  background: #F9FAFB;
}

.task-checkbox {
  margin-right: 12px;
}

.task-content {
  flex: 1;
  min-width: 0;
}

.task-title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

/* Kanban View */
.kanban-view {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px;
}

.kanban-column {
  background: #F9FAFB;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  min-height: 400px;
}

.column-header {
  padding: 12px 16px;
  border-bottom: 1px solid #E5E7EB;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.column-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #111827;
}

.column-count {
  font-size: 12px;
  color: #6B7280;
  background: white;
  padding: 2px 6px;
  border-radius: 4px;
}

.column-content {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-card-kanban {
  border: 1px solid #E5E7EB;
  border-radius: 6px;
  padding: 12px;
  background: white;
}

/* Empty State */
.empty-tasks {
  text-align: center;
  padding: 60px 20px;
  color: #6B7280;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-tasks h4 {
  font-size: 16px;
  font-weight: 600;
  color: #374151;
  margin: 0 0 8px 0;
}

.empty-tasks p {
  font-size: 14px;
  margin: 0 0 16px 0;
  line-height: 1.4;
}

@media (max-width: 768px) {
  .grid-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .quick-stats {
    flex-direction: row;
    gap: 12px;
  }
  
  .header-right {
    width: 100%;
    justify-content: space-between;
  }
  
  .task-filters {
    flex-direction: column;
    gap: 8px;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .grid-view {
    grid-template-columns: 1fr;
  }
  
  .kanban-view {
    grid-template-columns: 1fr;
  }
}
</style>