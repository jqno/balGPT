package database_test

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) LastYearMatchScores(homeTeamID, awayTeamID int) (int, int, error) {
	args := m.Called(homeTeamID, awayTeamID)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (m *MockDB) AverageGoalsInLastMatches(teamID, matches int) (float64, error) {
	args := m.Called(teamID, matches)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockDB) GetCurrentSeasonLeaderboard() (map[int]int, error) {
	args := m.Called()
	return args.Get(0).(map[int]int), args.Error(1)
}

func (m *MockDB) GetLastScrape() (time.Time, error) {
	args := m.Called()
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) InsertOrUpdateMatch(homeTeam, awayTeam string, homeGoals, awayGoals int, date time.Time) error {
	args := m.Called(homeTeam, awayTeam, homeGoals, awayGoals, date)
	return args.Error(0)
}

func (m *MockDB) UpdateLastScrape(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}
