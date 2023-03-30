package predictor

import (
	"github.com/jqno/balGPT/internal/database"
)

type Prediction struct {
	HomeGoals int
	AwayGoals int
}

type Predictor struct {
	DB *database.DB
}

func NewPredictor(db *database.DB) *Predictor {
	return &Predictor{DB: db}
}

func (p *Predictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	homeTeamAvgGoals := 0.0
	awayTeamAvgGoals := 0.0
	var err error

	if homeTeamID != -1 {
		homeTeamAvgGoals, err = p.DB.GetTeamAvgGoals(homeTeamID, true)
		if err != nil {
			return nil, err
		}
	}

	if awayTeamID != -1 {
		awayTeamAvgGoals, err = p.DB.GetTeamAvgGoals(awayTeamID, false)
		if err != nil {
			return nil, err
		}
	}

	prediction := &Prediction{
		HomeGoals: int(homeTeamAvgGoals + 0.5),
		AwayGoals: int(awayTeamAvgGoals + 0.5),
	}

	return prediction, nil
}
