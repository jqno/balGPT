package predictor

import (
	"math"
)

type AverageGoalsPredictor struct {
	db DB
}

func NewAverageGoalsPredictor(db DB) *AverageGoalsPredictor {
	return &AverageGoalsPredictor{db: db}
}

func (a *AverageGoalsPredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	homeAvgGoals, err := a.db.AverageGoalsInLastMatches(homeTeamID, 8)
	if err != nil {
		return nil, err
	}

	awayAvgGoals, err := a.db.AverageGoalsInLastMatches(awayTeamID, 8)
	if err != nil {
		return nil, err
	}

	homeGoals := int(math.Round(homeAvgGoals))
	awayGoals := int(math.Round(awayAvgGoals))

	return &Prediction{HomeGoals: homeGoals, AwayGoals: awayGoals}, nil
}
