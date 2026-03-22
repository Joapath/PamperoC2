<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <h2 class="text-3xl font-bold">C2 Agents</h2>
      <button @click="refreshAgents" class="btn btn-primary">Refrescar agents</button>
    </div>

    <div v-if="loading" class="text-center p-8">Cargando agentes...</div>

    <div v-else>
      <div v-if="agents.length === 0" class="text-center p-8 text-gray-500">No hay agentes conectados.</div>
      <table v-else class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Agent ID</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Host</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">User</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">OS</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Last Seen</th>
            <th class="px-6 py-3"></th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="agent in agents" :key="agent.id">
            <td class="px-6 py-4 text-sm font-mono">{{ agent.agent_id }}</td>
            <td class="px-6 py-4 text-sm">{{ agent.hostname }}</td>
            <td class="px-6 py-4 text-sm">{{ agent.username }}</td>
            <td class="px-6 py-4 text-sm">{{ agent.os }} / {{ agent.arch }}</td>
            <td class="px-6 py-4 text-sm">{{ agent.status }}</td>
            <td class="px-6 py-4 text-sm">{{ formatDate(agent.last_seen) }}</td>
            <td class="px-6 py-4 text-right">
              <router-link :to="`/agents/${agent.agent_id}`" class="text-blue-600 hover:text-blue-800">Ver</router-link>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useReportStore } from '../stores/reportStore'

const reportStore = useReportStore()
const loading = ref(false)

const formatDate = (v) => {
  if (!v) return '-'
  return new Date(v).toLocaleString()
}

const refreshAgents = async () => {
  loading.value = true
  try {
    await reportStore.loadAgents()
  } catch (err) {
    alert('Error cargando agentes: ' + err.message)
  } finally {
    loading.value = false
  }
}

onMounted(() => refreshAgents())
</script>
