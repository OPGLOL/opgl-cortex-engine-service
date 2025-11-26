package api

import (
	"github.com/gorilla/mux"
)

// SetupRouter configures all routes for the cortex engine
func SetupRouter(handler *Handler) *mux.Router {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", handler.HealthCheck).Methods("POST")

	// Analysis endpoint
	router.HandleFunc("/api/v1/analyze", handler.AnalyzePlayer).Methods("POST")

	return router
}
