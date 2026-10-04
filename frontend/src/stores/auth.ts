import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axiosInstance from '../api/axios'

export interface User {
  id: string
  name: string
  email: string
  tier: string
  created_at?: string
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  
  // Safely initialize cached user profile
  const storedUser = localStorage.getItem('user')
  const user = ref<User | null>(storedUser ? JSON.parse(storedUser) : null)

  const isAuthenticated = computed(() => !!token.value)
  const userTier = computed(() => user.value?.tier || 'Free')

  function setAuth(newToken: string, newUser: User) {
    token.value = newToken
    user.value = newUser
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  async function fetchCurrentUser() {
    if (!token.value) return null
    try {
      const res = await axiosInstance.get('/api/auth/me')
      if (res.data && res.data.user_id) {
        if (user.value) {
          user.value.tier = res.data.tier || user.value.tier
        }
      }
      return res.data
    } catch {
      // If token expired or invalid, log out
      logout()
      return null
    }
  }

  return {
    token,
    user,
    isAuthenticated,
    userTier,
    setAuth,
    logout,
    fetchCurrentUser,
  }
})
