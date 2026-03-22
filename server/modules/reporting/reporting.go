package reporting

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// GenerateBCRAReport genera un PDF completo según BCRA A 8398/2026
func GenerateBCRAReport(report *BCRAReport) ([]byte, error) {
	// Calcular estadísticas
	report.CalculateStats()
	report.CalculateRiskScores()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)

	// 1. Portada
	addCoverPage(pdf, report)

	// 2. Índice
	addTableOfContents(pdf, report)

	// 3. Resumen Ejecutivo
	addExecutiveSummary(pdf, report)

	// 4. Metodología
	addMethodology(pdf, report)

	// 5. Matriz de Riesgos
	addRiskMatrix(pdf, report)

	// 6. Hallazgos Detallados
	addFindings(pdf, report)

	// 7. Recomendaciones
	addRecommendations(pdf, report)

	// 8. Firma y Validación
	addSignaturePage(pdf, report)

	// Convertir a bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// addCoverPage agrega la portada del reporte
func addCoverPage(pdf *gofpdf.Fpdf, report *BCRAReport) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 24)

	// Logo/Título
	pdf.Ln(40)
	pdf.Cell(0, 10, "PamperoC2")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 16)
	pdf.MultiCell(0, 10, "Reporte de Evaluación de Seguridad\nComunicación BCRA A 8398/2026", "", "C", false)

	// Datos de la institución
	pdf.SetFont("Arial", "", 11)
	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(50, 8, "Institución:")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, report.InstitutionName)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(50, 8, "Tipo:")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, report.InstitutionType)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(50, 8, "Período Evaluado:")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, report.AssessmentPeriod)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(50, 8, "Fecha de Reporte:")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, report.ReportDate.Format("02/01/2006"))
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(50, 8, "Válido hasta:")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, report.ReportValidUntil.Format("02/01/2006"))
	pdf.Ln(8)

	// Estadísticas de riesgo en portada
	pdf.Ln(20)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 8, "Nivel de Riesgo General:")
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 14)
	getRiskColor := func(risk string) (uint8, uint8, uint8) {
		switch risk {
		case "CRITICAL":
			return 255, 0, 0 // Rojo
		case "HIGH":
			return 255, 165, 0 // Naranja
		case "MEDIUM":
			return 255, 255, 0 // Amarillo
		default:
			return 0, 128, 0 // Verde
		}
	}
	r, g, b := getRiskColor(report.OverallRiskLevel)
	pdf.SetTextColor(int(r), int(g), int(b))
	pdf.Cell(0, 10, report.OverallRiskLevel)
	pdf.Ln(10)
	pdf.SetTextColor(0, 0, 0)

	// Resumen de hallazgos
	pdf.SetFont("Arial", "", 10)
	pdf.Ln(10)
	pdf.Cell(50, 7, fmt.Sprintf("Hallazgos Críticos: %d", report.CriticalFindings))
	pdf.Ln(7)
	pdf.Cell(50, 7, fmt.Sprintf("Hallazgos Altos: %d", report.HighFindings))
	pdf.Ln(7)
	pdf.Cell(50, 7, fmt.Sprintf("Hallazgos Medios: %d", report.MediumFindings))
	pdf.Ln(7)
	pdf.Cell(50, 7, fmt.Sprintf("Hallazgos Bajos: %d", report.LowFindings))
	pdf.Ln(7)

	// Footer
	pdf.SetFont("Arial", "I", 8)
	pdf.Ln(30)
	pdf.MultiCell(0, 5, "Este documento contiene información confidencial y debe tratarse como tal.", "", "C", false)
	pdf.Ln(2)
	pdf.MultiCell(0, 5, "Clasificación: CONFIDENCIAL - Uso Interno Únicamente", "", "C", false)
}

// addTableOfContents agrega el índice
func addTableOfContents(pdf *gofpdf.Fpdf, report *BCRAReport) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Índice")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Ln(5)

	contents := []string{
		"1. Resumen Ejecutivo",
		"2. Metodología de Evaluación",
		"3. Matriz de Riesgos",
		"4. Hallazgos de Seguridad",
		"5. Recomendaciones y Remediación",
		"6. Información de Contacto",
	}

	for _, content := range contents {
		pdf.Cell(0, 8, content)
		pdf.Ln(8)
	}
}

// addExecutiveSummary agrega el resumen ejecutivo
func addExecutiveSummary(pdf *gofpdf.Fpdf, report *BCRAReport) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "1. Resumen Ejecutivo")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.Ln(5)

	summary := generateExecutiveSummary(report)
	pdf.MultiCell(0, 5, summary, "", "L", false)

	// Tabla de estadísticas
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, "Estadísticas de Hallazgos:")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.Ln(3)

	statsTable := [][]string{
		{"Criticidad", "Cantidad", "Porcentaje"},
		{"CRÍTICO", fmt.Sprintf("%d", report.CriticalFindings), fmt.Sprintf("%.1f%%", calcPercentage(report.CriticalFindings, len(report.Findings)))},
		{"ALTO", fmt.Sprintf("%d", report.HighFindings), fmt.Sprintf("%.1f%%", calcPercentage(report.HighFindings, len(report.Findings)))},
		{"MEDIO", fmt.Sprintf("%d", report.MediumFindings), fmt.Sprintf("%.1f%%", calcPercentage(report.MediumFindings, len(report.Findings)))},
		{"BAJO", fmt.Sprintf("%d", report.LowFindings), fmt.Sprintf("%.1f%%", calcPercentage(report.LowFindings, len(report.Findings)))},
	}

	addTable(pdf, statsTable, []float64{60, 60, 60})
}

// addMethodology agrega la sección de metodología
func addMethodology(pdf *gofpdf.Fpdf, report *BCRAReport) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "2. Metodología de Evaluación")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Ln(5)

	pdf.MultiCell(0, 5, fmt.Sprintf("Metodología utilizada: %s", report.MethodologyUsed), "", "L", false)
	pdf.Ln(3)
	pdf.MultiCell(0, 5, fmt.Sprintf("Total de controles evaluados: %d", report.ControlsAssessed), "", "L", false)
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, "Alcance de la Evaluación:")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5,
		"• Revisión de controles técnicos de seguridad\n"+
			"• Evaluación de procesos de gestión de riesgos\n"+
			"• Análisis conformidad normativa BCRA A 8398/2026\n"+
			"• Pruebas de penetración limitadas\n"+
			"• Revisión de políticas y procedimientos", "", "L", false)
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 5,
		"• BCRA Comunicación A 8398/2026: Gestión de Riesgos de Tecnología y Seguridad de la Información\n"+
			"• NIST Cybersecurity Framework\n"+
			"• ISO/IEC 27001:2022 - Seguridad de la Información\n"+
			"• CIS Controls v8", "", "L", false)
}

// addRiskMatrix agrega la matriz de riesgos
func addRiskMatrix(pdf *gofpdf.Fpdf, report *BCRAReport) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "3. Matriz de Riesgos")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.Ln(5)

	risks := report.RiskMatrix
	if len(risks) == 0 {
		pdf.MultiCell(0, 5, "No hay riesgos registrados.", "", "L", false)
		return
	}

	// Tabla de riesgos
	riskTable := [][]string{
		{"ID", "Categoría", "Riesgo", "Prob.", "Imp.", "Score"},
	}

	for _, risk := range risks {
		riskTable = append(riskTable, []string{
			risk.ID,
			truncateText(risk.Category, 15),
			truncateText(risk.RiskName, 20),
			fmt.Sprintf("%d", risk.Probability),
			fmt.Sprintf("%d", risk.Impact),
			fmt.Sprintf("%d", risk.RiskScore),
		})
	}

	addTable(pdf, riskTable, []float64{15, 30, 40, 15, 15, 20})
}

// addFindings agrega la sección de hallazgos detallados
func addFindings(pdf *gofpdf.Fpdf, report *BCRAReport) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "4. Hallazgos de Seguridad")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.Ln(5)

	for i, finding := range report.Findings {
		// Encabezado del hallazgo
		pdf.SetFont("Arial", "B", 11)
		header := fmt.Sprintf("%d. [%s] %s", i+1, finding.Risk, finding.ID)
		pdf.Cell(0, 8, header)
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 10)
		pdf.Ln(2)

		// Detalles
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(30, 6, "Título:")
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 6, finding.Title, "", "L", false)

		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(30, 6, "Descripción:")
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 6, finding.Description, "", "L", false)

		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(30, 6, "Impacto:")
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 6, finding.Impact, "", "L", false)

		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(30, 6, "Sistemas:")
		pdf.SetFont("Arial", "", 10)
		systems := ""
		for _, sys := range finding.AffectedSystems {
			systems += sys + ", "
		}
		pdf.MultiCell(0, 6, truncateText(systems, 100), "", "L", false)

		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(30, 6, "Evidencia:")
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 6, finding.Evidence, "", "L", false)

		pdf.Ln(5)
		pdf.SetDrawColor(200, 200, 200)
		pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
		pdf.Ln(3)

		// Salto de página si es necesario
		if pdf.GetY() > 250 {
			pdf.AddPage()
		}
	}
}

// addRecommendations agrega las recomendaciones
func addRecommendations(pdf *gofpdf.Fpdf, report *BCRAReport) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "5. Recomendaciones y Plan de Remediación")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.Ln(5)

	if len(report.Remediations) == 0 {
		pdf.MultiCell(0, 5, "No hay recomendaciones registradas.", "", "L", false)
		return
	}

	remTable := [][]string{
		{"ID", "Hallazgo", "Recomendación", "Timeline", "Responsable"},
	}

	for _, rem := range report.Remediations {
		remTable = append(remTable, []string{
			rem.ID,
			rem.FindingID,
			truncateText(rem.Recommendation, 30),
			rem.Timeline,
			truncateText(rem.Owner, 20),
		})
	}

	addTable(pdf, remTable, []float64{15, 15, 50, 25, 35})
}

// addSignaturePage agrega la página de firma
func addSignaturePage(pdf *gofpdf.Fpdf, report *BCRAReport) {
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Información de Contacto y Firma")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 10)
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(50, 8, "Equipo de Evaluación:")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, report.AssessmentTeam)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(50, 8, "Correo Electrónico:")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, report.AssessmentTeamEmail)
	pdf.Ln(8)

	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(50, 8, "Fecha de Firma:")
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, time.Now().Format("02/01/2006 15:04"))
	pdf.Ln(8)

	pdf.Ln(20)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(80, 6, "_________________________")
	pdf.Ln(6)
	pdf.Cell(80, 6, "Firma del Evaluador")
	pdf.Ln(6)

	pdf.Ln(15)

	pdf.SetFont("Arial", "I", 8)
	pdf.MultiCell(0, 4,
		"Este documento es confidencial y ha sido preparado únicamente para uso interno.\n"+
			"No puede ser reproducido o distribuido sin autorización explícita.\n"+
			"Generado por PamperoC2 - Framework C2 Argentino\n"+
			fmt.Sprintf("Reporte ID: %s", report.ReportID), "", "C", false)
}

// Funciones auxiliares

// addTable agrega una tabla simple al PDF
func addTable(pdf *gofpdf.Fpdf, data [][]string, colWidths []float64) {
	cellHeight := 7.0

	// Encabezado
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(200, 200, 200)

	for i, col := range data[0] {
		if i < len(colWidths) {
			pdf.CellFormat(colWidths[i], cellHeight, col, "1", 1, "CM", true, 0, "")
		}
	}
	pdf.Ln(cellHeight)

	// Datos
	pdf.SetFont("Arial", "", 8)
	pdf.SetFillColor(255, 255, 255)

	for _, row := range data[1:] {
		for i, col := range row {
			if i < len(colWidths) {
				pdf.CellFormat(colWidths[i], cellHeight, truncateText(col, 15), "1", 0, "L", false, 0, "")
			}
		}
		pdf.Ln(cellHeight)
	}
}

// truncateText trunca texto a longitud máxima
func truncateText(text string, maxLen int) string {
	if len(text) > maxLen {
		return text[:maxLen] + "..."
	}
	return text
}

// calcPercentage calcula porcentaje
func calcPercentage(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) * 100 / float64(total)
}

// generateExecutiveSummary genera un resumen ejecutivo automático
func generateExecutiveSummary(report *BCRAReport) string {
	if report.ExecutiveSummary != "" {
		return report.ExecutiveSummary
	}

	// Generar resumen automático si no hay uno personalizado
	total := len(report.Findings)
	if total == 0 {
		return "Se completó la evaluación de seguridad sin hallazgos significativos."
	}

	summary := fmt.Sprintf(
		"Se ha completado una evaluación integral de seguridad de la información en %s durante el período %s.\n\n"+
			"La evaluación identificó un total de %d hallazgos de seguridad:\n"+
			"- %d hallazgos críticos que requieren atención inmediata\n"+
			"- %d hallazgos de alto riesgo\n"+
			"- %d hallazgos de riesgo medio\n"+
			"- %d hallazgos de bajo riesgo\n\n"+
			"El nivel de riesgo general se clasifica como %s. Se recomienda implementar las medidas correctivas especificadas en este reporte dentro de los plazos indicados para cada hallazgo.",
		report.InstitutionName,
		report.AssessmentPeriod,
		total,
		report.CriticalFindings,
		report.HighFindings,
		report.MediumFindings,
		report.LowFindings,
		report.OverallRiskLevel,
	)

	return summary
}
