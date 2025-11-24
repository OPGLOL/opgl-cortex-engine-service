package models

import "time"

// Summoner represents a League of Legends player account
type Summoner struct {
	// Encrypted summoner ID returned by Riot API
	ID string `json:"id"`
	// Encrypted account ID
	AccountID string `json:"accountId"`
	// Encrypted PUUID (Player Universally Unique IDentifier)
	PUUID string `json:"puuid"`
	// Summoner name visible in game
	Name string `json:"name"`
	// Profile icon ID number
	ProfileIconID int `json:"profileIconId"`
	// Summoner level (non-ranked progression)
	SummonerLevel int64 `json:"summonerLevel"`
}

// Match represents a single League of Legends match
type Match struct {
	// Unique match identifier
	MatchID string `json:"matchId"`
	// Timestamp when the match started
	GameCreation time.Time `json:"gameCreation"`
	// Total duration of the match in seconds
	GameDuration int `json:"gameDuration"`
	// Game mode (e.g., CLASSIC, ARAM)
	GameMode string `json:"gameMode"`
	// Game type (e.g., MATCHED_GAME)
	GameType string `json:"gameType"`
	// List of all participants in the match
	Participants []Participant `json:"participants"`
}

// Participant represents a player's performance in a specific match
type Participant struct {
	// Player's PUUID
	PUUID string `json:"puuid"`
	// Summoner name at the time of the match
	SummonerName string `json:"summonerName"`
	// Champion ID played in this match
	ChampionID int `json:"championId"`
	// Champion name for easier reference
	ChampionName string `json:"championName"`
	// Number of enemy champions killed
	Kills int `json:"kills"`
	// Number of times the player died
	Deaths int `json:"deaths"`
	// Number of assists in killing enemy champions
	Assists int `json:"assists"`
	// Total gold earned during the match
	GoldEarned int `json:"goldEarned"`
	// Total damage dealt to champions
	TotalDamageDealtToChampions int `json:"totalDamageDealtToChampions"`
	// Total damage taken from all sources
	TotalDamageTaken int `json:"totalDamageTaken"`
	// Vision score (wards placed, destroyed, etc.)
	VisionScore int `json:"visionScore"`
	// Creep score (minions and monsters killed)
	TotalMinionsKilled int `json:"totalMinionsKilled"`
	// Whether the player's team won the match
	Win bool `json:"win"`
	// Player's role in the match (TOP, JUNGLE, MID, BOT, SUPPORT)
	TeamPosition string `json:"teamPosition"`
}

// PlayerStats represents aggregated statistics for a player
type PlayerStats struct {
	// Player's PUUID
	PUUID string `json:"puuid"`
	// Summoner name
	SummonerName string `json:"summonerName"`
	// Total number of matches analyzed
	TotalMatches int `json:"totalMatches"`
	// Overall win rate as a percentage
	WinRate float64 `json:"winRate"`
	// Average kills per game
	AverageKills float64 `json:"averageKills"`
	// Average deaths per game
	AverageDeaths float64 `json:"averageDeaths"`
	// Average assists per game
	AverageAssists float64 `json:"averageAssists"`
	// Kill/Death/Assist ratio
	KDA float64 `json:"kda"`
	// Average creep score (CS) per game
	AverageCS float64 `json:"averageCs"`
	// Average CS per minute
	CSPerMinute float64 `json:"csPerMinute"`
	// Average vision score per game
	AverageVisionScore float64 `json:"averageVisionScore"`
	// Average damage dealt to champions
	AverageDamage float64 `json:"averageDamage"`
	// Average gold earned per game
	AverageGold float64 `json:"averageGold"`
	// Most played champions with count
	ChampionPool map[string]int `json:"championPool"`
	// Role distribution (percentage of games in each role)
	RoleDistribution map[string]float64 `json:"roleDistribution"`
}

// ImprovementArea represents a specific area where the player can improve
type ImprovementArea struct {
	// Category of improvement (e.g., "CS", "Vision", "Deaths", "Damage")
	Category string `json:"category"`
	// Current performance metric value
	CurrentValue float64 `json:"currentValue"`
	// Average value for players at similar rank
	ExpectedValue float64 `json:"expectedValue"`
	// Difference between current and expected (negative means underperforming)
	Gap float64 `json:"gap"`
	// Priority level (HIGH, MEDIUM, LOW) based on impact
	Priority string `json:"priority"`
	// Specific recommendation text for the player
	Recommendation string `json:"recommendation"`
}

// AnalysisResult contains the complete analysis for a player
type AnalysisResult struct {
	// Player statistics summary
	PlayerStats PlayerStats `json:"playerStats"`
	// List of identified improvement areas
	ImprovementAreas []ImprovementArea `json:"improvementAreas"`
	// Timestamp of when the analysis was performed
	AnalyzedAt time.Time `json:"analyzedAt"`
}
