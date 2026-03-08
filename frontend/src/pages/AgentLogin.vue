<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
          Agent Access
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Enter your API key to access projects and tasks
        </p>
      </div>

      <div class="bg-blue-50 border border-blue-200 rounded-md p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-blue-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-blue-800">
              Agent Mode
            </h3>
            <div class="mt-2 text-sm text-blue-700">
              <p>Using API key authentication for read-only access to projects and tasks.</p>
            </div>
          </div>
        </div>
      </div>

      <form class="mt-8 space-y-6" @submit.prevent="handleApiKeyLogin">
        <div class="rounded-md shadow-sm space-y-4">
          <div>
            <label for="api-key" class="block text-sm font-medium text-gray-700">API Key</label>
            <input
              id="api-key"
              v-model="apiKey"
              name="api-key"
              type="password"
              autocomplete="off"
              required
              class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
              placeholder="Enter your agent API key"
            />
          </div>
        </div>

        <div v-if="error" class="rounded-md bg-red-50 p-4">
          <p class="text-sm text-red-800">{{ error }}</p>
        </div>

        <div>
          <button
            type="submit"
            :disabled="loading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ loading ? 'Verifying...' : 'Access as Agent' }}
          </button>
        </div>

        <div class="text-center">
          <router-link to="/login" class="font-medium text-blue-600 hover:text-blue-500 text-sm">
            Or login as a user
          </router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { projectService } from '../services/projectService'
import { setAgentMode } from '../utils/auth'

const router = useRouter()
const apiKey = ref('')
const loading = ref(false)
const error = ref('')

const handleApiKeyLogin = async () => {
  loading.value = true
  error.value = ''

  try {
    // Set agent mode
    setAgentMode(apiKey.value)

    // Test the API key by fetching projects
    await projectService.getProjectsForAgent()

    // If successful, redirect to dashboard
    router.push('/')
  } catch (err) {
    error.value = err.response?.data?.error || 'Invalid API key. Please check and try again.'
    // Clear the invalid API key
    localStorage.removeItem('agent_api_key')
  } finally {
    loading.value = false
  }
}
</script>
