import api from './api'

export const activityService = {
  async getActivity(limit = 50) {
    const response = await api.get(`/activity?limit=${limit}`)
    return response.data
  }
}
