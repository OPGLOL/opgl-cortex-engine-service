package services

import (
	"testing"

	"github.com/OPGLOL/opgl-cortex-engine-service/internal/models"
)

// TestNewAnalysisService tests the NewAnalysisService constructor
func TestNewAnalysisService(t *testing.T) {
	service := NewAnalysisService()

	if service == nil {
		t.Fatal("Expected service to not be nil")
	}
}

// TestAnalyzePlayer_EmptyMatches tests analysis with no matches
func TestAnalyzePlayer_EmptyMatches(t *testing.T) {
	service := NewAnalysisService()

	summoner := &models.Summoner{
		PUUID: "test-puuid",
		Name:  "TestPlayer",
	}
	matches := []models.Match{}

	result := service.AnalyzePlayer(summoner, matches)

	if result == nil {
		t.Fatal("Expected result to not be nil")
	}

	if result.PlayerStats.PUUID != summoner.PUUID {
		t.Errorf("Expected PUUID '%s', got '%s'", summoner.PUUID, result.PlayerStats.PUUID)
	}

	if result.PlayerStats.SummonerName != summoner.Name {
		t.Errorf("Expected SummonerName '%s', got '%s'", summoner.Name, result.PlayerStats.SummonerName)
	}

	if result.PlayerStats.TotalMatches != 0 {
		t.Errorf("Expected TotalMatches 0, got %d", result.PlayerStats.TotalMatches)
	}
}

// TestAnalyzePlayer_WithMatches tests analysis with match data
func TestAnalyzePlayer_WithMatches(t *testing.T) {
	service := NewAnalysisService()

	summoner := &models.Summoner{
		PUUID: "test-puuid",
		Name:  "TestPlayer",
	}

	matches := []models.Match{
		{
			MatchID:      "NA1_123",
			GameDuration: 1800, // 30 minutes
			Participants: []models.Participant{
				{
					PUUID:                       "test-puuid",
					ChampionName:                "Ahri",
					Kills:                       10,
					Deaths:                      2,
					Assists:                     8,
					TotalMinionsKilled:          200,
					VisionScore:                 50,
					TotalDamageDealtToChampions: 25000,
					GoldEarned:                  15000,
					Win:                         true,
					TeamPosition:                "MIDDLE",
				},
			},
		},
		{
			MatchID:      "NA1_124",
			GameDuration: 1500, // 25 minutes
			Participants: []models.Participant{
				{
					PUUID:                       "test-puuid",
					ChampionName:                "Zed",
					Kills:                       8,
					Deaths:                      4,
					Assists:                     5,
					TotalMinionsKilled:          180,
					VisionScore:                 35,
					TotalDamageDealtToChampions: 22000,
					GoldEarned:                  13000,
					Win:                         false,
					TeamPosition:                "MIDDLE",
				},
			},
		},
	}

	result := service.AnalyzePlayer(summoner, matches)

	if result == nil {
		t.Fatal("Expected result to not be nil")
	}

	if result.PlayerStats.TotalMatches != 2 {
		t.Errorf("Expected TotalMatches 2, got %d", result.PlayerStats.TotalMatches)
	}

	// Win rate should be 50% (1 win out of 2)
	expectedWinRate := 50.0
	if result.PlayerStats.WinRate != expectedWinRate {
		t.Errorf("Expected WinRate %.1f, got %.1f", expectedWinRate, result.PlayerStats.WinRate)
	}

	// Average kills should be 9 (10 + 8) / 2
	expectedAvgKills := 9.0
	if result.PlayerStats.AverageKills != expectedAvgKills {
		t.Errorf("Expected AverageKills %.1f, got %.1f", expectedAvgKills, result.PlayerStats.AverageKills)
	}

	// Champion pool should have 2 champions
	if len(result.PlayerStats.ChampionPool) != 2 {
		t.Errorf("Expected 2 champions in pool, got %d", len(result.PlayerStats.ChampionPool))
	}

	// Role distribution should be 100% MIDDLE
	if result.PlayerStats.RoleDistribution["MIDDLE"] != 100.0 {
		t.Errorf("Expected 100%% MIDDLE role, got %.1f%%", result.PlayerStats.RoleDistribution["MIDDLE"])
	}
}

// TestCalculateKDA tests the KDA calculation
func TestCalculateKDA(t *testing.T) {
	service := NewAnalysisService()

	testCases := []struct {
		name     string
		kills    float64
		deaths   float64
		assists  float64
		expected float64
	}{
		{"normal KDA", 10.0, 5.0, 10.0, 4.0},      // (10+10)/5 = 4.0
		{"perfect KDA (no deaths)", 10.0, 0.0, 5.0, 15.0}, // 10+5 = 15.0
		{"zero kills and assists", 0.0, 5.0, 0.0, 0.0},     // (0+0)/5 = 0
		{"equal KDA", 5.0, 5.0, 5.0, 2.0},         // (5+5)/5 = 2.0
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := service.calculateKDA(testCase.kills, testCase.deaths, testCase.assists)
			if result != testCase.expected {
				t.Errorf("Expected KDA %.2f, got %.2f", testCase.expected, result)
			}
		})
	}
}

// TestIdentifyImprovementAreas_LowCS tests improvement area detection for low CS
func TestIdentifyImprovementAreas_LowCS(t *testing.T) {
	service := NewAnalysisService()

	playerStats := &models.PlayerStats{
		CSPerMinute:        4.0, // Below 6.0 benchmark
		AverageVisionScore: 45.0,
		KDA:                3.5,
		AverageDeaths:      4.0,
		WinRate:            55.0,
	}

	areas := service.identifyImprovementAreas(playerStats)

	foundCSImprovement := false
	for _, area := range areas {
		if area.Category == "CS (Creep Score)" {
			foundCSImprovement = true
			if area.Priority != "HIGH" {
				t.Errorf("Expected CS improvement priority HIGH, got %s", area.Priority)
			}
			break
		}
	}

	if !foundCSImprovement {
		t.Error("Expected CS improvement area to be identified")
	}
}

// TestIdentifyImprovementAreas_LowVision tests improvement area detection for low vision
func TestIdentifyImprovementAreas_LowVision(t *testing.T) {
	service := NewAnalysisService()

	playerStats := &models.PlayerStats{
		CSPerMinute:        7.0,
		AverageVisionScore: 25.0, // Below 40.0 benchmark
		KDA:                3.5,
		AverageDeaths:      4.0,
		WinRate:            55.0,
	}

	areas := service.identifyImprovementAreas(playerStats)

	foundVisionImprovement := false
	for _, area := range areas {
		if area.Category == "Vision Control" {
			foundVisionImprovement = true
			if area.Priority != "MEDIUM" {
				t.Errorf("Expected Vision improvement priority MEDIUM, got %s", area.Priority)
			}
			break
		}
	}

	if !foundVisionImprovement {
		t.Error("Expected Vision improvement area to be identified")
	}
}

// TestIdentifyImprovementAreas_HighDeaths tests improvement area detection for high deaths
func TestIdentifyImprovementAreas_HighDeaths(t *testing.T) {
	service := NewAnalysisService()

	playerStats := &models.PlayerStats{
		CSPerMinute:        7.0,
		AverageVisionScore: 45.0,
		KDA:                3.5,
		AverageDeaths:      8.0, // Above 5.0 + 2.0 threshold
		WinRate:            55.0,
	}

	areas := service.identifyImprovementAreas(playerStats)

	foundDeathsImprovement := false
	for _, area := range areas {
		if area.Category == "Deaths" {
			foundDeathsImprovement = true
			if area.Priority != "HIGH" {
				t.Errorf("Expected Deaths improvement priority HIGH, got %s", area.Priority)
			}
			break
		}
	}

	if !foundDeathsImprovement {
		t.Error("Expected Deaths improvement area to be identified")
	}
}

// TestIdentifyImprovementAreas_LowWinRate tests improvement area detection for low win rate
func TestIdentifyImprovementAreas_LowWinRate(t *testing.T) {
	service := NewAnalysisService()

	playerStats := &models.PlayerStats{
		CSPerMinute:        7.0,
		AverageVisionScore: 45.0,
		KDA:                3.5,
		AverageDeaths:      4.0,
		WinRate:            40.0, // Below 45.0
	}

	areas := service.identifyImprovementAreas(playerStats)

	foundWinRateImprovement := false
	for _, area := range areas {
		if area.Category == "Win Rate" {
			foundWinRateImprovement = true
			if area.Priority != "HIGH" {
				t.Errorf("Expected Win Rate improvement priority HIGH, got %s", area.Priority)
			}
			break
		}
	}

	if !foundWinRateImprovement {
		t.Error("Expected Win Rate improvement area to be identified")
	}
}

// TestIdentifyImprovementAreas_NoImprovements tests when no improvements needed
func TestIdentifyImprovementAreas_NoImprovements(t *testing.T) {
	service := NewAnalysisService()

	playerStats := &models.PlayerStats{
		CSPerMinute:        7.0,  // Above 6.0
		AverageVisionScore: 45.0, // Above 30.0 (40-10)
		KDA:                4.0,  // Above 2.5 (3.0-0.5)
		AverageDeaths:      4.0,  // Below 6.0 (5.0+1.0)
		WinRate:            55.0, // Above 45.0
	}

	areas := service.identifyImprovementAreas(playerStats)

	// Should have at least one area (the positive feedback)
	if len(areas) == 0 {
		t.Error("Expected at least one improvement area (positive feedback)")
	}

	// Check for positive feedback
	foundPositive := false
	for _, area := range areas {
		if area.Category == "Overall Performance" && area.Priority == "LOW" {
			foundPositive = true
			break
		}
	}

	if !foundPositive {
		t.Error("Expected positive feedback when no improvements needed")
	}
}

// TestAnalysisServiceImplementsInterface verifies AnalysisService implements AnalysisServiceInterface
func TestAnalysisServiceImplementsInterface(t *testing.T) {
	var serviceInterface AnalysisServiceInterface = NewAnalysisService()

	if serviceInterface == nil {
		t.Error("AnalysisService should implement AnalysisServiceInterface")
	}
}

// TestAnalyzePlayer_MultipleChampions tests champion pool tracking
func TestAnalyzePlayer_MultipleChampions(t *testing.T) {
	service := NewAnalysisService()

	summoner := &models.Summoner{
		PUUID: "test-puuid",
		Name:  "TestPlayer",
	}

	matches := []models.Match{
		{
			MatchID:      "NA1_1",
			GameDuration: 1800,
			Participants: []models.Participant{
				{PUUID: "test-puuid", ChampionName: "Ahri", Win: true, TeamPosition: "MIDDLE"},
			},
		},
		{
			MatchID:      "NA1_2",
			GameDuration: 1800,
			Participants: []models.Participant{
				{PUUID: "test-puuid", ChampionName: "Ahri", Win: true, TeamPosition: "MIDDLE"},
			},
		},
		{
			MatchID:      "NA1_3",
			GameDuration: 1800,
			Participants: []models.Participant{
				{PUUID: "test-puuid", ChampionName: "Zed", Win: false, TeamPosition: "MIDDLE"},
			},
		},
	}

	result := service.AnalyzePlayer(summoner, matches)

	// Should have 2 unique champions
	if len(result.PlayerStats.ChampionPool) != 2 {
		t.Errorf("Expected 2 champions in pool, got %d", len(result.PlayerStats.ChampionPool))
	}

	// Ahri should have 2 games
	if result.PlayerStats.ChampionPool["Ahri"] != 2 {
		t.Errorf("Expected Ahri to have 2 games, got %d", result.PlayerStats.ChampionPool["Ahri"])
	}

	// Zed should have 1 game
	if result.PlayerStats.ChampionPool["Zed"] != 1 {
		t.Errorf("Expected Zed to have 1 game, got %d", result.PlayerStats.ChampionPool["Zed"])
	}
}

// TestAnalyzePlayer_MixedRoles tests role distribution tracking
func TestAnalyzePlayer_MixedRoles(t *testing.T) {
	service := NewAnalysisService()

	summoner := &models.Summoner{
		PUUID: "test-puuid",
		Name:  "TestPlayer",
	}

	matches := []models.Match{
		{
			MatchID:      "NA1_1",
			GameDuration: 1800,
			Participants: []models.Participant{
				{PUUID: "test-puuid", ChampionName: "Ahri", TeamPosition: "MIDDLE"},
			},
		},
		{
			MatchID:      "NA1_2",
			GameDuration: 1800,
			Participants: []models.Participant{
				{PUUID: "test-puuid", ChampionName: "Ashe", TeamPosition: "BOTTOM"},
			},
		},
		{
			MatchID:      "NA1_3",
			GameDuration: 1800,
			Participants: []models.Participant{
				{PUUID: "test-puuid", ChampionName: "LeeSin", TeamPosition: "JUNGLE"},
			},
		},
		{
			MatchID:      "NA1_4",
			GameDuration: 1800,
			Participants: []models.Participant{
				{PUUID: "test-puuid", ChampionName: "Zed", TeamPosition: "MIDDLE"},
			},
		},
	}

	result := service.AnalyzePlayer(summoner, matches)

	// MIDDLE should be 50% (2/4)
	expectedMiddle := 50.0
	if result.PlayerStats.RoleDistribution["MIDDLE"] != expectedMiddle {
		t.Errorf("Expected MIDDLE %.1f%%, got %.1f%%", expectedMiddle, result.PlayerStats.RoleDistribution["MIDDLE"])
	}

	// BOTTOM should be 25% (1/4)
	expectedBottom := 25.0
	if result.PlayerStats.RoleDistribution["BOTTOM"] != expectedBottom {
		t.Errorf("Expected BOTTOM %.1f%%, got %.1f%%", expectedBottom, result.PlayerStats.RoleDistribution["BOTTOM"])
	}

	// JUNGLE should be 25% (1/4)
	expectedJungle := 25.0
	if result.PlayerStats.RoleDistribution["JUNGLE"] != expectedJungle {
		t.Errorf("Expected JUNGLE %.1f%%, got %.1f%%", expectedJungle, result.PlayerStats.RoleDistribution["JUNGLE"])
	}
}
