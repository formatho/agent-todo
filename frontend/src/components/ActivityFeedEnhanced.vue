<template>
  <div class="activity-feed-enhanced">
    <div class="feed-header">
      <h3 class="text-lg font-semibold text-gray-900">Recent Activity</h3>
      <div class="feed-actions">
        <button @click="refreshFeed" class="btn-refresh" :disabled="loading">
          <svg v-if="loading" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <svg v-else class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
        <button @click="clearFilters" class="btn-clear" v-if="hasActiveFilters">
          Clear Filters
        </button>
      </div>
    </div>

    <!-- Filter Pills -->
    <div class="filter-pills" v-if="showFilters">
      <div class="filter-group">
        <span class="filter-label">Types:</span>
        <button 
          v-for="type in activityTypes" 
          :key="type"
          @click="toggleTypeFilter(type)"
          :class="['filter-pill', typeFilters.includes(type) ? 'active' : '']"
        >
          {{ type }}
        </button>
      </div>
      <div class="filter-group">
        <span class="filter-label">Status:</span>
        <button 
          v-for="status in activityStatuses" 
          :key="status"
          @click="toggleStatusFilter(status)"
          :class="['filter-pill', statusFilters.includes(status) ? 'active' : '']"
        >
          {{ status }}
        </button>
      </div>
    </div>

    <!-- Activity Items -->
    <div class="activity-list">
      <div 
        v-for="activity in filteredActivities" 
        :key="activity.id"
        class="activity-item"
        :class="getActivityClass(activity)"
      >
        <div class="activity-icon">
          <component :is="getActivityIcon(activity.type)" />
        </div>
        <div class="activity-content">
          <div class="activity-header">
            <span class="activity-title">{{ activity.title }}</span>
            <span class="activity-time">{{ formatTime(activity.timestamp) }}</span>
          </div>
          <div class="activity-description">{{ activity.description }}</div>
          <div class="activity-meta" v-if="activity.meta">
            <span v-if="activity.meta.project" class="meta-tag">
              📁 {{ activity.meta.project }}
            </span>
            <span v-if="activity.meta.agent" class="meta-tag">
              🤖 {{ activity.meta.agent }}
            </span>
            <span v-if="activity.meta.priority" class="meta-tag">
              ⚡ {{ activity.meta.priority }}
            </span>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="!loading && filteredActivities.length === 0" class="empty-activity">
        <div class="empty-icon">📊</div>
        <h4>No activity found</h4>
        <p>{{ getEmptyMessage() }}</p>
        <button @click="clearFilters" class="btn-clear-empty">
          Clear all filters
        </button>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-activity">
        <div v-for="i in 3" :key="i" class="skeleton-activity">
          <div class="skeleton-icon"></div>
          <div class="skeleton-content">
            <div class="skeleton-title"></div>
            <div class="skeleton-text"></div>
          </div>
        </div>
      </div>

      <!-- Load More -->
      <div v-if="!loading && hasMoreActivities" class="load-more">
        <button @click="loadMore" class="btn-load-more">
          Load More Activities
        </button>
      </div>
    </div>

    <!-- Stats Overview -->
    <div class="activity-stats" v-if="showStats">
      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-number">{{ todayActivities }}</div>
          <div class="stat-label">Today</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ weekActivities }}</div>
          <div class="stat-label">This Week</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ totalActivities }}</div>
          <div class="stat-label">Total</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ completedActivities }}</div>
          <div class="stat-label">Completed</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useActivityStore } from '../stores/activity'

const props = defineProps({
  filter: {
    type: String,
    default: 'all'
  }
})

const activityStore = useActivityStore()
const loading = ref(false)
const showFilters = ref(false)
const showStats = ref(true)

// Filter states
const typeFilters = ref([])
const statusFilters = ref([])

// Sample activity data - in real app, this would come from API
const sampleActivities = ref([
  {
    id: 1,
    type: 'task',
    status: 'completed',
    title: 'Task Completed: Setup Database',
    description: 'Agent "dev-assistant" successfully completed database setup',
    timestamp: new Date(Date.now() - 1000 * 60 * 5), // 5 minutes ago
    meta: {
      project: 'Backend Setup',
      agent: 'dev-assistant',
      priority: 'high'
    }
  },
  {
    id: 2,
    type: 'agent',
    status: 'started',
    title: 'Agent Started: code-review-bot',
    description: 'Code review agent has started monitoring repository',
    timestamp: new Date(Date.now() - 1000 * 60 * 15), // 15 minutes ago
    meta: {
      project: 'Code Quality',
      agent: 'code-review-bot'
    }
  },
  {
    id: 3,
    type: 'project',
    status: 'created',
    title: 'Project Created: Marketing Campaign',
    description: 'New project for Q2 marketing initiatives',
    timestamp: new Date(Date.now() - 1000 * 60 * 30), // 30 minutes ago
    meta: {
      project: 'Marketing Campaign'
    }
  },
  {
    id: 4,
    type: 'task',
    status: 'started',
    title: 'Task Started: API Documentation',
    description: 'Documentation agent began working on API docs',
    timestamp: new Date(Date.now() - 1000 * 60 * 45), // 45 minutes ago
    meta: {
      project: 'Documentation',
      agent: 'doc-writer',
      priority: 'medium'
    }
  }
])

const activities = ref(sampleActivities)

// Activity types and statuses for filtering
const activityTypes = ['task', 'agent', 'project', 'system']
const activityStatuses = ['created', 'started', 'completed', 'failed']

// Computed properties
const hasActiveFilters = computed(() => {
  return typeFilters.value.length > 0 || statusFilters.value.length > 0
})

const filteredActivities = computed(() => {
  let filtered = activities.value

  // Apply type filters
  if (typeFilters.value.length > 0) {
    filtered = filtered.filter(activity => typeFilters.value.includes(activity.type))
  }

  // Apply status filters
  if (statusFilters.value.length > 0) {
    filtered = filtered.filter(activity => statusFilters.value.includes(activity.status))
  }

  // Apply main filter (from props)
  if (props.filter !== 'all') {
    filtered = filtered.filter(activity => activity.type === props.filter)
  }

  // Sort by timestamp (newest first)
  return filtered.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp))
})

const todayActivities = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return activities.value.filter(activity => new Date(activity.timestamp) >= today).length
})

const weekActivities = computed(() => {
  const weekAgo = new Date()
  weekAgo.setDate(weekAgo.getDate() - 7)
  return activities.value.filter(activity => new Date(activity.timestamp) >= weekAgo).length
})

const totalActivities = computed(() => activities.value.length)
const completedActivities = computed(() => {
  return activities.value.filter(activity => activity.status === 'completed').length
})

const hasMoreActivities = computed(() => {
  return activities.value.length < 50 // Simulate having more data
})

// Methods
const refreshFeed = async () => {
  loading.value = true
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 1000))
  loading.value = false
}

const loadMore = () => {
  // Simulate loading more activities
  const newActivities = []
  for (let i = activities.value.length + 1; i <= activities.value.length + 5; i++) {
    newActivities.push({
      id: i,
      type: 'task',
      status: 'started',
      title: `Sample Task ${i}`,
      description: `This is a sample activity item for demonstration purposes`,
      timestamp: new Date(Date.now() - 1000 * 60 * i * 10),
      meta: {
        project: `Project ${i % 3 + 1}`,
        agent: `agent-${i}`,
        priority: i % 3 === 0 ? 'high' : i % 3 === 1 ? 'medium' : 'low'
      }
    })
  }
  activities.value.push(...newActivities)
}

const toggleTypeFilter = (type) => {
  const index = typeFilters.value.indexOf(type)
  if (index > -1) {
    typeFilters.value.splice(index, 1)
  } else {
    typeFilters.value.push(type)
  }
}

const toggleStatusFilter = (status) => {
  const index = statusFilters.value.indexOf(status)
  if (index > -1) {
    statusFilters.value.splice(index, 1)
  } else {
    statusFilters.value.push(status)
  }
}

const clearFilters = () => {
  typeFilters.value = []
  statusFilters.value = []
}

const getActivityIcon = (type) => {
  const icons = {
    task: '📋',
    agent: '🤖',
    project: '📁',
    system: '⚙️'
  }
  return icons[type] || '📝'
}

const getActivityClass = (activity) => {
  return {
    [`activity-${activity.type}`]: true,
    [`activity-${activity.status}`]: true,
    'activity-priority-high': activity.meta?.priority === 'high',
    'activity-priority-medium': activity.meta?.priority === 'medium',
    'activity-priority-low': activity.meta?.priority === 'low'
  }
}

const formatTime = (timestamp) => {
  const now = new Date()
  const diff = now - new Date(timestamp)
  
  if (diff < 60000) return 'Just now'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`
  return `${Math.floor(diff / 86400000)}d ago`
}

const getEmptyMessage = () => {
  if (hasActiveFilters.value) {
    return 'No activities match your current filters. Try adjusting them.'
  }
  if (props.filter !== 'all') {
    return `No ${props.filter} activities found yet.`
  }
  return 'No recent activity. Your first action will appear here!'
}

// Auto-refresh every 30 seconds
let refreshInterval
onMounted(() => {
  refreshInterval = setInterval(refreshFeed, 30000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.activity-feed-enhanced {
  background: white;
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  overflow: hidden;
}

.feed-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.feed-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.feed-actions {
  display: flex;
  gap: 8px;
}

.btn-refresh, .btn-clear {
  padding: 6px 12px;
  border: 1px solid #D1D5DB;
  background: white;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 4px;
}

.btn-refresh:hover, .btn-clear:hover {
  background: #F3F4F6;
  border-color: #9CA3AF;
}

.btn-refresh:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.filter-pills {
  padding: 16px 20px;
  border-bottom: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.filter-group:last-child {
  margin-bottom: 0;
}

.filter-label {
  font-size: 12px;
  font-weight: 500;
  color: #6B7280;
  min-width: 40px;
}

.filter-pill {
  padding: 4px 8px;
  border: 1px solid #D1D5DB;
  background: white;
  border-radius: 4px;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.filter-pill:hover {
  background: #F3F4F6;
  border-color: #9CA3AF;
}

.filter-pill.active {
  background: #3B82F6;
  color: white;
  border-color: #3B82F6;
}

.activity-list {
  max-height: 600px;
  overflow-y: auto;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid #F3F4F6;
  transition: all 0.2s ease;
}

.activity-item:hover {
  background: #F9FAFB;
}

.activity-item:last-child {
  border-bottom: none;
}

.activity-icon {
  font-size: 20px;
  flex-shrink: 0;
  width: 24px;
  text-align: center;
}

.activity-content {
  flex: 1;
  min-width: 0;
}

.activity-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 4px;
}

.activity-title {
  font-size: 14px;
  font-weight: 500;
  color: #111827;
  flex: 1;
  line-height: 1.4;
}

.activity-time {
  font-size: 11px;
  color: #9CA3AF;
  white-space: nowrap;
  flex-shrink: 0;
}

.activity-description {
  font-size: 13px;
  color: #6B7280;
  line-height: 1.4;
  margin-bottom: 8px;
}

.activity-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.meta-tag {
  font-size: 11px;
  padding: 2px 6px;
  background: #F3F4F6;
  color: #374151;
  border-radius: 3px;
  font-weight: 500;
}

/* Activity type specific styling */
.activity-task {
  border-left: 3px solid #3B82F6;
}

.activity-agent {
  border-left: 3px solid #10B981;
}

.activity-project {
  border-left: 3px solid #F59E0B;
}

.activity-system {
  border-left: 3px solid #8B5CF6;
}

/* Activity status specific styling */
.activity-completed {
  opacity: 0.8;
}

.activity-priority-high {
  border-left-color: #EF4444 !important;
}

.activity-priority-medium {
  border-left-color: #F59E0B !important;
}

.activity-priority-low {
  border-left-color: #10B981 !important;
}

.empty-activity {
  text-align: center;
  padding: 40px 20px;
  color: #6B7280;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-activity h4 {
  font-size: 16px;
  font-weight: 600;
  color: #374151;
  margin: 0 0 8px 0;
}

.empty-activity p {
  font-size: 14px;
  margin: 0 0 16px 0;
  line-height: 1.4;
}

.btn-clear-empty {
  padding: 8px 16px;
  background: #F3F4F6;
  color: #374151;
  border: 1px solid #D1D5DB;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-clear-empty:hover {
  background: #E5E7EB;
  border-color: #9CA3AF;
}

.loading-activity {
  padding: 20px;
}

.skeleton-activity {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid #F3F4F6;
}

.skeleton-activity:last-child {
  border-bottom: none;
}

.skeleton-icon {
  width: 24px;
  height: 24px;
  background: #E5E7EB;
  border-radius: 4px;
  flex-shrink: 0;
}

.skeleton-content {
  flex: 1;
}

.skeleton-title {
  height: 16px;
  background: #E5E7EB;
  border-radius: 4px;
  margin-bottom: 8px;
  width: 70%;
}

.skeleton-text {
  height: 14px;
  background: #E5E7EB;
  border-radius: 4px;
  width: 90%;
}

.load-more {
  padding: 16px 20px;
  text-align: center;
  border-top: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.btn-load-more {
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

.btn-load-more:hover {
  background: #2563EB;
}

.activity-stats {
  padding: 16px 20px;
  border-top: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 24px;
  font-weight: 700;
  color: #111827;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #6B7280;
  font-weight: 500;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }
  
  .filter-pills {
    padding: 12px 16px;
  }
  
  .activity-item {
    padding: 12px 16px;
  }
}
</style>