package storage

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bishopfox/sliver/server/modules/reporting"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// StoredReport es la entidad persistida en BD
type StoredReport struct {
	ID                  string    `gorm:"primaryKey" json:"id"`
	ReportID            string    `json:"report_id"`
	InstitutionName     string    `json:"institution_name"`
	InstitutionType     string    `json:"institution_type"`
	AssessmentPeriod    string    `json:"assessment_period"`
	OverallRiskLevel    string    `json:"overall_risk_level"`
	ComplianceStatus    string    `json:"compliance_status"`
	CriticalFindings    int       `json:"critical_findings"`
	HighFindings        int       `json:"high_findings"`
	MediumFindings      int       `json:"medium_findings"`
	LowFindings         int       `json:"low_findings"`
	FindingsCount       int       `json:"findings_count"`
	ReportDate          time.Time `json:"report_date"`
	PDFPath             string    `json:"pdf_path"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	ExecutiveSummary    string    `json:"executive_summary"`
	MethodologyUsed     string    `json:"methodology_used"`
	ControlsAssessed    int       `json:"controls_assessed"`
	AdditionalNotes     string    `json:"additional_notes"`
	AssessmentTeam      string    `json:"assessment_team"`
	AssessmentTeamEmail string    `json:"assessment_team_email"`
	ReportValidUntil    time.Time `json:"report_valid_until"`
	AIAnalysis          string    `json:"ai_analysis"` // JSON string
	AIAnalysisTimestamp time.Time `json:"ai_analysis_timestamp"`
	AIModelUsed         string    `json:"ai_model_used"`
}

// StoredAgent representa el implant en la DB
type StoredAgent struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	AgentID       string    `json:"agent_id"`
	Hostname      string    `json:"hostname"`
	Username      string    `json:"username"`
	OS            string    `json:"os"`
	Arch          string    `json:"arch"`
	BeaconHeaders string    `json:"beacon_headers"` // JSON encoded
	Status        string    `json:"status"`
	Tags          string    `json:"tags"` // comma-separated
	Zone          string    `json:"zone"`
	Profile       string    `json:"profile"`
	BeaconCount   int       `json:"beacon_count"`
	LastSeen      time.Time `json:"last_seen"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// StoredJob representa una tarea C2
type StoredJob struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	AgentID    string    `json:"agent_id"`
	Command    string    `json:"command"`
	Args       string    `json:"args"`   // JSON array string
	Status     string    `json:"status"` // pending,running,done,failed
	TimeoutSec int       `json:"timeout_sec"`
	ResultID   string    `json:"result_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// StoredJobResult representa salida de comando
type StoredJobResult struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	JobID     string    `json:"job_id"`
	AgentID   string    `json:"agent_id"`
	ExitCode  int       `json:"exit_code"`
	Stdout    string    `json:"stdout"`
	Stderr    string    `json:"stderr"`
	CreatedAt time.Time `json:"created_at"`
}

var DB *gorm.DB

// Init inicializa la base de datos SQLite
func Init(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("falló conectar a BD: %w", err)
	}
	err = DB.AutoMigrate(&StoredReport{}, &StoredAgent{}, &StoredJob{}, &StoredJobResult{})
	if err != nil {
		return fmt.Errorf("falló migración: %w", err)
	}
	return nil
}

// SaveReport guarda un reporte en BD
func SaveReport(bcraReport *reporting.BCRAReport, pdfPath string) (*StoredReport, error) {
	if DB == nil {
		return nil, fmt.Errorf("base de datos no inicializada")
	}
	stored := &StoredReport{
		ID:                  fmt.Sprintf("rpt_%d", time.Now().UnixNano()),
		ReportID:            bcraReport.ReportID,
		InstitutionName:     bcraReport.InstitutionName,
		InstitutionType:     bcraReport.InstitutionType,
		AssessmentPeriod:    bcraReport.AssessmentPeriod,
		OverallRiskLevel:    bcraReport.OverallRiskLevel,
		ComplianceStatus:    bcraReport.ComplianceStatus,
		CriticalFindings:    bcraReport.CriticalFindings,
		HighFindings:        bcraReport.HighFindings,
		MediumFindings:      bcraReport.MediumFindings,
		LowFindings:         bcraReport.LowFindings,
		FindingsCount:       len(bcraReport.Findings),
		ReportDate:          bcraReport.ReportDate,
		PDFPath:             pdfPath,
		ExecutiveSummary:    bcraReport.ExecutiveSummary,
		MethodologyUsed:     bcraReport.MethodologyUsed,
		ControlsAssessed:    bcraReport.ControlsAssessed,
		AdditionalNotes:     bcraReport.AdditionalNotes,
		AssessmentTeam:      bcraReport.AssessmentTeam,
		AssessmentTeamEmail: bcraReport.AssessmentTeamEmail,
		ReportValidUntil:    bcraReport.ReportValidUntil,
	}
	result := DB.Create(stored)
	return stored, result.Error
}

// GetReport obtiene un reporte por ID
func GetReport(id string) (*StoredReport, error) {
	var report StoredReport
	result := DB.First(&report, "id = ?", id)
	return &report, result.Error
}

// ListReports lista todos los reportes con paginación
func ListReports(page, pageSize int) ([]StoredReport, int64, error) {
	var reports []StoredReport
	var total int64
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	result := DB.Model(&StoredReport{}).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	result = DB.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&reports)
	return reports, total, result.Error
}

// SaveAgent guarda o actualiza información de un agente
func SaveAgent(agent *reporting.Agent) (*StoredAgent, error) {
	if DB == nil {
		return nil, fmt.Errorf("base de datos no inicializada")
	}
	beaconHeadersJSON, _ := json.Marshal(agent.BeaconHeaders)
	stored := &StoredAgent{
		ID:            fmt.Sprintf("agt_%d", time.Now().UnixNano()),
		AgentID:       agent.AgentID,
		Hostname:      agent.Hostname,
		Username:      agent.Username,
		OS:            agent.OS,
		Arch:          agent.Arch,
		BeaconHeaders: string(beaconHeadersJSON),
		Status:        agent.Status,
		Tags:          strings.Join(agent.Tags, ","),
		Zone:          agent.Zone,
		Profile:       agent.Profile,
		BeaconCount:   agent.BeaconCount,
		LastSeen:      agent.LastSeen,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := DB.Where("agent_id = ?", agent.AgentID).First(&StoredAgent{}).Error; err == nil {
		// ya existe -> update
		if err := DB.Model(&StoredAgent{}).Where("agent_id = ?", agent.AgentID).Updates(map[string]interface{}{
			"hostname":       agent.Hostname,
			"username":       agent.Username,
			"os":             agent.OS,
			"arch":           agent.Arch,
			"beacon_headers": string(beaconHeadersJSON),
			"status":         agent.Status,
			"tags":           strings.Join(agent.Tags, ","),
			"zone":           agent.Zone,
			"profile":        agent.Profile,
			"beacon_count":   agent.BeaconCount,
			"last_seen":      agent.LastSeen,
			"updated_at":     time.Now(),
		}).Error; err != nil {
			return nil, err
		}
		if err := DB.Where("agent_id = ?", agent.AgentID).First(stored).Error; err != nil {
			return nil, err
		}
		return stored, nil
	}

	// create nuevo
	if err := DB.Create(stored).Error; err != nil {
		return nil, err
	}
	return stored, nil
}

func GetAgent(agentID string) (*StoredAgent, error) {
	var agent StoredAgent
	if err := DB.First(&agent, "agent_id = ?", agentID).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

func ListAgents() ([]StoredAgent, error) {
	var agents []StoredAgent
	if err := DB.Order("last_seen desc").Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

func SaveJob(job *reporting.AgentJob) (*StoredJob, error) {
	if DB == nil {
		return nil, fmt.Errorf("base de datos no inicializada")
	}
	stored := &StoredJob{
		ID:         fmt.Sprintf("job_%d", time.Now().UnixNano()),
		AgentID:    job.AgentID,
		Command:    job.Command,
		Args:       strings.Join(job.Args, ","),
		Status:     "pending",
		TimeoutSec: job.TimeoutSec,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := DB.Create(stored).Error; err != nil {
		return nil, err
	}
	return stored, nil
}

func GetPendingJobs(agentID string) ([]StoredJob, error) {
	var jobs []StoredJob
	if err := DB.Where("agent_id = ? AND status = ?", agentID, "pending").Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func UpdateJobStatus(jobID, status, resultID string) error {
	return DB.Model(&StoredJob{}).Where("id = ?", jobID).Updates(map[string]interface{}{
		"status":     status,
		"result_id":  resultID,
		"updated_at": time.Now(),
	}).Error
}

func SaveJobResult(result *reporting.AgentJobResult) (*StoredJobResult, error) {
	if DB == nil {
		return nil, fmt.Errorf("base de datos no inicializada")
	}
	stored := &StoredJobResult{
		ID:        fmt.Sprintf("jr_%d", time.Now().UnixNano()),
		JobID:     result.JobID,
		AgentID:   result.AgentID,
		ExitCode:  result.ExitCode,
		Stdout:    result.Stdout,
		Stderr:    result.Stderr,
		CreatedAt: time.Now(),
	}
	if err := DB.Create(stored).Error; err != nil {
		return nil, err
	}
	return stored, nil
}

func GetJobResults(jobID string) ([]StoredJobResult, error) {
	var results []StoredJobResult
	if err := DB.Where("job_id = ?", jobID).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// DeleteReport elimina un reporte
func DeleteReport(id string) error {
	return DB.Delete(&StoredReport{}, "id = ?", id).Error
}

// GetStatistics obtiene estadísticas generales
func GetStatistics() (map[string]interface{}, error) {
	var total int64
	var critical, high, medium, low int64
	DB.Model(&StoredReport{}).Count(&total)
	DB.Model(&StoredReport{}).Select("SUM(critical_findings)").Row().Scan(&critical)
	DB.Model(&StoredReport{}).Select("SUM(high_findings)").Row().Scan(&high)
	DB.Model(&StoredReport{}).Select("SUM(medium_findings)").Row().Scan(&medium)
	DB.Model(&StoredReport{}).Select("SUM(low_findings)").Row().Scan(&low)
	return map[string]interface{}{
		"total_reports":  total,
		"total_critical": critical,
		"total_high":     high,
		"total_medium":   medium,
		"total_low":      low,
		"total_findings": critical + high + medium + low,
	}, nil
}

// GetAIAnalysis obtiene el análisis IA de un reporte
func GetAIAnalysis(reportID string) ([]reporting.AIAnalysisItem, error) {
	var report StoredReport
	if err := DB.Where("id = ?", reportID).First(&report).Error; err != nil {
		return nil, err
	}

	if report.AIAnalysis == "" {
		return nil, fmt.Errorf("no AI analysis found")
	}

	var aiAnalysis []reporting.AIAnalysisItem
	if err := json.Unmarshal([]byte(report.AIAnalysis), &aiAnalysis); err != nil {
		return nil, err
	}

	return aiAnalysis, nil
}

// SaveAIAnalysis guarda el análisis IA de un reporte
func SaveAIAnalysis(reportID string, aiAnalysis []reporting.AIAnalysisItem) error {
	aiAnalysisJSON, err := json.Marshal(aiAnalysis)
	if err != nil {
		return err
	}

	return DB.Model(&StoredReport{}).Where("id = ?", reportID).Updates(map[string]interface{}{
		"ai_analysis":           string(aiAnalysisJSON),
		"ai_analysis_timestamp": time.Now(),
		"ai_model_used":         "mistral",
	}).Error
}
