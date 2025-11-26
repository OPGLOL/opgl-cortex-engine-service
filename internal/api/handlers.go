package api

import (
	"encoding/json"
	"net/http"

	"github.com/OPGLOL/opgl-cortex-engine-service/internal/models"
	"github.com/OPGLOL/opgl-cortex-engine-service/internal/services"
)

// Handler manages HTTP request handlers for the cortex engine
type Handler struct {
	analysisService services.AnalysisServiceInterface
}

// NewHandler creates a new Handler instance
func NewHandler(analysisService services.AnalysisServiceInterface) *Handler {
	return &Handler{
		analysisService: analysisService,
	}
}

// HealthCheck handles health check requests
func (handler *Handler) HealthCheck(writer http.ResponseWriter, request *http.Request) {
	response := map[string]string{
		"status":  "healthy",
		"service": "opgl-cortex-engine",
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(response)
}

// AnalyzePlayer handles player analysis requests
func (handler *Handler) AnalyzePlayer(writer http.ResponseWriter, request *http.Request) {
	var analyzeRequest struct {
		Summoner *models.Summoner `json:"summoner"`
		Matches  []models.Match   `json:"matches"`
	}

	if err := json.NewDecoder(request.Body).Decode(&analyzeRequest); err != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	if analyzeRequest.Summoner == nil {
		http.Error(writer, "Summoner data is required", http.StatusBadRequest)
		return
	}

	analysisResult := handler.analysisService.AnalyzePlayer(analyzeRequest.Summoner, analyzeRequest.Matches)

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(analysisResult)
}
