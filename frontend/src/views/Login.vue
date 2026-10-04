<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import axiosInstance from '../api/axios'
import { KeyRound, Mail, Lock, User as UserIcon, ArrowRight, Loader2, Sparkles } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const isRegister = ref(false)
const name = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMessage = ref('')

const toggleMode = () => {
  isRegister.value = !isRegister.value
  errorMessage.value = ''
}

const handleSubmit = async () => {
  errorMessage.value = ''
  loading.value = true

  try {
    if (isRegister.value) {
      const res = await axiosInstance.post('/api/auth/register', {
        name: name.value,
        email: email.value,
        password: password.value,
      })
      authStore.setAuth(res.data.token, res.data.user)
    } else {
      const res = await axiosInstance.post('/api/auth/login', {
        email: email.value,
        password: password.value,
      })
      authStore.setAuth(res.data.token, res.data.user)
    }

    const redirectPath = (route.query.redirect as string) || '/dashboard'
    router.push(redirectPath)
  } catch (err: any) {
    if (err.response?.data?.error) {
      errorMessage.value = err.response.data.error
    } else {
      errorMessage.value = 'Failed to connect to authentication service. Please verify the backend is running.'
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-[#070b14] text-slate-100 flex items-center justify-center p-4 relative overflow-hidden font-sans">
    <!-- Ambient glowing backgrounds -->
    <div class="absolute -top-40 -left-40 w-96 h-96 bg-cyan-500/10 rounded-full blur-[128px] pointer-events-none"></div>
    <div class="absolute -bottom-40 -right-40 w-96 h-96 bg-indigo-500/15 rounded-full blur-[128px] pointer-events-none"></div>
    <div class="absolute inset-0 bg-[radial-gradient(#1e293b_1px,transparent_1px)] [background-size:24px_24px] opacity-25 pointer-events-none"></div>

    <div class="w-full max-w-md relative z-10">
      <!-- Brand header -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center h-12 w-12 rounded-2xl bg-gradient-to-tr from-cyan-500 to-indigo-600 shadow-xl shadow-cyan-500/25 ring-1 ring-white/20 mb-4">
          <KeyRound class="w-6 h-6 text-white" />
        </div>
        <h1 class="text-2xl font-bold tracking-tight text-white">
          PulseAPI Developer Portal
        </h1>
        <p class="text-sm text-slate-400 mt-1">
          {{ isRegister ? 'Create an account to manage API keys' : 'Sign in to access your developer console' }}
        </p>
      </div>

      <!-- Auth card -->
      <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-6 sm:p-8 backdrop-blur-xl shadow-2xl space-y-6">
        <!-- Error Banner -->
        <div 
          v-if="errorMessage" 
          class="p-3.5 bg-rose-950/60 border border-rose-500/30 rounded-xl text-rose-300 text-xs flex items-start space-x-2"
        >
          <span class="font-bold">Error:</span>
          <span>{{ errorMessage }}</span>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <!-- Name field (register only) -->
          <div v-if="isRegister" class="space-y-1.5">
            <label class="text-xs font-medium text-slate-300">Full Name</label>
            <div class="relative">
              <UserIcon class="w-4 h-4 text-slate-500 absolute left-3.5 top-1/2 -translate-y-1/2" />
              <input
                v-model="name"
                type="text"
                required
                placeholder="Alex Developer"
                class="w-full bg-slate-950/70 border border-slate-800 rounded-xl py-2.5 pl-10 pr-3.5 text-sm text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
              />
            </div>
          </div>

          <!-- Email field -->
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-slate-300">Email Address</label>
            <div class="relative">
              <Mail class="w-4 h-4 text-slate-500 absolute left-3.5 top-1/2 -translate-y-1/2" />
              <input
                v-model="email"
                type="email"
                required
                placeholder="developer@example.com"
                class="w-full bg-slate-950/70 border border-slate-800 rounded-xl py-2.5 pl-10 pr-3.5 text-sm text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
              />
            </div>
          </div>

          <!-- Password field -->
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-slate-300">Password</label>
            <div class="relative">
              <Lock class="w-4 h-4 text-slate-500 absolute left-3.5 top-1/2 -translate-y-1/2" />
              <input
                v-model="password"
                type="password"
                required
                minlength="6"
                placeholder="••••••••"
                class="w-full bg-slate-950/70 border border-slate-800 rounded-xl py-2.5 pl-10 pr-3.5 text-sm text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500 transition-colors"
              />
            </div>
            <p v-if="isRegister" class="text-[11px] text-slate-500">Minimum 6 characters</p>
          </div>

          <!-- Submit button -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full py-2.5 px-4 bg-gradient-to-r from-cyan-600 to-indigo-600 hover:from-cyan-500 hover:to-indigo-500 text-white rounded-xl text-sm font-semibold shadow-lg shadow-cyan-950/50 flex items-center justify-center space-x-2 transition-all disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer pt-3 pb-3"
          >
            <Loader2 v-if="loading" class="w-4 h-4 animate-spin" />
            <template v-else>
              <span>{{ isRegister ? 'Create Account' : 'Sign In' }}</span>
              <ArrowRight class="w-4 h-4" />
            </template>
          </button>
        </form>

        <!-- Mode Toggle footer -->
        <div class="pt-4 border-t border-slate-800/80 text-center">
          <p class="text-xs text-slate-400">
            {{ isRegister ? 'Already have an account?' : "Don't have an account yet?" }}
            <button
              type="button"
              @click="toggleMode"
              class="text-cyan-400 hover:text-cyan-300 font-medium ml-1 cursor-pointer focus:outline-none"
            >
              {{ isRegister ? 'Sign In instead' : 'Create Free Account' }}
            </button>
          </p>
        </div>
      </div>

      <!-- Tier feature notice -->
      <div class="mt-6 flex items-center justify-center space-x-2 text-xs text-slate-500">
        <Sparkles class="w-3.5 h-3.5 text-cyan-400" />
        <span>Free Tier includes 10 req/min &bull; Instant API Key Generation</span>
      </div>
    </div>
  </div>
</template>
