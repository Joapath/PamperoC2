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

  const findings = ref([])
  const loadingFindings = ref(false)
  const findingsError = ref(null)

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
})
