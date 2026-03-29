<template>
  <div class="agents-dashboard-compact">
    <div class="dashboard-header">
      <h3 class="text-lg font-semibold text-gray-900">Active Agents</h3>
      <button @click="showAddAgent = !showAddAgent" class="btn-add-agent">
        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Add Agent
      </button>
    </div>

    <!-- Add Agent Form -->
    <div v-if="showAddAgent" class="add-agent-form">
      <div class="form-group">
        <label class="form-label">Agent Name</label>
        <input 
          v-model="newAgent.name" 
          type="text" 
          class="form-input"
          placeholder="Enter agent name"
        />
      </div>
      <div class="form-group">
        <label class="form-label">Agent Type</label>
        <select v-model="newAgent.type" class="form-select">
          <option value="general">General</option>
          <option value="developer">Developer</option>
          <option value="designer">Designer</option>
          <option value="analyst">Analyst</option>
        </select>
      </div>
      <div class="form-actions">
        <button @click="addAgent" class="btn-primary" :disabled="!newAgent.name.trim()">
          Create Agent
        </button>
        <button @click="showAddAgent = false" class="btn-secondary">
          Cancel
        </button>
      </div>
    </div>

    <!-- Agent Stats -->
    <div class="agent-stats">
      <div class="stat-item">
        <div class="stat-icon">🤖</div>
        <div class="stat-content">
          <div class="stat-number">{{ totalAgents }}</div>
          <div class="stat-label">Total Agents</div>
        </div>
      </div>
      <div class="stat-item">
        <div class="stat-icon">⚡</div>
        <div class="stat-content">
          <div class="stat-number">{{ activeAgents }}</div>
          <div class="stat-label">Active</div>
        </div>
      </div>
      <div class="stat-item">
        <div class="stat-icon">💤</div>
        <div class="stat-content">
          <div class="stat-number">{{ idleAgents }}</div>
          <div class="stat-label">Idle</div>
        </div>
      </div>
    </div>

    <!-- Agent List -->
    <div class="agent-list">
      <div 
        v-for="agent in agents" 
        :key="agent.id"
        class="agent-item"
        :class="getAgentStatusClass(agent)"
      >
        <div class="agent-avatar">
          <span class="avatar-text">{{ agent.name.charAt(0).toUpperCase() }}</span>
          <div class="agent-status" :class="agent.status"></div>
        </div>
        <div class="agent-info">
          <div class="agent-name">{{ agent.name }}</div>
          <div class="agent-type">{{ agent.type }}</div>
          <div class="agent-meta">
            <span class="meta-item">
              📋 {{ agent.activeTasks || 0 }} tasks
            </span>
            <span v-if="agent.lastActive" class="meta-item">
              🕐 {{ formatTime(agent.lastActive) }}
            </span>
          </div>
        </div>
        <div class="agent-actions">
          <button 
            @click="toggleAgentStatus(agent)" 
            class="btn-toggle"
            :class="agent.status"
          >
            {{ agent.status === 'active' ? 'Pause' : 'Start' }}
          </button>
          <button @click="viewAgentDetails(agent)" class="btn-details">
            Details
          </button>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="agents.length === 0" class="empty-agents">
        <div class="empty-icon">🤖</div>
        <h4>No agents yet</h4>
        <p>Add your first AI agent to start automating tasks</p>
        <button @click="showAddAgent = true" class="btn-primary">
          Create First Agent
        </button>
      </div>
    </div>

    <!-- Performance Overview -->
    <div class="performance-overview" v-if="agents.length > 0">
      <h4 class="overview-title">Agent Performance</h4>
      <div class="performance-chart">
        <div class="chart-bar">
          <div class="bar-fill" :style="{ width: `${completionRate}%` }"></div>
          <div class="bar-label">{{ completionRate }}% Complete</div>
        </div>
        <div class="chart-stats">
          <span class="stat">Tasks Completed: {{ completedTasks }}</span>
          <span class="stat">Avg Response: {{ avgResponseTime }}s</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// Agent data management
const agents = ref([
  {
    id: 1,
    name: 'code-assistant',
    type: 'developer',
    status: 'active',
    activeTasks: 3,
    lastActive: new Date(Date.now() - 1000 * 60 * 5), // 5 minutes ago
    tasksCompleted: 24,
    avgResponseTime: 2.3
  },
  {
    id: 2,
    name: 'doc-writer',
    type: 'general',
    status: 'active',
    activeTasks: 1,
    lastActive: new Date(Date.now() - 1000 * 60 * 15), // 15 minutes ago
    tasksCompleted: 12,
    avgResponseTime: 3.7
  },
  {
    id: 3,
    name: 'code-review-bot',
    type: 'developer',
    status: 'idle',
    activeTasks: 0,
    lastActive: new Date(Date.now() - 1000 * 60 * 60), // 1 hour ago
    tasksCompleted: 18,
    avgResponseTime: 1.8
  },
  {
    id: 4,
    name: 'marketing-assistant',
    type: 'designer',
    status: 'active',
    activeTasks: 2,
    lastActive: new Date(Date.now() - 1000 * 60 * 8), // 8 minutes ago
    tasksCompleted: 8,
    avgResponseTime: 4.2
  }
])

const showAddAgent = ref(false)
const newAgent = ref({
  name: '',
  type: 'general'
})

// Computed properties
const totalAgents = computed(() => agents.value.length)
const activeAgents = computed(() => agents.value.filter(agent => agent.status === 'active').length)
const idleAgents = computed(() => agents.value.filter(agent => agent.status === 'idle').length)
const completedTasks = computed(() => agents.value.reduce((sum, agent) => sum + agent.tasksCompleted, 0))
const avgResponseTime = computed(() => {
  const total = agents.value.reduce((sum, agent) => sum + agent.avgResponseTime, 0)
  return total > 0 ? (total / agents.value.length).toFixed(1) : '0'
})
const completionRate = computed(() => {
  const totalActiveTasks = agents.value.reduce((sum, agent) => sum + agent.activeTasks, 0)
  const totalTasks = totalActiveTasks + completedTasks.value
  return totalTasks > 0 ? Math.round((completedTasks.value / totalTasks) * 100) : 0
})

// Methods
const addAgent = () => {
  if (!newAgent.value.name.trim()) return
  
  const newAgentObj = {
    id: agents.value.length + 1,
    name: newAgent.value.name.trim(),
    type: newAgent.value.type,
    status: 'active',
    activeTasks: 0,
    lastActive: new Date(),
    tasksCompleted: 0,
    avgResponseTime: 2.0
  }
  
  agents.value.push(newAgentObj)
  
  // Reset form
  newAgent.value = {
    name: '',
    type: 'general'
  }
  showAddAgent.value = false
}

const toggleAgentStatus = (agent) => {
  agent.status = agent.status === 'active' ? 'idle' : 'active'
  
  if (agent.status === 'idle') {
    agent.activeTasks = 0
  }
}

const viewAgentDetails = (agent) => {
  // In a real app, this would navigate to agent details page
  console.log('Viewing agent details:', agent)
}

const getAgentStatusClass = (agent) => {
  return {
    'agent-active': agent.status === 'active',
    'agent-idle': agent.status === 'idle'
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

// Simulate real-time updates
let updateInterval
onMounted(() => {
  updateInterval = setInterval(() => {
    // Simulate agent activity updates
    agents.value.forEach(agent => {
      if (agent.status === 'active' && Math.random() > 0.7) {
        agent.lastActive = new Date()
      }
    })
  }, 30000) // Update every 30 seconds
})

// Cleanup interval on unmount
onUnmounted(() => {
  if (updateInterval) {
    clearInterval(updateInterval)
  }
})
</script>

<style scoped>
.agents-dashboard-compact {
  background: white;
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  overflow: hidden;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.dashboard-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.btn-add-agent {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-add-agent:hover {
  background: #2563EB;
}

.add-agent-form {
  padding: 16px 20px;
  border-bottom: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.form-group {
  margin-bottom: 12px;
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
.form-select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #D1D5DB;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.2s ease;
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: #3B82F6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.btn-primary {
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

.btn-primary:hover:not(:disabled) {
  background: #2563EB;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  padding: 8px 16px;
  background: white;
  color: #374151;
  border: 1px solid #D1D5DB;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-secondary:hover {
  background: #F9FAFB;
  border-color: #9CA3AF;
}

.agent-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  padding: 16px 20px;
  border-bottom: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stat-icon {
  font-size: 20px;
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 18px;
  font-weight: 700;
  color: #111827;
  line-height: 1.2;
}

.stat-label {
  font-size: 11px;
  color: #6B7280;
  font-weight: 500;
}

.agent-list {
  max-height: 300px;
  overflow-y: auto;
}

.agent-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid #F3F4F6;
  transition: all 0.2s ease;
}

.agent-item:hover {
  background: #F9FAFB;
}

.agent-item:last-child {
  border-bottom: none;
}

.agent-avatar {
  position: relative;
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: #E5E7EB;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.avatar-text {
  font-size: 16px;
  font-weight: 600;
  color: #374151;
}

.agent-status {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  border: 2px solid white;
}

.agent-status.active {
  background: #10B981;
}

.agent-status.idle {
  background: #F59E0B;
}

.agent-info {
  flex: 1;
  min-width: 0;
}

.agent-name {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 2px;
}

.agent-type {
  font-size: 12px;
  color: #6B7280;
  margin-bottom: 4px;
}

.agent-meta {
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

.agent-actions {
  display: flex;
  gap: 6px;
}

.btn-toggle {
  padding: 4px 8px;
  border: none;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-toggle.active {
  background: #10B981;
  color: white;
}

.btn-toggle.idle {
  background: #F59E0B;
  color: white;
}

.btn-toggle:hover {
  opacity: 0.8;
}

.btn-details {
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

.btn-details:hover {
  background: #E5E7EB;
  border-color: #9CA3AF;
}

.empty-agents {
  text-align: center;
  padding: 40px 20px;
  color: #6B7280;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-agents h4 {
  font-size: 16px;
  font-weight: 600;
  color: #374151;
  margin: 0 0 8px 0;
}

.empty-agents p {
  font-size: 14px;
  margin: 0 0 16px 0;
  line-height: 1.4;
}

.performance-overview {
  padding: 16px 20px;
  border-top: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.overview-title {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 12px 0;
}

.performance-chart {
  margin-bottom: 12px;
}

.chart-bar {
  width: 100%;
  height: 8px;
  background: #E5E7EB;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 8px;
}

.bar-fill {
  height: 100%;
  background: #10B981;
  transition: width 0.5s ease;
}

.bar-label {
  font-size: 12px;
  color: #6B7280;
  text-align: right;
}

.chart-stats {
  display: flex;
  gap: 16px;
  font-size: 11px;
  color: #6B7280;
}

.stat {
  display: flex;
  align-items: center;
  gap: 4px;
}

@media (max-width: 768px) {
  .dashboard-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .agent-stats {
    grid-template-columns: 1fr;
    gap: 8px;
  }
  
  .agent-meta {
    flex-direction: column;
    gap: 2px;
  }
  
  .chart-stats {
    flex-direction: column;
    gap: 4px;
  }
}
</style>