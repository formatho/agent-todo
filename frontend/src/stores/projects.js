import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { projectService } from '../services/projectService'

export const useProjectStore = defineStore('projects', () => {
  const projects = ref([])
  const currentProject = ref(null)
  const loading = ref(false)
  const error = ref(null)

  const activeProjects = computed(() => {
    return projects.value.filter(p => p.status === 'active')
  })

  const fetchProjects = async (filters = {}) => {
    loading.value = true
    error.value = null
    try {
      projects.value = await projectService.getProjects(filters)
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch projects'
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchProject = async (id) => {
    loading.value = true
    error.value = null
    try {
      currentProject.value = await projectService.getProject(id)
      return currentProject.value
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch project'
      throw err
    } finally {
      loading.value = false
    }
  }

  const createProject = async (projectData) => {
    loading.value = true
    error.value = null
    try {
      const project = await projectService.createProject(projectData)
      projects.value.push(project)
      return project
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to create project'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateProject = async (id, updates) => {
    loading.value = true
    error.value = null
    try {
      const updated = await projectService.updateProject(id, updates)
      const index = projects.value.findIndex(p => p.id === id)
      if (index !== -1) {
        projects.value[index] = updated
      }
      if (currentProject.value?.id === id) {
        currentProject.value = updated
      }
      return updated
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to update project'
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteProject = async (id) => {
    loading.value = true
    error.value = null
    try {
      await projectService.deleteProject(id)
      projects.value = projects.value.filter(p => p.id !== id)
      if (currentProject.value?.id === id) {
        currentProject.value = null
      }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to delete project'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    projects,
    currentProject,
    loading,
    error,
    activeProjects,
    fetchProjects,
    fetchProject,
    createProject,
    updateProject,
    deleteProject
  }
})
