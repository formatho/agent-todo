import { defineStore } from 'pinia'
import { agentService } from '../services/agentService'

export const useAgentStore = defineStore('agents', {
  state: () => ({
    agents: [],
    loading: false,
    error: null
  }),

  actions: {
    async fetchAgents() {
      this.loading = true
      this.error = null
      try {
        const response = await agentService.getAgents()
        this.agents = Array.isArray(response) ? response : []
      } catch (error) {
        this.error = error.message
      } finally {
        this.loading = false
      }
    },

    async fetchAgentsWithTasks() {
      this.loading = true
      this.error = null
      try {
        const response = await agentService.getAgentsWithTasks()
        this.agents = Array.isArray(response) ? response : []
      } catch (error) {
        this.error = error.message
      } finally {
        this.loading = false
      }
    },

    async createAgent(agentData) {
      this.loading = true
      this.error = null
      try {
        const agent = await agentService.createAgent(agentData)
        this.agents.push(agent)
        return agent
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateAgent(id, agentData) {
      this.loading = true
      this.error = null
      try {
        const updatedAgent = await agentService.updateAgent(id, agentData)
        const index = this.agents.findIndex(a => a.id === id)
        if (index !== -1) {
          this.agents[index] = updatedAgent
        }
        return updatedAgent
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteAgent(id) {
      this.loading = true
      this.error = null
      try {
        await agentService.deleteAgent(id)
        this.agents = this.agents.filter(a => a.id !== id)
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})
