<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation -->
    <nav class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <router-link
              to="/projects"
              class="text-indigo-600 hover:text-indigo-900 text-sm font-medium"
            >
              &larr; Back to Projects
            </router-link>
          </div>
          <div class="flex items-center">
            <span class="text-gray-700 mr-4">{{ authStore.user?.email }}</span>
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

    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8" v-if="project">
      <!-- Project Header -->
      <div class="bg-white shadow rounded-lg mb-6">
        <div class="px-4 py-5 sm:px-6">
          <div class="md:flex md:items-center md:justify-between">
            <div class="flex-1 min-w-0">
              <div class="flex items-center space-x-3">
                <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl">
                  {{ project.name }}
                </h2>
                <span :class="getStatusClass(project.status)" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                  {{ formatStatus(project.status) }}
                </span>
              </div>
              <p class="mt-1 text-sm text-gray-500">{{ project.description || 'No description' }}</p>
            </div>
            <div class="mt-4 flex md:mt-0 md:ml-4 space-x-2">
              <button
                v-if="!isEditing"
                @click="startEditing"
                class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                </svg>
                Edit
              </button>
              <button
                v-if="isEditing"
                @click="cancelEditing"
                class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                v-if="isEditing"
                @click="handleSave"
                :disabled="!hasChanges"
                class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                </svg>
                Save Changes
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Project Details -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Main Content -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Description -->
          <div class="bg-white shadow rounded-lg">
            <div class="px-4 py-5 sm:px-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900">Description</h3>
            </div>
            <div class="px-4 py-5 sm:px-6">
              <div v-if="!isEditing">
                <p class="text-gray-700 whitespace-pre-wrap">{{ project.description || 'No description provided' }}</p>
              </div>
              <textarea
                v-else
                v-model="editForm.description"
                rows="4"
                class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md border p-2"
                placeholder="Project description..."
              ></textarea>
            </div>
          </div>

          <!-- LLM Context -->
          <div class="bg-white shadow rounded-lg">
            <div class="px-4 py-5 sm:px-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900">AI Agent Context</h3>
              <p class="text-sm text-gray-500">Instructions and guidelines for AI agents working on this project</p>
            </div>
            <div class="px-4 py-5 sm:px-6">
              <div v-if="!isEditing">
                <p class="text-gray-700 whitespace-pre-wrap">{{ project.llm_context || 'No AI context provided' }}</p>
              </div>
              <textarea
                v-else
                v-model="editForm.llm_context"
                rows="6"
                class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md border p-2 font-mono text-sm"
                placeholder="Enter instructions for AI agents..."
              ></textarea>
            </div>
          </div>

          <!-- Project Tasks -->
          <div class="bg-white shadow rounded-lg">
            <div class="px-4 py-5 sm:px-6 flex justify-between items-center">
              <h3 class="text-lg leading-6 font-medium text-gray-900">Project Tasks</h3>
              <router-link
                :to="`/tasks?project=${project.id}`"
                class="text-indigo-600 hover:text-indigo-900 text-sm font-medium"
              >
                View all tasks →
              </router-link>
            </div>
            <div class="px-4 py-5 sm:px-6">
              <div v-if="projectTasks.length === 0" class="text-gray-500 text-center py-4">
                No tasks in this project yet
              </div>
              <div v-else class="space-y-3">
                <div
                  v-for="task in projectTasks"
                  :key="task.id"
                  class="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition-colors"
                >
                  <div class="flex items-start justify-between">
                    <div class="flex-1">
                      <router-link
                        :to="`/tasks/${task.id}`"
                        class="text-sm font-medium text-gray-900 hover:text-indigo-600"
                      >
                        {{ task.title }}
                      </router-link>
                      <div class="mt-1 flex items-center space-x-2">
                        <span :class="getTaskStatusClass(task.status)" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                          {{ formatTaskStatus(task.status) }}
                        </span>
                        <span class="text-xs text-gray-500">{{ task.priority }} priority</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Sidebar -->
        <div class="lg:col-span-1 space-y-6">
          <!-- Project Info -->
          <div class="bg-white shadow rounded-lg">
            <div class="px-4 py-5 sm:px-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900">Project Info</h3>
            </div>
            <div class="px-4 py-5 sm:px-6">
              <dl class="space-y-4">
                <div>
                  <dt class="text-sm font-medium text-gray-500">Status</dt>
                  <dd v-if="!isEditing" class="mt-1">
                    <span :class="getStatusClass(project.status)" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                      {{ formatStatus(project.status) }}
                    </span>
                  </dd>
                  <select
                    v-else
                    v-model="editForm.status"
                    class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md border"
                  >
                    <option value="active">Active</option>
                    <option value="completed">Completed</option>
                    <option value="archived">Archived</option>
                  </select>
                </div>
                <div>
                  <dt class="text-sm font-medium text-gray-500">Created</dt>
                  <dd class="mt-1 text-sm text-gray-900">{{ formatDate(project.created_at) }}</dd>
                </div>
                <div>
                  <dt class="text-sm font-medium text-gray-500">Last Updated</dt>
                  <dd class="mt-1 text-sm text-gray-900">{{ formatDate(project.updated_at) }}</dd>
                </div>
              </dl>
            </div>
          </div>

          <!-- Links -->
          <div class="bg-white shadow rounded-lg">
            <div class="px-4 py-5 sm:px-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900">Links</h3>
            </div>
            <div class="px-4 py-5 sm:px-6">
              <dl class="space-y-4">
                <div>
                  <dt class="text-sm font-medium text-gray-500">Repository URL</dt>
                  <dd v-if="!isEditing" class="mt-1 text-sm">
                    <a
                      v-if="project.repository_url"
                      :href="project.repository_url"
                      target="_blank"
                      class="text-indigo-600 hover:text-indigo-900"
                    >
                      {{ project.repository_url }}
                    </a>
                    <span v-else class="text-gray-400">Not set</span>
                  </dd>
                  <input
                    v-else
                    v-model="editForm.repository_url"
                    type="url"
                    class="mt-1 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md border p-2"
                    placeholder="https://github.com/..."
                  />
                </div>
                <div>
                  <dt class="text-sm font-medium text-gray-500">Deployed URL</dt>
                  <dd v-if="!isEditing" class="mt-1 text-sm">
                    <a
                      v-if="project.deployed_url"
                      :href="project.deployed_url"
                      target="_blank"
                      class="text-indigo-600 hover:text-indigo-900"
                    >
                      {{ project.deployed_url }}
                    </a>
                    <span v-else class="text-gray-400">Not set</span>
                  </dd>
                  <input
                    v-else
                    v-model="editForm.deployed_url"
                    type="url"
                    class="mt-1 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md border p-2"
                    placeholder="https://..."
                  />
                </div>
                <div>
                  <dt class="text-sm font-medium text-gray-500">Documentation URL</dt>
                  <dd v-if="!isEditing" class="mt-1 text-sm">
                    <a
                      v-if="project.documentation_url"
                      :href="project.documentation_url"
                      target="_blank"
                      class="text-indigo-600 hover:text-indigo-900"
                    >
                      {{ project.documentation_url }}
                    </a>
                    <span v-else class="text-gray-400">Not set</span>
                  </dd>
                  <input
                    v-else
                    v-model="editForm.documentation_url"
                    type="url"
                    class="mt-1 shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md border p-2"
                    placeholder="https://docs..."
                  />
                </div>
              </dl>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-else class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <div class="text-center py-12">
        <div class="text-gray-500">Loading project...</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { projectService } from '../services/projectService'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const project = ref(null)
const projectTasks = ref([])
const isEditing = ref(false)
const editForm = ref({
  name: '',
  description: '',
  status: '',
  repository_url: '',
  deployed_url: '',
  documentation_url: '',
  llm_context: ''
})

onMounted(async () => {
  await loadProject()
})

const loadProject = async () => {
  try {
    const data = await projectService.getProject(route.params.id)
    project.value = data
    projectTasks.value = data.tasks || []
  } catch (error) {
    console.error('Failed to load project:', error)
    alert('Failed to load project')
    router.push('/projects')
  }
}

const startEditing = () => {
  isEditing.value = true
  editForm.value = {
    name: project.value.name,
    description: project.value.description || '',
    status: project.value.status,
    repository_url: project.value.repository_url || '',
    deployed_url: project.value.deployed_url || '',
    documentation_url: project.value.documentation_url || '',
    llm_context: project.value.llm_context || ''
  }
}

const cancelEditing = () => {
  isEditing.value = false
  editForm.value = {
    name: '',
    description: '',
    status: '',
    repository_url: '',
    deployed_url: '',
    documentation_url: '',
    llm_context: ''
  }
}

const hasChanges = computed(() => {
  return (
    editForm.value.name !== project.value.name ||
    editForm.value.description !== (project.value.description || '') ||
    editForm.value.status !== project.value.status ||
    editForm.value.repository_url !== (project.value.repository_url || '') ||
    editForm.value.deployed_url !== (project.value.deployed_url || '') ||
    editForm.value.documentation_url !== (project.value.documentation_url || '') ||
    editForm.value.llm_context !== (project.value.llm_context || '')
  )
})

const handleSave = async () => {
  try {
    await projectService.updateProject(project.value.id, editForm.value)
    isEditing.value = false
    await loadProject()
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to update project')
  }
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const formatDate = (dateStr) => {
  if (!dateStr) return 'N/A'
  return new Date(dateStr).toLocaleString()
}

const formatStatus = (status) => {
  return status.charAt(0).toUpperCase() + status.slice(1)
}

const getStatusClass = (status) => {
  const classes = {
    active: 'bg-green-100 text-green-800',
    completed: 'bg-blue-100 text-blue-800',
    archived: 'bg-gray-100 text-gray-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const formatTaskStatus = (status) => {
  return status.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const getTaskStatusClass = (status) => {
  const classes = {
    pending: 'bg-yellow-100 text-yellow-800',
    in_progress: 'bg-blue-100 text-blue-800',
    completed: 'bg-green-100 text-green-800',
    failed: 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}
</script>

<style scoped>
/* Add any scoped styles here if needed */
</style>
