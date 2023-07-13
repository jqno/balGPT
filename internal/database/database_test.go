package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetLastScrape(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database := DB{Conn: db}

	rows := sqlmock.NewRows([]string{"last_scrape"}).
		AddRow(time.Now())

	query := "SELECT last_scrape FROM stats ORDER BY id DESC LIMIT 1"
	mock.ExpectQuery(query).WillReturnRows(rows)

	_, err = database.GetLastScrape()
	assert.NoError(t, err)
}

func TestUpdateLastScrape(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database := DB{Conn: db}

	mock.ExpectExec("INSERT INTO stats \\(last_scrape\\) VALUES \\(\\$1\\)").
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = database.UpdateLastScrape(time.Now())
	assert.NoError(t, err)
}

func TestFetchTeamsFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database := DB{Conn: db}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "team1").
		AddRow(2, "team2")

	mock.ExpectQuery("SELECT id, name FROM teams").WillReturnRows(rows)

	teams, err := database.FetchTeamsFromDB()
	assert.NoError(t, err)
	assert.Len(t, teams, 2)
	assert.Equal(t, "team1", teams[0].Name)
	assert.Equal(t, "team2", teams[1].Name)
}

func TestInsertOrUpdateMatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database := DB{Conn: db}

	mock.ExpectQuery("SELECT id FROM teams WHERE name = \\$1").
		WithArgs("Home").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectQuery("SELECT id FROM teams WHERE name = \\$1").
		WithArgs("Away").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	mock.ExpectQuery("SELECT id FROM matches WHERE home_team = \\$1 AND away_team = \\$2 AND date = \\$3").
		WithArgs(1, 2, sqlmock.AnyArg()).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO matches \\(home_team, away_team, home_goals, away_goals, date\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)").
		WithArgs(1, 2, 3, 2, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = database.InsertOrUpdateMatch("Home", "Away", 3, 2, time.Now())
	assert.NoError(t, err)
}

func TestGetTeamID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database := DB{Conn: db}

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	mock.ExpectQuery("SELECT id FROM teams WHERE name = \\$1").WithArgs("team1").WillReturnRows(rows)

	id, err := database.GetTeamID("team1")
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestAverageGoalsInLastMatches(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database := DB{Conn: db}

	rows := sqlmock.NewRows([]string{"AVG(goals)"}).AddRow(1.5)

	mock.ExpectQuery(".*").WithArgs(1, 5).WillReturnRows(rows)

	avg, err := database.AverageGoalsInLastMatches(1, 5)
	assert.NoError(t, err)
	assert.Equal(t, 1.5, avg)
}

func TestLastYearMatchScores(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database := DB{Conn: db}

	rows := sqlmock.NewRows([]string{"home_goals", "away_goals"}).AddRow(2, 3)

	mock.ExpectQuery(".*").WithArgs(1, 2).WillReturnRows(rows)

	homeGoals, awayGoals, err := database.LastYearMatchScores(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, 2, homeGoals)
	assert.Equal(t, 3, awayGoals)
}

func TestGetCurrentSeasonLeaderboard(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database := DB{Conn: db}

	rows := sqlmock.NewRows([]string{"winner", "draw_team1", "draw_team2"}).
		AddRow(sql.NullInt64{Valid: true, Int64: 1}, sql.NullInt64{Valid: false, Int64: 0}, sql.NullInt64{Valid: false, Int64: 0}).
		AddRow(sql.NullInt64{Valid: false, Int64: 0}, sql.NullInt64{Valid: true, Int64: 2}, sql.NullInt64{Valid: true, Int64: 3})

	mock.ExpectQuery(".*").WillReturnRows(rows)

	result, err := database.GetCurrentSeasonLeaderboard()
	assert.NoError(t, err)
	assert.Equal(t, map[int]int{1: 3, 2: 1, 3: 1}, result)
}
