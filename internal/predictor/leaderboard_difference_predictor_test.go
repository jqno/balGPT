package predictor

import (
	"testing"

	"github.com/jqno/balGPT/internal/database_test"
	"github.com/stretchr/testify/assert"
)

func TestPredictWithEmptyLeaderboard(t *testing.T) {
	mockDB := new(database_test.MockDB)
	mockDB.On("GetCurrentSeasonLeaderboard").Return(map[int]int{}, nil)
	predictor := NewLeaderboardDifferencePredictor(mockDB)

	prediction, err := predictor.Predict(1, 2)

	assert.NoError(t, err)
	assert.Nil(t, prediction)
	mockDB.AssertExpectations(t)
}

func TestPredictWithHomeTeamLeading(t *testing.T) {
	mockDB := new(database_test.MockDB)
	mockDB.On("GetCurrentSeasonLeaderboard").Return(map[int]int{
		1: 25, // home team
		2: 20,
		3: 15,
		4: 10, // away team
	}, nil)
	predictor := NewLeaderboardDifferencePredictor(mockDB)

	prediction, err := predictor.Predict(1, 4)

	assert.NoError(t, err)
	assert.Equal(t, &Prediction{HomeGoals: 1, AwayGoals: 0}, prediction)
	mockDB.AssertExpectations(t)
}

func TestPredictWithAwayTeamLeading(t *testing.T) {
	mockDB := new(database_test.MockDB)
	mockDB.On("GetCurrentSeasonLeaderboard").Return(map[int]int{
		1: 10, // home team
		2: 15,
		3: 20,
		4: 25, // away team
	}, nil)
	predictor := NewLeaderboardDifferencePredictor(mockDB)

	prediction, err := predictor.Predict(1, 4)

	assert.NoError(t, err)
	assert.Equal(t, &Prediction{HomeGoals: 0, AwayGoals: 1}, prediction)
	mockDB.AssertExpectations(t)
}

func TestPredictWithEqualTeams(t *testing.T) {
	mockDB := new(database_test.MockDB)
	mockDB.On("GetCurrentSeasonLeaderboard").Return(map[int]int{
		1: 10, // home team
		2: 10, // away team
	}, nil)
	predictor := NewLeaderboardDifferencePredictor(mockDB)

	prediction, err := predictor.Predict(1, 2)

	assert.NoError(t, err)
	assert.Equal(t, &Prediction{HomeGoals: 0, AwayGoals: 0}, prediction)
	mockDB.AssertExpectations(t)
}

func TestPredictWithBothTeamsMinusOne(t *testing.T) {
	mockDB := new(database_test.MockDB)
	predictor := NewLeaderboardDifferencePredictor(mockDB)

	prediction, err := predictor.Predict(-1, -1)

	assert.NoError(t, err)
	assert.Equal(t, &Prediction{HomeGoals: 0, AwayGoals: 0}, prediction)
	mockDB.AssertExpectations(t)
}

func TestPredictWithHomeTeamMinusOne(t *testing.T) {
	mockDB := new(database_test.MockDB)
	predictor := NewLeaderboardDifferencePredictor(mockDB)

	prediction, err := predictor.Predict(-1, 4)

	assert.NoError(t, err)
	assert.Equal(t, &Prediction{HomeGoals: 0, AwayGoals: 1}, prediction)
	mockDB.AssertExpectations(t)
}

func TestPredictWithAwayTeamMinusOne(t *testing.T) {
	mockDB := new(database_test.MockDB)
	predictor := NewLeaderboardDifferencePredictor(mockDB)

	prediction, err := predictor.Predict(1, -1)

	assert.NoError(t, err)
	assert.Equal(t, &Prediction{HomeGoals: 1, AwayGoals: 0}, prediction)
	mockDB.AssertExpectations(t)
}
