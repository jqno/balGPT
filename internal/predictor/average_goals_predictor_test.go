package predictor

import (
	"errors"
	"testing"

	"github.com/jqno/balGPT/internal/database_test"
	"github.com/stretchr/testify/assert"
)

func TestNewAverageGoalsPredictor(t *testing.T) {
	mockDB := new(database_test.MockDB)
	predictor := NewAverageGoalsPredictor(mockDB)
	assert.NotNil(t, predictor)
}

func TestAverageGoalsPredictor_Predict(t *testing.T) {
	mockDB := new(database_test.MockDB)
	mockDB.On("AverageGoalsInLastMatches", 1, 8).Return(2.4, nil)
	mockDB.On("AverageGoalsInLastMatches", 2, 8).Return(1.6, nil)

	predictor := NewAverageGoalsPredictor(mockDB)
	prediction, err := predictor.Predict(1, 2)

	expectedPrediction := &Prediction{HomeGoals: 2, AwayGoals: 2}
	assert.Nil(t, err)
	assert.Equal(t, expectedPrediction, prediction)
}

func TestAverageGoalsPredictor_PredictDatabaseError(t *testing.T) {
	mockDB := new(database_test.MockDB)
	mockDB.On("AverageGoalsInLastMatches", 1, 8).Return(0.0, errors.New("some database error"))

	predictor := NewAverageGoalsPredictor(mockDB)
	prediction, err := predictor.Predict(1, 2)

	assert.Nil(t, prediction)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "some database error")
}
