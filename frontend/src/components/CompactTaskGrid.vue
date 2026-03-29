<template>
  <div class="compact-task-grid">
    <!-- Add Task Button (small) -->
    <div class="add-task-quick">
      <button 
        @click="$emit('add-task')" 
        class="btn-quick-add"
        title="Add New Task"
      >
        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
      </button>
    </div>

    <!-- Grid View -->
    <div v-if="view === 'grid'" class="grid-view">
      <div class="grid-header">
        <div class="sort-controls">
          <select v-model="sortBy" class="sort-select">
            <option value="priority">Priority</option>
            <option value="date">Date Created</option>
            <option value="due">Due Date</option>
            <option value="status">Status</option>
          </select>
          <button @click="toggleSortOrder" class="sort-btn">
            {{ sortOrder === 'asc' ? '↑' : '↓' }}
          </button>
        </div>
        <div class="view-options">
          <button 
            @click="setGrouping('none')" 
            :class="['group-btn', grouping === 'none' ? 'active' : '']"
            title="No Grouping"
          >
            <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
            </svg>
          </button>
          <button 
            @click="setGrouping('priority')" 
            :class="['group-btn', grouping === 'priority' ? 'active' : '']"
            title="Group by Priority"
          >
            ⚡
          </button>
          <button 
            @click="setGrouping('status')" 
            :class="['group-btn', grouping === 'status' ? 'active' : '']"
            title="Group by Status"
          >
            📊
          </button>
        </div>
      </div>

      <!-- Task Groups -->
      <div v-if="groupedTasks.length > 0" class="task-groups">
        <div v-for="group in groupedTasks" :key="group.name" class="task-group">
          <div class="group-header">
            <h4>{{ group.name }}</h4>
            <span class="group-count">{{ group.tasks.length }}</span>
          </div>
          <div class="task-cards">
            <div 
              v-for="task in group.tasks" 
              :key="task.id"
              class="compact-task-card"
              :class="getTaskClass(task)"
            >
              <div class="task-checkbox">
                <input 
                  type="checkbox" 
                  :checked="task.status === 'completed'"
                  @change="toggleTaskStatus(task)"
                  @click.stop
                />
              </div>
              <div class="task-content">
                <div class="task-title">{{ task.title }}</div>
                <div class="task-meta">
                  <span v-if="task.priority" class="meta-priority" :class="task.priority">
                    ⚡ {{ task.priority }}
                  </span>
                  <span v-if="task.assignee" class="meta-assignee">
                    🤖 {{ getAgentName(task.assignee) }}
                  </span>
                  <span class="meta-date">
                    📅 {{ formatDate(task.createdAt) }}
                  </span>
                </div>
              </div>
              <div class="task-actions">
                <button @click="editTask(task)" class="btn-compact-action" title="Edit">
                  ✏️
                </button>
                <button @click="deleteTask(task)" class="btn-compact-action danger" title="Delete">
                  🗑️
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="empty-grid">
        <div class="empty-icon">📋</div>
        <p>No tasks found</p>
        <button @click="$emit('add-task')" class="btn-primary">
          Create Your First Task
        </button>
      </div>
    </div>

    <!-- List View -->
    <div v-if="view === 'list'" class="list-view">
      <div class="list-header">
        <div class="row-checkbox">
          <input 
            type="checkbox" 
            :checked="allTasksCompleted"
            @change="toggleAllTasks"
          />
        </div>
        <div class="row-title">Task</div>
        <div class="row-priority">Priority</div>
        <div class="row-assignee">Assignee</div>
        <div class="row-date">Created</div>
        <div class="row-actions">Actions</div>
      </div>

      <div class="task-list">
        <div 
          v-for="task in sortedTasks" 
          :key="task.id"
          class="task-list-item"
          :class="getTaskClass(task)"
        >
          <div class="row-checkbox">
            <input 
              type="checkbox" 
              :checked="task.status === 'completed'"
              @change="toggleTaskStatus(task)"
              @click.stop
            />
          </div>
          <div class="row-content">
            <div class="task-title">{{ task.title }}</div>
            <div class="task-description">{{ task.description?.substring(0, 50) || 'No description' }}...</div>
          </div>
          <div class="row-priority">
            <span class="priority-badge" :class="task.priority">{{ task.priority }}</span>
          </div>
          <div class="row-assignee">
            <span v-if="task.assignee" class="assignee-badge">{{ getAgentName(task.assignee) }}</span>
            <span v-else class="unassigned">Unassigned</span>
          </div>
          <div class="row-date">
            {{ formatDate(task.createdAt) }}
          </div>
          <div class="row-actions">
            <button @click="editTask(task)" class="btn-list-action" title="Edit">
              ✏️
            </button>
            <button @click="deleteTask(task)" class="btn-list-action danger" title="Delete">
              🗑️
            </button>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="sortedTasks.length === 0" class="empty-list">
          <div class="empty-icon">📋</div>
          <p>No tasks found</p>
          <button @click="$emit('add-task')" class="btn-primary">
            Create Your First Task
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const props = defineProps({
  view: {
    type: String,
    default: 'grid',
    validator: (value) => ['grid', 'list'].includes(value)
  },
  filter: {
    type: String,
    default: 'all'
  }
})

const emit = defineEmits(['add-task', 'task-updated'])

// Mock task data
const tasks = ref([
  {
    id: 1,
    title: 'Setup database schema',
    description: 'Design and implement the database schema for the task management system',
    status: 'in-progress',
    priority: 'high',
    assignee: 1,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 2),
    completedAt: null
  },
  {
    id: 2,
    title: 'Create API endpoints',
    description: 'Implement REST API endpoints for task CRUD operations',
    status: 'pending',
    priority: 'high',
    assignee: 1,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 4),
    completedAt: null
  },
  {
    id: 3,
    title: 'Design user interface',
    description: 'Create wireframes and mockups for the dashboard',
    status: 'completed',
    priority: 'medium',
    assignee: 2,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24),
    completedAt: new Date(Date.now() - 1000 * 60 * 60 * 6)
  },
  {
    id: 4,
    title: 'Write documentation',
    description: 'Create comprehensive documentation for API usage',
    status: 'pending',
    priority: 'low',
    assignee: null,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 6),
    completedAt: null
  },
  {
    id: 5,
    title: 'Fix authentication bug',
    description: 'Resolve login issue with JWT token expiration',
    status: 'in-progress',
    priority: 'critical',
    assignee: 1,
    createdAt: new Date(Date.now() - 1000 * 60 * 30),
    completedAt: null
  }
])

// Mock agent data
const agents = ref([
  { id: 1, name: 'code-assistant', type: 'developer' },
  { id: 2, name: 'doc-writer', type: 'general' },
  { id: 3, name: 'design-bot', type: 'designer' }
])

// Sort and filter state
const sortBy = ref('priority')
const sortOrder = ref('desc')
const grouping = ref('none')

// Computed properties
const filteredTasks = computed(() => {
  let filtered = tasks.value

  // Apply filter
  if (props.filter === 'active') {
    filtered = filtered.filter(task => task.status !== 'completed')
  } else if (props.filter === 'completed') {
    filtered = filtered.filter(task => task.status === 'completed')
  } else if (props.filter === 'my-tasks') {
    filtered = filtered.filter(task => task.assignee === 1) // Assuming current user is agent 1
  }

  return filtered
})

const sortedTasks = computed(() => {
  const filtered = [...filteredTasks.value]
  
  return filtered.sort((a, b) => {
    let comparison = 0
    
    switch (sortBy.value) {
      case 'priority':
        const priorityOrder = { critical: 4, high: 3, medium: 2, low: 1 }
        comparison = priorityOrder[a.priority] - priorityOrder[b.priority]
        break
      case 'date':
        comparison = new Date(b.createdAt) - new Date(a.createdAt)
        break
      case 'status':
        comparison = a.status.localeCompare(b.status)
        break
      default:
        comparison = 0
    }
    
    return sortOrder.value === 'asc' ? comparison : -comparison
  })
})

const groupedTasks = computed(() => {
  if (grouping.value === 'none') {
    return [{ name: 'All Tasks', tasks: sortedTasks.value }]
  }

  const groups = {}
  const groupNames = {
    priority: { critical: 'Critical', high: 'High', medium: 'Medium', low: 'Low' },
    status: { pending: 'Pending', 'in-progress': 'In Progress', completed: 'Completed' }
  }

  sortedTasks.value.forEach(task => {
    const key = grouping.value
    const groupName = groupNames[key][task[key]] || task[key]
    
    if (!groups[groupName]) {
      groups[groupName] = []
    }
    groups[groupName].push(task)
  })

  return Object.entries(groups).map(([name, tasks]) => ({
    name,
    tasks
  }))
})

const allTasksCompleted = computed(() => {
  return filteredTasks.value.length > 0 && 
    filteredTasks.value.every(task => task.status === 'completed')
})

// Methods
const toggleSortOrder = () => {
  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
}

const setGrouping = (group) => {
  grouping.value = group
}

const toggleTaskStatus = (task) => {
  task.status = task.status === 'completed' ? 'pending' : 'completed'
  task.completedAt = task.status === 'completed' ? new Date() : null
  emit('task-updated', task)
}

const toggleAllTasks = () => {
  const allCompleted = allTasksCompleted.value
  filteredTasks.value.forEach(task => {
    if (task.status === completed !== allCompleted) {
      toggleTaskStatus(task)
    }
  })
}

const editTask = (task) => {
  // In a real app, this would open an edit modal
  console.log('Editing task:', task)
}

const deleteTask = (task) => {
  const index = tasks.value.findIndex(t => t.id === task.id)
  if (index > -1) {
    tasks.value.splice(index, 1)
    emit('task-updated', task)
  }
}

const getTaskClass = (task) => {
  return {
    [`task-${task.status}`]: true,
    [`task-${task.priority}`]: true,
    'task-completed': task.status === 'completed'
  }
}

const getAgentName = (agentId) => {
  const agent = agents.value.find(a => a.id === agentId)
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
</script>

<style scoped>
.compact-task-grid {
  position: relative;
}

.add-task-quick {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10;
}

.btn-quick-add {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #10B981;
  color: white;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.btn-quick-add:hover {
  background: #059669;
  transform: scale(1.05);
}

.grid-view {
  min-height: 400px;
}

.grid-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #E5E7EB;
  background: #F9FAFB;
  margin-bottom: 16px;
}

.sort-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sort-select {
  padding: 4px 8px;
  border: 1px solid #D1D5DB;
  border-radius: 4px;
  font-size: 12px;
  background: white;
  cursor: pointer;
}

.sort-btn {
  padding: 4px 8px;
  background: white;
  border: 1px solid #D1D5DB;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  min-width: 24px;
  text-align: center;
}

.sort-btn:hover {
  background: #F3F4F6;
}

.view-options {
  display: flex;
  gap: 4px;
}

.group-btn {
  padding: 6px 8px;
  border: 1px solid #D1D5DB;
  background: white;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 12px;
}

.group-btn:hover {
  background: #F3F4F6;
}

.group-btn.active {
  background: #10B981;
  color: white;
  border-color: #10B981;
}

.task-groups {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.task-group {
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  overflow: hidden;
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #F9FAFB;
  border-bottom: 1px solid #E5E7EB;
}

.group-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #111827;
}

.group-count {
  font-size: 12px;
  color: #6B7280;
  background: white;
  padding: 2px 6px;
  border-radius: 4px;
}

.task-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 8px;
  padding: 12px;
}

.compact-task-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border: 1px solid #E5E7EB;
  border-radius: 6px;
  background: white;
  transition: all 0.2s ease;
  cursor: pointer;
}

.compact-task-card:hover {
  background: #F9FAFB;
  border-color: #D1D5DB;
}

.compact-task-card.task-completed {
  opacity: 0.6;
}

.task-checkbox {
  flex-shrink: 0;
}

.task-checkbox input {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.task-content {
  flex: 1;
  min-width: 0;
}

.task-title {
  font-size: 13px;
  font-weight: 500;
  color: #111827;
  margin-bottom: 4px;
  line-height: 1.3;
}

.task-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  font-size: 11px;
}

.meta-priority {
  padding: 2px 6px;
  border-radius: 3px;
  font-weight: 500;
  text-transform: capitalize;
}

.meta-priority.high {
  background: #FEF3C7;
  color: #92400E;
}

.meta-priority.medium {
  background: #DBEAFE;
  color: #1E40AF;
}

.meta-priority.low {
  background: #D1FAE5;
  color: #065F46;
}

.meta-priority.critical {
  background: #FEE2E2;
  color: #991B1B;
}

.meta-assignee {
  color: #6B7280;
}

.meta-date {
  color: #9CA3AF;
}

.task-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.btn-compact-action {
  padding: 4px 6px;
  background: transparent;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s ease;
}

.btn-compact-action:hover {
  background: #F3F4F6;
}

.btn-compact-action.danger {
  color: #DC2626;
}

.btn-compact-action.danger:hover {
  background: #FEE2E2;
}

.empty-grid {
  text-align: center;
  padding: 60px 20px;
  color: #6B7280;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-grid p {
  font-size: 14px;
  margin: 0 0 16px 0;
}

/* List View */
.list-view {
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  overflow: hidden;
}

.list-header {
  display: grid;
  grid-template-columns: 40px 2fr 80px 120px 80px 80px;
  gap: 12px;
  padding: 12px 16px;
  background: #F9FAFB;
  border-bottom: 1px solid #E5E7EB;
  font-size: 12px;
  font-weight: 600;
  color: #6B7280;
}

.list-header div {
  display: flex;
  align-items: center;
}

.list-header .row-title {
  font-weight: 600;
  color: #374151;
}

.task-list {
  max-height: 500px;
  overflow-y: auto;
}

.task-list-item {
  display: grid;
  grid-template-columns: 40px 2fr 80px 120px 80px 80px;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid #F3F4F6;
  transition: all 0.2s ease;
}

.task-list-item:hover {
  background: #F9FAFB;
}

.task-list-item:last-child {
  border-bottom: none;
}

.task-list-item.task-completed {
  opacity: 0.6;
}

.row-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.task-description {
  font-size: 11px;
  color: #9CA3AF;
}

.priority-badge {
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 11px;
  font-weight: 500;
  text-transform: capitalize;
  text-align: center;
}

.priority-badge.high {
  background: #FEF3C7;
  color: #92400E;
}

.priority-badge.medium {
  background: #DBEAFE;
  color: #1E40AF;
}

.priority-badge.low {
  background: #D1FAE5;
  color: #065F46;
}

.priority-badge.critical {
  background: #FEE2E2;
  color: #991B1B;
}

.assignee-badge {
  padding: 2px 6px;
  background: #EDE9FE;
  color: #5B21B6;
  border-radius: 3px;
  font-size: 11px;
}

.unassigned {
  color: #9CA3AF;
  font-size: 11px;
}

.btn-list-action {
  padding: 4px 6px;
  background: transparent;
  border: none;
  border-radius: 3px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s ease;
}

.btn-list-action:hover {
  background: #F3F4F6;
}

.btn-list-action.danger {
  color: #DC2626;
}

.btn-list-action.danger:hover {
  background: #FEE2E2;
}

.empty-list {
  text-align: center;
  padding: 60px 20px;
  color: #6B7280;
}

.empty-list .empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-list p {
  font-size: 14px;
  margin: 0 0 16px 0;
}

.btn-primary {
  padding: 8px 16px;
  background: #10B981;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-primary:hover {
  background: #059669;
}

@media (max-width: 768px) {
  .grid-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .task-cards {
    grid-template-columns: 1fr;
  }
  
  .list-header {
    grid-template-columns: 30px 1fr 60px 80px;
    font-size: 11px;
    gap: 8px;
  }
  
  .task-list-item {
    grid-template-columns: 30px 1fr 60px 80px;
    font-size: 12px;
    gap: 8px;
    padding: 8px 12px;
  }
}
</style>