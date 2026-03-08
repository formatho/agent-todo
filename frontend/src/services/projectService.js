import api from './api'

export const projectService = {
  // User endpoints (JWT auth)
  async getProjects(filters = {}) {
    const params = new URLSearchParams()
    if (filters.status) params.append('status', filters.status)
    if (filters.search) params.append('search', filters.search)

    const response = await api.get(`/projects?${params}`)
    return response.data
  },

  async getProject(id) {
    const response = await api.get(`/projects/${id}`)
    return response.data
  },

  async createProject(projectData) {
    const response = await api.post('/projects', projectData)
    return response.data
  },

  async updateProject(id, updates) {
    const response = await api.patch(`/projects/${id}`, updates)
    return response.data
  },

  async deleteProject(id) {
    const response = await api.delete(`/projects/${id}`)
    return response.data
  },

  async getProjectTasks(projectId) {
    const response = await api.get(`/projects/${projectId}/tasks`)
    return response.data
  },

  // Agent endpoints (API Key auth)
  // These use /agent/projects routes which are read-only for agents
  async getProjectsForAgent(filters = {}) {
    const params = new URLSearchParams()
    if (filters.status) params.append('status', filters.status)
    if (filters.search) params.append('search', filters.search)

    const response = await api.get(`/agent/projects?${params}`)
    return response.data
  },

  async getProjectForAgent(id) {
    const response = await api.get(`/agent/projects/${id}`)
    return response.data
  }
}
