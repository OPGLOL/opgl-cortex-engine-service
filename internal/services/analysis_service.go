package services

import (
	"math"
	"time"

	"github.com/OPGLOL/opgl-cortex-engine-service/internal/models"
)

// AnalysisService performs player performance analysis
type AnalysisService struct{}

// NewAnalysisService creates a new AnalysisService instance
func NewAnalysisService() *AnalysisService {
	return &AnalysisService{}
}

// AnalyzePlayer performs comprehensive analysis on a player's match history
func (analysisService *AnalysisService) AnalyzePlayer(summoner *models.Summoner, matches []models.Match) *models.AnalysisResult {
	playerStats := analysisService.calculatePlayerStats(summoner, matches)
	improvementAreas := analysisService.identifyImprovementAreas(&playerStats)

	return &models.AnalysisResult{
		PlayerStats:      playerStats,
		ImprovementAreas: improvementAreas,
		AnalyzedAt:       time.Now(),
	}
}

// calculatePlayerStats aggregates statistics from match history
func (analysisService *AnalysisService) calculatePlayerStats(summoner *models.Summoner, matches []models.Match) models.PlayerStats {
	if len(matches) == 0 {
		return models.PlayerStats{
			PUUID:        summoner.PUUID,
			SummonerName: summoner.Name,
		}
	}

	var totalKills int
	var totalDeaths int
	var totalAssists int
	var totalCS int
	var totalVisionScore int
	var totalDamage int
	var totalGold int
	var totalGameDuration int
	var wins int

	championPool := make(map[string]int)
	roleDistribution := make(map[string]int)

	// Aggregate stats from all matches
	for _, match := range matches {
		// Find the player's participation in this match
		for _, participant := range match.Participants {
			if participant.PUUID == summoner.PUUID {
				totalKills += participant.Kills
				totalDeaths += participant.Deaths
				totalAssists += participant.Assists
				totalCS += participant.TotalMinionsKilled
				totalVisionScore += participant.VisionScore
				totalDamage += participant.TotalDamageDealtToChampions
				totalGold += participant.GoldEarned
				totalGameDuration += match.GameDuration

				if participant.Win {
					wins++
				}

				// Track champion pool
				championPool[participant.ChampionName]++

				// Track role distribution
				if participant.TeamPosition != "" {
					roleDistribution[participant.TeamPosition]++
				}

				break
			}
		}
	}

	matchCount := len(matches)
	matchCountFloat := float64(matchCount)

	// Calculate averages
	averageKills := float64(totalKills) / matchCountFloat
	averageDeaths := float64(totalDeaths) / matchCountFloat
	averageAssists := float64(totalAssists) / matchCountFloat
	averageCS := float64(totalCS) / matchCountFloat
	averageVisionScore := float64(totalVisionScore) / matchCountFloat
	averageDamage := float64(totalDamage) / matchCountFloat
	averageGold := float64(totalGold) / matchCountFloat

	// Calculate KDA ratio
	kda := analysisService.calculateKDA(averageKills, averageDeaths, averageAssists)

	// Calculate CS per minute
	averageGameDurationMinutes := float64(totalGameDuration) / matchCountFloat / 60.0
	csPerMinute := averageCS / averageGameDurationMinutes

	// Calculate win rate
	winRate := (float64(wins) / matchCountFloat) * 100.0

	// Convert role distribution to percentages
	rolePercentages := make(map[string]float64)
	for role, count := range roleDistribution {
		rolePercentages[role] = (float64(count) / matchCountFloat) * 100.0
	}

	return models.PlayerStats{
		PUUID:              summoner.PUUID,
		SummonerName:       summoner.Name,
		TotalMatches:       matchCount,
		WinRate:            winRate,
		AverageKills:       averageKills,
		AverageDeaths:      averageDeaths,
		AverageAssists:     averageAssists,
		KDA:                kda,
		AverageCS:          averageCS,
		CSPerMinute:        csPerMinute,
		AverageVisionScore: averageVisionScore,
		AverageDamage:      averageDamage,
		AverageGold:        averageGold,
		ChampionPool:       championPool,
		RoleDistribution:   rolePercentages,
	}
}

// calculateKDA computes the Kill/Death/Assist ratio
func (analysisService *AnalysisService) calculateKDA(kills float64, deaths float64, assists float64) float64 {
	// Avoid division by zero
	if deaths == 0 {
		return kills + assists
	}
	return (kills + assists) / deaths
}

// identifyImprovementAreas analyzes stats and identifies areas for improvement
func (analysisService *AnalysisService) identifyImprovementAreas(playerStats *models.PlayerStats) []models.ImprovementArea {
	var improvementAreas []models.ImprovementArea

	// Benchmark values for average players (these can be adjusted based on rank)
	benchmarkCSPerMinute := 6.0
	benchmarkVisionScore := 40.0
	benchmarkKDA := 3.0
	benchmarkDeaths := 5.0

	// CS per minute analysis
	csGap := playerStats.CSPerMinute - benchmarkCSPerMinute
	if csGap < -1.0 {
		priority := "HIGH"
		if csGap > -2.0 {
			priority = "MEDIUM"
		}

		improvementAreas = append(improvementAreas, models.ImprovementArea{
			Category:       "CS (Creep Score)",
			CurrentValue:   math.Round(playerStats.CSPerMinute*10) / 10,
			ExpectedValue:  benchmarkCSPerMinute,
			Gap:            math.Round(csGap*10) / 10,
			Priority:       priority,
			Recommendation: "Focus on last-hitting minions more consistently. Practice farming in training mode and aim to maintain CS during mid-game teamfights.",
		})
	}

	// Vision score analysis
	visionGap := playerStats.AverageVisionScore - benchmarkVisionScore
	if visionGap < -10.0 {
		priority := "HIGH"
		if visionGap > -20.0 {
			priority = "MEDIUM"
		}

		improvementAreas = append(improvementAreas, models.ImprovementArea{
			Category:       "Vision Control",
			CurrentValue:   math.Round(playerStats.AverageVisionScore*10) / 10,
			ExpectedValue:  benchmarkVisionScore,
			Gap:            math.Round(visionGap*10) / 10,
			Priority:       priority,
			Recommendation: "Purchase more control wards and place wards in key objectives (Dragon, Baron). Clear enemy wards when possible to increase vision score.",
		})
	}

	// KDA analysis
	kdaGap := playerStats.KDA - benchmarkKDA
	if kdaGap < -0.5 {
		priority := "MEDIUM"
		if kdaGap < -1.0 {
			priority = "HIGH"
		}

		improvementAreas = append(improvementAreas, models.ImprovementArea{
			Category:       "KDA Ratio",
			CurrentValue:   math.Round(playerStats.KDA*100) / 100,
			ExpectedValue:  benchmarkKDA,
			Gap:            math.Round(kdaGap*100) / 100,
			Priority:       priority,
			Recommendation: "Focus on safer positioning in teamfights. Prioritize assists over risky kills and avoid unnecessary deaths.",
		})
	}

	// Deaths analysis
	deathsGap := playerStats.AverageDeaths - benchmarkDeaths
	if deathsGap > 1.0 {
		priority := "MEDIUM"
		if deathsGap > 2.0 {
			priority = "HIGH"
		}

		improvementAreas = append(improvementAreas, models.ImprovementArea{
			Category:       "Deaths",
			CurrentValue:   math.Round(playerStats.AverageDeaths*10) / 10,
			ExpectedValue:  benchmarkDeaths,
			Gap:            math.Round(deathsGap*10) / 10,
			Priority:       priority,
			Recommendation: "Review your deaths to identify patterns. Common causes: overextending without vision, poor positioning in fights, or staying too long with low HP.",
		})
	}

	// Win rate analysis
	if playerStats.WinRate < 45.0 {
		improvementAreas = append(improvementAreas, models.ImprovementArea{
			Category:       "Win Rate",
			CurrentValue:   math.Round(playerStats.WinRate*10) / 10,
			ExpectedValue:  50.0,
			Gap:            math.Round((playerStats.WinRate-50.0)*10) / 10,
			Priority:       "HIGH",
			Recommendation: "Focus on macro gameplay: objective control, wave management, and better decision-making in mid-late game. Consider your champion pool and role effectiveness.",
		})
	}

	// If no improvement areas found, add positive feedback
	if len(improvementAreas) == 0 {
		improvementAreas = append(improvementAreas, models.ImprovementArea{
			Category:       "Overall Performance",
			CurrentValue:   0,
			ExpectedValue:  0,
			Gap:            0,
			Priority:       "LOW",
			Recommendation: "Your performance is above average! Continue maintaining good CS, vision control, and KDA. Focus on consistency and adapting to different team compositions.",
		})
	}

	return improvementAreas
}
