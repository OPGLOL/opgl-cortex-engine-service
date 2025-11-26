package services

import "github.com/OPGLOL/opgl-cortex-engine-service/internal/models"

// AnalysisServiceInterface defines the interface for analysis service operations
// This interface enables mocking in tests
type AnalysisServiceInterface interface {
	// AnalyzePlayer performs comprehensive analysis on a player's match history
	AnalyzePlayer(summoner *models.Summoner, matches []models.Match) *models.AnalysisResult
}
