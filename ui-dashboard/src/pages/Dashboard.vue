<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="mb-8">
      <h2 class="text-3xl font-bold text-gray-900 mb-2">Dashboard de Reportes BCRA</h2>
      <p class="text-gray-600">Estadísticas de evaluaciones de seguridad</p>
    </div>

    <!-- Estadísticas principales -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-4">
      <div class="card bg-blue-50 border-blue-200">
        <div class="text-3xl font-bold text-blue-600">{{ statistics.total_reports }}</div>
        <p class="text-gray-600 text-sm mt-2">Reportes Totales</p>
      </div>

      <div class="card bg-red-50 border-red-200">
        <div class="text-3xl font-bold text-red-600">{{ statistics.total_critical }}</div>
        <p class="text-gray-600 text-sm mt-2">Críticos</p>
      </div>

      <div class="card bg-orange-50 border-orange-200">
        <div class="text-3xl font-bold text-orange-600">{{ statistics.total_high }}</div>
        <p class="text-gray-600 text-sm mt-2">Altos</p>
      </div>

      <div class="card bg-yellow-50 border-yellow-200">
        <div class="text-3xl font-bold text-yellow-600">{{ statistics.total_medium }}</div>
        <p class="text-gray-600 text-sm mt-2">Medios</p>
      </div>

      <div class="card bg-green-50 border-green-200">
        <div class="text-3xl font-bold text-green-600">{{ statistics.total_low }}</div>
        <p class="text-gray-600 text-sm mt-2">Bajos</p>
      </div>

      <div class="card bg-purple-50 border-purple-200">
        <div class="text-3xl font-bold text-purple-600">{{ statistics.total_findings }}</div>
        <p class="text-gray-600 text-sm mt-2">Total Hallazgos</p>
      </div>
    </div>

    <!-- Últimos reportes -->
    <div class="card">
      <h3 class="text-xl font-bold mb-4">Últimos Reportes</h3>
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-sm font-medium text-gray-900">ID</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-gray-900">Institución</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-gray-900">Riesgo</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-gray-900">Hallazgos</th>
              <th class="px-6 py-3 text-left text-sm font-medium text-gray-900">Fecha</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-if="reports.length === 0" class="hover:bg-gray-50">
              <td colspan="5" class="px-6 py-4 text-center text-gray-500">Sin reportes aún</td>
            </tr>
            <tr v-for="report in reports.slice(0, 5)" :key="report.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 text-sm text-gray-900 font-mono">{{ report.report_id }}</td>
              <td class="px-6 py-4 text-sm text-gray-900">{{ report.institution }}</td>
              <td class="px-6 py-4 text-sm">
                <span :class="['badge', getRiskBadgeClass(report.risk_level)]">
                  {{ report.risk_level }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">
                <span class="badge badge-critical">{{ report.critical }}</span>
                <span class="badge badge-high">{{ report.high }}</span>
                <span class="badge badge-medium">{{ report.medium }}</span>
                <span class="badge badge-low">{{ report.low }}</span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ formatDate(report.created_at) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useReportStore } from '../stores/reportStore'

const reportStore = useReportStore()

const statistics = computed(() => reportStore.statistics)
const reports = computed(() => reportStore.reports)

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
</script>
