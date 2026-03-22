import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

export const useReportStore = defineStore('report', () => {
  const reports = ref([])
  const statistics = ref({
    total_reports: 0,
    total_critical: 0,
    total_high: 0,
    total_medium: 0,
    total_low: 0,
    total_findings: 0,
  })
  const loading = ref(false)
  const error = ref(null)

  const agents = ref([])
  const jobs = ref([])
  const findings = ref([])
  const loadingFindings = ref(false)
  const findingsError = ref(null)

  const apiClient = axios.create({
    baseURL: '/api/v1',
  })

  const loadReports = async (page = 1, pageSize = 10) => {
    loading.value = true
    error.value = null
    try {
      const response = await apiClient.get('/reports', {
        params: { page, page_size: pageSize },
      })
      reports.value = response.data.data
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const loadStatistics = async () => {
    try {
      const response = await apiClient.get('/statistics')
      statistics.value = response.data
    } catch (err) {
      console.error('Error cargando estadísticas:', err)
    }
  }

  const generateReport = async (reportData) => {
    loading.value = true
    error.value = null
    try {
      const response = await apiClient.post('/reports', reportData)
      reports.value.unshift(response.data)
      await loadStatistics()
      return response.data
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const downloadReport = async (reportId) => {
    try {
      const response = await apiClient.get(`/reports/${reportId}/download`, {
        responseType: 'blob',
      })
      const url = window.URL.createObjectURL(response.data)
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', `reporte_${reportId}.pdf`)
      document.body.appendChild(link)
      link.click()
      link.parentElement.removeChild(link)
    } catch (err) {
      error.value = 'Error descargando reporte'
      throw err
    }
  }

  const loadAgents = async () => {
    try {
      const response = await apiClient.get('/agents')
      agents.value = response.data.agents
      return agents.value
    } catch (err) {
      error.value = 'Error cargando agentes'
      throw err
    }
  }

  const loadFindings = async () => {
    loadingFindings.value = true
    findingsError.value = null
    try {
      const response = await apiClient.get('/findings')
      findings.value = response.data.data
    } catch (err) {
      findingsError.value = err.message
    } finally {
      loadingFindings.value = false
    }
  }

  const getAgent = async (agentId) => {
    try {
      const response = await apiClient.get(`/agents/${agentId}`)
      return response.data
    } catch (err) {
      error.value = 'Error cargando agente'
      throw err
    }
  }

  const executeAgentAction = async (agentId, action, params = {}) => {
    try {
      const response = await apiClient.post(`/agents/${agentId}/actions`, {
        action,
        params,
      })
      return response.data
    } catch (err) {
      error.value = 'Error ejecutando accion agente'
      throw err
    }
  }

  const getJobResults = async (jobId) => {
    try {
      const response = await apiClient.get(`/jobs/${jobId}/results`)
      return response.data.results
    } catch (err) {
      error.value = 'Error cargando resultados de job'
      throw err
    }
  }

  const reanalyzeFindings = async () => {
    loadingFindings.value = true
    findingsError.value = null
    try {
      const response = await apiClient.post('/findings/reanalyze')
      findings.value = response.data.data
    } catch (err) {
      findingsError.value = err.message
    } finally {
      loadingFindings.value = false
    }
  }

  const createAgentJob = async (agentId, command, args = [], timeoutSec = 300) => {
    try {
      const response = await apiClient.post(`/agents/${agentId}/jobs`, {
        command,
        args,
        timeout_sec: timeoutSec,
      })
      return response.data
    } catch (err) {
      error.value = 'Error creando job'
      throw err
    }
  }

  const getAgentJobs = async (agentId) => {
    try {
      const response = await apiClient.get(`/agents/${agentId}/jobs`)
      jobs.value = response.data.jobs
      return jobs.value
    } catch (err) {
      error.value = 'Error cargando jobs del agente'
      throw err
    }
  }

  const deleteReport = async (reportId) => {
    try {
      await apiClient.delete(`/reports/${reportId}`)
      reports.value = reports.value.filter(r => r.id !== reportId)
      await loadStatistics()
    } catch (err) {
      error.value = 'Error eliminando reporte'
      throw err
    }
  }

  return {
    reports,
    statistics,
    agents,
    jobs,
    findings,
    loading,
    error,
    loadingFindings,
    findingsError,
    loadReports,
    loadStatistics,
    generateReport,
    downloadReport,
    deleteReport,
    loadAgents,
    loadFindings,
    reanalyzeFindings,
    createAgentJob,
    getAgent,
    getAgentJobs,
    executeAgentAction,
    getJobResults,
  }
})
