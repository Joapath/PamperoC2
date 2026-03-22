# 🇦🇷 PamperoC2 - Sistema de Reportes BCRA con IA

> **Sistema completo de evaluación de seguridad financiera con análisis inteligente impulsado por IA local**

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![Vue.js](https://img.shields.io/badge/Vue.js-3.3+-green.svg)](https://vuejs.org/)
[![Ollama](https://img.shields.io/badge/Ollama-Mistral-orange.svg)](https://ollama.ai/)
[![License](https://img.shields.io/badge/License-GPLv3-red.svg)](LICENSE)

## 📋 Descripción

PamperoC2 es un sistema integral para la generación automática de reportes de cumplimiento BCRA (Banco Central de la República Argentina) con análisis inteligente de hallazgos de seguridad utilizando IA local. El sistema combina un backend robusto en Go con una interfaz web moderna en Vue.js, ofreciendo una solución completa para evaluaciones de seguridad financiera.

## ✨ Características Principales

### 🔒 **Módulo de Reportes BCRA**
- Generación automática de reportes PDF conforme a la normativa A 8398/2026
- Estructura completa de evaluación de seguridad financiera
- Matriz de riesgos integrada
- Sistema de remediaciones y seguimiento

### 🤖 **Análisis Inteligente con IA**
- Integración con Ollama para análisis local de hallazgos
- Identificación automática de vectores de ataque
- Generación de comandos técnicos específicos
- Sistema de prioridades y niveles de confianza
- Análisis no-bloqueante (funciona sin IA disponible)

### 🖥️ **Dashboard Web Moderno**
- Interfaz responsive con Vue.js 3 + Composition API
- Gestión completa de reportes (CRUD)
- Visualización de estadísticas en tiempo real
- Página dedicada para análisis IA vs hallazgos originales
- Tema profesional con Tailwind CSS

### 🏗️ **Arquitectura Robusta**
- Backend en Go con framework Gin
- Base de datos SQLite con migraciones automáticas
- API REST completa con documentación
- Arquitectura modular y escalable
- Sistema de logging y manejo de errores

## 🚀 Inicio Rápido

### Prerrequisitos

- **Go 1.21+** - [Descargar](https://golang.org/dl/)
- **Node.js 18+** - [Descargar](https://nodejs.org/)
- **Ollama** - [Instalar](https://ollama.ai/download)

### Instalación

1. **Clonar el repositorio**
   ```bash
   git clone https://github.com/Joapath/PamperoC2.git
   cd PamperoC2
   ```

2. **Instalar dependencias del backend**
   ```bash
   go mod download
   go mod vendor
   ```

3. **Instalar dependencias del frontend**
   ```bash
   cd ui-dashboard
   npm install
   cd ..
   ```

4. **Configurar Ollama**
   ```bash
   ollama pull mistral
   ```

### Ejecución

1. **Levantar el backend**
   ```bash
   cd server/modules/reporting/cmd/api
   go run main.go -db /tmp/pampero.db -port 3000
   ```

2. **Levantar el frontend** (en otra terminal)
   ```bash
   cd ui-dashboard
   npm run dev
   ```

3. **Acceder a la aplicación**
   - Dashboard: http://localhost:5173
   - API: http://localhost:3000
   - Health Check: http://localhost:3000/health

### Demo Automático

```bash
# Ejecutar demo completo
./run-demo.sh
```

## 📁 Estructura del Proyecto

```
PamperoC2/
├── server/modules/reporting/     # Backend Go
│   ├── ai/                       # Servicio Ollama
│   ├── api/                      # Handlers HTTP y servidor
│   ├── cmd/api/                  # Punto de entrada
│   ├── models.go                 # Estructuras de datos
│   ├── storage/                  # Base de datos SQLite
│   └── reporting.go              # Lógica de reportes
├── ui-dashboard/                 # Frontend Vue.js
│   ├── src/
│   │   ├── pages/                # Vistas principales
│   │   ├── stores/               # Estado Pinia
│   │   └── router/               # Configuración de rutas
│   └── package.json
├── docs/                         # Documentación
├── examples/                     # Ejemplos de uso
├── Dockerfile                    # Contenedor Docker
└── run-demo.sh                   # Script de demostración
```

## 🎯 Funcionalidades

### Reportes BCRA
- ✅ Generación automática de PDFs
- ✅ Conformidad con normativa A 8398/2026
- ✅ Matriz de riesgos integrada
- ✅ Sistema de hallazgos y remediaciones

### Análisis IA
- ✅ Integración con Ollama (Mistral)
- ✅ Análisis local (sin dependencias externas)
- ✅ Vectores de ataque identificados
- ✅ Comandos técnicos específicos
- ✅ Sistema de prioridades y confianza

### Dashboard Web
- ✅ Interfaz responsive y moderna
- ✅ Gestión completa de reportes
- ✅ Estadísticas en tiempo real
- ✅ Visualización de análisis IA
- ✅ Tema profesional

## 🔧 API Endpoints

### Reportes
- `POST /api/v1/reports` - Crear reporte
- `GET /api/v1/reports` - Listar reportes
- `GET /api/v1/reports/:id` - Obtener reporte
- `DELETE /api/v1/reports/:id` - Eliminar reporte
- `GET /api/v1/reports/:id/download` - Descargar PDF

### Análisis IA
- `GET /api/v1/reports/:id/ai-analysis` - Obtener análisis IA
- `POST /api/v1/reports/:id/reanalyze` - Re-analizar con IA

### Estadísticas
- `GET /api/v1/statistics` - Estadísticas generales
- `GET /health` - Health check

## 🧪 Testing

```bash
# Tests del backend
go test ./server/modules/reporting/...

# Tests del servicio IA
go test ./server/modules/reporting/ai/...

# Compilación del frontend
cd ui-dashboard && npm run build
```

## 🐳 Docker

```bash
# Construir imagen
docker build -t pampero-c2 .

# Ejecutar contenedor
docker run -p 3000:3000 -p 5173:5173 pampero-c2
```

## 📊 Arquitectura

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Frontend      │────│   API REST       │────│   Ollama AI     │
│   Vue.js 3      │    │   Gin Framework  │    │   (Mistral)     │
│                 │    │                  │    │                 │
│ • Dashboard     │    │ • Reportes       │    │ • Análisis      │
│ • Gestión       │    │ • Estadísticas   │    │ • Vectores      │
│ • Análisis IA   │    │ • IA Analysis    │    │ • Comandos      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   SQLite DB     │    │   PDF Reports    │    │   JSON Analysis │
│                 │    │                  │    │                 │
│ • Reportes      │    │ • BCRA Format    │    │ • Attack Vectors│
│ • Análisis IA   │    │ • Compliance     │    │ • Commands      │
│ • Estadísticas  │    │ • Downloadable   │    │ • Priority      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📝 Licencia

Este proyecto está bajo la Licencia GPL v3. Ver el archivo [LICENSE](LICENSE) para más detalles.

## 👥 Autores

- **Equipo PamperoC2** - Desarrollo inicial
- **Contribuciones** - Ver [CONTRIBUTORS](CONTRIBUTORS.md)

## 🙏 Agradecimientos

- [Bishop Fox](https://bishopfox.com/) - Por el framework Sliver base
- [Ollama](https://ollama.ai/) - Por el motor de IA local
- [Vue.js](https://vuejs.org/) - Por el framework frontend
- [Gin Web Framework](https://gin-gonic.com/) - Por el framework backend

## 📞 Contacto

- **Email**: team@pamperoc2.dev
- **GitHub**: [Joapath/PamperoC2](https://github.com/Joapath/PamperoC2)
- **Issues**: [Reportar problema](https://github.com/Joapath/PamperoC2/issues)

---

**🇦🇷 Hecho con ❤️ en Argentina para la seguridad financiera**
