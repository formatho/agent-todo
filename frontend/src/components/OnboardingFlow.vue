<template>
  <div class="onboarding-overlay" @click.self="closeOnboarding">
    <div class="onboarding-container">
      <!-- Progress Bar -->
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: `${(currentStep / totalSteps) * 100}%` }"></div>
      </div>

      <!-- Step Content -->
      <div class="step-content">
        <!-- Step 1: Welcome -->
        <div v-if="currentStep === 1" class="step">
          <div class="step-header">
            <div class="step-number">{{ currentStep }}</div>
            <h2>Welcome to Agent Todo!</h2>
          </div>
          <div class="step-body">
            <div class="welcome-message">
              <div class="icon">🚀</div>
              <p>Your AI-powered task management platform is ready to help you organize work and boost productivity.</p>
            </div>
            
            <div class="features-grid">
              <div class="feature-card">
                <div class="feature-icon">📁</div>
                <h4>Projects</h4>
                <p>Organize tasks into logical projects with custom contexts for AI agents</p>
              </div>
              <div class="feature-card">
                <div class="feature-icon">🤖</div>
                <h4>AI Agents</h4>
                <p>Assign tasks to AI agents that can work autonomously with your instructions</p>
              </div>
              <div class="feature-card">
                <div class="feature-icon">📊</div>
                <h4>Analytics</h4>
                <p>Track progress, completion rates, and team performance with detailed insights</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Step 2: Create Your First Project -->
        <div v-if="currentStep === 2" class="step">
          <div class="step-header">
            <div class="step-number">{{ currentStep }}</div>
            <h2>Create Your First Project</h2>
          </div>
          <div class="step-body">
            <p class="instruction">Every great task management system starts with a clear project structure. Let's create your first project together.</p>
            
            <div class="project-template">
              <div class="template-card">
                <div class="template-header">
                  <span class="template-icon">🎯</span>
                  <h4>Task Management</h4>
                </div>
                <p class="template-description">Organize your daily tasks and prioritize work effectively</p>
                <button @click="createTemplateProject('task-management')" class="template-select">
                  Use This Template
                </button>
              </div>
              
              <div class="template-card">
                <div class="template-header">
                  <span class="template-icon">💻</span>
                  <h4>Software Development</h4>
                </div>
                <p class="template-description">Manage development cycles, code reviews, and deployment tasks</p>
                <button @click="createTemplateProject('software-dev')" class="template-select">
                  Use This Template
                </button>
              </div>
              
              <div class="template-card">
                <div class="template-header">
                  <span class="template-icon">📝</span>
                  <h4>Content Creation</h4>
                </div>
                <p class="template-description">Plan and organize content creation workflows</p>
                <button @click="createTemplateProject('content-creation')" class="template-select">
                  Use This Template
                </button>
              </div>
            </div>
            
            <div class="custom-project-hint">
              <p>Or <button @click="showCustomForm = true" class="link-btn">create a custom project</button></p>
            </div>
            
            <!-- Custom Project Form -->
            <div v-if="showCustomForm" class="custom-form">
              <form @submit.prevent="createCustomProject">
                <div class="form-group">
                  <label>Project Name</label>
                  <input v-model="customProject.name" type="text" placeholder="e.g., My Awesome Project" required />
                </div>
                <div class="form-group">
                  <label>Description</label>
                  <textarea v-model="customProject.description" rows="3" placeholder="What's this project about?"></textarea>
                </div>
                <div class="form-actions">
                  <button type="button" @click="showCustomForm = false" class="btn-cancel">Cancel</button>
                  <button type="submit" class="btn-create">Create Project</button>
                </div>
              </form>
            </div>
          </div>
        </div>

        <!-- Step 3: Add Tasks -->
        <div v-if="currentStep === 3" class="step">
          <div class="step-header">
            <div class="step-number">{{ currentStep }}</div>
            <h2>Add Your First Tasks</h2>
          </div>
          <div class="step-body">
            <p class="instruction">Now let's populate your project with some sample tasks to get you started.</p>
            
            <div class="task-examples">
              <div class="task-example">
                <div class="task-icon">📋</div>
                <div class="task-content">
                  <h4>Setup Project Structure</h4>
                  <p>Configure project settings and add team members</p>
                  <div class="task-priority high">High Priority</div>
                </div>
              </div>
              
              <div class="task-example">
                <div class="task-icon">🔧</div>
                <div class="task-content">
                  <h4>Configure AI Agents</h4>
                  <p>Set up AI agents with proper context and permissions</p>
                  <div class="task-priority medium">Medium Priority</div>
                </div>
              </div>
              
              <div class="task-example">
                <div class="task-icon">📊</div>
                <div class="task-content">
                  <h4>Monitor Progress</h4>
                  <p>Check task completion rates and analytics</p>
                  <div class="task-priority low">Low Priority</div>
                </div>
              </div>
            </div>
            
            <div class="add-tasks-actions">
              <button @click="createSampleTasks" class="btn-primary">
                Add These Sample Tasks
              </button>
              <p class="skip-hint">
                Or <button @click="skipSampleTasks" class="link-btn">skip and add tasks manually later</button>
              </p>
            </div>
          </div>
        </div>

        <!-- Step 4: Success -->
        <div v-if="currentStep === 4" class="step">
          <div class="step-header">
            <div class="step-icon">✅</div>
            <h2>You're All Set!</h2>
          </div>
          <div class="step-body">
            <div class="success-message">
              <div class="success-icon">🎉</div>
              <p>You've successfully set up your Agent Todo workspace! You're now ready to manage tasks with AI-powered efficiency.</p>
            </div>
            
            <div class="next-steps">
              <h3>Next Steps:</h3>
              <ul>
                <li>Explore your dashboard to track progress</li>
                <li>Add AI agents to help with tasks</li>
                <li>Customize project settings and contexts</li>
                <li>Invite team members to collaborate</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <!-- Navigation -->
      <div class="step-navigation">
        <button 
          v-if="currentStep > 1" 
          @click="previousStep" 
          class="btn-nav btn-prev"
          :disabled="currentStep === 1"
        >
          ← Previous
        </button>
        
        <button 
          v-if="currentStep < totalSteps" 
          @click="nextStep" 
          class="btn-nav btn-next"
          :disabled="isCreating"
        >
          Next →
        </button>
        
        <button 
          v-if="currentStep === totalSteps" 
          @click="completeOnboarding" 
          class="btn-nav btn-complete"
        >
          Get Started →
        </button>
        
        <button 
          v-if="currentStep < totalSteps" 
          @click="skipOnboarding" 
          class="btn-nav btn-skip"
        >
          Skip Tutorial
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useProjectStore } from '../stores/projects'
import { useTaskStore } from '../stores/tasks'
import { useRouter } from 'vue-router'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'completed'])

const router = useRouter()
const projectStore = useProjectStore()
const taskStore = useTaskStore()

const currentStep = ref(1)
const totalSteps = 4
const showCustomForm = ref(false)
const isCreating = ref(false)

const customProject = ref({
  name: '',
  description: ''
})

const templateProjects = {
  'task-management': {
    name: 'Task Management',
    description: 'Organize and prioritize your daily tasks effectively',
    tasks: [
      { title: 'Set up daily planning routine', description: 'Establish a consistent workflow for task prioritization', priority: 'high' },
      { title: 'Review and organize tasks', description: 'Weekly review of task completion and priorities', priority: 'medium' },
      { title: 'Archive completed tasks', description: 'Maintain clean task lists by removing completed items', priority: 'low' }
    ]
  },
  'software-dev': {
    name: 'Software Development',
    description: 'Manage development lifecycle with automated AI assistance',
    tasks: [
      { title: 'Code review and validation', description: 'Review pull requests and ensure code quality standards', priority: 'high' },
      { title: 'Update documentation', description: 'Keep project documentation current and comprehensive', priority: 'medium' },
      { title: 'Performance optimization', description: 'Monitor and optimize application performance', priority: 'medium' },
      { title: 'Test coverage improvements', description: 'Increase test coverage and add new test suites', priority: 'low' }
    ]
  },
  'content-creation': {
    name: 'Content Creation',
    description: 'Streamline content production and publication workflows',
    tasks: [
      { title: 'Content planning and research', description: 'Plan content calendar and research topics', priority: 'high' },
      { title: 'Draft and edit content', description: 'Create high-quality drafts and edit for publication', priority: 'high' },
      { title: 'Publish and promote', description: 'Schedule publication and promote content across channels', priority: 'medium' },
      { title: 'Analytics and optimization', description: 'Review performance and optimize future content', priority: 'low' }
    ]
  }
}

const nextStep = () => {
  if (currentStep.value < totalSteps) {
    currentStep.value++
  }
}

const previousStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

const createTemplateProject = async (templateKey) => {
  isCreating.value = true
  try {
    const template = templateProjects[templateKey]
    await projectStore.createProject({
      name: template.name,
      description: template.description,
      status: 'active'
    })
    
    // Wait for project creation to complete
    await projectStore.fetchProjects()
    
    // Create sample tasks
    if (template.tasks) {
      const project = projectStore.projects.find(p => p.name === template.name)
      for (const task of template.tasks) {
        await taskStore.createTask({
          title: task.title,
          description: task.description,
          priority: task.priority,
          project_id: project.id
        })
      }
    }
    
    nextStep()
  } catch (error) {
    alert('Failed to create project. Please try again.')
  } finally {
    isCreating.value = false
  }
}

const createCustomProject = async () => {
  isCreating.value = true
  try {
    await projectStore.createProject({
      name: customProject.value.name,
      description: customProject.value.description,
      status: 'active'
    })
    
    await projectStore.fetchProjects()
    showCustomForm.value = false
    customProject.value = { name: '', description: '' }
    nextStep()
  } catch (error) {
    alert('Failed to create project. Please try again.')
  } finally {
    isCreating.value = false
  }
}

const createSampleTasks = async () => {
  isCreating.value = true
  try {
    // Get the first active project
    const projects = await projectStore.fetchProjects()
    const project = projects.find(p => p.status === 'active')
    
    if (project) {
      const sampleTasks = [
        { title: 'Setup Project Structure', description: 'Configure project settings and add team members', priority: 'high' },
        { title: 'Configure AI Agents', description: 'Set up AI agents with proper context and permissions', priority: 'medium' },
        { title: 'Monitor Progress', description: 'Check task completion rates and analytics', priority: 'low' }
      ]
      
      for (const task of sampleTasks) {
        await taskStore.createTask({
          title: task.title,
          description: task.description,
          priority: task.priority,
          project_id: project.id
        })
      }
    }
    
    nextStep()
  } catch (error) {
    alert('Failed to create tasks. Please try again.')
  } finally {
    isCreating.value = false
  }
}

const skipSampleTasks = () => {
  nextStep()
}

const skipOnboarding = () => {
  emit('close')
}

const closeOnboarding = () => {
  emit('close')
}

const completeOnboarding = () => {
  emit('completed')
  emit('close')
}
</script>

<style scoped>
.onboarding-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.onboarding-container {
  background: white;
  border-radius: 16px;
  padding: 32px;
  max-width: 800px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.3);
}

.progress-bar {
  height: 4px;
  background: #E5E7EB;
  border-radius: 2px;
  margin-bottom: 32px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #3B82F6 0%, #2563EB 100%);
  transition: width 0.3s ease;
}

.step-content {
  margin-bottom: 32px;
}

.step-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.step-number {
  width: 48px;
  height: 48px;
  background: #3B82F6;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 700;
}

.step-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #10B981 0%, #059669 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.step-header h2 {
  font-size: 24px;
  font-weight: 700;
  color: #111827;
  margin: 0;
}

/* Step 1: Welcome */
.welcome-message {
  text-align: center;
  margin-bottom: 32px;
}

.welcome-message .icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.feature-card {
  background: #F9FAFB;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  padding: 24px;
  text-align: center;
  transition: all 0.2s ease;
}

.feature-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.feature-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.feature-card h4 {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 8px 0;
}

.feature-card p {
  font-size: 13px;
  color: #6B7280;
  margin: 0;
  line-height: 1.5;
}

/* Step 2: Create Project */
.instruction {
  font-size: 16px;
  color: #374151;
  margin-bottom: 24px;
  line-height: 1.6;
}

.project-template {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.template-card {
  background: #F9FAFB;
  border: 2px solid #E5E7EB;
  border-radius: 12px;
  padding: 20px;
  transition: all 0.2s ease;
  cursor: pointer;
}

.template-card:hover {
  border-color: #3B82F6;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.1);
}

.template-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.template-icon {
  font-size: 24px;
}

.template-card h4 {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  margin: 0;
}

.template-description {
  font-size: 13px;
  color: #6B7280;
  margin-bottom: 16px;
  line-height: 1.5;
}

.template-select {
  width: 100%;
  padding: 10px 16px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.template-select:hover {
  background: #2563EB;
}

.custom-project-hint {
  text-align: center;
  padding: 16px;
  background: #EFF6FF;
  border-radius: 8px;
}

.link-btn {
  color: #3B82F6;
  text-decoration: none;
  font-weight: 600;
  cursor: pointer;
  background: none;
  border: none;
  padding: 0;
}

.link-btn:hover {
  text-decoration: underline;
}

.custom-form {
  background: #F9FAFB;
  border: 1px solid #E5E7EB;
  border-radius: 12px;
  padding: 24px;
  margin-top: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 6px;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #D1D5DB;
  border-radius: 6px;
  font-size: 14px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
}

/* Step 3: Add Tasks */
.task-examples {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 24px;
}

.task-example {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: #F9FAFB;
  border-radius: 12px;
  border: 1px solid #E5E7EB;
}

.task-icon {
  font-size: 24px;
  width: 40px;
  height: 40px;
  background: #E5E7EB;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.task-content h4 {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 4px 0;
}

.task-content p {
  font-size: 13px;
  color: #6B7280;
  margin: 0 0 8px 0;
}

.task-priority {
  font-size: 11px;
  font-weight: 600;
  padding: 4px 8px;
  border-radius: 6px;
  text-transform: uppercase;
}

.task-priority.high {
  background: #FEE2E2;
  color: #991B1B;
}

.task-priority.medium {
  background: #DBEAFE;
  color: #1E40AF;
}

.task-priority.low {
  background: #F3F4F6;
  color: #374151;
}

.add-tasks-actions {
  text-align: center;
}

.add-tasks-actions .btn-primary {
  padding: 12px 24px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 12px;
}

.add-tasks-actions .btn-primary:hover {
  background: #2563EB;
}

/* Step 4: Success */
.success-message {
  text-align: center;
  margin-bottom: 32px;
}

.success-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.success-message p {
  font-size: 16px;
  color: #374151;
  line-height: 1.6;
}

.next-steps {
  background: #F9FAFB;
  border-radius: 12px;
  padding: 24px;
}

.next-steps h3 {
  font-size: 18px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 16px 0;
}

.next-steps ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.next-steps li {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  color: #374151;
}

.next-steps li::before {
  content: "✓";
  background: #10B981;
  color: white;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
}

/* Navigation */
.step-navigation {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 24px;
  border-top: 1px solid #E5E7EB;
}

.btn-nav {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 14px;
}

.btn-nav:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-prev {
  background: white;
  color: #374151;
  border: 1px solid #D1D5DB;
}

.btn-prev:hover:not(:disabled) {
  background: #F9FAFB;
  border-color: #9CA3AF;
}

.btn-next {
  background: #3B82F6;
  color: white;
}

.btn-next:hover:not(:disabled) {
  background: #2563EB;
}

.btn-complete {
  background: #10B981;
  color: white;
}

.btn-complete:hover {
  background: #059669;
}

.btn-skip {
  background: white;
  color: #6B7280;
  border: 1px solid #D1D5DB;
}

.btn-skip:hover {
  background: #F9FAFB;
  color: #374151;
}

/* Responsive Design */
@media (max-width: 768px) {
  .onboarding-container {
    padding: 20px;
    margin: 20px;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
  
  .project-template {
    grid-template-columns: 1fr;
  }
  
  .step-navigation {
    flex-direction: column;
    gap: 12px;
  }
  
  .task-example {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>