<template>
  <div class="fixed z-10 inset-0 overflow-y-auto">
    <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="$emit('close')"></div>

      <span class="hidden sm:inline-block sm:align-middle sm:h-screen">&#8203;</span>

      <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
        <form @submit.prevent="handleSubmit">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg leading-6 font-medium text-gray-900">
                {{ task ? 'Edit Task' : 'Create New Task' }}
              </h3>
              <button 
                @click="$emit('close')" 
                class="text-gray-400 hover:text-gray-600 focus:outline-none"
              >
                <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <div class="space-y-6">
              <!-- Project Selection - Enhanced with guidance -->
              <div v-if="!task" class="project-selection">
                <div class="flex items-center justify-between mb-2">
                  <label class="block text-sm font-medium text-gray-700">
                    Project <span class="text-red-500">*</span>
                  </label>
                  <button 
                    @click="showProjectHelp = !showProjectHelp"
                    class="text-xs text-indigo-600 hover:text-indigo-800 flex items-center gap-1"
                  >
                    <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    Help
                  </button>
                </div>
                
                <div class="relative">
                  <select
                    v-model="form.project_id"
                    required
                    class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md border appearance-none"
                    :class="{ 
                      'border-red-300': !form.project_id && attemptedSubmit,
                      'border-indigo-500': form.project_id
                    }"
                  >
                    <option value="">Select a project...</option>
                    <option v-for="project in projectStore.activeProjects" :key="project.id" :value="project.id">
                      {{ project.name }}
                    </option>
                  </select>
                  <div class="absolute inset-y-0 right-0 flex items-center px-2 pointer-events-none">
                    <svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 20 20" stroke="currentColor">
                      <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
                    </svg>
                  </div>
                </div>
                
                <!-- Help and error messages -->
                <div v-if="showProjectHelp" class="mt-2 p-3 bg-blue-50 border border-blue-200 rounded-md">
                  <p class="text-xs text-blue-800">
                    📁 <strong>Projects help organize your tasks logically.</strong> Create a project first, then add tasks to it. This helps you group related work and track progress.
                  </p>
                </div>
                
                <p v-if="!form.project_id && attemptedSubmit" class="mt-1 text-sm text-red-500 flex items-center gap-1">
                  <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Please select a project for your task
                </p>
                
                <p v-if="projectStore.activeProjects.length === 0" class="mt-1 text-sm text-amber-600 bg-amber-50 p-2 rounded-md">
                  ⚠️ No active projects available. 
                  <router-link to="/projects" class="text-amber-700 hover:text-amber-900 font-medium underline">
                    Create your first project
                  </router-link>
                  to start adding tasks.
                </p>
                
                <p v-else-if="form.project_id" class="mt-1 text-sm text-green-600 flex items-center gap-1">
                  <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Project: {{ getSelectedProject()?.name }}
                </p>
              </div>

              <!-- Show project info for existing tasks -->
              <div v-else class="bg-gray-50 px-4 py-3 rounded-md border">
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-2">
                    <span class="text-sm text-gray-500">Project:</span>
                    <span class="text-sm font-medium text-gray-700">{{ task.project?.name || 'Unknown' }}</span>
                  </div>
                  <button 
                    @click="$emit('close'); $emit('project-selected', task.project_id)"
                    class="text-xs text-indigo-600 hover:text-indigo-800"
                  >
                    Change
                  </button>
                </div>
              </div>

              <!-- Task Title -->
              <div>
                <div class="flex items-center justify-between mb-2">
                  <label class="block text-sm font-medium text-gray-700">
                    Title <span class="text-red-500">*</span>
                  </label>
                  <span class="text-xs text-gray-500">{{ form.title.length }}/100</span>
                </div>
                <input
                  v-model="form.title"
                  type="text"
                  required
                  maxlength="100"
                  class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border p-2 transition-colors"
                  :class="{ 
                    'border-red-300': !form.title.trim() && attemptedSubmit,
                    'border-indigo-500': form.title.trim()
                  }"
                  placeholder="What needs to be done?"
                />
                <p v-if="!form.title.trim() && attemptedSubmit" class="mt-1 text-sm text-red-500 flex items-center gap-1">
                  <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Please enter a task title
                </p>
                <p class="mt-1 text-xs text-gray-500">
                  Be specific and clear about what needs to be accomplished.
                </p>
              </div>

              <!-- Task Description -->
              <div>
                <div class="flex items-center justify-between mb-2">
                  <label class="block text-sm font-medium text-gray-700">
                    Description
                  </label>
                  <button 
                    @click="showDescriptionHelp = !showDescriptionHelp"
                    class="text-xs text-indigo-600 hover:text-indigo-800 flex items-center gap-1"
                  >
                    Tips
                  </button>
                </div>
                <textarea
                  v-model="form.description"
                  rows="4"
                  maxlength="1000"
                  class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border p-2 transition-colors resize-none"
                  placeholder="Add details, requirements, or context for this task..."
                ></textarea>
                <div class="flex justify-between items-center mt-1">
                  <p v-if="showDescriptionHelp" class="text-xs text-blue-600">
                    💡 <strong>Pro tip:</strong> Include acceptance criteria, dependencies, or background information to help AI agents complete the task correctly.
                  </p>
                  <span class="text-xs text-gray-500">{{ form.description.length }}/1000</span>
                </div>
              </div>

              <!-- Priority and Due Date Row -->
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Priority <span class="text-red-500">*</span>
                  </label>
                  <select
                    v-model="form.priority"
                    required
                    class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md border transition-colors"
                    :class="{ 'border-indigo-500': form.priority }"
                  >
                    <option value="low">🟢 Low</option>
                    <option value="medium">🟡 Medium</option>
                    <option value="high">🟠 High</option>
                    <option value="critical">🔴 Critical</option>
                  </select>
                  <div class="mt-1 flex items-center gap-2">
                    <div class="w-2 h-2 rounded-full" :class="getPriorityColor(form.priority)"></div>
                    <span class="text-xs text-gray-500">{{ getPriorityDescription(form.priority) }}</span>
                  </div>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Due Date</label>
                  <input
                    v-model="form.due_date"
                    type="date"
                    class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border p-2 transition-colors"
                  />
                  <button 
                    @click="clearDueDate"
                    v-if="form.due_date"
                    class="mt-1 text-xs text-red-600 hover:text-red-800 flex items-center gap-1"
                  >
                    <svg class="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                    Clear
                  </button>
                </div>
              </div>

              <!-- Agent Assignment -->
              <div>
                <div class="flex items-center justify-between mb-2">
                  <label class="block text-sm font-medium text-gray-700">
                    Assign to AI Agent
                  </label>
                  <button 
                    @click="showAgentHelp = !showAgentHelp"
                    class="text-xs text-indigo-600 hover:text-indigo-800 flex items-center gap-1"
                  >
                    How it works
                  </button>
                </div>
                <select
                  v-model="form.assigned_agent_id"
                  class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md border transition-colors"
                  :class="{ 'border-indigo-500': form.assigned_agent_id }"
                >
                  <option value="">Unassigned (manual task)</option>
                  <option v-for="agent in agentStore.agents" :key="agent.id" :value="agent.id">
                    {{ agent.name }}
                  </option>
                </select>
                
                <!-- Agent help tooltip -->
                <div v-if="showAgentHelp" class="mt-2 p-3 bg-indigo-50 border border-indigo-200 rounded-md">
                  <p class="text-xs text-indigo-800">
                    🤖 <strong>AI Agents work autonomously!</strong> When you assign a task to an AI agent, it will use your project context and LLM instructions to complete the task automatically.
                  </p>
                </div>
                
                <!-- Selected agent info -->
                <div v-if="form.assigned_agent_id" class="mt-1 flex items-center gap-2 p-2 bg-blue-50 rounded-md">
                  <span class="text-xs text-blue-800">🤖 This task will be completed by AI automatically</span>
                </div>
              </div>

              <!-- Commit URL -->
              <div>
                <label class="block text-sm font-medium text-gray-700">Commit URL (Optional)</label>
                <div class="relative">
                  <input
                    v-model="form.commit_url"
                    type="url"
                    placeholder="https://github.com/org/repo/commit/abc123"
                    class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border p-2 transition-colors"
                  />
                  <button 
                    @click="clearCommitUrl"
                    v-if="form.commit_url"
                    class="absolute inset-y-0 right-3 pr-3 flex items-center text-gray-400 hover:text-red-600"
                  >
                    <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
                <p class="mt-1 text-xs text-gray-500">
                  Link to the Git commit (GitHub/GitLab) for this task. This helps track code changes.
                </p>
              </div>

              <!-- Task Template Suggestions -->
              <div v-if="!task && !form.description" class="task-templates">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Quick Templates
                </label>
                <div class="flex flex-wrap gap-2">
                  <button 
                    @click="fillTemplate('feature')"
                    class="template-btn"
                  >
                    🚀 Feature Task
                  </button>
                  <button 
                    @click="fillTemplate('bug')"
                    class="template-btn"
                  >
                    🐛 Bug Fix
                  </button>
                  <button 
                    @click="fillTemplate('refactor')"
                    class="template-btn"
                  >
                    🔧 Refactor
                  </button>
                  <button 
                    @click="fillTemplate('documentation')"
                    class="template-btn"
                  >
                    📚 Documentation
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Action buttons -->
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <div class="flex gap-2">
              <button
                type="button"
                @click="saveAsDraft"
                :disabled="loading"
                class="inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:text-sm disabled:opacity-50"
              >
                Save Draft
              </button>
              <button
                type="submit"
                :disabled="loading || (!task && !form.project_id) || !form.title.trim()"
                class="inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:text-sm disabled:opacity-50"
              >
                {{ loading ? 'Saving...' : (task ? 'Update Task' : 'Create Task') }}
              </button>
            </div>
            <button
              type="button"
              @click="$emit('close')"
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAgentStore } from '../stores/agents'
import { useProjectStore } from '../stores/projects'
import { taskService } from '../services/taskService'

const props = defineProps({
  task: {
    type: Object,
    default: null
  },
  preselectedProjectId: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['close', 'saved', 'project-selected'])

const agentStore = useAgentStore()
const projectStore = useProjectStore()
const loading = ref(false)
const attemptedSubmit = ref(false)
const showProjectHelp = ref(false)
const showDescriptionHelp = ref(false)
const showAgentHelp = ref(false)

const form = ref({
  title: '',
  description: '',
  priority: 'medium',
  due_date: '',
  commit_url: '',
  project_id: '',
  assigned_agent_id: ''
})

const templates = {
  feature: {
    title: 'Implement Feature X',
    description: `## Task Overview
Implement [Feature X] for the [Project Name] project.

### Acceptance Criteria
- [ ] Basic functionality implemented
- [ ] Unit tests written
- [ ] Integration tests pass
- [ ] Documentation updated

### Additional Information
- Follow existing code patterns
- Ensure accessibility compliance
- Update README if necessary`
  },
  bug: {
    title: 'Fix Bug in Component Y',
    description: `## Bug Report
Fix the bug reported in [Component Y] causing [specific issue].

### Bug Details
- **Issue**: [Describe the problem]
- **Environment**: [Browser/OS where it occurs]
- **Expected Behavior**: [What should happen]
- **Actual Behavior**: [What actually happens]

### Steps to Reproduce
1. [Step 1]
2. [Step 2]
3. [Step 3]

### Fix Requirements
- [ ] Implement proper error handling
- [ ] Add test case to prevent regression
- [ ] Update error messages if needed`
  },
  refactor: {
    title: 'Refactor Component Z',
    description: `## Refactoring Task
Improve the structure and maintainability of [Component Z].

### Current Issues
- [ ] Code duplication in [area]
- [ ] Poor separation of concerns
- [ ] Hard to test due to [reason]
- [ ] Performance bottlenecks in [area]

### Refactoring Goals
- [ ] Extract reusable components/functions
- [ ] Improve code readability
- [ ] Add proper error handling
- [ ] Ensure backward compatibility
- [ ] Update tests accordingly`
  },
  documentation: {
    title: 'Update Documentation for Module A',
    description: `## Documentation Update
Update and improve documentation for [Module A].

### Documentation Requirements
- [ ] API documentation generated
- [ ] Usage examples added
- [ ] Edge cases documented
- [ ] Contribution guide updated
- [ ] Troubleshooting section added

### Target Audience
- [ ] Developers using the API
- [ ] New team members onboarding
- [ ] External users/contributors`
  }
}

onMounted(async () => {
  await agentStore.fetchAgents()
  await projectStore.fetchProjects({ status: 'active' })

  if (props.task) {
    form.value = {
      title: props.task.title,
      description: props.task.description,
      priority: props.task.priority,
      due_date: props.task.due_date ? props.task.due_date.split('T')[0] : '',
      commit_url: props.task.commit_url || '',
      project_id: props.task.project_id || '',
      assigned_agent_id: props.task.assigned_agent_id || ''
    }
  } else if (props.preselectedProjectId) {
    form.value.project_id = props.preselectedProjectId
  }
})

const handleSubmit = async () => {
  attemptedSubmit.value = true

  if (!props.task && !form.value.project_id) {
    return
  }

  loading.value = true

  try {
    const taskData = {
      title: form.value.title.trim(),
      description: form.value.description.trim(),
      priority: form.value.priority,
      due_date: form.value.due_date ? new Date(form.value.due_date).toISOString() : null,
      commit_url: form.value.commit_url || '',
      assigned_agent_id: form.value.assigned_agent_id || null
    }

    // Only include project_id for new tasks
    if (!props.task) {
      taskData.project_id = form.value.project_id
    }

    if (props.task) {
      await taskService.updateTask(props.task.id, taskData)
    } else {
      await taskService.createTask(taskData)
    }

    emit('saved')
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to save task')
  } finally {
    loading.value = false
  }
}

const saveAsDraft = async () => {
  if (!form.value.title.trim()) {
    alert('Please enter a task title')
    return
  }

  loading.value = true
  try {
    await taskService.createTask({
      title: form.value.title.trim(),
      description: form.value.description.trim(),
      priority: form.value.priority,
      due_date: form.value.due_date ? new Date(form.value.due_date).toISOString() : null,
      commit_url: form.value.commit_url || '',
      assigned_agent_id: form.value.assigned_agent_id || null,
      project_id: form.value.project_id
    })
    
    emit('saved')
  } catch (error) {
    alert(error.response?.data?.error || 'Failed to save task')
  } finally {
    loading.value = false
  }
}

const fillTemplate = (templateType) => {
  const template = templates[templateType]
  if (template) {
    form.value.title = template.title
    form.value.description = template.description
  }
}

const clearDueDate = () => {
  form.value.due_date = ''
}

const clearCommitUrl = () => {
  form.value.commit_url = ''
}

const getSelectedProject = () => {
  return projectStore.projects.find(p => p.id === form.value.project_id)
}

const getPriorityColor = (priority) => {
  const colors = {
    low: 'bg-green-500',
    medium: 'bg-yellow-500',
    high: 'bg-orange-500',
    critical: 'bg-red-500'
  }
  return colors[priority] || 'bg-gray-500'
}

const getPriorityDescription = (priority) => {
  const descriptions = {
    low: 'Normal priority, can be done later',
    medium: 'Important, should be completed soon',
    high: 'Urgent, requires immediate attention',
    critical: 'Blocker, must be done now'
  }
  return descriptions[priority] || 'Priority not set'
}
</script>

<style scoped>
/* Project selection styles */
.project-selection select {
  background-image: none;
  transition: all 0.2s ease;
}

.project-selection select:hover {
  border-color: #9CA3AF;
}

.project-selection select:focus {
  outline: none;
  border-color: #6366F1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

/* Priority select custom styling */
select option {
  padding: 8px;
}

select:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

/* Template buttons */
.template-btn {
  padding: 6px 12px;
  background: #F9FAFB;
  border: 1px solid #E5E7EB;
  border-radius: 6px;
  font-size: 12px;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s ease;
}

.template-btn:hover {
  background: #E5E7EB;
  border-color: #D1D5DB;
  transform: translateY(-1px);
}

/* Form input styles */
input, textarea, select {
  transition: all 0.2s ease;
}

input:focus, textarea:focus, select:focus {
  outline: none;
  border-color: #6366F1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

/* Loading states */
button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

button:not(:disabled):hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* Responsive design */
@media (max-width: 640px) {
  .task-templates {
    flex-direction: column;
  }
  
  .template-btn {
    width: 100%;
    text-align: left;
  }
  
  .grid-cols-2 {
    grid-template-columns: 1fr;
  }
}

/* Modal animations */
@media (max-width: 640px) {
  .inline-block {
    margin: 0;
  }
}
</style>