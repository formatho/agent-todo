import { defineStore } from 'pinia'
import { activityService } from '../services/activityService'

export const useActivityStore = defineStore('activity', {
  state: () => ({
    activities: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchActivity(limit = 20) {
      this.loading = true
      this.error = null
      try {
        const response = await activityService.getActivity(limit)
        this.activities = Array.isArray(response) ? response : []
      } catch (error) {
        this.error = error.message
      } finally {
        this.loading = false
      }
    }
  }
})
