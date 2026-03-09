<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation -->
    <nav class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex">
            <div class="flex-shrink-0 flex items-center">
              <h1 class="text-xl font-bold text-gray-900">Agent Todo Platform</h1>
            </div>
            <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
              <router-link
                to="/tasks"
                class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
              >
                Tasks
              </router-link>
              <router-link
                to="/agents"
                class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
              >
                Agents
              </router-link>
              <router-link
                to="/projects"
                class="border-indigo-500 text-gray-900 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
              >
                Projects
              </router-link>
              <router-link
                to="/"
                class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
              >
                Dashboard
              </router-link>
            </div>
          </div>
          <div class="flex items-center">
            <span class="text-gray-700 mr-4">{{ authStore.user?.email }}</span>
            <button
              @click="handleLogout"
              class="bg-white py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 hover:bg-gray-50"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="flex justify-between items-center mb-6">
        <div>
          <h2 class="text-2xl font-bold text-gray-900">Projects</h2>
          <p class="text-sm text-gray-500 mt-1">Manage your projects and organize tasks</p>
        </div>
        <button @click="showCreateModal = true" class="btn-create">
          + New Project
        </button>
      </div>

      <!-- Filter Tabs -->
      <div class="flex gap-2 mb-6">
        <button 
          @click="statusFilter = ''" 
          :class="['tab-btn', { active: !statusFilter }]"
        >
          All
        </button>
        <button 
          @click="statusFilter = 'active'" 
          :class="['tab-btn', { active: statusFilter === 'active' }]"
        >
          Active
        </button>
        <button 
          @click="statusFilter = 'completed'" 
          :class="['tab-btn', { active: statusFilter === 'completed' }]"
        >
          Completed
        </button>
        <button 
          @click="statusFilter = 'archived'" 
          :class="['tab-btn', { active: statusFilter === 'archived' }]"
        >
          Archived
        </button>
      </div>

      <!-- Projects Grid -->
      <div v-if="filteredProjects.length > 0" class="projects-grid">
        <router-link
          v-for="project in filteredProjects"
          :key="project.id"
          :to="`/projects/${project.id}`"
          class="project-card cursor-pointer"
        >
          <div class="project-header">
            <div class="project-icon">📁</div>
            <div class="project-status" :class="project.status">
              {{ project.status }}
            </div>
          </div>
          
          <h3 class="project-name">{{ project.name }}</h3>
          <p class="project-description">{{ project.description || 'No description' }}</p>
          
          <!-- Project Links -->
          <div class="project-links" v-if="project.repository_url || project.deployed_url || project.documentation_url">
            <a v-if="project.repository_url" :href="project.repository_url" target="_blank" class="project-link" @click.stop>
              <span class="link-icon">📦</span>
              <span>Repository</span>
            </a>
            <a v-if="project.deployed_url" :href="project.deployed_url" target="_blank" class="project-link" @click.stop>
              <span class="link-icon">🚀</span>
              <span>Live App</span>
            </a>
            <a v-if="project.documentation_url" :href="project.documentation_url" target="_blank" class="project-link" @click.stop>
              <span class="link-icon">📚</span>
              <span>Docs</span>
            </a>
          </div>
          
          <!-- LLM Context Preview -->
          <div class="llm-context-preview" v-if="project.llm_context" @click="showLLMContext(project)">
            <span class="context-icon">🤖</span>
            <span class="context-text">LLM Context Available</span>
            <span class="view-context">View →</span>
          </div>
          
          <div class="project-stats">
            <div class="stat">
              <span class="stat-value">{{ project.taskCount || 0 }}</span>
              <span class="stat-label">Tasks</span>
            </div>
            <div class="stat">
              <span class="stat-value">{{ project.completedCount || 0 }}</span>
              <span class="stat-label">Completed</span>
            </div>
          </div>

          <div class="project-progress" v-if="project.taskCount > 0">
            <div class="progress-bar">
              <div 
                class="progress-fill" 
                :style="{ width: `${(project.completedCount / project.taskCount) * 100}%` }"
              ></div>
            </div>
            <span class="progress-text">
              {{ Math.round((project.completedCount / project.taskCount) * 100) }}%
            </span>
          </div>

          <div class="project-actions">
            <button @click.prevent="viewProjectTasks(project)" class="btn-action view">
              View Tasks
            </button>
            <button @click.prevent="editProject(project)" class="btn-action edit">
              Edit
            </button>
            <button 
              v-if="project.status === 'active'" 
              @click.prevent="archiveProject(project)" 
              class="btn-action archive"
            >
              Archive
            </button>
            <button 
              v-else-if="project.status === 'archived'" 
              @click="activateProject(project)" 
              class="btn-action activate"
            >
              Activate
            </button>
            <button @click="deleteProject(project)" class="btn-action delete">
              Delete
            </button>
          </div>

          <div class="project-meta">
            Created {{ formatDate(project.created_at) }}
          </div>
        </router-link>
      </div>

      <!-- Empty State -->
      <div v-else class="empty-state">
        <div class="empty-icon">📁</div>
        <h3>No projects found</h3>
        <p>{{ statusFilter ? 'No projects with this status.' : 'Create your first project to start organizing tasks.' }}</p>
        <button v-if="!statusFilter" @click="showCreateModal = true" class="btn-create-empty">
          Create Project
        </button>
      </div>
    </div>

    <!-- Create/Edit Project Modal -->
    <div v-if="showCreateModal || showEditModal" class="modal-overlay" @click.self="closeModals">
      <div class="modal-content">
        <h3 class="modal-title">{{ showEditModal ? 'Edit Project' : 'Create New Project' }}</h3>
        
        <form @submit.prevent="handleSubmit">
          <div class="form-group">
            <label class="form-label">Project Name <span class="required">*</span></label>
            <input
              v-model="form.name"
              type="text"
              required
              class="form-input"
              placeholder="e.g., Website Redesign"
            />
          </div>

          <div class="form-group">
            <label class="form-label">Description</label>
            <textarea
              v-model="form.description"
              rows="3"
              class="form-input"
              placeholder="Brief description of the project"
            ></textarea>
          </div>

          <div v-if="showEditModal" class="form-group">
            <label class="form-label">Status</label>
            <select v-model="form.status" class="form-input">
              <option value="active">Active</option>
              <option value="completed">Completed</option>
              <option value="archived">Archived</option>
            </select>
          </div>

          <div class="form-group">
            <label class="form-label">Repository URL</label>
            <input
              v-model="form.repository_url"
              type="url"
              class="form-input"
              placeholder="https://github.com/org/repo"
            />
          </div>

          <div class="form-group">
            <label class="form-label">Deployed URL</label>
            <input
              v-model="form.deployed_url"
              type="url"
              class="form-input"
              placeholder="https://app.example.com"
            />
          </div>

          <div class="form-group">
            <label class="form-label">Documentation URL</label>
            <input
              v-model="form.documentation_url"
              type="url"
              class="form-input"
              placeholder="https://docs.example.com"
            />
          </div>

          <div class="form-group">
            <label class="form-label">LLM Context</label>
            <textarea
              v-model="form.llm_context"
              rows="6"
              class="form-input font-mono"
              placeholder="Instructions, guidelines, and goals for AI agents working on this project..."
            ></textarea>
            <p class="form-hint">Markdown supported. This context will be available to AI agents working on project tasks.</p>
          </div>

          <div class="modal-actions">
            <button type="button" @click="closeModals" class="btn-cancel">
              Cancel
            </button>
            <button type="submit" :disabled="loading" class="btn-submit">
              {{ loading ? 'Saving...' : (showEditModal ? 'Update' : 'Create') }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useProjectStore } from '../stores/projects'
import { useTaskStore } from '../stores/tasks'

const router = useRouter()
const authStore = useAuthStore()
const projectStore = useProjectStore()
const taskStore = useTaskStore()

const showCreateModal = ref(false)
const showEditModal = ref(false)
const editingProject = ref(null)
const loading = ref(false)
const statusFilter = ref('')

const form = ref({
  name: '',
  description: '',
  status: 'active',
  repository_url: '',
  deployed_url: '',
  documentation_url: '',
  llm_context: ''
})

const filteredProjects = computed(() => {
  let projects = Array.isArray(projectStore.projects) ? projectStore.projects : []
  const tasks = Array.isArray(taskStore.tasks) ? taskStore.tasks : []
  
  // Add task counts to projects
  projects = projects.map(project => {
    const projectTasks = tasks.filter(t => t.project_id === project.id)
    return {
      ...project,
      taskCount: projectTasks.length,
      completedCount: projectTasks.filter(t => t.status === 'completed').length
    }
  })

  if (statusFilter.value) {
    projects = projects.filter(p => p.status === statusFilter.value)
  }

  return projects
})

onMounted(async () => {
  await projectStore.fetchProjects()
  await taskStore.fetchTasks()
})

watch(statusFilter, async () => {
  await projectStore.fetchProjects(statusFilter.value ? { status: statusFilter.value } : {})
})

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

const closeModals = () => {
  showCreateModal.value = false
  showEditModal.value = false
  editingProject.value = null
  form.value = { name: '', description: '', status: 'active', repository_url: '', deployed_url: '', documentation_url: '', llm_context: '' }
}

const editProject = (project) => {
  editingProject.value = project
  form.value = {
    name: project.name,
    description: project.description || '',
    status: project.status,
    repository_url: project.repository_url || '',
    deployed_url: project.deployed_url || '',
    documentation_url: project.documentation_url || '',
    llm_context: project.llm_context || ''
  }
  showEditModal.value = true
}

const handleSubmit = async () => {
  loading.value = true
  try {
    if (showEditModal.value && editingProject.value) {
      await projectStore.updateProject(editingProject.value.id, form.value)
    } else {
      await projectStore.createProject(form.value)
    }
    closeModals()
    await projectStore.fetchProjects()
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to save project')
  } finally {
    loading.value = false
  }
}

const archiveProject = async (project) => {
  if (confirm(`Archive project "${project.name}"?`)) {
    try {
      await projectStore.updateProject(project.id, { status: 'archived' })
      await projectStore.fetchProjects()
    } catch (error) {
      alert('Failed to archive project')
    }
  }
}

const activateProject = async (project) => {
  try {
    await projectStore.updateProject(project.id, { status: 'active' })
    await projectStore.fetchProjects()
  } catch (error) {
    alert('Failed to activate project')
  }
}

const showLLMContext = (project) => {
  alert(`LLM Context for ${project.name}:\n\n${project.llm_context}`)
}

const deleteProject = async (project) => {
  if (confirm(`Are you sure you want to delete project "${project.name}"? This action cannot be undone.`)) {
    try {
      await projectStore.deleteProject(project.id)
      await projectStore.fetchProjects(statusFilter.value ? { status: statusFilter.value } : {})
    } catch (error) {
      alert(error.response?.data?.error || 'Failed to delete project')
    }
  }
}

const viewProjectTasks = (project) => {
  router.push(`/tasks?project_id=${project.id}`)
}

const formatDate = (dateStr) => {
  if (!dateStr) return 'N/A'
  return new Date(dateStr).toLocaleDateString()
}
</script>

<style scoped>
.btn-create {
  padding: 10px 20px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-create:hover {
  background: #2563EB;
}

/* Tabs */
.tab-btn {
  padding: 8px 16px;
  background: white;
  border: 1px solid #E5E7EB;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  color: #6B7280;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-btn:hover {
  background: #F9FAFB;
  border-color: #D1D5DB;
}

.tab-btn.active {
  background: #3B82F6;
  border-color: #3B82F6;
  color: white;
}

/* Projects Grid */
.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.project-card {
  background: white;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  padding: 20px;
  transition: all 0.2s ease;
}

.project-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.project-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.project-icon {
  font-size: 24px;
}

.project-status {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
  text-transform: capitalize;
}

.project-status.active {
  background: #D1FAE5;
  color: #065F46;
}

.project-status.completed {
  background: #DBEAFE;
  color: #1E40AF;
}

.project-status.archived {
  background: #F3F4F6;
  color: #374151;
}

.project-name {
  font-size: 18px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 8px 0;
}

.project-description {
  font-size: 14px;
  color: #6B7280;
  margin: 0 0 16px 0;
  line-height: 1.5;
}

.project-links {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.project-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #F3F4F6;
  border-radius: 6px;
  font-size: 12px;
  color: #374151;
  text-decoration: none;
  transition: all 0.2s ease;
}

.project-link:hover {
  background: #E5E7EB;
  color: #111827;
}

.link-icon {
  font-size: 14px;
}

.llm-context-preview {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  background: linear-gradient(135deg, #EFF6FF 0%, #DBEAFE 100%);
  border: 1px solid #BFDBFE;
  border-radius: 8px;
  margin-bottom: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.llm-context-preview:hover {
  background: linear-gradient(135deg, #DBEAFE 0%, #BFDBFE 100%);
  border-color: #93C5FD;
}

.context-icon {
  font-size: 18px;
}

.context-text {
  flex: 1;
  font-size: 13px;
  font-weight: 500;
  color: #1E40AF;
}

.view-context {
  font-size: 12px;
  color: #3B82F6;
  font-weight: 600;
}

.project-stats {
  display: flex;
  gap: 24px;
  margin-bottom: 16px;
}

.stat {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #111827;
}

.stat-label {
  font-size: 12px;
  color: #6B7280;
}

.project-progress {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: #E5E7EB;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #10B981;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 12px;
  font-weight: 600;
  color: #374151;
}

.project-actions {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.btn-action {
  padding: 8px 12px;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-action.view {
  background: #3B82F6;
  color: white;
}

.btn-action.view:hover {
  background: #2563EB;
}

.btn-action.edit {
  background: #F3F4F6;
  color: #374151;
}

.btn-action.edit:hover {
  background: #E5E7EB;
}

.btn-action.archive {
  background: #FEF3C7;
  color: #92400E;
}

.btn-action.archive:hover {
  background: #FDE68A;
}

.btn-action.activate {
  background: #D1FAE5;
  color: #065F46;
}

.btn-action.activate:hover {
  background: #A7F3D0;
}

.btn-action.delete {
  background: #FEE2E2;
  color: #991B1B;
}

.btn-action.delete:hover {
  background: #FECACA;
}

.project-meta {
  font-size: 12px;
  color: #9CA3AF;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-state h3 {
  font-size: 20px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 8px 0;
}

.empty-state p {
  color: #6B7280;
  margin: 0 0 24px 0;
}

.btn-create-empty {
  padding: 12px 24px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-create-empty:hover {
  background: #2563EB;
}

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
}

.modal-content {
  background: white;
  border-radius: 12px;
  padding: 24px;
  width: 100%;
  max-width: 480px;
  margin: 20px;
}

.modal-title {
  font-size: 18px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 20px 0;
}

.form-group {
  margin-bottom: 16px;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 6px;
}

.required {
  color: #EF4444;
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #D1D5DB;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: #3B82F6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.font-mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.form-hint {
  margin-top: 6px;
  font-size: 12px;
  color: #6B7280;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.btn-cancel {
  padding: 10px 16px;
  background: white;
  border: 1px solid #D1D5DB;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  cursor: pointer;
}

.btn-cancel:hover {
  background: #F9FAFB;
}

.btn-submit {
  padding: 10px 16px;
  background: #3B82F6;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  color: white;
  cursor: pointer;
}

.btn-submit:hover {
  background: #2563EB;
}

.btn-submit:disabled {
  background: #9CA3AF;
  cursor: not-allowed;
}
</style>
