<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation -->
    <nav class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex">
            <div class="flex-shrink-0 flex items-center">
              <h1 class="text-xl font-bold text-gray-900">Agent Todo Platform</h1>
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
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Activity Feed -->
        <ActivityFeed />

        <!-- Agent Activity -->
        <AgentsDashboard />
      </div>

      <!-- Task Grid -->
      <TaskGrid />
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { isAgentMode, getAuthInfo } from '../utils/auth'
import TaskGrid from '../components/TaskGrid.vue'
import ActivityFeed from '../components/ActivityFeed.vue'
import AgentsDashboard from '../components/AgentsDashboard.vue'
import ThemeToggle from '../components/ThemeToggle.vue'

const router = useRouter()
const authStore = useAuthStore()
const agentMode = isAgentMode()
const authInfo = getAuthInfo()

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

