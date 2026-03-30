<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <h1 class="text-3xl font-bold text-gray-900">Beta Feedback Dashboard</h1>
        <p class="mt-1 text-sm text-gray-600">Review and manage beta tester feedback</p>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-5 gap-4 mb-8" v-if="stats">
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-2xl font-bold text-gray-900">{{ stats.total }}</div>
          <div class="text-sm text-gray-600">Total Feedback</div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-2xl font-bold text-blue-600">{{ stats.recent_count }}</div>
          <div class="text-sm text-gray-600">Last 7 Days</div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-2xl font-bold text-yellow-600">{{ stats.average_rating?.toFixed(1) || '-' }}</div>
          <div class="text-sm text-gray-600">Avg Rating</div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-2xl font-bold text-red-600">{{ getStatusCount('new') }}</div>
          <div class="text-sm text-gray-600">New</div>
        </div>
        <div class="bg-white rounded-lg shadow p-6">
          <div class="text-2xl font-bold text-green-600">{{ getStatusCount('resolved') + getStatusCount('closed') }}</div>
          <div class="text-sm text-gray-600">Resolved</div>
        </div>
      </div>

      <!-- Filters -->
      <div class="bg-white rounded-lg shadow mb-6 p-6">
        <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select v-model="filters.status" @change="loadFeedback"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
              <option value="">All</option>
              <option value="new">New</option>
              <option value="acknowledged">Acknowledged</option>
              <option value="in_progress">In Progress</option>
              <option value="resolved">Resolved</option>
              <option value="closed">Closed</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
            <select v-model="filters.feedbackType" @change="loadFeedback"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
              <option value="">All</option>
              <option value="bug">Bug</option>
              <option value="feature_request">Feature Request</option>
              <option value="improvement">Improvement</option>
              <option value="general">General</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Priority</label>
            <select v-model="filters.priority" @change="loadFeedback"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
              <option value="">All</option>
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
              <option value="critical">Critical</option>
            </select>
          </div>
          <div class="flex items-end">
            <button @click="resetFilters"
                    class="px-4 py-2 text-sm text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md transition-colors">
              Reset Filters
            </button>
          </div>
        </div>
      </div>

      <!-- Feedback List -->
      <div class="bg-white rounded-lg shadow overflow-hidden">
        <div v-if="loading" class="p-12 text-center">
          <svg class="animate-spin h-8 w-8 text-blue-600 mx-auto" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <p class="mt-2 text-gray-600">Loading feedback...</p>
        </div>

        <div v-else-if="feedback.length === 0" class="p-12 text-center">
          <svg class="h-12 w-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
          </svg>
          <p class="text-gray-600">No feedback found</p>
        </div>

        <div v-else class="divide-y divide-gray-200">
          <div v-for="item in feedback" :key="item.id"
               class="p-6 hover:bg-gray-50 cursor-pointer transition-colors"
               @click="openFeedbackModal(item)">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <div class="flex items-center space-x-3 mb-2">
                  <span class="px-2 py-1 text-xs font-medium rounded-full" :class="getTypeClass(item.feedback_type)">
                    {{ getTypeIcon(item.feedback_type) }} {{ formatType(item.feedback_type) }}
                  </span>
                  <span class="px-2 py-1 text-xs font-medium rounded-full" :class="getStatusClass(item.status)">
                    {{ formatStatus(item.status) }}
                  </span>
                  <span v-if="item.priority" class="px-2 py-1 text-xs font-medium rounded-full"
                        :class="getPriorityClass(item.priority)">
                    {{ item.priority }}
                  </span>
                </div>
                <h3 class="text-lg font-medium text-gray-900 mb-1">{{ item.title }}</h3>
                <p class="text-sm text-gray-600 line-clamp-2">{{ item.description }}</p>
                <div class="flex items-center space-x-4 mt-2 text-xs text-gray-500">
                  <span v-if="item.tester_name">👤 {{ item.tester_name }}</span>
                  <span v-if="item.tester_email">📧 {{ item.tester_email }}</span>
                  <span>📅 {{ formatDate(item.created_at) }}</span>
                  <span v-if="item.rating > 0">⭐ {{ item.rating }}/5</span>
                </div>
              </div>
              <svg class="h-5 w-5 text-gray-400 ml-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
              </svg>
            </div>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="total > limit" class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
          <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
            <div>
              <p class="text-sm text-gray-700">
                Showing <span class="font-medium">{{ offset + 1 }}</span> to
                <span class="font-medium">{{ Math.min(offset + limit, total) }}</span> of
                <span class="font-medium">{{ total }}</span> results
              </p>
            </div>
            <div>
              <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                <button @click="previousPage" :disabled="offset === 0"
                        class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed">
                  Previous
                </button>
                <button @click="nextPage" :disabled="offset + limit >= total"
                        class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed">
                  Next
                </button>
              </nav>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Feedback Detail Modal -->
    <div v-if="selectedFeedback" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
         @click.self="closeFeedbackModal">
      <div class="bg-white rounded-lg max-w-3xl w-full max-h-[90vh] overflow-y-auto">
        <div class="p-6">
          <div class="flex items-start justify-between mb-4">
            <h2 class="text-2xl font-bold text-gray-900">{{ selectedFeedback.title }}</h2>
            <button @click="closeFeedbackModal" class="text-gray-400 hover:text-gray-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          </div>

          <div class="space-y-4">
            <!-- Metadata -->
            <div class="flex items-center space-x-3">
              <span class="px-2 py-1 text-xs font-medium rounded-full"
                    :class="getTypeClass(selectedFeedback.feedback_type)">
                {{ getTypeIcon(selectedFeedback.feedback_type) }} {{ formatType(selectedFeedback.feedback_type) }}
              </span>
              <span class="px-2 py-1 text-xs font-medium rounded-full"
                    :class="getStatusClass(selectedFeedback.status)">
                {{ formatStatus(selectedFeedback.status) }}
              </span>
            </div>

            <!-- Description -->
            <div>
              <h3 class="text-sm font-medium text-gray-700 mb-2">Description</h3>
              <p class="text-gray-900 whitespace-pre-wrap">{{ selectedFeedback.description }}</p>
            </div>

            <!-- Submitter Info -->
            <div class="bg-gray-50 rounded-lg p-4">
              <h3 class="text-sm font-medium text-gray-700 mb-2">Submitter Information</h3>
              <div class="space-y-1 text-sm text-gray-600">
                <p v-if="selectedFeedback.tester_name">Name: {{ selectedFeedback.tester_name }}</p>
                <p v-if="selectedFeedback.tester_email">Email: {{ selectedFeedback.tester_email }}</p>
                <p>Submitted: {{ formatDate(selectedFeedback.created_at) }}</p>
                <p v-if="selectedFeedback.rating > 0">Rating: {{ selectedFeedback.rating }}/5</p>
                <p v-if="selectedFeedback.page">Page: {{ selectedFeedback.page }}</p>
              </div>
            </div>

            <!-- Update Status -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Update Status</label>
              <select v-model="statusUpdate.status" @change="updateFeedbackStatus"
                      class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
                <option value="new">New</option>
                <option value="acknowledged">Acknowledged</option>
                <option value="in_progress">In Progress</option>
                <option value="resolved">Resolved</option>
                <option value="closed">Closed</option>
              </select>
            </div>

            <!-- Admin Notes -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Admin Notes</label>
              <textarea v-model="statusUpdate.adminNotes" rows="3"
                        placeholder="Add internal notes..."
                        class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"></textarea>
              <button @click="updateFeedbackNotes"
                      class="mt-2 px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors">
                Save Notes
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useToast } from 'vue-toastification'

const toast = useToast()

const loading = ref(false)
const feedback = ref([])
const total = ref(0)
const limit = ref(50)
const offset = ref(0)
const stats = ref(null)
const selectedFeedback = ref(null)

const filters = ref({
  status: '',
  feedbackType: '',
  priority: ''
})

const statusUpdate = ref({
  status: '',
  adminNotes: ''
})

const loadFeedback = async () => {
  loading.value = true
  try {
    const params = new URLSearchParams()
    params.append('limit', limit.value)
    params.append('offset', offset.value)
    
    if (filters.value.status) params.append('status', filters.value.status)
    if (filters.value.feedbackType) params.append('feedback_type', filters.value.feedbackType)
    if (filters.value.priority) params.append('priority', filters.value.priority)

    const response = await fetch(`/api/feedback?${params}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (!response.ok) throw new Error('Failed to load feedback')

    const data = await response.json()
    feedback.value = data.feedback
    total.value = data.total
  } catch (error) {
    toast.error('Failed to load feedback')
    console.error(error)
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const response = await fetch('/api/feedback/stats', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (!response.ok) throw new Error('Failed to load stats')

    stats.value = await response.json()
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
}

const openFeedbackModal = (item) => {
  selectedFeedback.value = item
  statusUpdate.value.status = item.status
  statusUpdate.value.adminNotes = item.admin_notes || ''
}

const closeFeedbackModal = () => {
  selectedFeedback.value = null
}

const updateFeedbackStatus = async () => {
  if (!selectedFeedback.value) return

  try {
    const response = await fetch(`/api/feedback/${selectedFeedback.value.id}/status`, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(statusUpdate.value)
    })

    if (!response.ok) throw new Error('Failed to update status')

    toast.success('Status updated')
    await loadFeedback()
    await loadStats()
  } catch (error) {
    toast.error('Failed to update status')
    console.error(error)
  }
}

const updateFeedbackNotes = async () => {
  if (!selectedFeedback.value) return

  try {
    const response = await fetch(`/api/feedback/${selectedFeedback.value.id}/notes`, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ admin_notes: statusUpdate.value.adminNotes })
    })

    if (!response.ok) throw new Error('Failed to update notes')

    toast.success('Notes updated')
    await loadFeedback()
  } catch (error) {
    toast.error('Failed to update notes')
    console.error(error)
  }
}

const resetFilters = () => {
  filters.value = {
    status: '',
    feedbackType: '',
    priority: ''
  }
  offset.value = 0
  loadFeedback()
}

const previousPage = () => {
  if (offset.value > 0) {
    offset.value -= limit.value
    loadFeedback()
  }
}

const nextPage = () => {
  if (offset.value + limit.value < total.value) {
    offset.value += limit.value
    loadFeedback()
  }
}

const getStatusCount = (status) => {
  if (!stats.value?.by_status) return 0
  const stat = stats.value.by_status.find(s => s.Status === status)
  return stat ? stat.Count : 0
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatType = (type) => {
  const types = {
    bug: 'Bug Report',
    feature_request: 'Feature Request',
    improvement: 'Improvement',
    general: 'General'
  }
  return types[type] || type
}

const getTypeIcon = (type) => {
  const icons = {
    bug: '🐛',
    feature_request: '✨',
    improvement: '💡',
    general: '💬'
  }
  return icons[type] || '📝'
}

const formatStatus = (status) => {
  const statuses = {
    new: 'New',
    acknowledged: 'Acknowledged',
    in_progress: 'In Progress',
    resolved: 'Resolved',
    closed: 'Closed'
  }
  return statuses[status] || status
}

const getTypeClass = (type) => {
  const classes = {
    bug: 'bg-red-100 text-red-800',
    feature_request: 'bg-purple-100 text-purple-800',
    improvement: 'bg-blue-100 text-blue-800',
    general: 'bg-gray-100 text-gray-800'
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const getStatusClass = (status) => {
  const classes = {
    new: 'bg-yellow-100 text-yellow-800',
    acknowledged: 'bg-blue-100 text-blue-800',
    in_progress: 'bg-indigo-100 text-indigo-800',
    resolved: 'bg-green-100 text-green-800',
    closed: 'bg-gray-100 text-gray-800'
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

onMounted(() => {
  loadFeedback()
  loadStats()
})
</script>
