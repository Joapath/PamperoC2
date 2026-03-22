<template>
  <div class="space-y-6">
    <h2 class="text-3xl font-bold text-gray-900">Reportes BCRA</h2>

    <!-- Filtros y acciones -->
    <div class="flex gap-4 items-center">
      <input
        v-model="searchTerm"
        type="text"
        placeholder="Buscar por institución..."
        class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      />
      <router-link to="/crear" class="btn btn-primary">+ Crear Reporte</router-link>
    </div>

    <!-- Tabla de reportes -->
    <div class="card">
      <div v-if="loading" class="text-center py-8">
        <p class="text-gray-600">Cargando reportes...</p>
      </div>

      <div v-else-if="filteredReports.length === 0" class="text-center py-8">
        <p class="text-gray-500">No hay reportes disponibles</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-sm font-medium text-gray-900">ID</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-gray-900">Institución</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-gray-900">Tipo</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-gray-900">Riesgo</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-gray-900">Hallazgos</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-gray-900">Equipo</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-gray-900">Fecha</th>
              <th class="px-6 py-3 text-center text-sm font-medium text-gray-900">Acciones</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="report in filteredReports" :key="report.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 text-sm font-mono text-gray-900">{{ report.report_id.slice(0, 12) }}...</td>
              <td class="px-6 py-4 text-sm text-gray-900">{{ report.institution }}</td>
              <td class="px-6 py-4 text-sm text-gray-600">{{ report.institution_type }}</td>
              <td class="px-6 py-4 text-sm">
                <span :class="['badge', getRiskBadgeClass(report.risk_level)]">
                  {{ report.risk_level }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm">
                <div class="flex gap-1">
                  <span class="badge badge-critical text-xs">C:{{ report.critical }}</span>
                  <span class="badge badge-high text-xs">A:{{ report.high }}</span>
                  <span class="badge badge-medium text-xs">M:{{ report.medium }}</span>
                  <span class="badge badge-low text-xs">B:{{ report.low }}</span>
                </div>
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">{{ report.assessment_team }}</td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ formatDate(report.created_at) }}</td>
              <td class="px-6 py-4 text-center">
                <div class="flex gap-2 justify-center">
                  <button
                    @click="handleDownload(report.id)"
                    class="text-blue-600 hover:text-blue-800 text-sm font-medium"
                  >
                    📥
                  </button>
                  <button
                    @click="handleDelete(report.id)"
                    class="text-red-600 hover:text-red-800 text-sm font-medium"
                  >
                    🗑️
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useReportStore } from '../stores/reportStore'

const reportStore = useReportStore()
const searchTerm = ref('')

const loading = computed(() => reportStore.loading)
const allReports = computed(() => reportStore.reports)

const filteredReports = computed(() => {
  return allReports.value.filter(r =>
    r.institution.toLowerCase().includes(searchTerm.value.toLowerCase())
  )
})

onMounted(() => {
  reportStore.loadReports()
})

const getRiskBadgeClass = (risk) => {
  const classes = {
    CRITICAL: 'badge-critical',
    HIGH: 'badge-high',
    MEDIUM: 'badge-medium',
    LOW: 'badge-low',
  }
  return classes[risk] || 'badge-medium'
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('es-AR')
}

const handleDownload = async (reportId) => {
  try {
    await reportStore.downloadReport(reportId)
  } catch (error) {
    alert('Error descargando reporte: ' + error.message)
  }
}

const handleDelete = async (reportId) => {
  if (confirm('¿Eliminar este reporte?')) {
    try {
      await reportStore.deleteReport(reportId)
    } catch (error) {
      alert('Error eliminando reporte: ' + error.message)
    }
  }
}
</script>
