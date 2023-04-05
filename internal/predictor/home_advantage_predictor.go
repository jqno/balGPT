package predictor

type HomeAdvantagePredictor struct {
}

func NewHomeAdvantagePredictor() *HomeAdvantagePredictor {
	return &HomeAdvantagePredictor{}
}

func (h *HomeAdvantagePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	return &Prediction{HomeGoals: 1, AwayGoals: 0}, nil
}
