<template>
  <div class="fixed z-10 inset-0 overflow-y-auto">
    <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="$emit('close')"></div>

      <span class="hidden sm:inline-block sm:align-middle sm:h-screen">&#8203;</span>

      <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
        <form @submit.prevent="handleSubmit">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4">
              {{ task ? 'Edit Task' : 'Create New Task' }}
            </h3>

            <div class="space-y-4">
              <!-- Project Selection - Required for new tasks -->
              <div v-if="!task">
                <label class="block text-sm font-medium text-gray-700">
                  Project <span class="text-red-500">*</span>
                </label>
                <select
                  v-model="form.project_id"
                  required
                  class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md border"
                  :class="{ 'border-red-300': !form.project_id && attemptedSubmit }"
                >
                  <option value="">Select a project</option>
                  <option v-for="project in projectStore.activeProjects" :key="project.id" :value="project.id">
                    {{ project.name }}
                  </option>
                </select>
                <p v-if="!form.project_id && attemptedSubmit" class="mt-1 text-sm text-red-500">
                  Please select a project
                </p>
                <p v-if="projectStore.activeProjects.length === 0" class="mt-1 text-sm text-amber-600">
                  ⚠️ No active projects. <router-link to="/projects" class="underline">Create one first</router-link>
                </p>
              </div>

              <!-- Show project name for existing tasks -->
              <div v-else class="bg-gray-50 px-3 py-2 rounded-md">
                <span class="text-sm text-gray-500">Project:</span>
                <span class="text-sm font-medium text-gray-700 ml-1">
                  {{ task.project?.name || 'Unknown' }}
                </span>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700">
                  Title <span class="text-red-500">*</span>
                </label>
                <input
                  v-model="form.title"
                  type="text"
                  required
                  class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border p-2"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700">Description</label>
                <textarea
                  v-model="form.description"
                  rows="3"
                  class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border p-2"
                ></textarea>
              </div>

              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700">
                    Priority <span class="text-red-500">*</span>
                  </label>
                  <select
                    v-model="form.priority"
                    required
                    class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md border"
                  >
                    <option value="low">Low</option>
                    <option value="medium">Medium</option>
                    <option value="high">High</option>
                    <option value="critical">Critical</option>
                  </select>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700">Due Date</label>
                  <input
                    v-model="form.due_date"
                    type="date"
                    class="mt-1 focus:ring-indigo-500 focus:border-indigo-500 block w-full shadow-sm sm:text-sm border-gray-300 rounded-md border p-2"
                  />
                </div>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700">Assign to Agent</label>
                <select
                  v-model="form.assigned_agent_id"
                  class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md border"
                >
                  <option value="">Unassigned</option>
                  <option v-for="agent in agentStore.agents" :key="agent.id" :value="agent.id">
                    {{ agent.name }}
                  </option>
                </select>
              </div>
            </div>
          </div>

          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button
              type="submit"
              :disabled="loading || (!task && !form.project_id)"
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
            >
              {{ loading ? 'Saving...' : (task ? 'Update' : 'Create') }}
            </button>
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
import { ref, onMounted } from 'vue'
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

const emit = defineEmits(['close', 'saved'])

const agentStore = useAgentStore()
const projectStore = useProjectStore()
const loading = ref(false)
const attemptedSubmit = ref(false)

const form = ref({
  title: '',
  description: '',
  priority: 'medium',
  due_date: '',
  project_id: '',
  assigned_agent_id: ''
})

onMounted(async () => {
  await agentStore.fetchAgents()
  await projectStore.fetchProjects({ status: 'active' })

  if (props.task) {
    form.value = {
      title: props.task.title,
      description: props.task.description,
      priority: props.task.priority,
      due_date: props.task.due_date ? props.task.due_date.split('T')[0] : '',
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
      title: form.value.title,
      description: form.value.description,
      priority: form.value.priority,
      due_date: form.value.due_date ? new Date(form.value.due_date).toISOString() : null,
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
</script>
