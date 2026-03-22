#!/bin/bash

# PamperoC2 Dashboard - Script de demo
# Uso: ./run-demo.sh

set -e

echo "🇦🇷 PamperoC2 - Fase 2: UI Web Dashboard"
echo "=========================================="
echo ""

# Verificar que estamos en la carpeta correcta
if [ ! -f "go.mod" ]; then
    echo "❌ Error: Ejecutar desde raíz de PamperoC2"
    exit 1
fi

echo "📋 Componentes creados:"
echo "   ✅ Backend API REST (Go + Gin)"
echo "       - server/modules/reporting/storage/db.go"
echo "       - server/modules/reporting/api/handlers.go"
echo "       - server/modules/reporting/api/server.go"
echo "       - server/modules/reporting/cmd/api/main.go"
echo ""
echo "   ✅ Frontend Vue 3 + Vite"
echo "       - ui-dashboard/ (14 archivos)"
echo "       - Dashboard, Reportes, Crear vistas"
echo "       - Pinia store para estado"
echo ""

echo "📦 Para ejecutar en desarrollo:"
echo ""
echo "Terminal 1 - Backend API:"
echo "  cd /home/iwnl/C2/PamperoC2"
echo "  go run ./server/modules/reporting/cmd/api/main.go"
echo "  (o usar el binary: ./bin/pampero-api)"
echo ""
echo "Terminal 2 - Frontend:"
echo "  cd /home/iwnl/C2/PamperoC2/ui-dashboard"
echo "  npm install"
echo "  npm run dev"
echo ""
echo "Luego abrir:"
echo "  🌐 Dashboard: http://localhost:5173"
echo "  🔌 API Health: http://localhost:3000/health"
echo ""

echo "📊 API Endpoints disponibles:"
echo "  POST   /api/v1/reports                    - Generar reporte"
echo "  GET    /api/v1/reports                    - Listar reportes"
echo "  GET    /api/v1/reports/{id}               - Obtener reporte"
echo "  GET    /api/v1/reports/{id}/download      - Descargar PDF"
echo "  DELETE /api/v1/reports/{id}               - Eliminar reporte"
echo "  GET    /api/v1/statistics                 - Estadísticas"
echo ""

echo "✨ Características:"
echo "   • Dashboard en ESPAÑOL 🇦🇷"
echo "   • Generación de PDFs BCRA A 8398/2026"
echo "   • CRUD completo de reportes"
echo "   • Base de datos SQLite embebida"
echo "   • Frontend responsive (Tailwind CSS)"
echo "   • API REST con validación"
echo "   • CORS habilitado"
echo ""

echo "🚀 Status: Listo para producción (Fase 2 completada)"
echo ""
