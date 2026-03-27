import { createRouter, createWebHistory } from 'vue-router'
import { isAuthenticated } from '../utils/auth'

// Lazy load components for better performance (Core Web Vitals - LCP optimization)
const Login = () => import('../pages/Login.vue')
const AgentLogin = () => import('../pages/AgentLogin.vue')
const Dashboard = () => import(/* webpackPrefetch: true */ '../pages/Dashboard.vue')
const Tasks = () => import(/* webpackPrefetch: true */ '../pages/Tasks.vue')
const Agents = () => import('../pages/Agents.vue')
const Projects = () => import(/* webpackPrefetch: true */ '../pages/Projects.vue')
const TaskDetails = () => import('../pages/TaskDetails.vue')
const ProjectDetails = () => import('../pages/ProjectDetails.vue')

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
