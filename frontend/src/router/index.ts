import { createRouter, createWebHistory } from 'vue-router'
import DashboardLayout from '../layouts/DashboardLayout.vue'

const routes = [
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { guestOnly: true },
  },
  {
    path: '/dashboard',
    component: DashboardLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'DashboardOverview',
        component: () => import('../views/DashboardOverview.vue'),
      },
      {
        path: 'keys',
        name: 'ApiKeys',
        component: () => import('../views/ApiKeys.vue'),
      },
      {
        path: 'logs',
        name: 'AuditLogs',
        component: () => import('../views/AuditLogs.vue'),
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation Guard: Protect private dashboard routes and redirect guests
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')

  if (to.matched.some((record) => record.meta.requiresAuth)) {
    if (!token) {
      return next({ path: '/login', query: { redirect: to.fullPath } })
    }
  }

  if (to.matched.some((record) => record.meta.guestOnly)) {
    if (token) {
      return next({ path: '/dashboard' })
    }
  }

  next()
})

export default router
