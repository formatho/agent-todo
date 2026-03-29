<template>
  <div class="compact-activity-feed">
    <div class="feed-header">
      <h3 class="text-md font-semibold text-gray-900">Recent Activity</h3>
      <button @click="refreshFeed" class="btn-refresh" :disabled="loading">
        <svg v-if="loading" class="animate-spin h-3 w-3" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <svg v-else class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </button>
    </div>

    <div class="activity-list">
      <div 
        v-for="activity in recentActivities" 
        :key="activity.id"
        class="activity-item"
      >
        <div class="activity-icon" :class="activity.type">
          <component :is="getActivityIcon(activity.type)" />
        </div>
        <div class="activity-content">
          <div class="activity-title">{{ activity.title }}</div>
          <div class="activity-description">{{ activity.description }}</div>
          <div class="activity-time">{{ formatTime(activity.timestamp) }}</div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="recentActivities.length === 0" class="empty-activity">
        <div class="empty-icon">📊</div>
        <p>No recent activity</p>
      </div>

      <!-- Load More -->
      <div v-if="recentActivities.length > 0" class="load-more">
        <button @click="loadMore" class="btn-load-more">
          View All Activity
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const loading = ref(false)
const recentActivities = ref([
  {
    id: 1,
    type: 'task',
    title: 'Task Completed: Database Schema',
    description: 'code-assistant finished the database setup',
    timestamp: new Date(Date.now() - 1000 * 60 * 5)
  },
  {
    id: 2,
    type: 'agent',
    title: 'Agent Started: code-review-bot',
    description: 'Code review agent is now active',
    timestamp: new Date(Date.now() - 1000 * 60 * 15)
  },
  {
    id: 3,
    type: 'project',
    title: 'Project Created: Marketing Campaign',
    description: 'New project for Q2 initiatives',
    timestamp: new Date(Date.now() - 1000 * 60 * 30)
  },
  {
    id: 4,
    type: 'task',
    title: 'Task Started: API Documentation',
    description: 'doc-writer began working on documentation',
    timestamp: new Date(Date.now() - 1000 * 60 * 45)
  },
  {
    id: 5,
    type: 'system',
    title: 'System Update: Security Patch',
    description: 'Security patches applied successfully',
    timestamp: new Date(Date.now() - 1000 * 60 * 60)
  }
])

const getActivityIcon = (type) => {
  const icons = {
    task: '📋',
    agent: '🤖',
    project: '📁',
    system: '⚙️'
  }
  return icons[type] || '📝'
}

const formatTime = (timestamp) => {
  const now = new Date()
  const diff = now - new Date(timestamp)
  
  if (diff < 60000) return 'Just now'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}m ago`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}h ago`
  return `${Math.floor(diff / 86400000)}d ago`
}

const refreshFeed = async () => {
  loading.value = true
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 1000))
  loading.value = false
}

const loadMore = () => {
  // In a real app, this would load more activities
  console.log('Loading more activities...')
}

// Auto-refresh every 2 minutes
let refreshInterval
onMounted(() => {
  refreshInterval = setInterval(refreshFeed, 120000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.compact-activity-feed {
  background: white;
  border-radius: 8px;
  border: 1px solid #E5E7EB;
  overflow: hidden;
}

.feed-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.feed-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #111827;
}

.btn-refresh {
  padding: 4px 8px;
  background: transparent;
  border: 1px solid #D1D5DB;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 4px;
}

.btn-refresh:hover {
  background: #F3F4F6;
}

.btn-refresh:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.activity-list {
  max-height: 400px;
  overflow-y: auto;
}

.activity-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 16px;
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
  font-size: 16px;
  flex-shrink: 0;
  width: 20px;
  text-align: center;
}

.activity-icon.task {
  color: #3B82F6;
}

.activity-icon.agent {
  color: #10B981;
}

.activity-icon.project {
  color: #F59E0B;
}

.activity-icon.system {
  color: #8B5CF6;
}

.activity-content {
  flex: 1;
  min-width: 0;
}

.activity-title {
  font-size: 12px;
  font-weight: 500;
  color: #111827;
  margin-bottom: 2px;
  line-height: 1.3;
}

.activity-description {
  font-size: 11px;
  color: #6B7280;
  margin-bottom: 4px;
  line-height: 1.3;
}

.activity-time {
  font-size: 10px;
  color: #9CA3AF;
}

.empty-activity {
  text-align: center;
  padding: 40px 20px;
  color: #6B7280;
}

.empty-icon {
  font-size: 32px;
  margin-bottom: 8px;
  opacity: 0.5;
}

.empty-activity p {
  font-size: 12px;
  margin: 0;
}

.load-more {
  padding: 12px 16px;
  text-align: center;
  border-top: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.btn-load-more {
  padding: 6px 12px;
  background: #F3F4F6;
  color: #374151;
  border: 1px solid #D1D5DB;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-load-more:hover {
  background: #E5E7EB;
  border-color: #9CA3AF;
}
</style>