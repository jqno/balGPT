package predictor

import (
	"database/sql"

	"github.com/jqno/balGPT/internal/database"
)

type LastYearMatchPredictor struct {
	db           *database.DB
	flippedTeams bool
}

func NewLastYearMatchPredictor(db *database.DB) *LastYearMatchPredictor {
	return &LastYearMatchPredictor{db: db, flippedTeams: false}
}

func NewFlippedLastYearMatchPredictor(db *database.DB) *LastYearMatchPredictor {
	return &LastYearMatchPredictor{db: db, flippedTeams: true}
}
func (l *LastYearMatchPredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	if l.flippedTeams {
		homeTeamID, awayTeamID = awayTeamID, homeTeamID
	}

	homeGoals, awayGoals, err := l.db.LastYearMatchScores(homeTeamID, awayTeamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if l.flippedTeams {
		homeGoals, awayGoals = awayGoals, homeGoals
	}

	return &Prediction{HomeGoals: homeGoals, AwayGoals: awayGoals}, nil
}
