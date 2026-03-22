# Módulo Reporting - PamperoC2

Generador automático de reportes de seguridad en PDF conforme a **BCRA A 8398/2026**.

## Características

✅ Generación de PDF profesionales en formato BCRA  
✅ Compatibilidad con estándares internacionales (NIST, ISO 27001, CIS)  
✅ Cálculo automático de matrices de riesgo  
✅ Soporte para múltiples niveles de severidad  
✅ Tests unitarios incluidos  
✅ Bajo consumo de recursos (sin dependencias externas pesadas)  

## Uso Rápido

```go
package main

import (
	"os"
	"time"
	"github.com/bishopfox/sliver/server/modules/reporting"
)

func main() {
	// Crear reporte básico
	report := &reporting.BCRAReport{
		ReportID:         "BCRA-2026-001",
		InstitutionName:  "Mi Banco",
		InstitutionType:  "Banco",
		AssessmentPeriod: "Enero 2026",
		ReportDate:       time.Now(),
		ReportValidUntil: time.Now().AddDate(0, 0, 365),
		OverallRiskLevel: "HIGH",
		ComplianceStatus: "PARTIAL_COMPLIANT",
		MethodologyUsed:  "TLPT + NIST",
		ControlsAssessed: 42,
		Findings: []reporting.Finding{
			{
				ID:          "BCRA-001",
				Title:       "RDP Expuesto",
				Description: "RDP sin MFA",
				Risk:        "CRITICAL",
				Impact:      "Acceso directo a sistemas",
				Evidence:    "Puerto 3389 abierto",
				AffectedSystems: []string{"SRV-01"},
			},
		},
	}

	// Generar PDF
	pdf, _ := reporting.GenerateBCRAReport(report)
	os.WriteFile("reporte.pdf", pdf, 0644)
}
```

## Estructura

```
server/modules/reporting/
├── models.go              # Estructuras de datos
├── reporting.go           # Lógica de generación PDF
├── reporting_test.go      # Tests unitarios
└── README.md              # Este archivo
```

## Testing

```bash
# Ejecutar todos los tests
go test -v ./server/modules/reporting

# Genera un PDF de prueba en /tmp/
# Recomendado para validación manual
```

## Dependencias

- `github.com/jung-kurt/gofpdf/v2` - Generación de PDF

## API Reference

### GenerateBCRAReport
Genera PDF del reporte BCRA.

```go
func GenerateBCRAReport(report *BCRAReport) ([]byte, error)
```

### BCRAReport
Estructura principal del reporte.

```go
type BCRAReport struct {
	ReportID         string
	InstitutionName  string
	Findings         []Finding
	RiskMatrix       []RiskItem
	Remediations     []RemediationItem
	// ... más campos
}
```

## Contacto

PamperoC2 Red Team  
redteam@pampero.ar

## Licencia

MIT
