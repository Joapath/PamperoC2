<template>
  <div class="space-y-6">
    <h2 class="text-3xl font-bold text-gray-900">Crear Nuevo Reporte BCRA</h2>

    <form @submit.prevent="submitForm" class="space-y-6">
      <!-- Datos Institucionales -->
      <div class="card">
        <h3 class="text-lg font-semibold mb-4">Datos Institucionales</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <input
            v-model="form.institution_name"
            type="text"
            placeholder="Nombre de Institución *"
            required
            class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <input
            v-model="form.institution_type"
            type="text"
            placeholder="Tipo (Fintech, Banco, etc.)"
            class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <input
            v-model="form.institution_address"
            type="text"
            placeholder="Dirección"
            class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <input
            v-model="form.assessment_period"
            type="text"
            placeholder="Período (marzo 2026)"
            class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <input
            v-model="form.assessment_team"
            type="text"
            placeholder="Equipo evaluador"
            class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <input
            v-model="form.assessment_team_email"
            type="email"
            placeholder="Email del equipo"
            class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <!-- Reporte -->
      <div class="card">
        <h3 class="text-lg font-semibold mb-4">Información del Reporte</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <select
            v-model="form.overall_risk_level"
            class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="CRITICAL">Riesgo Crítico</option>
            <option value="HIGH">Riesgo Alto</option>
            <option value="MEDIUM">Riesgo Medio</option>
            <option value="LOW">Riesgo Bajo</option>
          </select>
          <select
            v-model="form.compliance_status"
            class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="COMPLIANT">Cumplimiento Completo</option>
            <option value="PARTIAL_COMPLIANT">Cumplimiento Parcial</option>
            <option value="NON_COMPLIANT">No Cumplimiento</option>
          </select>
          <input
            v-model="form.controls_assessed"
            type="number"
            placeholder="Controles evaluados"
            class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <input
            v-model="form.methodology_used"
            type="text"
            placeholder="Metodología (NIST, ISO, etc.)"
            class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <textarea
          v-model="form.executive_summary"
          placeholder="Resumen Ejecutivo"
          rows="4"
          class="w-full mt-4 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
        ></textarea>
      </div>

      <!-- Vista previa de hallazgos -->
      <div class="card">
        <h3 class="text-lg font-semibold mb-4">Hallazgos de Ejemplo</h3>
        <p class="text-gray-600 text-sm mb-4">
          Ahora mismo se generará un reporte con hallazgos de ejemplo. En integración final, se importarán desde la BD de Sliver.
        </p>
      </div>

      <!-- Botones -->
      <div class="flex gap-4">
        <button type="submit" :disabled="loading" class="btn btn-primary">
          {{ loading ? 'Generando...' : 'Generar Reporte' }}
        </button>
        <router-link to="/reportes" class="btn btn-secondary">Cancelar</router-link>
      </div>

      <div v-if="error" class="p-4 bg-red-100 border border-red-300 rounded text-red-800">
        {{ error }}
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useReportStore } from '../stores/reportStore'

const router = useRouter()
const reportStore = useReportStore()

const loading = ref(false)
const error = ref(null)

const form = ref({
  institution_name: '',
  institution_type: 'Fintech',
  institution_address: '',
  assessment_period: 'Marzo 2026',
  assessment_team: 'PamperoC2 Red Team',
  assessment_team_email: 'redteam@pampero.ar',
  overall_risk_level: 'HIGH',
  compliance_status: 'PARTIAL_COMPLIANT',
  controls_assessed: 56,
  methodology_used: 'TLPT (Threat Led Penetration Testing)',
  executive_summary: '',
})

const submitForm = async () => {
  if (!form.value.institution_name) {
    error.value = 'Falta nombre de institución'
    return
  }

  loading.value = true
  error.value = null

  try {
    // Hallazgos de ejemplo
    const findings = [
      {
        id: 'FIN-001',
        title: 'API sin autenticación',
        description: 'Endpoint expuesto sin requerir tokens',
        risk: 'CRITICAL',
        impact: 'Fuga de datos sensibles',
        evidence: 'GET /api/data retorna JSON sin auth',
        affected_systems: ['API Server 1', 'Load Balancer'],
        discovered_date: new Date().toISOString(),
      },
      {
        id: 'FIN-002',
        title: 'SQL Injection',
        description: 'Parámetro vulnerable a inyección SQL',
        risk: 'CRITICAL',
        impact: 'Acceso no autorizado a BD',
        evidence: 'customer_id=1\' OR \'1\'=\'1',
        affected_systems: ['PostgreSQL', 'Aplicación Web'],
        discovered_date: new Date().toISOString(),
      },
      {
        id: 'FIN-003',
        title: 'TLS 1.0 habilitado',
        description: 'Servidor acepta conexiones TLS 1.0 deprecado',
        risk: 'HIGH',
        impact: 'Vulnerabilidad a ataques BEAST',
        evidence: 'nmap mostró TLSv1.0 activo',
        affected_systems: ['HTTPS Gateway'],
        discovered_date: new Date().toISOString(),
      },
    ]

    const riskMatrix = [
      {
        id: 'RK-001',
        category: 'Seguridad de Aplicaciones',
        risk_name: 'API sin autenticación',
        probability: 9,
        impact: 10,
        mitigation_id: 'REM-001',
      },
      {
        id: 'RK-002',
        category: 'BD de Datos',
        risk_name: 'Inyección SQL',
        probability: 8,
        impact: 9,
        mitigation_id: 'REM-002',
      },
    ]

    const remediations = [
      {
        id: 'REM-001',
        finding_id: 'FIN-001',
        recommendation: 'Implementar OAuth 2.0 con JWT',
        priority: 'CRITICAL',
        timeline: '7 días',
        owner: 'Ing. Backend Lead',
        status: 'Pendiente',
        estimated_cost: '$3000 USD',
      },
    ]

    const reportData = {
      institution_name: form.value.institution_name,
      institution_type: form.value.institution_type,
      institution_address: form.value.institution_address,
      assessment_period: form.value.assessment_period,
      assessment_team: form.value.assessment_team,
      assessment_team_email: form.value.assessment_team_email,
      overall_risk_level: form.value.overall_risk_level,
      compliance_status: form.value.compliance_status,
      methodology_used: form.value.methodology_used,
      controls_assessed: form.value.controls_assessed,
      executive_summary: form.value.executive_summary,
      findings,
      risk_matrix: riskMatrix,
      remediations,
    }

    const response = await reportStore.generateReport(reportData)
    alert(`Reporte creado: ${response.report_id}`)
    router.push('/reportes')
  } catch (err) {
    error.value = `Error: ${err.message}`
  } finally {
    loading.value = false
  }
}
</script>
