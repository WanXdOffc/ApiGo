<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { 
  LayoutDashboard, 
  KeyRound, 
  FileText, 
  LogOut, 
  Menu, 
  X, 
  Zap, 
  ExternalLink,
  ShieldCheck
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const mobileMenuOpen = ref(false)

const user = computed(() => authStore.user)
const userTier = computed(() => authStore.userTier)

const navItems = [
  {
    name: 'Overview',
    path: '/dashboard',
    icon: LayoutDashboard,
    exact: true,
  },
  {
    name: 'API Keys',
    path: '/dashboard/keys',
    icon: KeyRound,
    exact: false,
  },
  {
    name: 'Audit Logs',
    path: '/dashboard/logs',
    icon: FileText,
    exact: false,
  },
]

const isActive = (item: typeof navItems[0]) => {
  if (item.exact) {
    return route.path === item.path
  }
  return route.path.startsWith(item.path)
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen bg-[#070b14] text-slate-100 flex font-sans selection:bg-cyan-500/20 selection:text-cyan-300">
    <!-- Mobile sidebar overlay -->
    <div 
      v-if="mobileMenuOpen" 
      @click="mobileMenuOpen = false" 
      class="fixed inset-0 bg-black/60 z-40 lg:hidden backdrop-blur-sm"
    ></div>

    <!-- Sidebar navigation -->
    <aside 
      class="fixed inset-y-0 left-0 z-50 w-64 bg-slate-950/80 border-r border-slate-800/80 backdrop-blur-xl flex flex-col transition-transform duration-300 lg:translate-x-0"
      :class="mobileMenuOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'"
    >
      <!-- Brand header -->
      <div class="h-16 px-6 flex items-center justify-between border-b border-slate-800/80">
        <router-link to="/dashboard" class="flex items-center space-x-3 group">
          <div class="h-9 w-9 rounded-xl bg-gradient-to-tr from-cyan-500 to-indigo-600 flex items-center justify-center shadow-lg shadow-cyan-500/20 ring-1 ring-white/20">
            <Zap class="w-5 h-5 text-white" />
          </div>
          <div>
            <span class="font-bold text-base tracking-tight bg-gradient-to-r from-white to-slate-300 bg-clip-text text-transparent">
              PulseAPI
            </span>
            <span class="block text-[10px] text-slate-500 font-mono">Developer SaaS</span>
          </div>
        </router-link>

        <button 
          @click="mobileMenuOpen = false" 
          class="lg:hidden p-1.5 text-slate-400 hover:text-white rounded-lg"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <!-- Navigation links -->
      <nav class="flex-1 px-4 py-6 space-y-1.5 overflow-y-auto">
        <div class="px-3 pb-2 text-[10px] font-mono uppercase tracking-wider text-slate-500 font-semibold">
          Platform Console
        </div>

        <router-link
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          @click="mobileMenuOpen = false"
          class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-all group"
          :class="isActive(item) 
            ? 'bg-gradient-to-r from-cyan-500/15 to-indigo-500/10 text-cyan-300 border border-cyan-500/30 shadow-sm shadow-cyan-950/40' 
            : 'text-slate-400 hover:text-slate-200 hover:bg-slate-900/60'"
        >
          <component 
            :is="item.icon" 
            class="w-4 h-4 transition-colors" 
            :class="isActive(item) ? 'text-cyan-400' : 'text-slate-500 group-hover:text-slate-300'"
          />
          <span>{{ item.name }}</span>
        </router-link>

        <div class="pt-6 px-3 pb-2 text-[10px] font-mono uppercase tracking-wider text-slate-500 font-semibold">
          Developer Resources
        </div>

        <a 
          href="http://localhost:8080/api/health" 
          target="_blank" 
          class="flex items-center justify-between px-3.5 py-2.5 rounded-xl text-sm font-medium text-slate-400 hover:text-slate-200 hover:bg-slate-900/60 transition-all group"
        >
          <div class="flex items-center space-x-3">
            <ShieldCheck class="w-4 h-4 text-slate-500 group-hover:text-slate-300" />
            <span>Health Check</span>
          </div>
          <ExternalLink class="w-3.5 h-3.5 text-slate-600 group-hover:text-slate-400" />
        </a>
      </nav>

      <!-- User Profile / Footer area -->
      <div class="p-4 border-t border-slate-800/80 bg-slate-950/40 space-y-3">
        <div class="flex items-center justify-between">
          <div class="min-w-0 pr-2">
            <p class="text-xs font-semibold text-white truncate">
              {{ user?.name || 'Developer' }}
            </p>
            <p class="text-[11px] text-slate-400 truncate">
              {{ user?.email || 'dev@example.com' }}
            </p>
          </div>
          <span 
            class="px-2 py-0.5 rounded-full text-[10px] font-mono font-medium border"
            :class="userTier === 'Premium' 
              ? 'bg-amber-950/50 text-amber-300 border-amber-500/30' 
              : 'bg-cyan-950/50 text-cyan-300 border-cyan-500/30'"
          >
            {{ userTier }}
          </span>
        </div>

        <button
          @click="handleLogout"
          class="w-full py-2 px-3 bg-slate-900/80 hover:bg-rose-950/40 text-slate-400 hover:text-rose-300 border border-slate-800 hover:border-rose-500/30 rounded-xl text-xs font-medium flex items-center justify-center space-x-2 transition-all cursor-pointer"
        >
          <LogOut class="w-3.5 h-3.5" />
          <span>Sign Out</span>
        </button>
      </div>
    </aside>

    <!-- Main Content Area -->
    <div class="flex-1 lg:pl-64 flex flex-col min-w-0">
      <!-- Top header bar -->
      <header class="h-16 border-b border-slate-800/80 bg-slate-950/60 backdrop-blur-md sticky top-0 z-30 px-4 sm:px-8 flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <button 
            @click="mobileMenuOpen = true" 
            class="lg:hidden p-2 text-slate-400 hover:text-white rounded-lg border border-slate-800"
          >
            <Menu class="w-5 h-5" />
          </button>
          
          <div class="flex items-center space-x-2 text-xs">
            <span class="text-slate-500">Dashboard</span>
            <span class="text-slate-600">/</span>
            <span class="text-white font-medium capitalize">{{ route.name?.toString() || 'Console' }}</span>
          </div>
        </div>

        <!-- Live Server Status Badge -->
        <div class="flex items-center space-x-2 px-3 py-1 rounded-full text-xs font-medium bg-emerald-950/40 text-emerald-300 border border-emerald-500/30">
          <span class="relative flex h-2 w-2">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2 w-2 bg-emerald-400"></span>
          </span>
          <span class="font-mono text-[11px]">Core API :8080 Active</span>
        </div>
      </header>

      <!-- Page View Render -->
      <main class="flex-1 p-4 sm:p-8 max-w-7xl w-full mx-auto">
        <router-view />
      </main>
    </div>
  </div>
</template>
