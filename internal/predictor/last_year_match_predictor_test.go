package predictor

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/jqno/balGPT/internal/database_test"
	"github.com/stretchr/testify/assert"
)

func TestNewLastYearMatchPredictor(t *testing.T) {
	mockDB := new(database_test.MockDB)
	predictor := NewLastYearMatchPredictor(mockDB)
	assert.False(t, predictor.flippedTeams)
}

func TestNewFlippedLastYearMatchPredictor(t *testing.T) {
	mockDB := new(database_test.MockDB)
	predictor := NewFlippedLastYearMatchPredictor(mockDB)
	assert.True(t, predictor.flippedTeams)
}

func TestPredictWithFlippedTeams(t *testing.T) {
	mockDB := new(database_test.MockDB)
	mockDB.On("LastYearMatchScores", 2, 1).Return(1, 2, nil)
	predictor := NewFlippedLastYearMatchPredictor(mockDB)
	prediction, err := predictor.Predict(1, 2)

	assert.Nil(t, err)
	assert.Equal(t, &Prediction{HomeGoals: 2, AwayGoals: 1}, prediction)
}

func TestPredictWithoutFlippedTeams(t *testing.T) {
	mockDB := new(database_test.MockDB)
	mockDB.On("LastYearMatchScores", 1, 2).Return(1, 2, nil)
	predictor := NewLastYearMatchPredictor(mockDB)
	prediction, err := predictor.Predict(1, 2)

	assert.Nil(t, err)
	assert.Equal(t, &Prediction{HomeGoals: 1, AwayGoals: 2}, prediction)
}

func TestPredictDatabaseError(t *testing.T) {
	mockDB := new(database_test.MockDB)
	mockDB.On("LastYearMatchScores", 1, 2).Return(0, 0, errors.New("some database error"))
	predictor := NewLastYearMatchPredictor(mockDB)
	prediction, err := predictor.Predict(1, 2)

	assert.Nil(t, prediction)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "some database error")
}

func TestPredictNoRowsError(t *testing.T) {
	mockDB := new(database_test.MockDB)
	mockDB.On("LastYearMatchScores", 1, 2).Return(0, 0, sql.ErrNoRows)
	predictor := NewLastYearMatchPredictor(mockDB)
	prediction, err := predictor.Predict(1, 2)

	assert.Nil(t, prediction)
	assert.Nil(t, err)
}
