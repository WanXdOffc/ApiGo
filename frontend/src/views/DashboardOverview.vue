<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import axiosInstance from '../api/axios'
import { 
  KeyRound, 
  Activity, 
  Zap, 
  ShieldAlert, 
  ArrowUpRight, 
  Terminal, 
  Copy, 
  Check 
} from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()

const keyCount = ref<number | null>(null)
const logCount = ref<number | null>(null)
const loading = ref(true)
const copied = ref(false)

const user = computed(() => authStore.user)
const tier = computed(() => authStore.userTier)
const rateLimit = computed(() => tier.value === 'Premium' ? '1,000 req / min' : '10 req / min')

const sampleCurl = computed(() => {
  return `curl -X GET http://localhost:8080/api/v1/data \\\n  -H "x-api-key: YOUR_API_KEY"`
})

const copySnippet = async () => {
  try {
    await navigator.clipboard.writeText(sampleCurl.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {}
}

const loadStats = async () => {
  loading.value = true
  try {
    const [keysRes, logsRes] = await Promise.allSettled([
      axiosInstance.get('/api/keys'),
      axiosInstance.get('/api/logs'),
    ])

    if (keysRes.status === 'fulfilled') {
      keyCount.value = keysRes.value.data.total ?? keysRes.value.data.keys?.length ?? 0
    }
    if (logsRes.status === 'fulfilled') {
      logCount.value = logsRes.value.data.total ?? logsRes.value.data.logs?.length ?? 0
    }
  } catch {
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStats()
})
</script>

<template>
  <div class="space-y-8">
    <!-- Welcome Banner -->
    <div class="relative overflow-hidden rounded-2xl bg-gradient-to-r from-cyan-950/40 via-slate-900 to-indigo-950/30 border border-slate-800 p-6 sm:p-8 backdrop-blur-xl">
      <div class="relative z-10 max-w-2xl space-y-2">
        <div class="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-cyan-950/80 border border-cyan-800/40 text-[11px] font-mono text-cyan-400">
          <span>Active Plan: <strong>{{ tier }}</strong></span>
        </div>
        <h1 class="text-2xl sm:text-3xl font-bold tracking-tight text-white">
          Welcome back, {{ user?.name || 'Developer' }}
        </h1>
        <p class="text-sm text-slate-400 leading-relaxed">
          Monitor your developer API credentials, audit request latency, and test live REST endpoints across your distributed platform.
        </p>
      </div>
    </div>

    <!-- Telemetry Metric Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Card 1 -->
      <div class="bg-slate-900/60 border border-slate-800 rounded-2xl p-5 hover:border-slate-700 transition-all">
        <div class="flex items-center justify-between text-slate-400 mb-3">
          <span class="text-xs font-medium">Active API Keys</span>
          <div class="p-2 rounded-xl bg-cyan-500/10 text-cyan-400">
            <KeyRound class="w-4 h-4" />
          </div>
        </div>
        <div class="text-2xl font-bold text-white font-mono">
          {{ loading ? '...' : (keyCount ?? 0) }}
        </div>
        <p class="text-[11px] text-slate-500 mt-1">Managed credentials</p>
      </div>

      <!-- Card 2 -->
      <div class="bg-slate-900/60 border border-slate-800 rounded-2xl p-5 hover:border-slate-700 transition-all">
        <div class="flex items-center justify-between text-slate-400 mb-3">
          <span class="text-xs font-medium">Total Requests Logged</span>
          <div class="p-2 rounded-xl bg-indigo-500/10 text-indigo-400">
            <Activity class="w-4 h-4" />
          </div>
        </div>
        <div class="text-2xl font-bold text-white font-mono">
          {{ loading ? '...' : (logCount ?? 0) }}
        </div>
        <p class="text-[11px] text-slate-500 mt-1">Audit log records</p>
      </div>

      <!-- Card 3 -->
      <div class="bg-slate-900/60 border border-slate-800 rounded-2xl p-5 hover:border-slate-700 transition-all">
        <div class="flex items-center justify-between text-slate-400 mb-3">
          <span class="text-xs font-medium">Rate Limiting Window</span>
          <div class="p-2 rounded-xl bg-emerald-500/10 text-emerald-400">
            <Zap class="w-4 h-4" />
          </div>
        </div>
        <div class="text-base sm:text-lg font-bold text-white font-mono">
          {{ rateLimit }}
        </div>
        <p class="text-[11px] text-slate-500 mt-1">Sliding-window quota</p>
      </div>

      <!-- Card 4 -->
      <div class="bg-slate-900/60 border border-slate-800 rounded-2xl p-5 hover:border-slate-700 transition-all">
        <div class="flex items-center justify-between text-slate-400 mb-3">
          <span class="text-xs font-medium">Protection Status</span>
          <div class="p-2 rounded-xl bg-cyan-500/10 text-cyan-400">
            <ShieldAlert class="w-4 h-4" />
          </div>
        </div>
        <div class="text-base sm:text-lg font-bold text-emerald-400 font-mono">
          Anti-DDoS Active
        </div>
        <p class="text-[11px] text-slate-500 mt-1">Redis sliding window</p>
      </div>
    </div>

    <!-- Quick Navigation & Code Snippet -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Fast Actions -->
      <div class="space-y-4">
        <h2 class="text-sm font-semibold text-white">Quick Actions</h2>
        
        <div 
          @click="router.push('/dashboard/keys')"
          class="bg-slate-900/60 hover:bg-slate-800/80 border border-slate-800 rounded-2xl p-5 cursor-pointer transition-all flex items-center justify-between group"
        >
          <div class="flex items-center space-x-3.5">
            <div class="p-2.5 rounded-xl bg-cyan-500/10 text-cyan-400 group-hover:bg-cyan-500 group-hover:text-white transition-colors">
              <KeyRound class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-sm font-semibold text-white">Generate API Key</h3>
              <p class="text-xs text-slate-400">Issue a new SHA-256 protected key</p>
            </div>
          </div>
          <ArrowUpRight class="w-4 h-4 text-slate-500 group-hover:text-white transition-colors" />
        </div>

        <div 
          @click="router.push('/dashboard/logs')"
          class="bg-slate-900/60 hover:bg-slate-800/80 border border-slate-800 rounded-2xl p-5 cursor-pointer transition-all flex items-center justify-between group"
        >
          <div class="flex items-center space-x-3.5">
            <div class="p-2.5 rounded-xl bg-indigo-500/10 text-indigo-400 group-hover:bg-indigo-500 group-hover:text-white transition-colors">
              <Activity class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-sm font-semibold text-white">View Audit Trail</h3>
              <p class="text-xs text-slate-400">Inspect request latency & status</p>
            </div>
          </div>
          <ArrowUpRight class="w-4 h-4 text-slate-500 group-hover:text-white transition-colors" />
        </div>
      </div>

      <!-- Quickstart Code Snippet -->
      <div class="lg:col-span-2 bg-slate-900/70 border border-slate-800 rounded-2xl p-6 backdrop-blur-xl flex flex-col justify-between">
        <div>
          <div class="flex items-center justify-between pb-4 border-b border-slate-800">
            <div class="flex items-center space-x-2">
              <Terminal class="w-4 h-4 text-cyan-400" />
              <span class="text-xs font-semibold text-white">Live Endpoint Quickstart</span>
            </div>
            <button
              @click="copySnippet"
              class="flex items-center space-x-1.5 px-2.5 py-1 rounded-lg bg-slate-800 hover:bg-slate-700 text-xs text-slate-300 transition-colors cursor-pointer"
            >
              <Check v-if="copied" class="w-3.5 h-3.5 text-emerald-400" />
              <Copy v-else class="w-3.5 h-3.5" />
              <span>{{ copied ? 'Copied' : 'Copy cURL' }}</span>
            </button>
          </div>

          <p class="text-xs text-slate-400 my-4">
            Invoke the protected Consumer API endpoint using your generated API key in the <code>x-api-key</code> header:
          </p>

          <pre class="bg-slate-950 p-4 rounded-xl border border-slate-800/90 text-xs font-mono text-cyan-300 overflow-x-auto selection:bg-cyan-900 selection:text-white">{{ sampleCurl }}</pre>
        </div>

        <div class="mt-4 pt-4 border-t border-slate-800/80 flex items-center justify-between text-xs text-slate-500">
          <span>Target: <code>http://localhost:8080/api/v1/data</code></span>
          <span class="text-emerald-400">Response: 200 OK with Rate Limits</span>
        </div>
      </div>
    </div>
  </div>
</template>
