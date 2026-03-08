<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation -->
    <nav class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <router-link
              to="/"
              class="text-indigo-600 hover:text-indigo-900 text-sm font-medium"
            >
              &larr; Back to Tasks
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

    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8" v-if="task">
      <!-- Task Header -->
      <div class="bg-white shadow rounded-lg mb-6">
        <div class="px-4 py-5 sm:px-6">
          <div class="md:flex md:items-center md:justify-between">
            <div class="flex-1 min-w-0">
              <h2 class="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
                {{ task.title }}
              </h2>
              <div class="mt-2 flex items-center space-x-4">
                <span :class="getStatusClass(task.status)" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                  {{ formatStatus(task.status) }}
                </span>
                <span :class="getPriorityClass(task.priority)" class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full">
                  {{ task.priority }}
                </span>
              </div>
            </div>
            <div class="mt-4 flex md:mt-0 md:ml-4 space-x-2">
              <button
                @click="showEditModal = true"
                class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
              >
                Edit
              </button>
              <button
                @click="handleDeleteTask"
                class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700"
              >
                Delete
              </button>
              <button
                v-if="!task.assigned_agent_id"
                @click="showAssignModal = true"
                class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700"
              >
                Assign Agent
              </button>
              <button
                v-else
                @click="handleUnassign"
                class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
              >
                Unassign
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Task Details -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Main Content -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Description -->
          <div class="bg-white shadow rounded-lg">
            <div class="px-4 py-5 sm:px-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900">Description</h3>
            </div>
            <div class="px-4 py-5 sm:px-6">
              <p class="text-gray-700">{{ task.description || 'No description provided' }}</p>
            </div>
          </div>

          <!-- Comments -->
          <div class="bg-white shadow rounded-lg">
            <div class="px-4 py-5 sm:px-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900">Comments</h3>
            </div>
            <div class="px-4 py-5 sm:px-6 space-y-4">
              <div v-if="comments.length === 0" class="text-gray-500 text-center py-4">
                No comments yet
              </div>

              <div
                v-for="comment in comments"
                :key="comment.id"
                class="border-b border-gray-200 pb-4 last:border-0"
              >
                <div class="flex items-center justify-between mb-2">
                  <span class="font-medium text-gray-900">{{ comment.author_name }}</span>
                  <span class="text-sm text-gray-500">{{ formatDate(comment.created_at) }}</span>
                </div>
                <p class="text-gray-700">{{ comment.content }}</p>
              </div>

              <!-- Add Comment -->
              <div class="mt-4">
                <textarea
                  v-model="newComment"
                  rows="3"
                  placeholder="Add a comment..."
                  class="shadow-sm focus:ring-indigo-500 focus:border-indigo-500 block w-full sm:text-sm border-gray-300 rounded-md border p-2"
                ></textarea>
                <div class="mt-2 flex justify-end">
                  <button
                    @click="handleAddComment"
                    :disabled="!newComment.trim()"
                    class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50"
                  >
                    Add Comment
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Event History -->
          <div class="bg-white shadow rounded-lg">
            <div class="px-4 py-5 sm:px-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900">History</h3>
            </div>
            <div class="px-4 py-5 sm:px-6">
              <div v-if="!task.events || task.events.length === 0" class="text-gray-500 text-center py-4">
                No history available
              </div>

              <div v-else class="space-y-3">
                <div
                  v-for="event in task.events"
                  :key="event.id"
                  class="flex items-start space-x-3 text-sm"
                >
                  <div class="flex-shrink-0">
                    <div class="h-2 w-2 rounded-full bg-indigo-600 mt-2"></div>
                  </div>
                  <div class="flex-1">
                    <p class="text-gray-900">
                      <span class="font-medium">{{ event.changed_by }}</span>
                      {{ formatEvent(event) }}
                    </p>
                    <p class="text-gray-500 text-xs">{{ formatDate(event.created_at) }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Sidebar -->
        <div class="space-y-6">
          <!-- Details Card -->
          <div class="bg-white shadow rounded-lg">
            <div class="px-4 py-5 sm:px-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">Details</h3>
              <dl class="space-y-3">
                <div v-if="task.project">
                  <dt class="text-sm font-medium text-gray-500">Project</dt>
                  <dd class="mt-1 text-sm text-gray-900 flex items-center gap-2">
                    <span>📁</span>
                    <router-link :to="`/projects`" class="text-indigo-600 hover:underline">
                      {{ task.project.name }}
                    </router-link>
                  </dd>
                </div>
                <div>
                  <dt class="text-sm font-medium text-gray-500">Status</dt>
                  <dd class="mt-1 text-sm text-gray-900">{{ formatStatus(task.status) }}</dd>
                </div>
                <div>
                  <dt class="text-sm font-medium text-gray-500">Priority</dt>
                  <dd class="mt-1 text-sm text-gray-900">{{ task.priority }}</dd>
                </div>
                <div>
                  <dt class="text-sm font-medium text-gray-500">Due Date</dt>
                  <dd class="mt-1 text-sm text-gray-900">{{ formatDate(task.due_date) }}</dd>
                </div>
                <div>
                  <dt class="text-sm font-medium text-gray-500">Created By</dt>
                  <dd class="mt-1 text-sm text-gray-900">{{ task.created_by?.email }}</dd>
                </div>
                <div>
                  <dt class="text-sm font-medium text-gray-500">Assigned Agent</dt>
                  <dd class="mt-1 text-sm text-gray-900">
                    {{ task.assigned_agent?.name || 'Unassigned' }}
                  </dd>
                </div>
                <div>
                  <dt class="text-sm font-medium text-gray-500">Created</dt>
                  <dd class="mt-1 text-sm text-gray-900">{{ formatDate(task.created_at) }}</dd>
                </div>
              </dl>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Task Modal -->
    <TaskModal
      v-if="showEditModal"
      :task="task"
      @close="showEditModal = false"
      @saved="handleTaskUpdated"
    />

    <!-- Assign Agent Modal -->
    <div v-if="showAssignModal" class="fixed z-10 inset-0 overflow-y-auto">
      <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="showAssignModal = false"></div>

        <span class="hidden sm:inline-block sm:align-middle sm:h-screen">&#8203;</span>

        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">
              Assign Agent
            </h3>

            <div class="space-y-2">
              <button
                v-for="agent in agentStore.agents"
                :key="agent.id"
                @click="handleAssign(agent.id)"
                class="w-full text-left px-4 py-3 border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
              >
                <div class="font-medium text-gray-900">{{ agent.name }}</div>
                <div class="text-sm text-gray-500">{{ agent.description || 'No description' }}</div>
              </button>
            </div>
          </div>

          <div class="bg-gray-50 px-4 py-3 sm:px-6">
            <button
              @click="showAssignModal = false"
              class="w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none sm:text-sm"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useTaskStore } from '../stores/tasks'
import { useAgentStore } from '../stores/agents'
import { taskService } from '../services/taskService'
import TaskModal from '../components/TaskModal.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const taskStore = useTaskStore()
const agentStore = useAgentStore()

const task = ref(null)
const comments = ref([])
const newComment = ref('')
const showEditModal = ref(false)
const showAssignModal = ref(false)

onMounted(async () => {
  await loadTask()
  await loadComments()
  await agentStore.fetchAgents()
})

const loadTask = async () => {
  task.value = await taskService.getTask(route.params.id)
}

const loadComments = async () => {
  comments.value = await taskService.getComments(route.params.id)
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const handleTaskUpdated = () => {
  showEditModal.value = false
  loadTask()
}

const handleDeleteTask = async () => {
  if (!confirm('Are you sure you want to delete this task?')) return

  try {
    await taskStore.deleteTask(route.params.id)
    router.push('/')
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to delete task')
  }
}

const handleAssign = async (agentId) => {
  try {
    await taskService.assignAgent(route.params.id, agentId)
    showAssignModal.value = false
    await loadTask()
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to assign agent')
  }
}

const handleUnassign = async () => {
  if (!confirm('Are you sure you want to unassign this agent?')) return

  try {
    await taskService.unassignAgent(route.params.id)
    await loadTask()
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to unassign agent')
  }
}

const handleAddComment = async () => {
  if (!newComment.value.trim()) return

  try {
    await taskService.addComment(route.params.id, newComment.value)
    newComment.value = ''
    await loadComments()
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to add comment')
  }
}

const formatStatus = (status) => {
  return status.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())
}

const getStatusClass = (status) => {
  const classes = {
    pending: 'bg-yellow-100 text-yellow-800',
    in_progress: 'bg-blue-100 text-blue-800',
    completed: 'bg-green-100 text-green-800',
    failed: 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getPriorityClass = (priority) => {
  const classes = {
    low: 'bg-gray-100 text-gray-800',
    medium: 'bg-blue-100 text-blue-800',
    high: 'bg-orange-100 text-orange-800',
    critical: 'bg-red-100 text-red-800'
  }
  return classes[priority] || 'bg-gray-100 text-gray-800'
}

const formatEvent = (event) => {
  const messages = {
    created: 'created this task',
    updated: 'updated this task',
    status_changed: `changed status from ${event.previous_state} to ${event.new_state}`,
    assigned: `assigned to ${event.new_state}`,
    unassigned: `unassigned ${event.previous_state}`
  }
  return messages[event.event_type] || event.event_type
}

const formatDate = (dateStr) => {
  if (!dateStr) return 'N/A'
  return new Date(dateStr).toLocaleString()
}
</script>
