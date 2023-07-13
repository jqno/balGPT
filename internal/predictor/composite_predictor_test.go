package predictor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPredictor struct {
	mock.Mock
}

func (m *MockPredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	args := m.Called(homeTeamID, awayTeamID)
	return args.Get(0).(*Prediction), args.Error(1)
}

func TestCompositePredictorWithNoPredictors(t *testing.T) {
	c := NewCompositePredictor()

	prediction, err := c.Predict(1, 2)
	assert.Nil(t, prediction)
	assert.Equal(t, errors.New("No predictors provided"), err)
}

func TestCompositePredictorWithOnePredictor(t *testing.T) {
	mockPredictor := new(MockPredictor)
	mockPredictor.On("Predict", 1, 2).Return(&Prediction{HomeGoals: 3, AwayGoals: 1}, nil)

	c := NewCompositePredictor(mockPredictor)

	prediction, err := c.Predict(1, 2)
	assert.Nil(t, err)
	assert.Equal(t, &Prediction{HomeGoals: 3, AwayGoals: 1}, prediction)
}

func TestCompositePredictorWithMultiplePredictors(t *testing.T) {
	mockPredictor1 := new(MockPredictor)
	mockPredictor1.On("Predict", 1, 2).Return(&Prediction{HomeGoals: 3, AwayGoals: 1}, nil)

	mockPredictor2 := new(MockPredictor)
	mockPredictor2.On("Predict", 1, 2).Return(&Prediction{HomeGoals: 2, AwayGoals: 2}, nil)

	c := NewCompositePredictor(mockPredictor1, mockPredictor2)

	prediction, err := c.Predict(1, 2)
	assert.Nil(t, err)
	assert.Equal(t, &Prediction{HomeGoals: 2, AwayGoals: 1}, prediction) // Medians of 3,2 and 1,2 are 2 and 1 respectively.
}
