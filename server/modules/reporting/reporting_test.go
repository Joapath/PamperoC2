package reporting

import (
	"os"
	"testing"
	"time"
)

// TestGenerateBCRAReport_BasicStructure verifica que el PDF se genera correctamente
func TestGenerateBCRAReport_BasicStructure(t *testing.T) {
	report := createMockBCRAReport()

	pdf, err := GenerateBCRAReport(report)
	if err != nil {
		t.Errorf("Error generando PDF: %v", err)
	}

	if len(pdf) == 0 {
		t.Error("PDF vacío - no se generó contenido")
	}

	// Verificar que el PDF tiene la firma correcta
	if pdf[0] != 0x25 || pdf[1] != 0x50 || pdf[2] != 0x44 || pdf[3] != 0x46 { // %PDF
		t.Error("Archivo generado no es un PDF válido")
	}
}

// TestCalculateStats verifica el cálculo de estadísticas
func TestCalculateStats(t *testing.T) {
	report := createMockBCRAReport()
	report.CalculateStats()

	if report.CriticalFindings != 1 {
		t.Errorf("Esperaba 1 hallazgo crítico, obtuvo %d", report.CriticalFindings)
	}

	if report.HighFindings != 1 {
		t.Errorf("Esperaba 1 hallazgo alto, obtuvo %d", report.HighFindings)
	}

	if report.MediumFindings != 1 {
		t.Errorf("Esperaba 1 hallazgo medio, obtuvo %d", report.MediumFindings)
	}

	if report.LowFindings != 1 {
		t.Errorf("Esperaba 1 hallazgo bajo, obtuvo %d", report.LowFindings)
	}
}

// TestCalculateRiskScores verifica el cálculo de riesgos
func TestCalculateRiskScores(t *testing.T) {
	report := createMockBCRAReport()
	report.CalculateRiskScores()

	for _, risk := range report.RiskMatrix {
		expected := risk.Probability * risk.Impact
		if risk.RiskScore != expected {
			t.Errorf("Risk score incorrecto: esperaba %d, obtuvo %d", expected, risk.RiskScore)
		}
	}
}

// TestGenerateBCRAReport_PDFSize verifica que el PDF tenga tamaño razonable
func TestGenerateBCRAReport_PDFSize(t *testing.T) {
	report := createMockBCRAReport()

	pdf, err := GenerateBCRAReport(report)
	if err != nil {
		t.Errorf("Error generando PDF: %v", err)
	}

	// El PDF debe tener al menos 500 bytes
	if len(pdf) < 500 {
		t.Errorf("PDF demasiado pequeño: %d bytes", len(pdf))
	}

	// El PDF no debe exceder 5MB (sanidad check)
	if len(pdf) > 5*1024*1024 {
		t.Errorf("PDF demasiado grande: %d bytes", len(pdf))
	}
}

// TestGenerateBCRAReport_SaveTestPDF genera un PDF de prueba
func TestGenerateBCRAReport_SaveTestPDF(t *testing.T) {
	report := createMockBCRAReport()

	pdf, err := GenerateBCRAReport(report)
	if err != nil {
		t.Fatalf("Error generando PDF: %v", err)
	}

	// Guardar PDF de prueba en /tmp/
	testFile := "/tmp/test_bcra_report_" + time.Now().Format("20060102_150405") + ".pdf"
	err = os.WriteFile(testFile, pdf, 0644)
	if err != nil {
		t.Logf("No se pudo guardar PDF de prueba: %v", err)
	} else {
		t.Logf("PDF de prueba guardado en: %s", testFile)
	}
}

// TestTruncateText verifica que el truncamiento funciona correctamente
func TestTruncateText(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"Hello World", 5, "Hello..."},
		{"Hi", 5, "Hi"},
		{"Exact", 5, "Exact"},
	}

	for _, test := range tests {
		result := truncateText(test.input, test.maxLen)
		if result != test.expected {
			t.Errorf("truncateText(%q, %d) = %q, expected %q", test.input, test.maxLen, result, test.expected)
		}
	}
}

// Funciones auxiliares

// createMockBCRAReport crea un reporte de prueba con datos mock
func createMockBCRAReport() *BCRAReport {
	now := time.Now()
	validUntil := now.AddDate(0, 0, 365)

	return &BCRAReport{
		ReportID:           "BCRA-2026-001",
		InstitutionName:    "Banco Tecnológico Argentino S.A.",
		InstitutionType:    "Banco",
		InstitutionAddress: "Avenida Libertador 101, Buenos Aires",
		AssessmentPeriod:   "Enero-Febrero 2026",
		ReportDate:         now,
		ReportValidUntil:   validUntil,

		AssessmentTeam:      "PamperoC2 Red Team",
		AssessmentTeamEmail: "redteam@pampero.ar",

		ExecutiveSummary: "Se completó una evaluación integral de seguridad de la información...",
		OverallRiskLevel: "HIGH",
		ComplianceStatus: "PARTIAL_COMPLIANT",

		MethodologyUsed:  "TLPT + NIST CSF",
		ControlsAssessed: 42,

		Findings: []Finding{
			{
				ID:              "BCRA-001",
				Title:           "RDP Expuesto en Internet",
				Description:     "Se identificó que el servidor de RDP está accesible directamente desde internet.",
				Risk:            "CRITICAL",
				Impact:          "Un atacante podría obtener acceso directo a sistemas internos.",
				Evidence:        "Puerto 3389 abierto en IP pública",
				AffectedSystems: []string{"SRV-PROD-01"},
				DiscoveredDate:  now.AddDate(0, 0, -5),
			},
			{
				ID:              "BCRA-002",
				Title:           "Contraseñas Débiles",
				Description:     "Múltiples usuarios tienen contraseñas débiles.",
				Risk:            "HIGH",
				Impact:          "Compromiso de cuentas de usuario.",
				Evidence:        "Auditoría de contraseñas: 15 usuarios con passwords débiles",
				AffectedSystems: []string{"AD"},
				DiscoveredDate:  now.AddDate(0, 0, -4),
			},
			{
				ID:              "BCRA-003",
				Title:           "Software Desactualizado",
				Description:     "Varios servidores ejecutan versiones antiguas de software.",
				Risk:            "MEDIUM",
				Impact:          "Exposición a exploit de vulnerabilidades conocidas.",
				Evidence:        "Exchange 2013 (EOL) detectado",
				AffectedSystems: []string{"SRV-MAIL-01"},
				DiscoveredDate:  now.AddDate(0, 0, -3),
			},
			{
				ID:              "BCRA-004",
				Title:           "Falta de Monitoreo",
				Description:     "No existe auditoría centralizada de logs de acceso.",
				Risk:            "LOW",
				Impact:          "Dificultad en detectar accesos no autorizados.",
				Evidence:        "No hay eventos forwarding",
				AffectedSystems: []string{"General"},
				DiscoveredDate:  now.AddDate(0, 0, -2),
			},
		},

		RiskMatrix: []RiskItem{
			{
				ID:           "RK-001",
				Category:     "Gestión de Tecnología",
				RiskName:     "Falta de parches de seguridad",
				Probability:  8,
				Impact:       9,
				MitigationId: "REM-001",
			},
			{
				ID:           "RK-002",
				Category:     "Seguridad de la Información",
				RiskName:     "Acceso no autorizado",
				Probability:  7,
				Impact:       10,
				MitigationId: "REM-002",
			},
			{
				ID:           "RK-003",
				Category:     "Continuidad",
				RiskName:     "Falta de recuperación",
				Probability:  6,
				Impact:       8,
				MitigationId: "REM-003",
			},
		},

		Remediations: []RemediationItem{
			{
				ID:             "REM-001",
				FindingID:      "BCRA-001",
				Recommendation: "Bloquear acceso RDP desde internet.",
				Priority:       "CRITICAL",
				Timeline:       "7 días",
				Owner:          "Ing. Juan Pérez",
				Status:         "Pendiente",
				EstimatedCost:  "$2000 USD",
			},
			{
				ID:             "REM-002",
				FindingID:      "BCRA-002",
				Recommendation: "Enforcer política de contraseñas.",
				Priority:       "HIGH",
				Timeline:       "14 días",
				Owner:          "Lic. María García",
				Status:         "En Progreso",
				EstimatedCost:  "$500 USD",
			},
		},

		AdditionalNotes: "Se recomienda implementar programa de capacitación en seguridad.",
	}
}
