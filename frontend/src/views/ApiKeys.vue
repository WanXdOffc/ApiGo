<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axiosInstance from '../api/axios'
import { 
  KeyRound, 
  Plus, 
  Trash2, 
  Copy, 
  Check, 
  AlertTriangle, 
  X, 
  Loader2 
} from 'lucide-vue-next'

interface APIKeyItem {
  id: string
  name: string
  masked_key: string
  created_at: string
  is_active: boolean
}

const keys = ref<APIKeyItem[]>([])
const loading = ref(true)
const creating = ref(false)
const revokingId = ref<string | null>(null)
const errorMessage = ref('')

// Modal state for creating key
const showCreateModal = ref(false)
const newKeyName = ref('')

// Modal state for revealing raw key (CRITICAL ONE-TIME REVEAL)
const showRevealModal = ref(false)
const newlyGeneratedRawKey = ref('')
const newlyGeneratedKeyName = ref('')
const rawKeyCopied = ref(false)

const fetchKeys = async () => {
  loading.value = true
  errorMessage.value = ''
  try {
    const res = await axiosInstance.get('/api/keys')
    keys.value = res.data.keys || []
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Failed to fetch API keys'
  } finally {
    loading.value = false
  }
}

const openCreateModal = () => {
  newKeyName.value = ''
  showCreateModal.value = true
}

const handleCreateKey = async () => {
  if (!newKeyName.value.trim()) return
  creating.value = true
  errorMessage.value = ''

  try {
    const res = await axiosInstance.post('/api/keys', {
      name: newKeyName.value.trim(),
    })

    showCreateModal.value = false
    newlyGeneratedRawKey.value = res.data.raw_key
    newlyGeneratedKeyName.value = res.data.name || newKeyName.value
    rawKeyCopied.value = false
    showRevealModal.value = true

    // Refresh list
    await fetchKeys()
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Failed to create API key'
  } finally {
    creating.value = false
  }
}

const copyRawKey = async () => {
  try {
    await navigator.clipboard.writeText(newlyGeneratedRawKey.value)
    rawKeyCopied.value = true
    setTimeout(() => { rawKeyCopied.value = false }, 3000)
  } catch {}
}

const handleRevokeKey = async (id: string, name: string) => {
  if (!confirm(`Are you sure you want to revoke "${name}"? This action cannot be undone.`)) {
    return
  }

  revokingId.value = id
  try {
    await axiosInstance.delete(`/api/keys/${id}`)
    await fetchKeys()
  } catch (err: any) {
    alert(err.response?.data?.error || 'Failed to revoke API key')
  } finally {
    revokingId.value = null
  }
}

const formatDate = (isoString: string) => {
  if (!isoString) return '—'
  const d = new Date(isoString)
  return d.toLocaleDateString(undefined, { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric', 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

onMounted(() => {
  fetchKeys()
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header with action button -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-slate-800">
      <div>
        <h1 class="text-xl sm:text-2xl font-bold tracking-tight text-white flex items-center space-x-2.5">
          <KeyRound class="w-6 h-6 text-cyan-400" />
          <span>API Key Management</span>
        </h1>
        <p class="text-xs sm:text-sm text-slate-400 mt-1">
          Create, view, and revoke secret API keys used to authenticate against the Consumer REST API.
        </p>
      </div>

      <button
        @click="openCreateModal"
        class="inline-flex items-center space-x-2 px-4 py-2.5 bg-gradient-to-r from-cyan-600 to-indigo-600 hover:from-cyan-500 hover:to-indigo-500 text-white rounded-xl text-xs font-semibold shadow-lg shadow-cyan-950/50 transition-all cursor-pointer"
      >
        <Plus class="w-4 h-4" />
        <span>Generate New Key</span>
      </button>
    </div>

    <!-- Error message alert -->
    <div 
      v-if="errorMessage" 
      class="p-4 bg-rose-950/60 border border-rose-500/30 rounded-2xl text-rose-300 text-xs flex items-center justify-between"
    >
      <span>{{ errorMessage }}</span>
      <button @click="fetchKeys" class="underline hover:text-white ml-2">Retry</button>
    </div>

    <!-- Keys Table Card -->
    <div class="bg-slate-900/70 border border-slate-800 rounded-2xl backdrop-blur-xl overflow-hidden shadow-2xl">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-slate-800 bg-slate-950/60 text-[11px] font-mono uppercase tracking-wider text-slate-400">
              <th class="py-3.5 px-6 font-semibold">Key Name</th>
              <th class="py-3.5 px-6 font-semibold">Secret Identifier</th>
              <th class="py-3.5 px-6 font-semibold">Status</th>
              <th class="py-3.5 px-6 font-semibold">Created At</th>
              <th class="py-3.5 px-6 font-semibold text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/80 text-sm">
            <!-- Loading state -->
            <tr v-if="loading">
              <td colspan="5" class="py-12 text-center text-slate-500">
                <Loader2 class="w-6 h-6 animate-spin mx-auto text-cyan-400 mb-2" />
                <span class="text-xs">Loading active API keys...</span>
              </td>
            </tr>

            <!-- Empty state -->
            <tr v-else-if="keys.length === 0">
              <td colspan="5" class="py-12 text-center text-slate-500">
                <KeyRound class="w-8 h-8 mx-auto text-slate-600 mb-2" />
                <p class="text-sm font-medium text-slate-300">No API keys generated yet</p>
                <p class="text-xs text-slate-500 mt-1">Generate your first key to start consuming the API.</p>
                <button
                  @click="openCreateModal"
                  class="mt-4 px-3.5 py-1.5 bg-slate-800 hover:bg-slate-700 text-cyan-300 rounded-lg text-xs font-medium border border-slate-700 transition-colors"
                >
                  Create API Key
                </button>
              </td>
            </tr>

            <!-- Keys list rows -->
            <tr 
              v-else 
              v-for="item in keys" 
              :key="item.id"
              class="hover:bg-slate-800/30 transition-colors"
            >
              <!-- Name -->
              <td class="py-4 px-6 font-medium text-white">
                {{ item.name }}
              </td>

              <!-- Masked Key -->
              <td class="py-4 px-6 font-mono text-xs text-slate-300">
                <span class="bg-slate-950/80 px-2.5 py-1 rounded-lg border border-slate-800">
                  {{ item.masked_key }}
                </span>
              </td>

              <!-- Status -->
              <td class="py-4 px-6">
                <span 
                  class="inline-flex items-center space-x-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium border"
                  :class="item.is_active 
                    ? 'bg-emerald-950/50 text-emerald-300 border-emerald-500/30' 
                    : 'bg-rose-950/50 text-rose-300 border-rose-500/30'"
                >
                  <span 
                    class="w-1.5 h-1.5 rounded-full" 
                    :class="item.is_active ? 'bg-emerald-400' : 'bg-rose-400'"
                  ></span>
                  <span>{{ item.is_active ? 'Active' : 'Revoked' }}</span>
                </span>
              </td>

              <!-- Created At -->
              <td class="py-4 px-6 text-xs text-slate-400">
                {{ formatDate(item.created_at) }}
              </td>

              <!-- Actions -->
              <td class="py-4 px-6 text-right">
                <button
                  v-if="item.is_active"
                  @click="handleRevokeKey(item.id, item.name)"
                  :disabled="revokingId === item.id"
                  class="p-2 text-slate-400 hover:text-rose-400 hover:bg-rose-950/30 rounded-lg transition-colors cursor-pointer border border-transparent hover:border-rose-900/40"
                  title="Revoke Key"
                >
                  <Loader2 v-if="revokingId === item.id" class="w-4 h-4 animate-spin text-rose-400" />
                  <Trash2 v-else class="w-4 h-4" />
                </button>
                <span v-else class="text-xs text-slate-600 italic">Revoked</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal 1: Generate Key Modal -->
    <div 
      v-if="showCreateModal" 
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
    >
      <div class="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-5">
        <div class="flex items-center justify-between pb-3 border-b border-slate-800">
          <h3 class="text-base font-semibold text-white flex items-center space-x-2">
            <KeyRound class="w-5 h-5 text-cyan-400" />
            <span>Generate New API Key</span>
          </h3>
          <button 
            @click="showCreateModal = false" 
            class="text-slate-500 hover:text-white"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="handleCreateKey" class="space-y-4">
          <div class="space-y-1.5">
            <label class="text-xs font-medium text-slate-300">Key Name / Description</label>
            <input
              v-model="newKeyName"
              type="text"
              required
              placeholder="e.g., CI/CD Integration Key, Production Backend"
              class="w-full bg-slate-950 border border-slate-800 rounded-xl py-2.5 px-3.5 text-sm text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 focus:ring-1 focus:ring-cyan-500"
            />
          </div>

          <p class="text-xs text-slate-400 leading-relaxed">
            The key will be generated with 256-bit entropy. Only a SHA-256 hash will be stored on our servers.
          </p>

          <div class="flex items-center justify-end space-x-3 pt-3">
            <button
              type="button"
              @click="showCreateModal = false"
              class="px-4 py-2 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl text-xs font-medium transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="creating || !newKeyName.trim()"
              class="px-4 py-2 bg-gradient-to-r from-cyan-600 to-indigo-600 hover:from-cyan-500 hover:to-indigo-500 text-white rounded-xl text-xs font-semibold shadow-lg shadow-cyan-950/40 flex items-center space-x-2 disabled:opacity-50"
            >
              <Loader2 v-if="creating" class="w-4 h-4 animate-spin" />
              <span>{{ creating ? 'Generating...' : 'Create Key' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Modal 2: CRITICAL ONE-TIME RAW KEY REVEAL MODAL -->
    <div 
      v-if="showRevealModal" 
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 backdrop-blur-md"
    >
      <div class="bg-slate-900 border border-amber-500/40 rounded-2xl max-w-lg w-full p-6 sm:p-8 shadow-2xl space-y-6">
        <!-- Warning alert banner -->
        <div class="flex items-start space-x-3 p-4 bg-amber-950/50 border border-amber-500/40 rounded-xl text-amber-300 text-xs">
          <AlertTriangle class="w-5 h-5 shrink-0 text-amber-400" />
          <div class="space-y-1">
            <p class="font-bold">Save your secret API key now!</p>
            <p class="text-amber-200/80 leading-relaxed">
              This is the <strong>only time</strong> your raw API key will ever be displayed. If you lose it, you will need to revoke it and generate a new one.
            </p>
          </div>
        </div>

        <div>
          <span class="text-xs text-slate-400 block mb-1">Key Name:</span>
          <span class="text-sm font-semibold text-white">{{ newlyGeneratedKeyName }}</span>
        </div>

        <!-- Raw Key display box with Copy button -->
        <div class="space-y-2">
          <div class="flex items-center justify-between text-xs text-slate-400">
            <span>Raw Secret API Key:</span>
            <span class="font-mono text-[11px] text-cyan-400">SHA-256 Hashed on Server</span>
          </div>
          <div class="relative flex items-center bg-slate-950 border border-slate-800 rounded-xl p-3 font-mono text-xs text-cyan-300 break-all select-all">
            <span>{{ newlyGeneratedRawKey }}</span>
          </div>
        </div>

        <!-- Action buttons -->
        <div class="flex flex-col sm:flex-row items-center gap-3 pt-2">
          <button
            @click="copyRawKey"
            class="w-full sm:flex-1 py-2.5 px-4 bg-cyan-600 hover:bg-cyan-500 text-white rounded-xl text-xs font-semibold flex items-center justify-center space-x-2 transition-all cursor-pointer shadow-lg shadow-cyan-950/50"
          >
            <Check v-if="rawKeyCopied" class="w-4 h-4 text-white" />
            <Copy v-else class="w-4 h-4" />
            <span>{{ rawKeyCopied ? 'Copied to Clipboard!' : 'Copy API Key' }}</span>
          </button>

          <button
            @click="showRevealModal = false"
            class="w-full sm:w-auto py-2.5 px-5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded-xl text-xs font-medium transition-colors cursor-pointer"
          >
            I've Saved It
          </button>
        </div>
      </div>
    </div>

  </div>
</template>
