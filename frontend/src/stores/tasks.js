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
        const tasks = await taskService.getTasks(this.filters)
        // Ensure all tasks have comments and events as arrays
        this.tasks = tasks.map(task => ({
          ...task,
          comments: Array.isArray(task.comments) ? task.comments : [],
          events: Array.isArray(task.events) ? task.events : [],
          project: task.project || null
        }))
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
        const task = await taskService.getTask(id)
        // Ensure task has comments and events as arrays
        this.currentTask = {
          ...task,
          comments: Array.isArray(task.comments) ? task.comments : [],
          events: Array.isArray(task.events) ? task.events : [],
          project: task.project || null
        }
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
        // Ensure task has comments and events as arrays
        const normalizedTask = {
          ...task,
          comments: Array.isArray(task.comments) ? task.comments : [],
          events: Array.isArray(task.events) ? task.events : [],
          project: task.project || null
        }
        this.tasks.unshift(normalizedTask)
        return normalizedTask
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
        // Ensure task has comments and events as arrays
        const normalizedTask = {
          ...task,
          comments: Array.isArray(task.comments) ? task.comments : [],
          events: Array.isArray(task.events) ? task.events : [],
          project: task.project || null
        }
        const index = this.tasks.findIndex(t => t.id === id)
        if (index !== -1) {
          this.tasks[index] = normalizedTask
        }
        if (this.currentTask?.id === id) {
          this.currentTask = normalizedTask
        }
        return normalizedTask
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
