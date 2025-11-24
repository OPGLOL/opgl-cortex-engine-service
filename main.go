package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/OPGLOL/opgl-cortex-engine/internal/api"
	"github.com/OPGLOL/opgl-cortex-engine/internal/services"
)

func main() {
	// Get port from environment variable (default: 8082)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Initialize analysis service
	analysisService := services.NewAnalysisService()

	// Initialize HTTP handler
	handler := api.NewHandler(analysisService)

	// Set up router
	router := api.SetupRouter(handler)

	// Start server
	serverAddress := fmt.Sprintf(":%s", port)
	log.Printf("OPGL Cortex Engine starting on port %s", port)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}
