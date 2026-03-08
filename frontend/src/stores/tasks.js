import { defineStore } from 'pinia'
import { taskService } from '../services/taskService'

export const useTaskStore = defineStore('tasks', {
  state: () => ({
    tasks: [],
    currentTask: null,
    loading: false,
    error: null,
    filters: {
      status: '',
      agent_id: '',
      priority: '',
      search: ''
    }
  }),

  actions: {
    async fetchTasks() {
      this.loading = true
      this.error = null
      try {
        this.tasks = await taskService.getTasks(this.filters)
      } catch (error) {
        this.error = error.message
      } finally {
        this.loading = false
      }
    },

    async fetchTask(id) {
      this.loading = true
      this.error = null
      try {
        this.currentTask = await taskService.getTask(id)
      } catch (error) {
        this.error = error.message
      } finally {
        this.loading = false
      }
    },

    async createTask(taskData) {
      this.loading = true
      this.error = null
      try {
        const task = await taskService.createTask(taskData)
        this.tasks.unshift(task)
        return task
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateTask(id, updates) {
      this.loading = true
      this.error = null
      try {
        const task = await taskService.updateTask(id, updates)
        const index = this.tasks.findIndex(t => t.id === id)
        if (index !== -1) {
          this.tasks[index] = task
        }
        if (this.currentTask?.id === id) {
          this.currentTask = task
        }
        return task
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteTask(id) {
      this.loading = true
      this.error = null
      try {
        await taskService.deleteTask(id)
        this.tasks = this.tasks.filter(t => t.id !== id)
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.loading = false
      }
    },

    setFilters(filters) {
      this.filters = { ...this.filters, ...filters }
    }
  }
})
