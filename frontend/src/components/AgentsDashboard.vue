<template>
  <div class="agents-dashboard">
    <h3 class="agents-title">Agent Activity</h3>

    <div v-if="loading" class="loading">Loading agents...</div>

    <div v-else-if="agents.length === 0" class="empty-state">
      <p>No agents found</p>
    </div>

    <div v-else class="agents-list">
      <div v-for="agent in agents" :key="agent.id" class="agent-card">
        <div class="agent-header">
          <div class="agent-info">
            <h4 class="agent-name">{{ agent.name }}</h4>
            <span class="agent-status" :class="{ 'active': agent.enabled, 'inactive': !agent.enabled }">
              {{ agent.enabled ? 'Active' : 'Disabled' }}
            </span>
          </div>
        </div>

        <div class="agent-tasks">
          <div v-if="agent.active_tasks && agent.active_tasks.length > 0">
            <p class="tasks-label">Currently working on ({{ agent.active_tasks.length }}):</p>
            <div class="task-list">
              <div v-for="task in agent.active_tasks" :key="task.id" class="task-item">
                <span class="task-priority" :class="task.priority">{{ task.priority }}</span>
                <span class="task-title">{{ task.title }}</span>
              </div>
            </div>
          </div>
          <p v-else class="no-tasks">No active tasks</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAgentStore } from '../stores/agents'

const agentStore = useAgentStore()
const agents = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    await agentStore.fetchAgents()
    agents.value = agentStore.agents
  } catch (error) {
    console.error('Failed to load agents:', error)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.agents-dashboard {
  background: white;
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  padding: 24px;
  margin-bottom: 24px;
}

.agents-title {
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

.agents-list {
  display: grid;
  gap: 16px;
}

.agent-card {
  background: #F9FAFB;
  border-radius: 8px;
  padding: 16px;
  border: 1px solid #E5E7EB;
}

.agent-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.agent-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.agent-name {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  margin: 0;
}

.agent-status {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.agent-status.active {
  background: #DEF7EC;
  color: #03543F;
}

.agent-status.inactive {
  background: #FEF2F2;
  color: #9B2C2C;
}

.agent-tasks {
  margin-top: 8px;
}

.tasks-label {
  font-size: 14px;
  font-weight: 500;
  color: #6B7280;
  margin: 0 0 8px 0;
}

.task-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-item {
  display: flex;
  gap: 8px;
  align-items: center;
  padding: 8px;
  background: white;
  border-radius: 6px;
  border: 1px solid #E5E7EB;
}

.task-priority {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
}

.task-priority.high {
  background: #FEE2E2;
  color: #991B1B;
}

.task-priority.medium {
  background: #FEF3C7;
  color: #92400E;
}

.task-priority.low {
  background: #DBEAFE;
  color: #1E40AF;
}

.task-priority.critical {
  background: #FCA5A5;
  color: #7F1D1D;
}

.task-title {
  font-size: 14px;
  color: #374151;
  flex: 1;
}

.no-tasks {
  font-size: 14px;
  color: #9CA3AF;
  font-style: italic;
  margin: 0;
}
</style>
