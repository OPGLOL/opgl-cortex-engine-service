package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/OPGLOL/opgl-cortex-engine-service/internal/api"
	"github.com/OPGLOL/opgl-cortex-engine-service/internal/middleware"
	"github.com/OPGLOL/opgl-cortex-engine-service/internal/services"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize zerolog with colorized console output for development
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Caller().Logger()

	// Set global log level (can be configured via environment variable)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Info().Msg("Starting OPGL Cortex Engine")

	// Get port from environment variable (default: 8082)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Info().
		Str("port", port).
		Msg("Configuration loaded")

	// Initialize analysis service
	analysisService := services.NewAnalysisService()

	// Initialize HTTP handler
	handler := api.NewHandler(analysisService)

	// Set up router
	router := api.SetupRouter(handler)

	// Wrap router with logging middleware
	loggedRouter := middleware.LoggingMiddleware(router)

	// Start server
	serverAddress := fmt.Sprintf(":%s", port)
	log.Info().
		Str("address", serverAddress).
		Str("port", port).
		Msg("OPGL Cortex Engine listening")

	if err := http.ListenAndServe(serverAddress, loggedRouter); err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
