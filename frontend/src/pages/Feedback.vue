<template>
  <div class="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-3xl mx-auto">
      <!-- Header -->
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-gray-900 mb-2">Beta Feedback</h1>
        <p class="text-lg text-gray-600">Help us improve Formatho Agent Todo</p>
      </div>

      <!-- Success Message -->
      <div v-if="submitted" class="bg-green-50 border border-green-200 rounded-lg p-6 mb-6">
        <div class="flex items-center">
          <svg class="w-6 h-6 text-green-600 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
          </svg>
          <div>
            <h3 class="text-green-800 font-semibold text-lg">Thank you for your feedback!</h3>
            <p class="text-green-700">We appreciate your input and will review it shortly.</p>
          </div>
        </div>
        <button @click="resetForm" class="mt-4 text-green-700 underline hover:text-green-800">
          Submit another feedback
        </button>
      </div>

      <!-- Feedback Form -->
      <form v-else @submit.prevent="submitFeedback" class="bg-white shadow-sm rounded-lg p-8">
        <!-- Feedback Type -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Feedback Type <span class="text-red-500">*</span>
          </label>
          <select v-model="form.feedbackType" required
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
            <option value="">Select type...</option>
            <option value="bug">🐛 Bug Report</option>
            <option value="feature_request">✨ Feature Request</option>
            <option value="improvement">💡 Improvement Suggestion</option>
            <option value="general">💬 General Feedback</option>
          </select>
        </div>

        <!-- Title -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Title <span class="text-red-500">*</span>
          </label>
          <input v-model="form.title" type="text" required
                 placeholder="Brief summary of your feedback"
                 class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
        </div>

        <!-- Description -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Description <span class="text-red-500">*</span>
          </label>
          <textarea v-model="form.description" required rows="5"
                    placeholder="Please provide detailed information..."
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"></textarea>
        </div>

        <!-- Optional Fields -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <!-- Name -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Your Name (Optional)
            </label>
            <input v-model="form.testerName" type="text"
                   placeholder="John Doe"
                   class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
          </div>

          <!-- Email -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Email (Optional, for follow-up)
            </label>
            <input v-model="form.testerEmail" type="email"
                   placeholder="john@example.com"
                   class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
          </div>
        </div>

        <!-- Priority -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Priority (Optional)
          </label>
          <select v-model="form.priority"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
            <option value="">Not specified</option>
            <option value="low">Low</option>
            <option value="medium">Medium</option>
            <option value="high">High</option>
            <option value="critical">Critical</option>
          </select>
        </div>

        <!-- Rating -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Overall Experience Rating (Optional)
          </label>
          <div class="flex items-center space-x-2">
            <button v-for="i in 5" :key="i"
                    type="button"
                    @click="form.rating = i"
                    class="text-3xl transition-transform hover:scale-110">
              <span :class="i <= form.rating ? 'text-yellow-400' : 'text-gray-300'">★</span>
            </button>
            <span v-if="form.rating" class="ml-2 text-sm text-gray-600">{{ form.rating }}/5</span>
          </div>
        </div>

        <!-- Current Page Info -->
        <input type="hidden" v-model="form.page">

        <!-- Error Message -->
        <div v-if="error" class="mb-4 p-4 bg-red-50 border border-red-200 rounded-md">
          <p class="text-red-700">{{ error }}</p>
        </div>

        <!-- Submit Button -->
        <button type="submit"
                :disabled="loading"
                class="w-full bg-blue-600 text-white py-3 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">
          <span v-if="loading" class="flex items-center justify-center">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Submitting...
          </span>
          <span v-else>Submit Feedback</span>
        </button>
      </form>

      <!-- Footer -->
      <div class="mt-8 text-center text-sm text-gray-500">
        <p>Your feedback is valuable to us and helps shape the future of Formatho Agent Todo.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const loading = ref(false)
const submitted = ref(false)
const error = ref('')

const form = ref({
  feedbackType: '',
  title: '',
  description: '',
  testerName: '',
  testerEmail: '',
  priority: '',
  rating: 0,
  page: window.location.href
})

const submitFeedback = async () => {
  loading.value = true
  error.value = ''

  try {
    const response = await fetch('/api/feedback', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(form.value)
    })

    if (!response.ok) {
      const data = await response.json()
      throw new Error(data.error || 'Failed to submit feedback')
    }

    submitted.value = true
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  form.value = {
    feedbackType: '',
    title: '',
    description: '',
    testerName: '',
    testerEmail: '',
    priority: '',
    rating: 0,
    page: window.location.href
  }
  submitted.value = false
  error.value = ''
}

onMounted(() => {
  // Pre-fill page from query parameter if provided
  if (route.query.page) {
    form.value.page = route.query.page
  }
})
</script>
