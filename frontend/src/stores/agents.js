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
        this.agents = await agentService.getAgents()
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
