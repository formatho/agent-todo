import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import Login from '../pages/Login.vue'
import Register from '../pages/Register.vue'
import Dashboard from '../pages/Dashboard.vue'
import Agents from '../pages/Agents.vue'
import TaskDetails from '../pages/TaskDetails.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/register',
    name: 'Register',
    component: Register
  },
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/agents',
    name: 'Agents',
    component: Agents,
    meta: { requiresAuth: true }
  },
  {
    path: '/tasks/:id',
    name: 'TaskDetails',
    component: TaskDetails,
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if ((to.name === 'Login' || to.name === 'Register') && authStore.isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router
