<template>
  <div class="compact-project-list">
    <div class="list-header">
      <h3 class="text-sm font-semibold text-gray-900">My Projects</h3>
      <button @click="showCreateModal = true" class="btn-add-project">
        <svg class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
      </button>
    </div>

    <!-- Create Project Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click="showCreateModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h4>Create New Project</h4>
          <button @click="showCreateModal = false" class="btn-close">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">Project Name</label>
            <input 
              v-model="newProject.name" 
              type="text" 
              class="form-input"
              placeholder="Enter project name"
            />
          </div>
          <div class="form-group">
            <label class="form-label">Description</label>
            <textarea 
              v-model="newProject.description" 
              class="form-textarea"
              placeholder="Enter project description"
              rows="2"
            ></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showCreateModal = false" class="btn-secondary">Cancel</button>
          <button @click="createProject" class="btn-primary" :disabled="!newProject.name.trim()">
            Create
          </button>
        </div>
      </div>
    </div>

    <!-- Project List -->
    <div class="projects-container">
      <div 
        v-for="project in projects" 
        :key="project.id"
        class="project-item"
        @click="selectProject(project)"
      >
        <div class="project-main">
          <div class="project-icon">
            <span class="icon-text">{{ getProjectIcon(project.name) }}</span>
          </div>
          <div class="project-info">
            <div class="project-name">{{ project.name }}</div>
            <div class="project-meta">
              <span class="meta-item">
                📋 {{ project.taskCount || 0 }} tasks
              </span>
              <span class="meta-item">
                🤖 {{ project.agentCount || 0 }} agents
              </span>
            </div>
          </div>
          <div class="project-status">
            <div class="status-indicator" :class="project.status"></div>
          </div>
        </div>
        <div class="project-progress">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: `${project.progress}%` }"></div>
          </div>
          <div class="progress-text">{{ project.progress }}% complete</div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-if="projects.length === 0" class="empty-projects">
        <div class="empty-icon">📁</div>
        <h4>No projects yet</h4>
        <p>Create your first project to organize your tasks</p>
        <button @click="showCreateModal = true" class="btn-primary">
          Create Project
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const showCreateModal = ref(false)
const newProject = ref({
  name: '',
  description: ''
})

const projects = ref([
  {
    id: 1,
    name: 'Backend Development',
    description: 'API and database implementation',
    status: 'active',
    progress: 75,
    taskCount: 12,
    agentCount: 2,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 7)
  },
  {
    id: 2,
    name: 'Frontend Design',
    description: 'User interface and experience design',
    status: 'active',
    progress: 60,
    taskCount: 8,
    agentCount: 1,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 5)
  },
  {
    id: 3,
    name: 'Documentation',
    description: 'User guides and API documentation',
    status: 'planning',
    progress: 25,
    taskCount: 5,
    agentCount: 1,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 3)
  },
  {
    id: 4,
    name: 'Marketing Campaign',
    description: 'Q2 marketing initiatives',
    status: 'completed',
    progress: 100,
    taskCount: 15,
    agentCount: 3,
    createdAt: new Date(Date.now() - 1000 * 60 * 60 * 24 * 14)
  }
])

const createProject = () => {
  if (!newProject.value.name.trim()) return
  
  const newProjectObj = {
    id: projects.value.length + 1,
    name: newProject.value.name.trim(),
    description: newProject.value.description,
    status: 'planning',
    progress: 0,
    taskCount: 0,
    agentCount: 0,
    createdAt: new Date()
  }
  
  projects.value.unshift(newProjectObj)
  
  // Reset form
  newProject.value = {
    name: '',
    description: ''
  }
  showCreateModal.value = false
}

const selectProject = (project) => {
  console.log('Selected project:', project)
  // In a real app, this would navigate to project details
}

const getProjectIcon = (name) => {
  const icons = {
    'backend': '⚙️',
    'frontend': '🎨',
    'design': '🎨',
    'docs': '📚',
    'doc': '📚',
    'marketing': '📢',
    'campaign': '📢',
    'api': '🔌',
    'database': '🗄️',
    'development': '💻',
    'dev': '💻'
  }
  
  const lowerName = name.toLowerCase()
  for (const [key, icon] of Object.entries(icons)) {
    if (lowerName.includes(key)) {
      return icon
    }
  }
  
  return '📁'
}
</script>

<style scoped>
.compact-project-list {
  background: white;
  border-radius: 8px;
  border: 1px solid #E5E7EB;
  overflow: hidden;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #E5E7EB;
  background: #F9FAFB;
}

.list-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #111827;
}

.btn-add-project {
  padding: 6px 8px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-add-project:hover {
  background: #2563EB;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  border: 1px solid #E5E7EB;
  width: 90%;
  max-width: 400px;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #E5E7EB;
}

.modal-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #111827;
}

.btn-close {
  background: none;
  border: none;
  font-size: 18px;
  color: #6B7280;
  cursor: pointer;
  padding: 0;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 3px;
}

.btn-close:hover {
  background: #F3F4F6;
  color: #374151;
}

.modal-body {
  padding: 16px;
}

.form-group {
  margin-bottom: 12px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 11px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 4px;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #D1D5DB;
  border-radius: 4px;
  font-size: 12px;
  transition: all 0.2s ease;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #3B82F6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 60px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid #E5E7EB;
}

.projects-container {
  max-height: 300px;
  overflow-y: auto;
}

.project-item {
  padding: 12px 16px;
  border-bottom: 1px solid #F3F4F6;
  cursor: pointer;
  transition: all 0.2s ease;
}

.project-item:hover {
  background: #F9FAFB;
}

.project-item:last-child {
  border-bottom: none;
}

.project-main {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.project-icon {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: #E5E7EB;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.icon-text {
  font-size: 16px;
}

.project-info {
  flex: 1;
  min-width: 0;
}

.project-name {
  font-size: 13px;
  font-weight: 500;
  color: #111827;
  margin-bottom: 4px;
}

.project-meta {
  display: flex;
  gap: 8px;
  font-size: 11px;
  color: #6B7280;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 2px;
}

.project-status {
  flex-shrink: 0;
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-indicator.active {
  background: #10B981;
}

.status-indicator.planning {
  background: #F59E0B;
}

.status-indicator.completed {
  background: #6B7280;
}

.project-progress {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.progress-bar {
  width: 100%;
  height: 4px;
  background: #E5E7EB;
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #10B981;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 10px;
  color: #6B7280;
  text-align: right;
}

.empty-projects {
  text-align: center;
  padding: 40px 20px;
  color: #6B7280;
}

.empty-icon {
  font-size: 32px;
  margin-bottom: 12px;
  opacity: 0.5;
}

.empty-projects h4 {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin: 0 0 8px 0;
}

.empty-projects p {
  font-size: 12px;
  margin: 0 0 16px 0;
  line-height: 1.4;
}

.btn-primary {
  padding: 8px 16px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-primary:hover:not(:disabled) {
  background: #2563EB;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  padding: 8px 16px;
  background: white;
  color: #374151;
  border: 1px solid #D1D5DB;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-secondary:hover {
  background: #F9FAFB;
  border-color: #9CA3AF;
}
</style>