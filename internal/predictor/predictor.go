package predictor

type Prediction struct {
	HomeGoals int
	AwayGoals int
}

type Predictor interface {
	Predict(homeTeamID, awayTeamID int) (*Prediction, error)
}
