package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OPGLOL/opgl-cortex-engine-service/internal/models"
)

// MockAnalysisService is a mock implementation of AnalysisServiceInterface for testing
type MockAnalysisService struct {
	AnalyzePlayerFunc func(summoner *models.Summoner, matches []models.Match) *models.AnalysisResult
}

func (m *MockAnalysisService) AnalyzePlayer(summoner *models.Summoner, matches []models.Match) *models.AnalysisResult {
	if m.AnalyzePlayerFunc != nil {
		return m.AnalyzePlayerFunc(summoner, matches)
	}
	return nil
}

// TestNewHandler tests the NewHandler constructor
func TestNewHandler(t *testing.T) {
	mockService := &MockAnalysisService{}
	handler := NewHandler(mockService)

	if handler == nil {
		t.Fatal("Expected handler to not be nil")
	}

	if handler.analysisService != mockService {
		t.Error("Expected analysisService to be set correctly")
	}
}

// TestHealthCheck tests the health check endpoint
func TestHealthCheck(t *testing.T) {
	handler := &Handler{analysisService: nil}

	request, err := http.NewRequest("POST", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	responseRecorder := httptest.NewRecorder()
	handler.HealthCheck(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, responseRecorder.Code)
	}

	var response map[string]string
	err = json.NewDecoder(responseRecorder.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}

	if response["service"] != "opgl-cortex-engine" {
		t.Errorf("Expected service 'opgl-cortex-engine', got '%s'", response["service"])
	}
}

// TestHealthCheckContentType tests that health check returns JSON content type
func TestHealthCheckContentType(t *testing.T) {
	handler := &Handler{analysisService: nil}

	request, err := http.NewRequest("POST", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	responseRecorder := httptest.NewRecorder()
	handler.HealthCheck(responseRecorder, request)

	contentType := responseRecorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

// TestAnalyzePlayer_Success tests successful player analysis
func TestAnalyzePlayer_Success(t *testing.T) {
	expectedResult := &models.AnalysisResult{
		PlayerStats: models.PlayerStats{
			PUUID:        "test-puuid",
			SummonerName: "TestPlayer",
			TotalMatches: 10,
			WinRate:      55.0,
		},
		ImprovementAreas: []models.ImprovementArea{
			{
				Category:       "CS (Creep Score)",
				Priority:       "MEDIUM",
				Recommendation: "Focus on last-hitting",
			},
		},
		AnalyzedAt: time.Now(),
	}

	mockService := &MockAnalysisService{
		AnalyzePlayerFunc: func(summoner *models.Summoner, matches []models.Match) *models.AnalysisResult {
			if summoner.PUUID != "test-puuid" {
				t.Errorf("Expected PUUID 'test-puuid', got '%s'", summoner.PUUID)
			}
			return expectedResult
		},
	}

	handler := NewHandler(mockService)

	requestBody := map[string]interface{}{
		"summoner": map[string]interface{}{
			"puuid": "test-puuid",
			"name":  "TestPlayer",
		},
		"matches": []map[string]interface{}{
			{"matchId": "NA1_123"},
		},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	request, err := http.NewRequest("POST", "/api/v1/analyze", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	responseRecorder := httptest.NewRecorder()
	handler.AnalyzePlayer(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, responseRecorder.Code)
	}

	var response models.AnalysisResult
	err = json.NewDecoder(responseRecorder.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.PlayerStats.PUUID != expectedResult.PlayerStats.PUUID {
		t.Errorf("Expected PUUID '%s', got '%s'", expectedResult.PlayerStats.PUUID, response.PlayerStats.PUUID)
	}
}

// TestAnalyzePlayer_InvalidJSON tests invalid JSON request body
func TestAnalyzePlayer_InvalidJSON(t *testing.T) {
	handler := NewHandler(&MockAnalysisService{})

	request, err := http.NewRequest("POST", "/api/v1/analyze", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	responseRecorder := httptest.NewRecorder()
	handler.AnalyzePlayer(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
}

// TestAnalyzePlayer_MissingSummoner tests missing summoner data
func TestAnalyzePlayer_MissingSummoner(t *testing.T) {
	handler := NewHandler(&MockAnalysisService{})

	requestBody := map[string]interface{}{
		"matches": []map[string]interface{}{
			{"matchId": "NA1_123"},
		},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	request, _ := http.NewRequest("POST", "/api/v1/analyze", bytes.NewBuffer(bodyBytes))
	request.Header.Set("Content-Type", "application/json")

	responseRecorder := httptest.NewRecorder()
	handler.AnalyzePlayer(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
}

// TestAnalyzePlayer_NullSummoner tests null summoner value
func TestAnalyzePlayer_NullSummoner(t *testing.T) {
	handler := NewHandler(&MockAnalysisService{})

	requestBody := map[string]interface{}{
		"summoner": nil,
		"matches":  []map[string]interface{}{},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	request, _ := http.NewRequest("POST", "/api/v1/analyze", bytes.NewBuffer(bodyBytes))
	request.Header.Set("Content-Type", "application/json")

	responseRecorder := httptest.NewRecorder()
	handler.AnalyzePlayer(responseRecorder, request)

	if responseRecorder.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, responseRecorder.Code)
	}
}

// TestAnalyzePlayer_EmptyMatches tests analysis with no matches
func TestAnalyzePlayer_EmptyMatches(t *testing.T) {
	expectedResult := &models.AnalysisResult{
		PlayerStats: models.PlayerStats{
			PUUID:        "test-puuid",
			SummonerName: "TestPlayer",
		},
		ImprovementAreas: []models.ImprovementArea{},
		AnalyzedAt:       time.Now(),
	}

	mockService := &MockAnalysisService{
		AnalyzePlayerFunc: func(summoner *models.Summoner, matches []models.Match) *models.AnalysisResult {
			if len(matches) != 0 {
				t.Errorf("Expected 0 matches, got %d", len(matches))
			}
			return expectedResult
		},
	}

	handler := NewHandler(mockService)

	requestBody := map[string]interface{}{
		"summoner": map[string]interface{}{
			"puuid": "test-puuid",
			"name":  "TestPlayer",
		},
		"matches": []map[string]interface{}{},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	request, _ := http.NewRequest("POST", "/api/v1/analyze", bytes.NewBuffer(bodyBytes))
	request.Header.Set("Content-Type", "application/json")

	responseRecorder := httptest.NewRecorder()
	handler.AnalyzePlayer(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, responseRecorder.Code)
	}
}

// TestAnalyzePlayer_ContentType tests that analyze returns JSON content type
func TestAnalyzePlayer_ContentType(t *testing.T) {
	mockService := &MockAnalysisService{
		AnalyzePlayerFunc: func(summoner *models.Summoner, matches []models.Match) *models.AnalysisResult {
			return &models.AnalysisResult{}
		},
	}

	handler := NewHandler(mockService)

	requestBody := map[string]interface{}{
		"summoner": map[string]interface{}{
			"puuid": "test-puuid",
		},
		"matches": []map[string]interface{}{},
	}
	bodyBytes, _ := json.Marshal(requestBody)

	request, _ := http.NewRequest("POST", "/api/v1/analyze", bytes.NewBuffer(bodyBytes))
	request.Header.Set("Content-Type", "application/json")

	responseRecorder := httptest.NewRecorder()
	handler.AnalyzePlayer(responseRecorder, request)

	contentType := responseRecorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}
