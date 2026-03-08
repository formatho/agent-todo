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
                class="border-indigo-500 text-gray-900 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
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
                class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
              >
                Dashboard
              </router-link>
            </div>
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

    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="md:flex md:items-center md:justify-between">
        <div class="flex-1 min-w-0">
          <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
            Agents
          </h2>
        </div>
        <div class="mt-4 flex md:mt-0 md:ml-4">
          <button
            @click="showCreateModal = true"
            type="button"
            class="ml-3 inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            Create Agent
          </button>
        </div>
      </div>

      <!-- Agents Grid -->
      <div class="mt-6 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        <div
          v-if="agentStore.loading"
          class="col-span-full text-center py-12 text-gray-500"
        >
          Loading agents...
        </div>

        <div
          v-else-if="agentStore.agents.length === 0"
          class="col-span-full text-center py-12 text-gray-500"
        >
          No agents found. Create your first agent to get started.
        </div>

        <AgentCard
          v-else
          v-for="agent in agentStore.agents"
          :key="agent.id"
          :agent="agent"
          :tasks="getAgentTasks(agent.id)"
          @edit="handleEdit"
          @delete="handleDelete"
        />
      </div>
    </div>

    <!-- Create Agent Modal -->
    <div v-if="showCreateModal" class="fixed z-10 inset-0 overflow-y-auto">
      <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="showCreateModal = false"></div>

        <span class="hidden sm:inline-block sm:align-middle sm:h-screen">&#8203;</span>

        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form @submit.prevent="handleCreateAgent">
            <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">
                Create New Agent
              </h3>

              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700">Name</label>
                  <input
                    v-model="agentForm.name"
                    type="text"
                    required
                    class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border p-2"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700">Description</label>
                  <textarea
                    v-model="agentForm.description"
                    rows="3"
                    class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border p-2"
                  ></textarea>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700">Role</label>
                  <select
                    v-model="agentForm.role"
                    class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border p-2"
                  >
                    <option value="regular">Regular - Can only update own tasks</option>
                    <option value="supervisor">Supervisor - Can update any task & create agents</option>
                    <option value="admin">Admin - Full system access</option>
                  </select>
                  <p class="mt-1 text-xs text-gray-500">
                    Regular agents can only update tasks assigned to them. Supervisors can update any task and create new agents.
                  </p>
                </div>
              </div>
            </div>

            <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
              <button
                type="submit"
                :disabled="loading"
                class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
              >
                {{ loading ? 'Creating...' : 'Create' }}
              </button>
              <button
                type="button"
                @click="showCreateModal = false"
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Edit Agent Modal -->
    <div v-if="showEditModal && editingAgent" class="fixed z-10 inset-0 overflow-y-auto">
      <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="showEditModal = false"></div>

        <span class="hidden sm:inline-block sm:align-middle sm:h-screen">&#8203;</span>

        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form @submit.prevent="handleUpdateAgent">
            <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">
                Edit Agent
              </h3>

              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700">Name</label>
                  <input
                    v-model="editingAgent.name"
                    type="text"
                    required
                    class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border p-2"
                  />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700">Description</label>
                  <textarea
                    v-model="editingAgent.description"
                    rows="3"
                    class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border p-2"
                  ></textarea>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700">Role</label>
                  <select
                    v-model="editingAgent.role"
                    class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border p-2"
                  >
                    <option value="regular">Regular - Can only update own tasks</option>
                    <option value="supervisor">Supervisor - Can update any task & create agents</option>
                    <option value="admin">Admin - Full system access</option>
                  </select>
                </div>

                <div>
                  <label class="flex items-center">
                    <input
                      type="checkbox"
                      v-model="editingAgent.enabled"
                      class="rounded border-gray-300 text-indigo-600 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                    />
                    <span class="ml-2 text-sm text-gray-700">Agent Enabled</span>
                  </label>
                  <p class="mt-1 text-xs text-gray-500">
                    Disabled agents cannot access the API.
                  </p>
                </div>
              </div>
            </div>

            <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
              <button
                type="submit"
                :disabled="loading"
                class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
              >
                {{ loading ? 'Updating...' : 'Update' }}
              </button>
              <button
                type="button"
                @click="showEditModal = false"
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useAgentStore } from '../stores/agents'
import { useTaskStore } from '../stores/tasks'
import AgentCard from '../components/AgentCard.vue'

const router = useRouter()
const authStore = useAuthStore()
const agentStore = useAgentStore()
const taskStore = useTaskStore()

const showCreateModal = ref(false)
const showEditModal = ref(false)
const loading = ref(false)
const editingAgent = ref(null)

const agentForm = ref({
  name: '',
  description: '',
  role: 'regular'
})

onMounted(async () => {
  await agentStore.fetchAgents()
  await taskStore.fetchTasks()
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const handleCreateAgent = async () => {
  loading.value = true

  try {
    await agentStore.createAgent(agentForm.value)
    showCreateModal.value = false
    agentForm.value = { name: '', description: '', role: 'regular' }
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to create agent')
  } finally {
    loading.value = false
  }
}

const handleUpdateAgent = async () => {
  loading.value = true

  try {
    await agentStore.updateAgent(editingAgent.value.id, {
      name: editingAgent.value.name,
      description: editingAgent.value.description,
      role: editingAgent.value.role,
      enabled: editingAgent.value.enabled
    })
    showEditModal.value = false
    editingAgent.value = null
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to update agent')
  } finally {
    loading.value = false
  }
}

const handleEdit = (agent) => {
  editingAgent.value = { ...agent }
  showEditModal.value = true
}

const handleDelete = async (id) => {
  try {
    await agentStore.deleteAgent(id)
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to delete agent')
  }
}

const getAgentTasks = (agentId) => {
  return taskStore.tasks.filter(task => task.assigned_agent_id === agentId)
}
</script>
