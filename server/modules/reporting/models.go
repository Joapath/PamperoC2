package reporting

import "time"

// Finding representa un hallazgo de seguridad identificado en la evaluación
type Finding struct {
	ID              string
	Title           string
	Description     string
	Risk            string
	Impact          string
	Evidence        string
	AffectedSystems []string
	DiscoveredDate  time.Time
}

// RiskItem representa un riesgo en la matriz de riesgos
type RiskItem struct {
	ID           string
	Category     string
	RiskName     string
	Probability  int
	Impact       int
	RiskScore    int
	MitigationId string
}

// RemediationItem representa una acción de remediación
type RemediationItem struct {
	ID              string
	FindingID       string
	Recommendation  string
	Priority        string
	Timeline        string
	Owner           string
	Status          string
	EstimatedCost   string
	ImplementedDate time.Time
}

// AIAnalysisItem representa el análisis inteligente de un hallazgo
type AIAnalysisItem struct {
	FindingID         string
	AIEnrichment      string
	AttackVectors     []string
	TechnicalCommands []string
	Priority          string
	Confidence        float64
}

// BCRAReport es la estructura principal del reporte BCRA A 8398/2026
type BCRAReport struct {
	ReportID            string
	InstitutionName     string
	InstitutionType     string
	InstitutionAddress  string
	AssessmentPeriod    string
	ReportDate          time.Time
	ReportValidUntil    time.Time
	AssessmentTeam      string
	AssessmentTeamEmail string
	ExecutiveSummary    string
	OverallRiskLevel    string
	ComplianceStatus    string
	Findings            []Finding
	RiskMatrix          []RiskItem
	Remediations        []RemediationItem
	AIAnalysis          []AIAnalysisItem
	MethodologyUsed     string
	ControlsAssessed    int
	AdditionalNotes     string
	CriticalFindings    int
	HighFindings        int
	MediumFindings      int
	LowFindings         int
}

// CalculateStats calcula estadísticas del reporte
func (r *BCRAReport) CalculateStats() {
	r.CriticalFindings = 0
	r.HighFindings = 0
	r.MediumFindings = 0
	r.LowFindings = 0

	for _, finding := range r.Findings {
		switch finding.Risk {
		case "CRITICAL":
			r.CriticalFindings++
		case "HIGH":
			r.HighFindings++
		case "MEDIUM":
			r.MediumFindings++
		case "LOW":
			r.LowFindings++
		}
	}
}

// CalculateRiskScores calcula las puntuaciones de riesgo
func (r *BCRAReport) CalculateRiskScores() {
	for i := range r.RiskMatrix {
		r.RiskMatrix[i].RiskScore = r.RiskMatrix[i].Probability * r.RiskMatrix[i].Impact
	}
}
