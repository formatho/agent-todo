import api from './api'

export const subtaskService = {
  // Get all subtasks for a task
  async getSubtasks(taskId) {
    const response = await api.get(`/tasks/${taskId}/subtasks`)
    return response.data
  },

  // Create a new subtask
  async createSubtask(taskId, title) {
    const response = await api.post(`/tasks/${taskId}/subtasks`, { title })
    return response.data
  },

  // Update a subtask
  async updateSubtask(subtaskId, updates) {
    const response = await api.patch(`/subtasks/${subtaskId}`, updates)
    return response.data
  },

  // Delete a subtask
  async deleteSubtask(subtaskId) {
    await api.delete(`/subtasks/${subtaskId}`)
  },

  // Toggle subtask status
  async toggleSubtask(subtaskId, currentStatus) {
    const newStatus = currentStatus === 'pending' ? 'completed' : 'pending'
    const response = await api.patch(`/subtasks/${subtaskId}`, { status: newStatus })
    return response.data
  },

  // Reorder subtasks
  async reorderSubtasks(taskId, subtaskIds) {
    const response = await api.post(`/tasks/${taskId}/subtasks/reorder`, { subtask_ids: subtaskIds })
    return response.data
  }
}
