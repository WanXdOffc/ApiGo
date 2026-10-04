<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axiosInstance from '../api/axios'
import { 
  FileText, 
  RefreshCw, 
  Terminal, 
  Loader2 
} from 'lucide-vue-next'

interface AuditLogItem {
  id: string
  api_key_id: string
  endpoint: string
  method: string
  status_code: number
  latency: number
  timestamp: string
}

const logs = ref<AuditLogItem[]>([])
const loading = ref(true)
const refreshing = ref(false)
const errorMessage = ref('')

const fetchLogs = async (isManualRefresh = false) => {
  if (isManualRefresh) refreshing.value = true
  else loading.value = true
  errorMessage.value = ''

  try {
    const res = await axiosInstance.get('/api/logs')
    logs.value = res.data.logs || []
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Failed to fetch audit logs'
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

const formatTimestamp = (isoString: string) => {
  if (!isoString) return '—'
  const d = new Date(isoString)
  return d.toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  })
}

const getMethodColor = (method: string) => {
  switch (method.toUpperCase()) {
    case 'GET':
      return 'bg-cyan-950/60 text-cyan-300 border-cyan-800/40'
    case 'POST':
      return 'bg-emerald-950/60 text-emerald-300 border-emerald-800/40'
    case 'DELETE':
      return 'bg-rose-950/60 text-rose-300 border-rose-800/40'
    case 'PUT':
    case 'PATCH':
      return 'bg-amber-950/60 text-amber-300 border-amber-800/40'
    default:
      return 'bg-slate-900 text-slate-300 border-slate-700'
  }
}

const getStatusColor = (code: number) => {
  if (code >= 200 && code < 300) {
    return 'bg-emerald-950/50 text-emerald-300 border-emerald-500/30'
  } else if (code === 429) {
    return 'bg-purple-950/50 text-purple-300 border-purple-500/30'
  } else if (code >= 400 && code < 500) {
    return 'bg-amber-950/50 text-amber-300 border-amber-500/30'
  } else {
    return 'bg-rose-950/50 text-rose-300 border-rose-500/30'
  }
}

const getLatencyColor = (latency: number) => {
  if (latency < 100) return 'text-emerald-400'
  if (latency < 300) return 'text-cyan-400'
  if (latency < 1000) return 'text-amber-400'
  return 'text-rose-400'
}

onMounted(() => {
  fetchLogs()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header with refresh button -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-slate-800">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-white flex items-center space-x-2.5">
          <FileText class="w-6 h-6 text-indigo-400" />
          <span>Real-Time Audit Trail</span>
        </h1>
        <p class="text-xs sm:text-sm text-slate-400 mt-1">
          Historical log of every authenticated inbound request processed through your API credentials.
        </p>
      </div>

      <button
        @click="() => fetchLogs(true)"
        :disabled="loading || refreshing"
        class="inline-flex items-center space-x-2 px-3.5 py-2 bg-slate-900 hover:bg-slate-800 border border-slate-800 text-slate-300 hover:text-white rounded-xl text-xs font-medium transition-all cursor-pointer disabled:opacity-50"
      >
        <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': refreshing }" />
        <span>{{ refreshing ? 'Refreshing...' : 'Refresh Logs' }}</span>
      </button>
    </div>

    <!-- Error alert -->
    <div 
      v-if="errorMessage" 
      class="p-4 bg-rose-950/60 border border-rose-500/30 rounded-2xl text-rose-300 text-xs flex items-center justify-between"
    >
      <span>{{ errorMessage }}</span>
      <button @click="() => fetchLogs(true)" class="underline hover:text-white ml-2">Retry</button>
    </div>

    <!-- Logs Table Card -->
    <div class="bg-slate-900/70 border border-slate-800 rounded-2xl backdrop-blur-xl overflow-hidden shadow-2xl">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-slate-800 bg-slate-950/60 text-[11px] font-mono uppercase tracking-wider text-slate-400">
              <th class="py-3.5 px-6 font-semibold">Method</th>
              <th class="py-3.5 px-6 font-semibold">Endpoint</th>
              <th class="py-3.5 px-6 font-semibold">Status Code</th>
              <th class="py-3.5 px-6 font-semibold">Latency</th>
              <th class="py-3.5 px-6 font-semibold text-right">Timestamp</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/80 text-sm font-mono">
            <!-- Loading state -->
            <tr v-if="loading">
              <td colspan="5" class="py-12 text-center text-slate-500 font-sans">
                <Loader2 class="w-6 h-6 animate-spin mx-auto text-indigo-400 mb-2" />
                <span class="text-xs">Fetching audit logs from database...</span>
              </td>
            </tr>

            <!-- Empty state -->
            <tr v-else-if="logs.length === 0">
              <td colspan="5" class="py-12 text-center text-slate-500 font-sans">
                <Terminal class="w-8 h-8 mx-auto text-slate-600 mb-2" />
                <p class="text-sm font-medium text-slate-300">No requests recorded yet</p>
                <p class="text-xs text-slate-500 mt-1 max-w-sm mx-auto">
                  Make a request to <code>GET http://localhost:8080/api/v1/data</code> using your API key to see telemetry appear here in real-time.
                </p>
              </td>
            </tr>

            <!-- Log rows -->
            <tr 
              v-else 
              v-for="log in logs" 
              :key="log.id"
              class="hover:bg-slate-800/30 transition-colors"
            >
              <!-- HTTP Method -->
              <td class="py-3.5 px-6">
                <span 
                  class="px-2.5 py-0.5 rounded-md text-[11px] font-semibold border inline-block"
                  :class="getMethodColor(log.method)"
                >
                  {{ log.method }}
                </span>
              </td>

              <!-- Endpoint -->
              <td class="py-3.5 px-6 text-xs text-white font-medium">
                {{ log.endpoint }}
              </td>

              <!-- Status Code -->
              <td class="py-3.5 px-6">
                <span 
                  class="px-2.5 py-0.5 rounded-full text-xs font-semibold border inline-flex items-center space-x-1"
                  :class="getStatusColor(log.status_code)"
                >
                  <span>{{ log.status_code }}</span>
                </span>
              </td>

              <!-- Latency -->
              <td class="py-3.5 px-6 text-xs font-medium" :class="getLatencyColor(log.latency)">
                {{ log.latency }} ms
              </td>

              <!-- Timestamp -->
              <td class="py-3.5 px-6 text-right text-xs text-slate-400 font-sans">
                {{ formatTimestamp(log.timestamp) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
