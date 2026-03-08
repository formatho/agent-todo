import api from './api'

export const taskService = {
  async getTasks(filters = {}) {
    const params = new URLSearchParams()
    if (filters.status) params.append('status', filters.status)
    if (filters.agent_id) params.append('agent_id', filters.agent_id)
    if (filters.priority) params.append('priority', filters.priority)
    if (filters.search) params.append('search', filters.search)

    const response = await api.get(`/tasks?${params}`)
    return response.data
  },

  async getTask(id) {
    const response = await api.get(`/tasks/${id}`)
    return response.data
  },

  async createTask(taskData) {
    const response = await api.post('/tasks', taskData)
    return response.data
  },

  async updateTask(id, updates) {
    const response = await api.patch(`/tasks/${id}`, updates)
    return response.data
  },

  async deleteTask(id) {
    const response = await api.delete(`/tasks/${id}`)
    return response.data
  },

  async assignAgent(taskId, agentId) {
    const response = await api.patch(`/tasks/${taskId}/assign`, { agent_id: agentId })
    return response.data
  },

  async unassignAgent(taskId) {
    const response = await api.patch(`/tasks/${taskId}/unassign`)
    return response.data
  },

  async getComments(taskId) {
    const response = await api.get(`/tasks/${taskId}/comments`)
    return response.data
  },

  async addComment(taskId, content) {
    const response = await api.post(`/tasks/${taskId}/comments`, { content })
    return response.data
  }
}
