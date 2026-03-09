import { createRouter, createWebHistory } from 'vue-router'
import { isAuthenticated } from '../utils/auth'
import Login from '../pages/Login.vue'
// import Register from '../pages/Register.vue' // Disabled: registration is closed
import AgentLogin from '../pages/AgentLogin.vue'
import Dashboard from '../pages/Dashboard.vue'
import Tasks from '../pages/Tasks.vue'
import Agents from '../pages/Agents.vue'
import Projects from '../pages/Projects.vue'
import TaskDetails from '../pages/TaskDetails.vue'
import ProjectDetails from '../pages/ProjectDetails.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  // Registration disabled
  // {
  //   path: '/register',
  //   name: 'Register',
  //   component: Register
  // },
  {
    path: '/register',
    redirect: '/login' // Redirect registration attempts to login
  },
  {
    path: '/agent/login',
    name: 'AgentLogin',
    component: AgentLogin
  },
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/tasks',
    name: 'Tasks',
    component: Tasks,
    meta: { requiresAuth: true }
  },
  {
    path: '/agents',
    name: 'Agents',
    component: Agents,
    meta: { requiresAuth: true }
  },
  {
    path: '/projects',
    name: 'Projects',
    component: Projects,
    meta: { requiresAuth: true }
  },
  {
    path: '/tasks/:id',
    name: 'TaskDetails',
    component: TaskDetails,
    meta: { requiresAuth: true }
  },
  {
    path: '/projects/:id',
    name: 'ProjectDetails',
    component: ProjectDetails,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const authenticated = isAuthenticated()

  if (to.meta.requiresAuth && !authenticated) {
    next('/login')
  } else if (to.name === 'Login' && authenticated) {
    next('/')
  } else {
    next()
  }
})

export default router
