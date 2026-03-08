import api from './api'

export const agentService = {
  async getAgents() {
    const response = await api.get('/agents')
    return response.data
  },

  async getAgentsWithTasks() {
    const response = await api.get('/agents/activity')
    return response.data
  },

  async getAgent(id) {
    const response = await api.get(`/agents/${id}`)
    return response.data
  },

  async createAgent(agentData) {
    const response = await api.post('/agents', agentData)
    return response.data
  },

  async updateAgent(id, agentData) {
    const response = await api.patch(`/agents/${id}`, agentData)
    return response.data
  },

  async deleteAgent(id) {
    const response = await api.delete(`/agents/${id}`)
    return response.data
  }
}
