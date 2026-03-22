<template>
  <div class="space-y-6">
    <button @click="goBack" class="text-blue-600 hover:text-blue-800">← Volver a Agentes</button>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div class="col-span-2 p-6 bg-white shadow rounded">
        <h2 class="text-xl font-bold mb-4">Agent {{ agent.agent_id || '...' }}</h2>
        <div class="mb-3"><strong>Host:</strong> {{ agent.hostname }}</div>
        <div class="mb-3"><strong>User:</strong> {{ agent.username }}</div>
        <div class="mb-3"><strong>OS:</strong> {{ agent.os }} / {{ agent.arch }}</div>
        <div class="mb-3"><strong>Status:</strong> {{ agent.status }}</div>
        <div class="mb-3"><strong>Profile:</strong> {{ agent.profile }}</div>
        <div class="mb-3"><strong>Last seen:</strong> {{ formatDate(agent.last_seen) }}</div>
      </div>

      <div class="p-6 bg-white shadow rounded">
        <h3 class="text-lg font-semibold mb-3">Acciones rápidas</h3>
        <div class="space-y-3">
          <button class="btn btn-secondary w-full" @click="runAction('recon-auto')">Recon Auto</button>
          <button class="btn btn-secondary w-full" @click="runAction('lateral-move', { target: '10.0.0.55' })">Lateral Move</button>
          <button class="btn btn-secondary w-full" @click="runAction('phish-whatsapp', { target: '+5491123456789' })">Phish WhatsApp</button>
          <button class="btn btn-secondary w-full" @click="runAction('generate-evasion-kaspersky')">Generate Evasion Kaspersky</button>
        </div>
      </div>
    </div>

    <div class="p-6 bg-white shadow rounded">
      <h3 class="text-xl font-bold mb-4">Enviar comando libre</h3>
      <div class="flex flex-wrap gap-2 items-end">
        <input v-model="command" type="text" class="input" placeholder="comando" />
        <input v-model="args" type="text" class="input" placeholder="args (separados por coma)" />
        <button class="btn btn-primary" @click="sendCustomCommand">Enviar</button>
      </div>
    </div>

    <div class="p-6 bg-white shadow rounded">
      <h3 class="text-xl font-bold mb-4">Jobs / resultados</h3>
      <div v-if="jobLoading">Cargando jobs...</div>
      <div v-else>
        <div v-if="jobs.length === 0" class="text-gray-500">No hay jobs</div>
        <ul class="space-y-3">
          <li v-for="job in jobs" :key="job.id" class="border rounded p-3">
            <div><strong>{{ job.command }}</strong> <span class="text-xs text-gray-500">{{ job.status }}</span></div>
            <div class="text-xs text-gray-500">Creado: {{ formatDate(job.created_at) }}</div>
            <button class="text-sm text-blue-600" @click="loadJobResults(job.id)">Ver resultados</button>
            <div v-if="jobResults[job.id]" class="mt-2 bg-gray-100 p-2 rounded">
              <div><strong>Stdout:</strong> <pre>{{ jobResults[job.id].stdout }}</pre></div>
              <div><strong>Stderr:</strong> <pre>{{ jobResults[job.id].stderr }}</pre></div>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useReportStore } from '../stores/reportStore'

const route = useRoute()
const router = useRouter()
const reportStore = useReportStore()

const agent = ref({})
const jobs = ref([])
const jobResults = ref({})
const command = ref('')
const args = ref('')
const jobLoading = ref(false)

const formatDate = (v) => {
  if (!v) return '-'
  return new Date(v).toLocaleString()
}

const loadAgent = async () => {
  const id = route.params.id
  try {
    const response = await reportStore.getAgent(id)
    agent.value = response
  } catch (err) {
    alert('No se encontró el agente')
    router.push('/agents')
  }
}

const loadJobs = async () => {
  jobLoading.value = true
  try {
    jobs.value = await reportStore.getAgentJobs(route.params.id)
  } catch (_) {
    jobs.value = []
  } finally {
    jobLoading.value = false
  }
}

const loadJobResults = async (jobID) => {
  try {
    const response = await reportStore.getJobResults(jobID)
    jobResults.value[jobID] = response
  } catch (err) {
    console.error(err)
  }
}

const runAction = async (action, params = {}) => {
  try {
    await reportStore.executeAgentAction(agent.value.agent_id, action, params)
    await loadJobs()
  } catch (err) {
    alert('Error ejecutando acción')
  }
}

const sendCustomCommand = async () => {
  const cmd = command.value.trim()
  if (!cmd) {
    alert('Ingrese un comando')
    return
  }
  const argList = args.value.split(',').map((a) => a.trim()).filter(Boolean)
  await runAction(cmd, { cmd })
  command.value = ''
  args.value = ''
}

const goBack = () => router.push('/agents')

onMounted(async () => {
  await loadAgent()
  await loadJobs()
})
</script>

<style scoped>
.input {
  border: 1px solid #d1d5db;
  padding: 0.5rem;
  border-radius: 0.375rem;
}

.btn {
  @apply px-4 py-2 rounded bg-blue-600 text-white;
}

.btn-primary {
  @apply bg-blue-600 hover:bg-blue-700;
}

.btn-secondary {
  @apply bg-gray-100 text-gray-800 hover:bg-gray-200;
}
</style>