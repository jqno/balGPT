package predictor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHomeAdvantagePredictor(t *testing.T) {
	predictor := NewHomeAdvantagePredictor()
	assert.NotNil(t, predictor)
}

func TestHomeAdvantagePredictor_Predict(t *testing.T) {
	predictor := NewHomeAdvantagePredictor()

	prediction, err := predictor.Predict(1, 2)

	expectedPrediction := &Prediction{HomeGoals: 1, AwayGoals: 0}
	assert.Nil(t, err)
	assert.Equal(t, expectedPrediction, prediction)
}
