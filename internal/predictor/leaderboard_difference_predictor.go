package predictor

import (
	"log"
	"sort"
)

type LeaderboardDifferencePredictor struct {
	db DB
}

type leaderboardEntry struct {
	teamID int
	points int
}

func NewLeaderboardDifferencePredictor(db DB) *LeaderboardDifferencePredictor {
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

	sortedLeaderboard := l.sortLeaderboard(leaderboard)
	l.logLeaderboard(sortedLeaderboard)

	homePosition, awayPosition := l.getTeamPositions(homeTeamID, awayTeamID, sortedLeaderboard)
	positionDifference := abs(homePosition-awayPosition) / 2

	if homePosition < awayPosition {
		return &Prediction{HomeGoals: positionDifference, AwayGoals: 0}, nil
	} else if homePosition > awayPosition {
		return &Prediction{HomeGoals: 0, AwayGoals: positionDifference}, nil
	} else {
		return &Prediction{HomeGoals: 0, AwayGoals: 0}, nil
	}
}

func (l *LeaderboardDifferencePredictor) sortLeaderboard(leaderboard map[int]int) []leaderboardEntry {

	sortedLeaderboard := make([]leaderboardEntry, 0, len(leaderboard))
	for teamID, points := range leaderboard {
		sortedLeaderboard = append(sortedLeaderboard, leaderboardEntry{teamID: teamID, points: points})
	}

	sort.Slice(sortedLeaderboard, func(i, j int) bool {
		return sortedLeaderboard[i].points > sortedLeaderboard[j].points
	})

	return sortedLeaderboard
}

func (l *LeaderboardDifferencePredictor) logLeaderboard(sortedLeaderboard []leaderboardEntry) {
	log.Println("Current Season Leaderboard:")
	for _, entry := range sortedLeaderboard {
		log.Printf(" - Team ID: %d, Points: %d\n", entry.teamID, entry.points)
	}
}

func (l *LeaderboardDifferencePredictor) getTeamPositions(homeTeamID, awayTeamID int, sortedLeaderboard []leaderboardEntry) (int, int) {
	var homePosition, awayPosition int

	for position, entry := range sortedLeaderboard {
		if entry.teamID == homeTeamID {
			homePosition = position + 1
		} else if entry.teamID == awayTeamID {
			awayPosition = position + 1
		}

		if homePosition > 0 && awayPosition > 0 {
			break
		}
	}

	return homePosition, awayPosition
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
