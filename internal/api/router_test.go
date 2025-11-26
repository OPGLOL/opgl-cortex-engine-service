package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OPGLOL/opgl-cortex-engine-service/internal/models"
)

// TestSetupRouter tests that all routes are registered correctly
func TestSetupRouter(t *testing.T) {
	mockService := &MockAnalysisService{}
	handler := NewHandler(mockService)
	router := SetupRouter(handler)

	if router == nil {
		t.Fatal("Expected router to not be nil")
	}
}

// TestRouterHealthEndpoint tests that the health endpoint is registered
func TestRouterHealthEndpoint(t *testing.T) {
	mockService := &MockAnalysisService{}
	handler := NewHandler(mockService)
	router := SetupRouter(handler)

	request, _ := http.NewRequest("POST", "/health", nil)
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, responseRecorder.Code)
	}
}

// TestRouterHealthEndpointMethodNotAllowed tests that GET is not allowed for health
func TestRouterHealthEndpointMethodNotAllowed(t *testing.T) {
	mockService := &MockAnalysisService{}
	handler := NewHandler(mockService)
	router := SetupRouter(handler)

	request, _ := http.NewRequest("GET", "/health", nil)
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d for GET /health, got %d", http.StatusMethodNotAllowed, responseRecorder.Code)
	}
}

// TestRouterAnalyzeEndpoint tests that the analyze endpoint is registered
func TestRouterAnalyzeEndpoint(t *testing.T) {
	mockService := &MockAnalysisService{
		AnalyzePlayerFunc: func(summoner *models.Summoner, matches []models.Match) *models.AnalysisResult {
			return &models.AnalysisResult{}
		},
	}
	handler := NewHandler(mockService)
	router := SetupRouter(handler)

	// Send invalid JSON body to trigger BadRequest (proves endpoint is registered)
	request, _ := http.NewRequest("POST", "/api/v1/analyze", bytes.NewBufferString("invalid"))
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	// Endpoint should be reachable (returns 400 due to invalid body, not 404)
	if responseRecorder.Code == http.StatusNotFound {
		t.Error("Expected /api/v1/analyze endpoint to be registered")
	}
}

// TestRouterAnalyzeEndpointMethodNotAllowed tests that GET is not allowed for analyze
func TestRouterAnalyzeEndpointMethodNotAllowed(t *testing.T) {
	mockService := &MockAnalysisService{}
	handler := NewHandler(mockService)
	router := SetupRouter(handler)

	request, _ := http.NewRequest("GET", "/api/v1/analyze", nil)
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d for GET /api/v1/analyze, got %d", http.StatusMethodNotAllowed, responseRecorder.Code)
	}
}

// TestRouterNonExistentEndpoint tests that non-existent endpoints return 404
func TestRouterNonExistentEndpoint(t *testing.T) {
	mockService := &MockAnalysisService{}
	handler := NewHandler(mockService)
	router := SetupRouter(handler)

	request, _ := http.NewRequest("POST", "/api/v1/nonexistent", nil)
	responseRecorder := httptest.NewRecorder()

	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for non-existent endpoint, got %d", http.StatusNotFound, responseRecorder.Code)
	}
}

// TestRouterAllEndpointsUsePOST verifies all endpoints use POST method
func TestRouterAllEndpointsUsePOST(t *testing.T) {
	mockService := &MockAnalysisService{}
	handler := NewHandler(mockService)
	router := SetupRouter(handler)

	endpoints := []string{
		"/health",
		"/api/v1/analyze",
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			// Test GET method should not be allowed
			request, _ := http.NewRequest("GET", endpoint, nil)
			responseRecorder := httptest.NewRecorder()

			router.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected GET %s to return %d, got %d", endpoint, http.StatusMethodNotAllowed, responseRecorder.Code)
			}
		})
	}
}
