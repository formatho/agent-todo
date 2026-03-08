<template>
  <div class="activity-feed">
    <h3 class="feed-title">Activity Feed</h3>

    <div v-if="loading" class="loading">Loading activity...</div>

    <div v-else-if="activities.length === 0" class="empty-state">
      <p>No recent activity</p>
    </div>

    <div v-else class="activity-list">
      <div v-for="activity in activities" :key="activity.id" class="activity-item">
        <div class="activity-icon">
          {{ getActivityIcon(activity.event_type) }}
        </div>

        <div class="activity-content">
          <p class="activity-description">{{ activity.description }}</p>
          <div class="activity-meta">
            <span class="activity-time">{{ formatTime(activity.created_at) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useActivityStore } from '../stores/activity'

const activityStore = useActivityStore()
const activities = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    await activityStore.fetchActivity()
    activities.value = activityStore.activities
  } catch (error) {
    console.error('Failed to load activity:', error)
  } finally {
    loading.value = false
  }
})

const getActivityIcon = (eventType) => {
  const icons = {
    'created': '✨',
    'status_changed': '🔄',
    'assigned': '👤',
    'unassigned': '👤',
    'comment_added': '💬'
  }
  return icons[eventType] || '📋'
}

const formatTime = (timestamp) => {
  const date = new Date(timestamp)
  const now = new Date()
  const diff = now - date
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (minutes < 1) return 'Just now'
  if (minutes < 60) return `${minutes}m ago`
  if (hours < 24) return `${hours}h ago`
  if (days < 7) return `${days}d ago`

  return date.toLocaleDateString()
}
</script>

<style scoped>
.activity-feed {
  background: white;
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  padding: 24px;
  margin-bottom: 24px;
}

.feed-title {
  font-size: 20px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 16px 0;
}

.loading,
.empty-state {
  text-align: center;
  padding: 20px;
  color: #6B7280;
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.activity-item {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: #F9FAFB;
  border-radius: 8px;
  transition: background 0.2s;
}

.activity-item:hover {
  background: #F3F4F6;
}

.activity-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.activity-content {
  flex: 1;
}

.activity-description {
  margin: 0 0 4px 0;
  font-size: 14px;
  color: #374151;
  line-height: 1.5;
}

.activity-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #9CA3AF;
}

.activity-time {
  font-weight: 500;
}
</style>
