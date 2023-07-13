package predictor

type DB interface {
	LastYearMatchScores(homeTeamID, awayTeamID int) (int, int, error)
	AverageGoalsInLastMatches(teamID, matches int) (float64, error)
	GetCurrentSeasonLeaderboard() (map[int]int, error)
}
