package predictor

import (
	"errors"
	"log"
	"sort"
	"time"
)

type CompositePredictor struct {
	predictors []Predictor
}

func NewCompositePredictor(predictors ...Predictor) *CompositePredictor {
	return &CompositePredictor{predictors: predictors}
}

func (c *CompositePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	if len(c.predictors) == 0 {
		return nil, errors.New("No predictors provided")
	}

	homeGoals := make([]int, 0, len(c.predictors))
	awayGoals := make([]int, 0, len(c.predictors))

	for _, predictor := range c.predictors {
		predictionStartTime := time.Now()
		prediction, err := predictor.Predict(homeTeamID, awayTeamID)
		predictionEndTime := time.Now()

		if err != nil {
			return nil, err
		}
		if prediction != nil {
			log.Printf("Prediction from %T: %+v, runtime: %v", predictor, prediction, predictionEndTime.Sub(predictionStartTime))
			homeGoals = append(homeGoals, prediction.HomeGoals)
			awayGoals = append(awayGoals, prediction.AwayGoals)
		}
	}

	if len(homeGoals) == 0 || len(awayGoals) == 0 {
		return nil, errors.New("No predictions available")
	}

	sort.Ints(homeGoals)
	sort.Ints(awayGoals)

	medianHomeGoals := calculateMedian(homeGoals)
	medianAwayGoals := calculateMedian(awayGoals)

	return &Prediction{HomeGoals: medianHomeGoals, AwayGoals: medianAwayGoals}, nil
}

func calculateMedian(sortedValues []int) int {
	middle := len(sortedValues) / 2

	if len(sortedValues)%2 == 0 {
		return (sortedValues[middle-1] + sortedValues[middle]) / 2
	}

	return sortedValues[middle]
}
