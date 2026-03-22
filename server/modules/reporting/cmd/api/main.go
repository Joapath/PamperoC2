package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/bishopfox/sliver/server/modules/reporting/api"
	"github.com/bishopfox/sliver/server/modules/reporting/storage"
)

func main() {
	dbPath := flag.String("db", "/tmp/pampero.db", "Ruta a base de datos SQLite")
	port := flag.String("port", "3000", "Puerto API")
	flag.Parse()

	fmt.Println("🇦🇷 PamperoC2 Dashboard API")
	fmt.Println("============================")

	// Inicializar BD
	fmt.Printf("📁 Inicializando BD: %s\n", *dbPath)
	if err := storage.Init(*dbPath); err != nil {
		log.Fatalf("❌ Error inicializando BD: %v", err)
	}
	fmt.Println("✅ BD inicializada")

	// Iniciar servidor
	fmt.Printf("🚀 Iniciando API en puerto %s\n", *port)
	fmt.Printf("📡 Health check: http://localhost:%s/health\n", *port)
	fmt.Printf("📊 API Base: http://localhost:%s/api/v1\n", *port)
	fmt.Println("Press CTRL+C para detener...")

	if err := api.StartServer(*port); err != nil {
		log.Fatalf("❌ Error servidor: %v", err)
	}
}
