<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation -->
    <nav class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex">
            <div class="flex-shrink-0 flex items-center">
              <h1 class="text-xl font-bold text-gray-900">Agent Todo Platform</h1>
              <span class="ml-3 px-2 py-1 text-xs font-medium bg-green-100 text-green-800 rounded-full">
                V2 - Compact Grid
              </span>
            </div>
            <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
              <router-link
                to="/tasks"
                class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
              >
                Tasks
              </router-link>
              <router-link
                to="/agents"
                class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
              >
                Agents
              </router-link>
              <router-link
                to="/projects"
                class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
              >
                Projects
              </router-link>
              <router-link
                to="/"
                class="border-indigo-500 text-gray-900 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
              >
                Dashboard
              </router-link>
            </div>
          </div>
          <div class="flex items-center">
            <span v-if="agentMode" class="text-blue-700 mr-4 font-medium flex items-center gap-1">
              <svg class="h-5 w-5 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
              {{ authInfo?.identifier }}
            </span>
            <span v-else class="text-gray-700 mr-4">{{ authStore.user?.email }}</span>
            <ThemeToggle class="mr-2" />
            <button
              @click="handleLogout"
              class="bg-white py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 hover:bg-gray-50"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <!-- Welcome Banner for New Users -->
      <div v-if="showWelcomeBanner" class="welcome-banner">
        <div class="banner-content">
          <div class="banner-icon">👋</div>
          <div class="banner-text">
            <h2>Welcome to Agent Todo!</h2>
            <p>Get started quickly with our guided onboarding experience.</p>
          </div>
          <button @click="startOnboarding" class="btn-start-onboarding">
            Start Tutorial
          </button>
          <button @click="dismissWelcomeBanner" class="btn-dismiss">
            ×
          </button>
        </div>
      </div>

      <!-- Quick Stats Bar -->
      <div class="quick-stats-bar">
        <div class="stat-item">
          <div class="stat-icon">📊</div>
          <div class="stat-content">
            <div class="stat-number">{{ totalTasks }}</div>
            <div class="stat-label">Total Tasks</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon">⚡</div>
          <div class="stat-content">
            <div class="stat-number">{{ activeTasks }}</div>
            <div class="stat-label">Active</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon">🤖</div>
          <div class="stat-content">
            <div class="stat-number">{{ activeAgents }}</div>
            <div class="stat-label">Active Agents</div>
          </div>
        </div>
        <div class="stat-item">
          <div class="stat-icon">📈</div>
          <div class="stat-content">
            <div class="stat-number">{{ completionRate }}%</div>
            <div class="stat-label">Completion Rate</div>
          </div>
        </div>
      </div>

      <!-- Compact Task Grid Focus -->
      <div class="compact-task-section">
        <div class="section-header">
          <div class="header-left">
            <h2 class="text-xl font-bold text-gray-900">Task Management</h2>
            <div class="view-toggle">
              <button 
                @click="setTaskView('grid')" 
                :class="['view-btn', taskView === 'grid' ? 'active' : '']"
              >
                Grid
              </button>
              <button 
                @click="setTaskView('list')" 
                :class="['view-btn', taskView === 'list' ? 'active' : '']"
              >
                List
              </button>
            </div>
          </div>
          <div class="header-right">
            <button @click="showAddTask = true" class="btn-add-task">
              <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              Add Task
            </button>
          </div>
        </div>

        <!-- Quick Filters -->
        <div class="quick-filters">
          <button 
            v-for="filter in quickFilters" 
            :key="filter"
            @click="setQuickFilter(filter)"
            :class="['filter-btn', quickFilter === filter ? 'active' : '']"
          >
            {{ filter }}
          </button>
        </div>

        <!-- Compact Task Grid -->
        <CompactTaskGrid 
          :view="taskView"
          :filter="quickFilter"
          @add-task="showAddTask = true"
        />
      </div>

      <!-- Side Panel with Activity -->
      <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <!-- Compact Activity Feed -->
        <div class="lg:col-span-1">
          <CompactActivityFeed />
        </div>
        
        <!-- Main Content Area -->
        <div class="lg:col-span-3">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Project Overview -->
            <div class="project-overview">
              <h3 class="text-lg font-semibold text-gray-900 mb-4">Projects Overview</h3>
              <CompactProjectList />
            </div>

            <!-- Agent Status -->
            <div class="agent-status">
              <h3 class="text-lg font-semibold text-gray-900 mb-4">Agent Status</h3>
              <CompactAgentList />
            </div>
          </div>

          <!-- Recent Activity -->
          <div class="recent-activity">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Recent Activity</h3>
            <CompactRecentActivity />
          </div>
        </div>
      </div>
      
      <!-- Empty State with Enhanced Guidance -->
      <div v-else class="enhanced-empty-state">
        <div class="empty-content">
          <div class="empty-icon">🚀</div>
          <h2>Ready to Get Organized?</h2>
          <p>Your Agent Todo workspace is ready. Let's set up your first project to start managing tasks efficiently.</p>
          
          <div class="empty-steps">
            <div class="step">
              <div class="step-number">1</div>
              <div class="step-content">
                <h4>Create a Project</h4>
                <p>Organize your tasks into logical projects</p>
                <button @click="navigateToProjects" class="btn-step-action">
                  Create Project
                </button>
              </div>
            </div>
            
            <div class="step">
              <div class="step-number">2</div>
              <div class="step-content">
                <h4>Add Tasks</h4>
                <p>Create tasks and assign them to AI agents</p>
                <button @click="navigateToTasks" class="btn-step-action">
                  Create Tasks
                </button>
              </div>
            </div>
            
            <div class="step">
              <div class="step-number">3</div>
              <div class="step-content">
                <h4>Track Progress</h4>
                <p>Monitor task completion and analytics</p>
                <router-link to="/tasks" class="btn-step-action">
                  View Dashboard
                </router-link>
              </div>
            </div>
          </div>
          
          <div class="empty-actions">
            <button @click="startOnboarding" class="btn-primary-onboarding">
              Take Guided Tour
            </button>
            <button @click="dismissWelcomeBanner" class="btn-secondary">
              Skip for Now
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Task Modal -->
    <TaskQuickAdd 
      v-if="showAddTask" 
      @close="showAddTask = false"
      @task-created="handleTaskCreated"
    />

    <!-- Onboarding Flow Component -->
    <OnboardingFlow 
      v-if="showOnboarding" 
      @close="closeOnboarding" 
      @completed="onboardingCompleted"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useProjectStore } from '../stores/projects'
import { useTaskStore } from '../stores/tasks'
import { isAgentMode, getAuthInfo } from '../utils/auth'
import CompactTaskGrid from './CompactTaskGrid.vue'
import CompactActivityFeed from './CompactActivityFeed.vue'
import CompactProjectList from './CompactProjectList.vue'
import CompactAgentList from './CompactAgentList.vue'
import CompactRecentActivity from './CompactRecentActivity.vue'
import TaskQuickAdd from './TaskQuickAdd.vue'
import ThemeToggle from './ThemeToggle.vue'
import OnboardingFlow from './OnboardingFlow.vue'

const router = useRouter()
const authStore = useAuthStore()
const projectStore = useProjectStore()
const taskStore = useTaskStore()

const agentMode = isAgentMode()
const authInfo = getAuthInfo()
const showWelcomeBanner = ref(true)
const showQuickActions = ref(true)
const showOnboarding = ref(false)
const showAddTask = ref(false)
const hasProjects = computed(() => Array.isArray(projectStore.projects) && projectStore.projects.length > 0)

// Task view state
const taskView = ref('grid')
const quickFilter = ref('all')

// Quick filters
const quickFilters = ['all', 'active', 'completed', 'my-tasks']

// Stats computed properties
const totalTasks = computed(() => taskStore.tasks?.length || 0)
const activeTasks = computed(() => taskStore.tasks?.filter(task => task.status !== 'completed').length || 0)
const activeAgents = computed(() => projectStore.agents?.filter(agent => agent.status === 'active').length || 0)
const completionRate = computed(() => {
  const total = totalTasks.value
  const completed = totalTasks.value - activeTasks.value
  return total > 0 ? Math.round((completed / total) * 100) : 0
})

const setTaskView = (view) => {
  taskView.value = view
}

const setQuickFilter = (filter) => {
  quickFilter.value = filter
}

const startOnboarding = () => {
  showOnboarding.value = true
  showWelcomeBanner.value = false
  showQuickActions.value = false
}

const closeOnboarding = () => {
  showOnboarding.value = false
}

const onboardingCompleted = () => {
  showOnboarding.value = false
  // Refresh data to show newly created projects and tasks
  projectStore.fetchProjects()
  taskStore.fetchTasks()
}

const dismissWelcomeBanner = () => {
  showWelcomeBanner.value = false
  showQuickActions.value = false
}

const navigateToProjects = () => {
  router.push('/projects')
}

const navigateToTasks = () => {
  router.push('/tasks')
}

const navigateToAgents = () => {
  router.push('/agents')
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const handleTaskCreated = () => {
  showAddTask.value = false
  // Refresh tasks
  taskStore.fetchTasks()
}

onMounted(async () => {
  // Check if this is a new user (no projects and recent login)
  await projectStore.fetchProjects()
  await taskStore.fetchTasks()
  
  // Hide welcome banner if user has projects or after a certain time
  const lastVisit = localStorage.getItem('lastVisit')
  if (hasProjects.value || lastVisit) {
    showWelcomeBanner.value = false
    showQuickActions.value = false
  }
  
  // Store last visit time
  localStorage.setItem('lastVisit', new Date().toISOString())
})
</script>

<style scoped>
.welcome-banner {
  background: linear-gradient(135deg, #EFF6FF 0%, #DBEAFE 100%);
  border: 1px solid #BFDBFE;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
}

.banner-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.banner-icon {
  font-size: 48px;
  margin-right: 24px;
}

.banner-text {
  flex: 1;
}

.banner-text h2 {
  font-size: 20px;
  font-weight: 700;
  color: #1E40AF;
  margin: 0 0 4px 0;
}

.banner-text p {
  color: #1E40AF;
  margin: 0;
  font-size: 14px;
  opacity: 0.8;
}

.btn-start-onboarding {
  padding: 12px 24px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-start-onboarding:hover {
  background: #2563EB;
}

.btn-dismiss {
  background: none;
  border: none;
  font-size: 24px;
  color: #6B7280;
  cursor: pointer;
  padding: 0;
  margin-left: 16px;
}

.btn-dismiss:hover {
  color: #374151;
}

.quick-stats-bar {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: white;
  border-radius: 8px;
  border: 1px solid #E5E7EB;
}

.stat-icon {
  font-size: 24px;
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 20px;
  font-weight: 700;
  color: #111827;
  line-height: 1.2;
}

.stat-label {
  font-size: 12px;
  color: #6B7280;
  font-weight: 500;
}

.compact-task-section {
  background: white;
  border-radius: 12px;
  border: 1px solid #E5E7EB;
  margin-bottom: 24px;
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.section-header h2 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.view-toggle {
  display: flex;
  gap: 2px;
  padding: 2px;
  background: #E5E7EB;
  border-radius: 6px;
}

.view-btn {
  padding: 6px 12px;
  border: none;
  background: transparent;
  border-radius: 4px;
  font-size: 12px;
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
  background: #10B981;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-add-task:hover {
  background: #059669;
}

.quick-filters {
  padding: 12px 20px;
  border-bottom: 1px solid #E5E7EB;
  background: #F9FAFB;
  display: flex;
  gap: 8px;
}

.filter-btn {
  padding: 6px 12px;
  border: 1px solid #D1D5DB;
  background: white;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.filter-btn:hover {
  background: #F3F4F6;
  border-color: #9CA3AF;
}

.filter-btn.active {
  background: #10B981;
  color: white;
  border-color: #10B981;
}

.enhanced-empty-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 16px;
  border: 1px solid #E5E7EB;
  margin-top: 24px;
}

.empty-content {
  max-width: 600px;
  margin: 0 auto;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 24px;
}

.enhanced-empty-state h2 {
  font-size: 28px;
  font-weight: 700;
  color: #111827;
  margin: 0 0 16px 0;
}

.enhanced-empty-state p {
  font-size: 16px;
  color: #6B7280;
  margin: 0 0 40px 0;
  line-height: 1.6;
}

.empty-steps {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 24px;
  margin-bottom: 40px;
}

.step {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  text-align: left;
}

.step-number {
  width: 32px;
  height: 32px;
  background: #10B981;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
  flex-shrink: 0;
}

.step-content h4 {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 4px 0;
}

.step-content p {
  font-size: 13px;
  color: #6B7280;
  margin: 0 0 12px 0;
}

.btn-step-action {
  padding: 8px 16px;
  background: #F3F4F6;
  color: #374151;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-step-action:hover {
  background: #E5E7EB;
  color: #111827;
}

.empty-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
}

.btn-primary-onboarding {
  padding: 14px 28px;
  background: #10B981;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-primary-onboarding:hover {
  background: #059669;
}

.btn-secondary {
  padding: 14px 28px;
  background: white;
  color: #374151;
  border: 1px solid #D1D5DB;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-secondary:hover {
  background: #F9FAFB;
  border-color: #9CA3AF;
}

@media (max-width: 768px) {
  .banner-content {
    flex-direction: column;
    text-align: center;
    gap: 16px;
  }
  
  .banner-icon {
    margin-right: 0;
    margin-bottom: 0;
  }
  
  .quick-stats-bar {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }
  
  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  
  .view-toggle {
    align-self: flex-start;
  }
  
  .empty-steps {
    grid-template-columns: 1fr;
  }
  
  .empty-actions {
    flex-direction: column;
  }
}
</style>