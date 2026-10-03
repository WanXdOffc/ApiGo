<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface HealthResponse {
  status: string
  [key: string]: unknown
}

const status = ref<string | null>(null)
const responseData = ref<HealthResponse | null>(null)
const loading = ref<boolean>(true)
const error = ref<string | null>(null)
const latency = ref<number | null>(null)
const lastChecked = ref<string | null>(null)

const API_BASE_URL = 'http://localhost:8080'

const fetchHealth = async () => {
  loading.value = true
  error.value = null
  const startTime = performance.now()

  try {
    const res = await fetch(`${API_BASE_URL}/api/health`, {
      method: 'GET',
      headers: {
        'Accept': 'application/json',
      },
    })

    const duration = Math.round(performance.now() - startTime)
    latency.value = duration

    if (!res.ok) {
      throw new Error(`HTTP ${res.status}: ${res.statusText}`)
    }

    const data = await res.json()
    status.value = data.status || 'unknown'
    responseData.value = data
    lastChecked.value = new Date().toLocaleTimeString()
  } catch (err: unknown) {
    status.value = 'offline'
    if (err instanceof Error) {
      error.value = err.message
    } else {
      error.value = 'Failed to connect to backend server'
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchHealth()
})
</script>

<template>
  <div class="min-h-screen bg-[#070b14] text-slate-100 flex flex-col font-sans selection:bg-cyan-500/20 selection:text-cyan-300">
    <!-- Background glowing accents -->
    <div class="fixed inset-0 overflow-hidden pointer-events-none -z-10">
      <div class="absolute -top-40 -left-40 w-96 h-96 bg-cyan-500/10 rounded-full blur-[128px]"></div>
      <div class="absolute top-1/3 -right-40 w-96 h-96 bg-indigo-500/10 rounded-full blur-[128px]"></div>
      <div class="absolute -bottom-40 left-1/3 w-96 h-96 bg-emerald-500/10 rounded-full blur-[128px]"></div>
      <div class="absolute inset-0 bg-[radial-gradient(#1e293b_1px,transparent_1px)] [background-size:24px_24px] opacity-30"></div>
    </div>

    <!-- Navigation Header -->
    <header class="border-b border-slate-800/80 bg-slate-950/60 backdrop-blur-md sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="h-9 w-9 rounded-xl bg-gradient-to-tr from-cyan-500 to-indigo-600 flex items-center justify-center shadow-lg shadow-cyan-500/20 ring-1 ring-white/20">
            <svg class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div>
            <div class="flex items-center space-x-2">
              <span class="font-bold text-lg tracking-tight bg-gradient-to-r from-white via-slate-100 to-slate-400 bg-clip-text text-transparent">
                PulseAPI Platform
              </span>
              <span class="text-[10px] font-mono uppercase px-2 py-0.5 rounded-full bg-cyan-950/80 text-cyan-400 border border-cyan-700/50">
                Phase 1 Active
              </span>
            </div>
            <p class="text-xs text-slate-400 hidden sm:block">B2D High-Performance API SaaS Platform</p>
          </div>
        </div>

        <div class="flex items-center space-x-3">
          <!-- Live Service Status Pill in Header -->
          <div 
            class="flex items-center space-x-2 px-3 py-1.5 rounded-full text-xs font-medium border backdrop-blur-md transition-all duration-300"
            :class="status === 'ok' 
              ? 'bg-emerald-950/40 text-emerald-300 border-emerald-500/30' 
              : loading 
              ? 'bg-amber-950/40 text-amber-300 border-amber-500/30' 
              : 'bg-rose-950/40 text-rose-300 border-rose-500/30'"
          >
            <span class="relative flex h-2 w-2">
              <span 
                v-if="status === 'ok'"
                class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"
              ></span>
              <span 
                class="relative inline-flex rounded-full h-2 w-2"
                :class="status === 'ok' ? 'bg-emerald-400' : loading ? 'bg-amber-400' : 'bg-rose-400'"
              ></span>
            </span>
            <span>{{ loading ? 'Checking API...' : status === 'ok' ? 'API Online' : 'API Offline' }}</span>
          </div>

          <button
            @click="fetchHealth"
            :disabled="loading"
            class="p-2 text-slate-400 hover:text-white bg-slate-900/60 hover:bg-slate-800/80 border border-slate-800 rounded-lg transition-all focus:outline-none focus:ring-2 focus:ring-cyan-500/50"
            title="Refresh Health Status"
          >
            <svg 
              class="w-4 h-4 transition-transform duration-500" 
              :class="{ 'animate-spin': loading }" 
              fill="none" 
              viewBox="0 0 24 24" 
              stroke="currentColor"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content Area -->
    <main class="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-10 space-y-10">
      
      <!-- Hero Banner -->
      <section class="text-center max-w-3xl mx-auto space-y-4 pt-4">
        <div class="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-slate-900/80 border border-slate-800 text-xs text-slate-300 backdrop-blur-sm shadow-inner">
          <span class="w-1.5 h-1.5 rounded-full bg-cyan-400"></span>
          <span>Full-Stack Monorepo Initialized: <strong>Vue 3 (Vite)</strong> + <strong>Golang (Fiber)</strong></span>
        </div>
        <h1 class="text-3xl sm:text-5xl font-extrabold tracking-tight text-white leading-tight">
          Next-Gen <span class="bg-gradient-to-r from-cyan-400 via-sky-300 to-indigo-400 bg-clip-text text-transparent">B2D API Platform</span>
        </h1>
        <p class="text-sm sm:text-base text-slate-400 max-w-2xl mx-auto">
          Monorepo development environment running concurrently. The frontend communicates directly with the high-performance Go Fiber Core API.
        </p>
      </section>

      <!-- Health Verification & Live Telemetry Card -->
      <section class="max-w-2xl mx-auto">
        <div class="relative group">
          <!-- Subtle glow ring -->
          <div 
            class="absolute -inset-0.5 rounded-2xl blur-lg transition duration-500 opacity-70 group-hover:opacity-100"
            :class="status === 'ok' ? 'bg-gradient-to-r from-emerald-500/20 to-cyan-500/20' : 'bg-gradient-to-r from-rose-500/20 to-amber-500/20'"
          ></div>

          <div class="relative rounded-2xl bg-slate-900/90 border border-slate-800/90 p-6 sm:p-8 backdrop-blur-xl shadow-2xl">
            <div class="flex items-center justify-between pb-5 border-b border-slate-800">
              <div class="flex items-center space-x-3">
                <div 
                  class="w-10 h-10 rounded-xl flex items-center justify-center border"
                  :class="status === 'ok' 
                    ? 'bg-emerald-500/10 border-emerald-500/30 text-emerald-400' 
                    : 'bg-rose-500/10 border-rose-500/30 text-rose-400'"
                >
                  <svg v-if="status === 'ok'" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                </div>
                <div>
                  <h2 class="text-lg font-semibold text-white">Backend Health Check</h2>
                  <p class="text-xs text-slate-400 font-mono">GET /api/health</p>
                </div>
              </div>

              <button
                @click="fetchHealth"
                :disabled="loading"
                class="px-4 py-2 bg-gradient-to-r from-cyan-600 to-indigo-600 hover:from-cyan-500 hover:to-indigo-500 text-white rounded-lg text-xs font-medium shadow-md shadow-cyan-900/30 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
              >
                <svg 
                  class="w-3.5 h-3.5" 
                  :class="{ 'animate-spin': loading }" 
                  fill="none" 
                  viewBox="0 0 24 24" 
                  stroke="currentColor"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                <span>{{ loading ? 'Testing...' : 'Test Endpoint' }}</span>
              </button>
            </div>

            <!-- Stats grid -->
            <div class="grid grid-cols-2 sm:grid-cols-3 gap-4 my-6">
              <div class="bg-slate-950/60 border border-slate-800/80 rounded-xl p-3.5">
                <span class="text-[11px] text-slate-400 block mb-1">Health Status</span>
                <span 
                  class="text-sm font-semibold tracking-wide uppercase font-mono"
                  :class="status === 'ok' ? 'text-emerald-400' : 'text-rose-400'"
                >
                  {{ status || 'CHECKING...' }}
                </span>
              </div>
              <div class="bg-slate-950/60 border border-slate-800/80 rounded-xl p-3.5">
                <span class="text-[11px] text-slate-400 block mb-1">Response Latency</span>
                <span class="text-sm font-mono font-medium text-slate-200">
                  {{ latency !== null ? `${latency} ms` : '—' }}
                </span>
              </div>
              <div class="col-span-2 sm:col-span-1 bg-slate-950/60 border border-slate-800/80 rounded-xl p-3.5">
                <span class="text-[11px] text-slate-400 block mb-1">Last Timestamp</span>
                <span class="text-sm font-mono text-slate-300">
                  {{ lastChecked || 'Never' }}
                </span>
              </div>
            </div>

            <!-- Response Raw Payload Preview -->
            <div class="space-y-2">
              <div class="flex items-center justify-between text-xs text-slate-400 font-mono">
                <span>Response Payload:</span>
                <span class="text-[10px] text-slate-500">application/json</span>
              </div>
              <div class="bg-slate-950 rounded-xl p-4 border border-slate-800 font-mono text-xs overflow-x-auto text-emerald-400">
                <pre v-if="responseData">{{ JSON.stringify(responseData, null, 2) }}</pre>
                <div v-else-if="loading" class="text-slate-500 italic">Fetching from {{ API_BASE_URL }}/api/health...</div>
                <div v-else-if="error" class="text-rose-400 font-sans text-xs">
                  <p class="font-semibold mb-1">Connection Error:</p>
                  <p class="text-rose-300/80 font-mono">{{ error }}</p>
                  <p class="mt-2 text-slate-400">Verify that the Go Fiber backend is running on <code>http://localhost:8080</code>.</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Monorepo Architecture Overview Cards -->
      <section class="space-y-4 pt-4">
        <h2 class="text-center text-xs font-mono uppercase tracking-widest text-slate-400">
          Architecture & Roadmap (Phase 1 Completed)
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
          <!-- Card 1: Backend Stack -->
          <div class="bg-slate-900/60 border border-slate-800/90 rounded-2xl p-5 hover:border-slate-700/80 transition-all">
            <div class="flex items-center justify-between mb-3">
              <div class="h-8 w-8 rounded-lg bg-cyan-500/10 text-cyan-400 flex items-center justify-center font-mono font-bold text-xs">
                GO
              </div>
              <span class="text-[10px] font-mono px-2 py-0.5 rounded-full bg-cyan-950 text-cyan-400 border border-cyan-800/40">
                Port :8080
              </span>
            </div>
            <h3 class="text-base font-semibold text-white mb-1">Golang Fiber Backend</h3>
            <p class="text-xs text-slate-400 leading-relaxed mb-4">
              High-throughput HTTP engine with CORS middleware, recovery handlers, and Air hot-reloading.
            </p>
            <div class="text-[11px] font-mono text-slate-400 bg-slate-950/80 rounded-lg p-2.5 border border-slate-800/80">
              <span class="text-cyan-400">module:</span> api-platform<br />
              <span class="text-cyan-400">framework:</span> fiber/v2<br />
              <span class="text-cyan-400">hot-reload:</span> .air.toml
            </div>
          </div>

          <!-- Card 2: Frontend Stack -->
          <div class="bg-slate-900/60 border border-slate-800/90 rounded-2xl p-5 hover:border-slate-700/80 transition-all">
            <div class="flex items-center justify-between mb-3">
              <div class="h-8 w-8 rounded-lg bg-emerald-500/10 text-emerald-400 flex items-center justify-center font-mono font-bold text-xs">
                VUE
              </div>
              <span class="text-[10px] font-mono px-2 py-0.5 rounded-full bg-emerald-950 text-emerald-400 border border-emerald-800/40">
                Port :5173
              </span>
            </div>
            <h3 class="text-base font-semibold text-white mb-1">Vue 3 + Vite Frontend</h3>
            <p class="text-xs text-slate-400 leading-relaxed mb-4">
              TypeScript-powered developer portal styled with Tailwind CSS, utilizing Composition API.
            </p>
            <div class="text-[11px] font-mono text-slate-400 bg-slate-950/80 rounded-lg p-2.5 border border-slate-800/80">
              <span class="text-emerald-400">tooling:</span> Vite 6 / TS<br />
              <span class="text-emerald-400">styling:</span> Tailwind CSS<br />
              <span class="text-emerald-400">runtime:</span> Vue 3 Composition
            </div>
          </div>

          <!-- Card 3: Monorepo Orchestration -->
          <div class="bg-slate-900/60 border border-slate-800/90 rounded-2xl p-5 hover:border-slate-700/80 transition-all">
            <div class="flex items-center justify-between mb-3">
              <div class="h-8 w-8 rounded-lg bg-indigo-500/10 text-indigo-400 flex items-center justify-center font-mono font-bold text-xs">
                DEV
              </div>
              <span class="text-[10px] font-mono px-2 py-0.5 rounded-full bg-indigo-950 text-indigo-400 border border-indigo-800/40">
                make dev
              </span>
            </div>
            <h3 class="text-base font-semibold text-white mb-1">Unified Monorepo Workflow</h3>
            <p class="text-xs text-slate-400 leading-relaxed mb-4">
              Root Makefile orchestrating both servers concurrently with hot-reloading across platforms.
            </p>
            <div class="text-[11px] font-mono text-slate-400 bg-slate-950/80 rounded-lg p-2.5 border border-slate-800/80">
              <span class="text-indigo-400">command:</span> make dev<br />
              <span class="text-indigo-400">concurrent:</span> yes<br />
              <span class="text-indigo-400">ready for:</span> Phase 2 (MongoDB)
            </div>
          </div>
        </div>
      </section>

    </main>

    <!-- Footer -->
    <footer class="border-t border-slate-800/80 bg-slate-950/50 py-6 text-center text-xs text-slate-500">
      <div class="max-w-7xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-3">
        <span>B2D SaaS API Platform &mdash; Phase 1 (Monorepo Scaffolding)</span>
        <span class="font-mono text-[11px] text-slate-400">Frontend: :5173 &bull; Backend: :8080</span>
      </div>
    </footer>
  </div>
</template>
