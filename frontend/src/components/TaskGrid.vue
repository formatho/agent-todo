<template>
  <div class="task-grid-container">
    <!-- Grid Header -->
    <div class="grid-header">
      <div>
        <h2 class="grid-title">
          Tasks
          <span class="task-count">({{ filteredTasks.length }})</span>
        </h2>
        <div v-if="filters.project_id" class="project-filter-indicator">
          <span class="filter-icon">📁</span>
          <span class="filter-text">
            Filtering by: {{ getProjectName(filters.project_id) }}
          </span>
          <button @click="clearProjectFilter" class="btn-clear-filter">×</button>
        </div>
        <div v-if="filters.agent_id" class="project-filter-indicator agent-filter">
          <span class="filter-icon">🤖</span>
          <span class="filter-text">
            Filtering by: {{ getAgentName(filters.agent_id) }}
          </span>
          <button @click="clearAgentFilter" class="btn-clear-filter">×</button>
        </div>
      </div>

      <ViewToggle v-model="viewMode" :options="viewOptions" @change="handleViewChange" />

      <div class="create-task-wrapper">
        <button 
          @click="handleCreateTask" 
          class="btn-create" 
          :disabled="projectStore.activeProjects.length === 0"
          :class="{ 'has-tips': showTooltips }"
        >
          + Create Task
        </button>
        
        <!-- Tooltip for new users -->
        <div v-if="showTooltips && projectStore.activeProjects.length === 0" class="tooltip">
          <div class="tooltip-content">
            <div class="tooltip-header">
              <span class="tooltip-icon">📝</span>
              <span class="tooltip-title">Get Started</span>
            </div>
            <p class="tooltip-text">Create a project first to organize your tasks, then add tasks to it.</p>
            <button @click="navigateToProjects" class="tooltip-action">Create Project</button>
          </div>
          <div class="tooltip-arrow"></div>
        </div>
        
        <!-- Tips for creating tasks -->
        <div v-if="showTooltips && projectStore.activeProjects.length > 0" class="tips">
          <div class="tips-content">
            <span class="tips-icon">💡</span>
            <p class="tips-text">
              Click here to create tasks. Assign them to AI agents for automatic completion!
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Warning if no projects -->
    <div v-if="projectStore.activeProjects.length === 0" class="no-projects-warning">
      <span class="warning-icon">⚠️</span>
      <span>No active projects. You need a project to create tasks.</span>
      <router-link to="/projects" class="btn-create-project">Create Project</router-link>
    </div>

    <!-- Filters Bar -->
    <div class="filters-bar">
      <div class="search-box">
        <span class="search-icon">🔍</span>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search tasks..."
          class="search-input"
        />
      </div>

      <label class="toggle-completed">
        <input type="checkbox" v-model="showCompleted" />
        <span>Show completed tasks</span>
      </label>

      <button @click="showFilters = !showFilters" class="btn-filters">
        <span>⚙️</span>
        <span>Filters</span>
        <span v-if="hasActiveFilters" class="filter-indicator">{{ activeFilterCount }}</span>
      </button>
    </div>

    <!-- Expanded Filters -->
    <div v-if="showFilters" class="filters-panel">
      <div class="filter-group">
        <label>Project</label>
        <select v-model="filters.project_id" @change="applyFilters" class="filter-select">
          <option value="">All Projects</option>
          <option v-for="project in projectStore.activeProjects" :key="project.id" :value="project.id">
            {{ project.name }}
          </option>
        </select>
      </div>

      <div class="filter-group">
        <label>Status</label>
        <select v-model="filters.status" @change="applyFilters" class="filter-select">
          <option value="">All Statuses</option>
          <option value="pending">Pending</option>
          <option value="in_progress">In Progress</option>
          <option value="blocked">Blocked</option>
          <option value="completed">Completed</option>
          <option value="failed">Failed</option>
        </select>
      </div>

      <div class="filter-group">
        <label>Priority</label>
        <select v-model="filters.priority" @change="applyFilters" class="filter-select">
          <option value="">All Priorities</option>
          <option value="critical">Critical</option>
          <option value="high">High</option>
          <option value="medium">Medium</option>
          <option value="low">Low</option>
        </select>
      </div>

      <div class="filter-group">
        <label>Agent</label>
        <select v-model="filters.agent_id" @change="applyFilters" class="filter-select">
          <option value="">All Agents</option>
          <option v-for="agent in agents" :key="agent.id" :value="agent.id">
            {{ agent.name }}
          </option>
        </select>
      </div>

      <button @click="clearFilters" class="btn-clear">Clear All</button>
    </div>

    <!-- Grid View -->
    <transition-group
      v-if="viewMode === 'grid'"
      name="task-card"
      tag="div"
      class="tasks-grid"
    >
      <TaskCard
        v-for="task in filteredTasks"
        :key="task.id"
        :task="task"
      />
    </transition-group>

    <!-- List View (Table) -->
    <div v-else-if="viewMode === 'list'" class="tasks-list">
      <table class="task-table">
        <thead>
          <tr>
            <th>Task</th>
            <th>Project</th>
            <th>Status</th>
            <th>Priority</th>
            <th>Agent</th>
            <th>Due Date</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="task in filteredTasks" :key="task.id">
            <td>
              <div class="task-cell-title">{{ task.title }}</div>
              <div class="task-cell-desc">{{ truncatedDesc(task.description) }}</div>
            </td>
            <td>
              <span class="project-badge">{{ task.project?.name || '-' }}</span>
            </td>
            <td>
              <span :class="['status-badge', task.status]">
                {{ formatStatus(task.status) }}
              </span>
            </td>
            <td>
              <span :class="['priority-badge', task.priority]">
                {{ task.priority }}
              </span>
            </td>
            <td>
              <AgentAvatar
                v-if="task.assigned_agent"
                :agent="task.assigned_agent"
                size="small"
              />
              <span v-else class="unassigned">-</span>
            </td>
            <td>{{ formatDate(task.due_date) }}</td>
            <td>
              <router-link :to="`/tasks/${task.id}`" class="btn-view-small">
                View
              </router-link>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Kanban Board View -->
    <KanbanBoard v-else-if="viewMode === 'board'" :projectFilter="filters.project_id" />

    <!-- Empty State -->
    <div v-if="filteredTasks.length === 0" class="empty-state">
      <div class="empty-icon">📋</div>
      <h3>No tasks found</h3>
      <p>{{ emptyMessage }}</p>
      <button v-if="!hasActiveFilters && projectStore.activeProjects.length > 0" @click="handleCreateTask" class="btn-create-empty">
        Create your first task
      </button>
    </div>

    <!-- Create Task Modal -->
    <TaskModal
      v-if="showCreateModal"
      @close="showCreateModal = false"
      @saved="handleTaskCreated"
      :preselectedProjectId="filters.project_id"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTaskStore } from '../stores/tasks'
import { useAgentStore } from '../stores/agents'
import { useProjectStore } from '../stores/projects'
import TaskCard from './TaskCard.vue'
import AgentAvatar from './AgentAvatar.vue'
import TaskModal from './TaskModal.vue'
import ViewToggle from './ViewToggle.vue'
import KanbanBoard from './KanbanBoard.vue'

const route = useRoute()
const router = useRouter()
const taskStore = useTaskStore()
const agentStore = useAgentStore()
const projectStore = useProjectStore()

const viewMode = ref('grid')
const searchQuery = ref('')
const showCompleted = ref(true) // Show completed tasks by default
const showTooltips = ref(true) // Enable tooltips for new users

const viewOptions = [
  { value: 'grid', label: 'Grid View', icon: '⊞' },
  { value: 'list', label: 'List View', icon: '☰' },
  { value: 'board', label: 'Board View', icon: '▦' }
]
const showFilters = ref(false)
const showCreateModal = ref(false)

const filters = ref({
  status: '',
  priority: '',
  agent_id: '',
  project_id: '',
  search: ''
})

// Priority order mapping (higher number = higher priority)
const priorityOrder = {
  critical: 4,
  high: 3,
  medium: 2,
  low: 1
}

const filteredTasks = computed(() => {
  let tasks = Array.isArray(taskStore.tasks) ? taskStore.tasks : []

  // Filter out completed tasks if toggle is off
  if (!showCompleted.value) {
    tasks = tasks.filter(task => task.status !== 'completed')
  }

  // Sort by priority (highest first)
  tasks.sort((a, b) => {
    const priorityA = priorityOrder[a.priority] || 0
    const priorityB = priorityOrder[b.priority] || 0
    return priorityB - priorityA
  })

  return tasks
})

const agents = computed(() => Array.isArray(agentStore.agents) ? agentStore.agents : [])

const hasActiveFilters = computed(() => {
  return filters.value.status ||
         filters.value.priority ||
         filters.value.agent_id ||
         filters.value.project_id ||
         searchQuery.value
})

const activeFilterCount = computed(() => {
  let count = 0
  if (filters.value.status) count++
  if (filters.value.priority) count++
  if (filters.value.agent_id) count++
  if (filters.value.project_id) count++
  if (searchQuery.value) count++
  return count
})

const emptyMessage = computed(() => {
  if (hasActiveFilters.value) {
    return 'No tasks match your current filters. Try adjusting your search criteria.'
  }
  if (projectStore.activeProjects.length === 0) {
    return 'Create a project first, then add tasks.'
  }
  return 'Get started by creating your first task.'
})

onMounted(async () => {
  await projectStore.fetchProjects({ status: 'active' })
  await agentStore.fetchAgents()

  // Check for query parameters and apply them
  const projectIdFromUrl = route.query.project_id
  const agentIdFromUrl = route.query.agent_id
  const statusFromUrl = route.query.status
  const priorityFromUrl = route.query.priority
  const searchFromUrl = route.query.search

  // Apply filters from URL
  if (projectIdFromUrl) filters.value.project_id = projectIdFromUrl
  if (agentIdFromUrl) filters.value.agent_id = agentIdFromUrl
  if (statusFromUrl) filters.value.status = statusFromUrl
  if (priorityFromUrl) filters.value.priority = priorityFromUrl
  if (searchFromUrl) {
    searchQuery.value = searchFromUrl
    filters.value.search = searchFromUrl
  }

  // Build params for API call
  const params = {}
  if (projectIdFromUrl) params.project_id = projectIdFromUrl
  if (agentIdFromUrl) params.agent_id = agentIdFromUrl

  if (Object.keys(params).length > 0) {
    await taskStore.fetchTasks(params)
  } else {
    await taskStore.fetchTasks()
  }
})

// Watch for route query changes
watch(
  () => route.query,
  async (newQuery) => {
    // Update filters from URL
    filters.value.project_id = newQuery.project_id || ''
    filters.value.agent_id = newQuery.agent_id || ''
    filters.value.status = newQuery.status || ''
    filters.value.priority = newQuery.priority || ''
    
    if (newQuery.search && newQuery.search !== searchQuery.value) {
      searchQuery.value = newQuery.search
      filters.value.search = newQuery.search
    }

    // Fetch tasks with appropriate filters
    const params = {}
    if (newQuery.project_id) params.project_id = newQuery.project_id
    if (newQuery.agent_id) params.agent_id = newQuery.agent_id

    if (Object.keys(params).length > 0) {
      await taskStore.fetchTasks(params)
    } else {
      await taskStore.fetchTasks()
    }
  },
  { deep: true }
)

watch(searchQuery, (newVal) => {
  filters.value.search = newVal
  applyFilters()
})

const handleViewChange = (newView) => {
  viewMode.value = newView
}

const handleToggleCompleted = () => {
  // No need to fetch, just update local filter
  // The computed property will automatically update
}

const getProjectName = (projectId) => {
  const project = projectStore.projects.find(p => p.id === projectId)
  return project ? project.name : 'Unknown Project'
}

const getAgentName = (agentId) => {
  const agent = agents.value.find(a => a.id === agentId)
  return agent ? agent.name : 'Unknown Agent'
}

const clearProjectFilter = async () => {
  filters.value.project_id = ''
  // Update URL to remove only the project_id query parameter
  const query = { ...route.query }
  delete query.project_id
  await router.push({ path: route.path, query })
  await taskStore.fetchTasks()
}

const clearAgentFilter = async () => {
  filters.value.agent_id = ''
  // Update URL to remove only the agent_id query parameter
  const query = { ...route.query }
  delete query.agent_id
  await router.push({ path: route.path, query })
  await taskStore.fetchTasks()
}

const applyFilters = async () => {
  // Update URL with current filters
  const query = {}
  
  if (filters.value.status) query.status = filters.value.status
  if (filters.value.priority) query.priority = filters.value.priority
  if (filters.value.agent_id) query.agent_id = filters.value.agent_id
  if (filters.value.project_id) query.project_id = filters.value.project_id
  if (searchQuery.value) query.search = searchQuery.value
  
  // Update URL without reloading
  await router.push({ path: route.path, query })
  
  taskStore.setFilters(filters.value)
  taskStore.fetchTasks()
}

const clearFilters = async () => {
  filters.value = {
    status: '',
    priority: '',
    agent_id: '',
    project_id: '',
    search: ''
  }
  searchQuery.value = ''
  
  // Clear URL query parameters
  await router.push({ path: route.path })
  
  taskStore.setFilters(filters.value)
  taskStore.fetchTasks()
}

const handleCreateTask = () => {
  if (projectStore.activeProjects.length === 0) {
    return
  }
  showCreateModal.value = true
}

const navigateToProjects = () => {
  router.push('/projects')
}

const handleTaskCreated = () => {
  showCreateModal.value = false
  taskStore.fetchTasks()
}

const formatStatus = (status) => {
  return status.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const formatDate = (dateStr) => {
  if (!dateStr) return 'N/A'
  return new Date(dateStr).toLocaleDateString()
}

const truncatedDesc = (desc) => {
  if (!desc) return ''
  return desc.length > 50 ? desc.substring(0, 50) + '...' : desc
}
</script>

<style scoped>
.task-grid-container {
  width: 100%;
}

/* Grid Header */
.grid-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.grid-title {
  font-size: 28px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 8px 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.task-count {
  font-size: 18px;
  font-weight: 400;
  color: #6B7280;
}

/* Project Filter Indicator */
.project-filter-indicator {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: #E0E7FF;
  border: 1px solid #C7D2FE;
  border-radius: 20px;
  font-size: 13px;
  color: #3730A3;
  margin-top: 8px;
}

.project-filter-indicator.agent-filter {
  background: #DBEAFE;
  border-color: #BFDBFE;
  color: #1E40AF;
}

.filter-icon {
  font-size: 14px;
}

.filter-text {
  font-weight: 500;
}

.btn-clear-filter {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  background: #6366F1;
  color: white;
  border: none;
  border-radius: 50%;
  font-size: 18px;
  line-height: 1;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-clear-filter:hover {
  background: #4F46E5;
  transform: scale(1.1);
}

.btn-create {
  padding: 10px 20px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-create:hover:not(:disabled) {
  background: #2563EB;
  transform: translateY(-1px);
}

.btn-create:disabled {
  background: #9CA3AF;
  cursor: not-allowed;
}

/* No Projects Warning */
.no-projects-warning {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #FEF3C7;
  border: 1px solid #FCD34D;
  border-radius: 8px;
  margin-bottom: 16px;
  color: #92400E;
}

.warning-icon {
  font-size: 20px;
}

.btn-create-project {
  padding: 6px 12px;
  background: #F59E0B;
  color: white;
  border-radius: 6px;
  text-decoration: none;
  font-weight: 600;
  font-size: 13px;
  transition: all 0.2s ease;
}

.btn-create-project:hover {
  background: #D97706;
}

/* Filters Bar */
.filters-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
  align-items: center;
}

.toggle-completed {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: white;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.toggle-completed:hover {
  background: #F9FAFB;
  border-color: #D1D5DB;
}

.toggle-completed input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.search-box {
  flex: 1;
  min-width: 280px;
  position: relative;
  display: flex;
  align-items: center;
  background: white;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  padding: 10px 16px;
}

.search-icon {
  font-size: 18px;
  margin-right: 10px;
}

.search-input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 14px;
}

.search-input::placeholder {
  color: #9CA3AF;
}

.btn-filters {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  background: white;
  border: 1px solid #E5E7EB;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-filters:hover {
  background: #F9FAFB;
  border-color: #D1D5DB;
}

.filter-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 20px;
  background: #EF4444;
  color: white;
  border-radius: 10px;
  font-size: 11px;
  font-weight: 700;
}

/* Filters Panel */
.filters-panel {
  background: white;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 24px;
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-group label {
  font-size: 13px;
  font-weight: 600;
  color: #374151;
}

.filter-select {
  padding: 8px 12px;
  border: 1px solid #D1D5DB;
  border-radius: 6px;
  font-size: 14px;
  background: white;
  cursor: pointer;
  min-width: 160px;
}

.btn-clear {
  padding: 8px 16px;
  background: #FEE2E2;
  color: #991B1B;
  border: 1px solid #FECACA;
  border-radius: 6px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  align-self: flex-end;
}

.btn-clear:hover {
  background: #FECACA;
}

/* Tasks Grid */
.tasks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

@media (max-width: 768px) {
  .tasks-grid {
    grid-template-columns: 1fr;
  }
}

/* Card Animations */
.task-card-enter-active {
  animation: slideIn 0.3s ease;
}

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

/* List View */
.tasks-list {
  background: white;
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  overflow: hidden;
}

.task-table {
  width: 100%;
  border-collapse: collapse;
}

.task-table th {
  background: #F9FAFB;
  padding: 12px 16px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  color: #6B7280;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 1px solid #E5E7EB;
}

.task-table td {
  padding: 16px;
  border-bottom: 1px solid #F3F4F6;
}

.task-table tr:last-child td {
  border-bottom: none;
}

.task-cell-title {
  font-weight: 600;
  color: #111827;
  margin-bottom: 4px;
}

.task-cell-desc {
  font-size: 13px;
  color: #6B7280;
}

.project-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  background: #E0E7FF;
  color: #3730A3;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
}

.status-badge,
.priority-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
}

.status-badge.pending {
  background: #FEF3C7;
  color: #92400E;
}

.status-badge.in_progress {
  background: #DBEAFE;
  color: #1E40AF;
}

.status-badge.blocked {
  background: #EDE9FE;
  color: #5B21B6;
}

.status-badge.completed {
  background: #D1FAE5;
  color: #065F46;
}

.status-badge.failed {
  background: #FEE2E2;
  color: #991B1B;
}

.priority-badge.critical {
  background: #FECACA;
  color: #991B1B;
}

.priority-badge.high {
  background: #FED7AA;
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

.unassigned {
  color: #9CA3AF;
  font-style: italic;
}

.btn-view-small {
  color: #3B82F6;
  text-decoration: none;
  font-weight: 600;
  font-size: 13px;
}

.btn-view-small:hover {
  text-decoration: underline;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-state h3 {
  font-size: 20px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 8px 0;
}

.empty-state p {
  color: #6B7280;
  margin: 0 0 24px 0;
}

.btn-create-empty {
  padding: 12px 24px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-create-empty:hover {
  background: #2563EB;
}

/* Create Task Wrapper with Tooltips */
.create-task-wrapper {
  position: relative;
  display: inline-block;
}

.btn-create.has-tips {
  position: relative;
}

/* Tooltip Styles */
.tooltip {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  z-index: 10;
}

.tooltip-content {
  background: white;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  max-width: 280px;
  min-width: 280px;
}

.tooltip-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.tooltip-icon {
  font-size: 18px;
}

.tooltip-title {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
}

.tooltip-text {
  font-size: 12px;
  color: #6B7280;
  line-height: 1.4;
  margin-bottom: 12px;
}

.tooltip-action {
  width: 100%;
  padding: 8px 12px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tooltip-action:hover {
  background: #2563EB;
}

.tooltip-arrow {
  position: absolute;
  top: -6px;
  right: 20px;
  width: 12px;
  height: 12px;
  background: white;
  border-right: 1px solid #E5E7EB;
  border-top: 1px solid #E5E7EB;
  transform: rotate(45deg);
}

/* Tips Styles */
.tips {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  z-index: 10;
}

.tips-content {
  background: linear-gradient(135deg, #EFF6FF 0%, #DBEAFE 100%);
  border: 1px solid #BFDBFE;
  border-radius: 8px;
  padding: 10px 12px;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.15);
  max-width: 250px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.tips-icon {
  font-size: 14px;
  flex-shrink: 0;
}

.tips-text {
  font-size: 11px;
  color: #1E40AF;
  line-height: 1.3;
  font-weight: 500;
}

/* Responsive Design */
@media (max-width: 768px) {
  .tooltip-content,
  .tips-content {
    max-width: 240px;
  }
  
  .tooltip {
    right: auto;
    left: 50%;
    transform: translateX(-50%);
  }
  
  .tooltip-arrow {
    right: 50%;
    left: auto;
    transform: translateX(50%) rotate(45deg);
  }
}
</style>
