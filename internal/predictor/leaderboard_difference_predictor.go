package predictor

import (
	"log"
	"sort"

	"github.com/jqno/balGPT/internal/database"
)

type LeaderboardDifferencePredictor struct {
	db *database.DB
}

func NewLeaderboardDifferencePredictor(db *database.DB) *LeaderboardDifferencePredictor {
	return &LeaderboardDifferencePredictor{db: db}
}

func (l *LeaderboardDifferencePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	leaderboard, err := l.db.GetCurrentSeasonLeaderboard()
	if err != nil {
		return nil, err
	}

	if len(leaderboard) == 0 {
		return nil, nil
	}

	l.logLeaderboard(leaderboard)

	homePoints := leaderboard[homeTeamID]
	awayPoints := leaderboard[awayTeamID]

	positionDifference := abs(homePoints-awayPoints) / 2

	if homePoints > awayPoints {
		return &Prediction{HomeGoals: positionDifference, AwayGoals: 0}, nil
	} else if homePoints < awayPoints {
		return &Prediction{HomeGoals: 0, AwayGoals: positionDifference}, nil
	} else {
		return &Prediction{HomeGoals: 0, AwayGoals: 0}, nil
	}
}

func (l *LeaderboardDifferencePredictor) logLeaderboard(leaderboard map[int]int) {
	type leaderboardEntry struct {
		teamID int
		points int
	}

	sortedLeaderboard := make([]leaderboardEntry, 0, len(leaderboard))
	for teamID, points := range leaderboard {
		sortedLeaderboard = append(sortedLeaderboard, leaderboardEntry{teamID: teamID, points: points})
	}

	sort.Slice(sortedLeaderboard, func(i, j int) bool {
		return sortedLeaderboard[i].points > sortedLeaderboard[j].points
	})

	log.Println("Current Season Leaderboard:")
	for _, entry := range sortedLeaderboard {
		log.Printf("Team ID: %d, Points: %d\n", entry.teamID, entry.points)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}