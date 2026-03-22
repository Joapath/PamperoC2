<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white border-b border-gray-200 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 bg-blue-600 rounded flex items-center justify-center">
            <span class="text-white font-bold text-sm">🇦🇷</span>
          </div>
          <h1 class="text-2xl font-bold text-gray-900">PamperoC2</h1>
        </div>
        <nav class="flex gap-6">
          <router-link to="/" class="text-gray-600 hover:text-gray-900 font-medium">Dashboard</router-link>
          <router-link to="/reportes" class="text-gray-600 hover:text-gray-900 font-medium">Reportes</router-link>
          <router-link to="/crear" class="text-gray-600 hover:text-gray-900 font-medium">Crear</router-link>
        </nav>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <router-view />
    </main>

    <!-- Footer -->
    <footer class="bg-gray-900 text-gray-400 py-8 mt-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <p>PamperoC2 v1.0 — Framework C2 Argentino | © 2026</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useReportStore } from './stores/reportStore'

const reportStore = useReportStore()

onMounted(() => {
  reportStore.loadStatistics()
})
</script>

<style scoped>
</style>
