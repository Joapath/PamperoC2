package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bishopfox/sliver/server/modules/reporting"
)

// ExampleBCRAReportGeneration muestra cómo generar un reporte BCRA
// Ejecutar: go run examples/bcra_report_example.go
func main() {
	fmt.Println("🇦🇷 PamperoC2 - Generador de Reportes BCRA")
	fmt.Println("=========================================\n")

	// Crear reporte de ejemplo
	report := createExampleReport()

	// Generar PDF
	fmt.Println("Generando reporte BCRA...")
	pdf, err := reporting.GenerateBCRAReport(report)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Guardar archivo
	filename := fmt.Sprintf("reporte_BCRA_%s.pdf", time.Now().Format("20060102_150405"))
	err = os.WriteFile(filename, pdf, 0644)
	if err != nil {
		fmt.Printf("❌ Error al guardar: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Reporte generado: %s\n", filename)
	fmt.Printf("📊 Tamaño: %.2f KB\n", float64(len(pdf))/1024)
	fmt.Printf("📄 Hallazgos: %d (Críticos: %d, Altos: %d, Medios: %d, Bajos: %d)\n",
		len(report.Findings), report.CriticalFindings, report.HighFindings,
		report.MediumFindings, report.LowFindings)
}

func createExampleReport() *reporting.BCRAReport {
	now := time.Now()

	return &reporting.BCRAReport{
		ReportID:            "BCRA-2026-PamperoC2-001",
		InstitutionName:     "Fintech Argentina Demo",
		InstitutionType:     "Fintech",
		InstitutionAddress:  "Palermo, Buenos Aires",
		AssessmentPeriod:    "Febrero 2026",
		ReportDate:          now,
		ReportValidUntil:    now.AddDate(0, 0, 365),
		AssessmentTeam:      "PamperoC2 Red Team",
		AssessmentTeamEmail: "redteam@pampero.ar",

		ExecutiveSummary: "Evaluación de seguridad realizada del 01-15 de febrero 2026. Se identificaron vulnerabilidades críticas en infraestructura externa y deficiencias en controles de acceso.",
		OverallRiskLevel: "HIGH",
		ComplianceStatus: "PARTIAL_COMPLIANT",

		MethodologyUsed:  "TLPT (Threat Led Penetration Testing) + NIST CSF",
		ControlsAssessed: 56,

		Findings: []reporting.Finding{
			{
				ID:              "BCRA-2026-001",
				Title:           "ServiceStack API sin autenticación",
				Description:     "Endpoint interno expuesto que retorna datos sensibles sin requerir tokens de autenticación. Afecta transacciones de los últimos 30 días.",
				Risk:            "CRITICAL",
				Impact:          "Exposición de datos de clientes, transacciones fraudulentas, violación de BCRA A 8398/2026 punto 4.3",
				Evidence:        "GET /api/transactions retorna JSON sin authentication. Response headers: Server: ServiceStack/5.11. Datos: 125 transacciones de muestra.",
				AffectedSystems: []string{"API Server 1", "API Server 2", "Load Balancer"},
				DiscoveredDate:  now.AddDate(0, 0, -15),
			},
			{
				ID:              "BCRA-2026-002",
				Title:           "SQL Injection en módulo de reportes",
				Description:     "Parámetro 'customer_id' vulnerable a inyección SQL. Permite lectura de tablas no autorizadas.",
				Risk:            "CRITICAL",
				Impact:          "Fuga de datos PII, IBAN, estados de cuenta. Cumplimiento: VIOLADOR",
				Evidence:        "POST /reports/generate - customer_id=1' OR '1'='1. SQLi confirmado en tabla 'accounts'",
				AffectedSystems: []string{"Reporting Backend", "PostgreSQL DB"},
				DiscoveredDate:  now.AddDate(0, 0, -14),
			},
			{
				ID:              "BCRA-2026-003",
				Title:           "Contraseña admin por defecto",
				Description:     "Consola de administración (Jenkins, Grafana) con credenciales por defecto admin:admin",
				Risk:            "HIGH",
				Impact:          "Ejecución de código en servidores de CI/CD. Compromiso de infraestructura completa",
				Evidence:        "Jenkins en 203.0.113.50:8080 - Acceso exitoso con admin:admin. Perfil: 90 trabajos sensibles",
				AffectedSystems: []string{"Jenkins", "Grafana", "Prometheus"},
				DiscoveredDate:  now.AddDate(0, 0, -13),
			},
			{
				ID:              "BCRA-2026-004",
				Title:           "TLS 1.0 habilitado",
				Description:     "Servidores aceptan conexiones con TLS 1.0 (deprecado 2018). Vulnerable a ataques BEAST/downgrade",
				Risk:            "HIGH",
				Impact:          "Interceptación de sesiones HTTPS encriptadas",
				Evidence:        "nmap --script ssl-enum-ciphers: TLSv1.0 habilitado en puerto 443",
				AffectedSystems: []string{"HTTPS Gateway", "API Servers"},
				DiscoveredDate:  now.AddDate(0, 0, -10),
			},
			{
				ID:              "BCRA-2026-005",
				Title:           "Falta de WAF",
				Description:     "No existe Web Application Firewall. Expuesto a OWASP Top 10",
				Risk:            "MEDIUM",
				Impact:          "Ataques XSS, CSRF sin detección. Cumplimiento parcial",
				Evidence:        "XSS reflejado en /search?q=<script>alert(1)</script>. Headers WAF: ninguno",
				AffectedSystems: []string{"Web Application"},
				DiscoveredDate:  now.AddDate(0, 0, -8),
			},
		},

		RiskMatrix: []reporting.RiskItem{
			{
				ID:           "RK-001",
				Category:     "Seguridad de Aplicaciones",
				RiskName:     "Inyección SQL en reportes",
				Probability:  9,
				Impact:       10,
				MitigationId: "REM-001",
			},
			{
				ID:           "RK-002",
				Category:     "Gestión de Acceso",
				RiskName:     "Credenciales por defecto en administración",
				Probability:  8,
				Impact:       9,
				MitigationId: "REM-002",
			},
			{
				ID:           "RK-003",
				Category:     "Tecnología y Encriptación",
				RiskName:     "TLS 1.0 expuesto",
				Probability:  7,
				Impact:       7,
				MitigationId: "REM-003",
			},
			{
				ID:           "RK-004",
				Category:     "Seguridad de Aplicaciones",
				RiskName:     "Falta de WAF",
				Probability:  8,
				Impact:       6,
				MitigationId: "REM-004",
			},
		},

		Remediations: []reporting.RemediationItem{
			{
				ID:             "REM-001",
				FindingID:      "BCRA-2026-001",
				Recommendation: "Implementar autenticación OAuth 2.0 con tokens JWT. Requerir Authorization header en todos los endpoints.",
				Priority:       "CRITICAL",
				Timeline:       "7 días",
				Owner:          "Ing. Backend Lead",
				Status:         "Pendiente",
				EstimatedCost:  "$3000 USD",
			},
			{
				ID:             "REM-002",
				FindingID:      "BCRA-2026-002",
				Recommendation: "Usar prepared statements en todas las queries. Integrar ORM (GORM/SQLAlchemy) y validar inputs",
				Priority:       "CRITICAL",
				Timeline:       "5 días",
				Owner:          "Ing. Security + Backend",
				Status:         "En Progreso",
				EstimatedCost:  "$2000 USD",
			},
			{
				ID:             "REM-003",
				FindingID:      "BCRA-2026-003",
				Recommendation: "Cambiar credenciales admin en Jenkins, Grafana, Prometheus. Implementar LDAP/AD. Deshabilitar acceso local.",
				Priority:       "CRITICAL",
				Timeline:       "1 día",
				Owner:          "DevOps Lead",
				Status:         "Completado",
				EstimatedCost:  "$500 USD",
			},
			{
				ID:             "REM-004",
				FindingID:      "BCRA-2026-004",
				Recommendation: "Deshabilitar TLS 1.0/1.1. Requerir TLS 1.2+ con suites de cifrado modernas (ECDHE, AES-256)",
				Priority:       "HIGH",
				Timeline:       "3 días",
				Owner:          "Ing. Infraestructura",
				Status:         "Pendiente",
				EstimatedCost:  "$800 USD",
			},
			{
				ID:             "REM-005",
				FindingID:      "BCRA-2026-005",
				Recommendation: "Implementar WAF (AWS WAF o ModSecurity). Configurar reglas OWASP Core Rule Set v3.3",
				Priority:       "HIGH",
				Timeline:       "14 días",
				Owner:          "Ing. Security",
				Status:         "Pendiente",
				EstimatedCost:  "$5000 USD",
			},
		},

		AdditionalNotes: "Costo total estimado: $11,300 USD. Timeline crítico: 7 días. Se recomienda implementar programa de capacitación en OWASP Top 10 para equipo de desarrollo. Próxima evaluación: Marzo 2026.",
	}
}
