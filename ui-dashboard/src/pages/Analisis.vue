<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="flex justify-between items-center mb-8">
      <div>
        <h2 class="text-3xl font-bold text-gray-900">Análisis de Hallazgos</h2>
        <p class="text-gray-600">Comparación entre hallazgo original y análisis de IA</p>
      </div>
      <button
        @click="reanalyze"
        :disabled="loadingFindings"
        class="btn-primary"
      >
        <span v-if="loadingFindings">Analizando...</span>
        <span v-else>Re-analizar</span>
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loadingFindings" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      <span class="ml-4 text-gray-600">Analizando hallazgos...</span>
    </div>

    <!-- Error State -->
    <div v-else-if="findingsError" class="card bg-red-50 border-red-200">
      <div class="text-red-600">
        <h3 class="font-bold">Error al cargar análisis</h3>
        <p>{{ findingsError }}</p>
      </div>
    </div>

    <!-- Findings Analysis -->
    <div v-else-if="findings.length > 0" class="space-y-6">
      <div
        v-for="finding in findings"
        :key="finding.id"
        class="grid grid-cols-1 lg:grid-cols-2 gap-6"
      >
        <!-- Original Finding -->
        <div class="card">
          <h3 class="text-xl font-bold mb-4 text-gray-800">Hallazgo Original</h3>
          <div class="space-y-3">
            <div>
              <strong>Título:</strong> {{ finding.title }}
            </div>
            <div>
              <strong>Descripción:</strong> {{ finding.description }}
            </div>
            <div>
              <strong>Riesgo:</strong>
              <span
                :class="getRiskClass(finding.risk)"
                class="px-2 py-1 rounded text-sm font-medium"
              >
                {{ finding.risk }}
              </span>
            </div>
          </div>
        </div>

        <!-- AI Analysis -->
        <div class="card">
          <h3 class="text-xl font-bold mb-4 text-blue-800">Análisis IA</h3>
          <div v-if="finding.aiAnalysis" class="space-y-3">
            <div>
              <strong>Attack Vectors:</strong>
              <ul class="list-disc list-inside ml-4">
                <li v-for="vector in finding.aiAnalysis.attackVectors" :key="vector">
                  {{ vector }}
                </li>
              </ul>
            </div>
            <div>
              <strong>Técnicas:</strong>
              <ul class="list-disc list-inside ml-4">
                <li v-for="technique in finding.aiAnalysis.techniques" :key="technique">
                  {{ technique }}
                </li>
              </ul>
            </div>
            <div>
              <strong>Comandos:</strong>
              <pre class="bg-gray-100 p-2 rounded text-sm overflow-x-auto">{{ finding.aiAnalysis.commands }}</pre>
            </div>
            <div>
              <strong>Prioridad:</strong>
              <span
                :class="getPriorityClass(finding.aiAnalysis.priority)"
                class="px-2 py-1 rounded text-sm font-medium"
              >
                {{ finding.aiAnalysis.priority }}
              </span>
            </div>
          </div>
          <div v-else class="text-gray-500">
            Análisis no disponible
          </div>
        </div>
      </div>
    </div>

    <!-- No Findings -->
    <div v-else class="card text-center py-12">
      <p class="text-gray-600">No hay hallazgos para analizar</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useReportStore } from '../stores/reportStore'

// Store
const reportStore = useReportStore()

// Computed from store
const { findings, loadingFindings, findingsError, loadFindings, reanalyzeFindings } = reportStore

// Methods
const reanalyze = async () => {
  await reanalyzeFindings()
}

const getRiskClass = (risk) => {
  const classes = {
    Critical: 'bg-red-100 text-red-800',
    High: 'bg-orange-100 text-orange-800',
    Medium: 'bg-yellow-100 text-yellow-800',
    Low: 'bg-green-100 text-green-800'
  }
  return classes[risk] || 'bg-gray-100 text-gray-800'
}

const getPriorityClass = (priority) => {
  const classes = {
    Critical: 'bg-red-100 text-red-800',
    High: 'bg-orange-100 text-orange-800',
    Medium: 'bg-yellow-100 text-yellow-800',
    Low: 'bg-green-100 text-green-800'
  }
  return classes[priority] || 'bg-gray-100 text-gray-800'
}

// Lifecycle
onMounted(() => {
  loadFindings()
})
</script>

<style scoped>
.card {
  @apply bg-white border border-gray-200 rounded-lg p-6 shadow-sm;
}

.btn-primary {
  @apply bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed;
}
</style>