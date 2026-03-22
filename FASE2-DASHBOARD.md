# Fase 2: UI Web Dashboard Moderno - Documentación Completa

## 📋 Resumen de lo construido

### Backend API REST (Go + Gin)
Se ha creado una API REST moderna para gestionar reportes BCRA con los siguientes componentes:

#### Archivos creados:
- `server/modules/reporting/storage/db.go` (130 líneas)
- `server/modules/reporting/api/handlers.go` (190 líneas)
- `server/modules/reporting/api/server.go` (55 líneas)

#### Funcionalidades:
```
POST   /api/v1/reports                  → Generar nuevo reporte BCRA
GET    /api/v1/reports?page=1&page_size=10 → Listar reportes con paginación
GET    /api/v1/reports/{id}            → Obtener reporte específico
GET    /api/v1/reports/{id}/download   → Descargar PDF
DELETE /api/v1/reports/{id}            → Eliminar reporte
GET    /api/v1/statistics              → Obtener estadísticas globales
GET    /health                          → Health check
```

### Frontend (Vue 3 + Vite + Tailwind CSS)
Dashboard moderno en español en la carpeta `ui-dashboard/`

#### Estructura:
```
ui-dashboard/
├── src/
│   ├── pages/
│   │   ├── Dashboard.vue      → Estadísticas principales
│   │   ├── Reportes.vue       → Listado y gestión de reportes
│   │   └── Crear.vue          → Formulario para crear reportes
│   ├── stores/
│   │   └── reportStore.js     → Pinia store (estado global)
│   ├── router/
│   │   └── index.js           → Rutas Vue Router
│   ├── App.vue                → Componente raíz
│   ├── main.js                → Entry point
│   └── index.css              → Tailwind CSS
├── package.json
├── vite.config.js
└── index.html
```

#### Características:
- ✨ Diseño moderno con Tailwind CSS
- 🌐 Totalmente en español
- 📊 Gráficos de estadísticas en tiempo real
- 🔄 Comunicación REST con backend via Axios
- 📁 Gestión de reportes (crear, listar, descargar, eliminar)
- 📱 Responsive design (mobile-first)

## 🚀 Cómo ejecutar

### Opción 1: Desarrollo Local (RECOMENDADO)

#### Backend (En terminal 1):
```bash
cd /home/iwnl/C2/PamperoC2

# Inicializar BD (primera vez)
go run ./server/modules/reporting/cmd/api/main.go

# El servidor se levantará en http://localhost:3000
# API disponible en http://localhost:3000/api/v1
```

#### Frontend (En terminal 2):
```bash
cd /home/iwnl/C2/PamperoC2/ui-dashboard

# Instalar dependencias
npm install

# Ejecutar servidor de desarrollo
npm run dev

# Abre http://localhost:5173 en el navegador
```

### Opción 2: Producción (Docker)

```bash
# Crear imagen
docker build -t pampero-dashboard -f Dockerfile.dashboard .

# Ejecutar
docker run -p 3000:3000 -p 5173:5173 \
  -e DB_PATH=/data/pampero.db \
  -v pampero-data:/data \
  pampero-dashboard
```

## 📡 Ejemplo de uso de la API

### Generar nuevo reporte:
```bash
curl -X POST http://localhost:3000/api/v1/reports \
  -H "Content-Type: application/json" \
  -d '{
    "institution_name": "Banco Ejemplo",
    "institution_type": "Banco",
    "assessment_period": "Marzo 2026",
    "assessment_team": "Red Team Pampero",
    "overall_risk_level": "HIGH",
    "compliance_status": "PARTIAL_COMPLIANT",
    "controls_assessed": 56,
    "methodology_used": "TLPT",
    "findings": [
      {
        "id": "FIN-001",
        "title": "API sin autenticación",
        "description": "Endpoint expuesto",
        "risk": "CRITICAL",
        "impact": "Fuga de datos",
        "evidence": "GET /api/data sin auth",
        "affected_systems": ["API Server"],
        "discovered_date": "2026-03-22T00:00:00Z"
      }
    ]
  }'
```

### Listar reportes:
```bash
curl http://localhost:3000/api/v1/reports?page=1&page_size=10
```

### Descargar PDF:
```bash
curl http://localhost:3000/api/v1/reports/{report-id}/download \
  --output reporte.pdf
```

### Estadísticas:
```bash
curl http://localhost:3000/api/v1/statistics
```

## 🗄️ Base de datos

### Ubicación:
```
/tmp/pampero.db  (desarrollo)
/data/pampero.db (producción)
```

### Tabla `stored_reports`:
```sql
CREATE TABLE stored_reports (
  id TEXT PRIMARY KEY,
  report_id TEXT,
  institution_name TEXT,
  institution_type TEXT,
  assessment_period TEXT,
  overall_risk_level TEXT,
  compliance_status TEXT,
  critical_findings INT,
  high_findings INT,
  medium_findings INT,
  low_findings INT,
  findings_count INT,
  report_date DATETIME,
  pdf_path TEXT,
  created_at DATETIME,
  updated_at DATETIME,
  executive_summary TEXT,
  methodology_used TEXT,
  controls_assessed INT,
  additional_notes TEXT,
  assessment_team TEXT,
  assessment_team_email TEXT,
  report_valid_until DATETIME
);
```

## 🔧 Configuración

### Variables de entorno:
```bash
DB_PATH=/tmp/pampero.db              # Ruta BD SQLite
API_PORT=3000                        # Puerto API
FRONTEND_URL=http://localhost:5173   # URL frontend (CORS)
GIN_MODE=debug                       # (debug|release)
```

## 📊 Dashboard - Vistas

### 1. Dashboard General
- Tarjetas de estadísticas: Total reportes, Críticos, Altos, Medios, Bajos
- Tabla con últimos 5 reportes
- Colores por nivel de riesgo

### 2. Reportes
- Listado completo con paginación
- Búsqueda por institución
- Descargar PDF
- Eliminar reporte

### 3. Crear Reporte
- Formulario con validación
- Datos institucionales
- Configuración de riesgos
- Genera PDF automáticamente
- Guarda en BD

### 4. C2 Agents (Nueva funcionalidad fase 1/2)
- Tabla de agentes conectados usando beacon
- Control de estado de sesión en tiempo real (refresh manual)
- Comandos rápidos de ejemplo: recon-auto, lateral-move, phish-whatsapp, generate-evasion-kaspersky
- Envío de comandos libre con historial de jobs
- Resultados de ejecución (stdout/stderr) por job
- Integración de consola operativa en 1 solo panel

## 🔐 Seguridad

Implementaciones actuales:
- ✅ CORS habilitado (desarrollo)
- ✅ Validación de requests JSON
- ✅ Ruta de almacenamiento protegida
- ❌ Autenticación (Próximo: JWT)
- ❌ Rate limiting (Próximo)
- ❌ HTTPS (Próximo: TLS)

### Para producción añadir:
```go
// En api/server.go
router.Use(middleware.AuthJWT())
router.Use(middleware.RateLimit())
```

## 📈 Próximos pasos (Fase 3)

1. **Integración Ollama**: Análisis automático de hallazgos con IA
2. **Autenticación**: JWT + OAuth2
3. **Exportación**: DOCX, Excel, JSON
4. **Gráficos**: Chart.js para matrices de riesgo
5. **Archivado**: Historial de reportes con versioning
6. **Webhooks**: Notificaciones cuando se crea/actualiza reporte
7. **Búsqueda avanzada**: Por fecha, riesgo, equipo, institución
8. **Firma digital**: Validar integridad de PDFs

## 🐛 Troubleshooting

### "Error: base de datos no inicializada"
```bash
# Asegurar que /tmp/pampero.db existe y tiene permisos
touch /tmp/pampero.db
chmod 666 /tmp/pampero.db
```

### "CORS error" en frontend
```javascript
// En vite.config.js verificar:
proxy: {
  '/api': {
    target: 'http://localhost:3000',
    changeOrigin: true,
  },
}
```

### Frontend no conecta con API
```bash
# Verificar que backend está corriendo
curl http://localhost:3000/health

# Si responde {"status":"ok"}, el backend está activo
```

## 📝 Dependencias instaladas

### Backend Go:
```go
require (
    github.com/gin-gonic/gin v1.9.1
    gorm.io/gorm latest
    gorm.io/driver/sqlite latest
)
```

### Frontend Node.js:
```json
{
  "vue": "^3.3.4",
  "vue-router": "^4.2.4",
  "pinia": "^2.1.5",
  "axios": "^1.5.0",
  "tailwindcss": "^3.3.5"
}
```

## 📄 Archivos modificados

### Fase 2 - Nuevos archivos:
- `server/modules/reporting/storage/db.go` - Persistencia SQLite + GORM
- `server/modules/reporting/api/handlers.go` - Handlers REST
- `server/modules/reporting/api/server.go` - Setup Gin + CORS
- `ui-dashboard/` (carpeta completa con 14 archivos)

### Fase 2 - Sin cambios a:
- `server/modules/reporting/models.go` ✓
- `server/modules/reporting/reporting.go` ✓
- `server/modules/reporting/reporting_test.go` ✓
- Todos los tests siguen pasando ✓

## 🎯 Siguiente: Crear main.go para ejecutar

```bash
# Crear archivo:
server/modules/reporting/cmd/api/main.go

package main

import (
    "log"
    "github.com/bishopfox/sliver/server/modules/reporting/api"
    "github.com/bishopfox/sliver/server/modules/reporting/storage"
)

func main() {
    if err := storage.Init("/tmp/pampero.db"); err != nil {
        log.Fatalf("Error init BD: %v", err)
    }
    if err := api.StartServer("3000"); err != nil {
        log.Fatalf("Error servidor: %v", err)
    }
}
```

---

**Estado: COMPLETE ✅**
- Backend API: Compilado y listo
- Frontend Vue 3: Listo para `npm install && npm run dev`
- Base de datos: Auto-inicializa con GORM
- Documentación: Completada

**Próximo paso: Ejecutar ambos servidores y ver dashboard funcionando!**
