package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jqno/balGPT/internal/team"
	_ "github.com/lib/pq"
)

type DB struct {
	Conn *sql.DB
}

func New(connectionString string, appBaseDir string) *DB {
	conn, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	pgDriver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		panic(fmt.Errorf("unable to create postgres driver instance: %v", err))
	}

	if err := RunMigrations(pgDriver, appBaseDir); err != nil {
		panic(fmt.Errorf("failed to run migrations: %v", err))
	}

	return &DB{Conn: conn}
}

func (db *DB) GetLastScrape() (time.Time, error) {
	var lastScrape time.Time
	err := db.Conn.QueryRow("SELECT last_scrape FROM stats ORDER BY id DESC LIMIT 1").Scan(&lastScrape)
	if err != nil {
		if err == sql.ErrNoRows {
			return time.Time{}, nil
		}
		return time.Time{}, err
	}

	return lastScrape, nil
}

func (db *DB) UpdateLastScrape(scrapeTime time.Time) error {
	_, err := db.Conn.Exec("INSERT INTO stats (last_scrape) VALUES ($1)", scrapeTime)
	return err
}

func (db *DB) FetchTeamsFromDB() ([]team.Team, error) {
	rows, err := db.Conn.Query("SELECT id, name FROM teams")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teams := []team.Team{}
	for rows.Next() {
		var team team.Team
		if err := rows.Scan(&team.ID, &team.Name); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (db *DB) InsertOrUpdateMatch(homeTeam, awayTeam string, homeGoals, awayGoals int, date time.Time) error {
	// Insert or update the home team
	homeTeamID, err := db.insertOrUpdateTeam(homeTeam)
	if err != nil {
		return err
	}

	// Insert or update the away team
	awayTeamID, err := db.insertOrUpdateTeam(awayTeam)
	if err != nil {
		return err
	}

	// Check if the match already exists
	var matchID int
	err = db.Conn.QueryRow("SELECT id FROM matches WHERE home_team = $1 AND away_team = $2 AND date = $3",
		homeTeamID, awayTeamID, date).Scan(&matchID)

	switch {
	case err == sql.ErrNoRows:
		// Insert a new match
		_, err := db.Conn.Exec("INSERT INTO matches (home_team, away_team, home_goals, away_goals, date) VALUES ($1, $2, $3, $4, $5)",
			homeTeamID, awayTeamID, homeGoals, awayGoals, date)
		return err
	case err != nil:
		return err
	default:
		// Do nothing if the match already exists
		return nil
	}
}

func (db *DB) insertOrUpdateTeam(name string) (int, error) {
	var teamID int
	err := db.Conn.QueryRow("SELECT id FROM teams WHERE name = $1", name).Scan(&teamID)

	switch {
	case err == sql.ErrNoRows:
		// Insert a new team
		err := db.Conn.QueryRow("INSERT INTO teams (name) VALUES ($1) RETURNING id", name).Scan(&teamID)
		if err != nil {
			return 0, err
		}
	case err != nil:
		return 0, err
	}

	return teamID, nil
}

func (db *DB) GetTeamID(teamName string) (int, error) {
	var teamID int
	err := db.Conn.QueryRow("SELECT id FROM teams WHERE name = $1", teamName).Scan(&teamID)
	if err != nil {
		return 0, err
	}

	return teamID, nil
}

func (db *DB) GetTeamAvgGoals(teamID int, isHomeTeam bool) (float64, error) {
	var avgGoals float64
	var err error

	if isHomeTeam {
		err = db.Conn.QueryRow("SELECT AVG(home_goals) FROM matches WHERE home_team = $1", teamID).Scan(&avgGoals)
	} else {
		err = db.Conn.QueryRow("SELECT AVG(away_goals) FROM matches WHERE away_team = $1", teamID).Scan(&avgGoals)
	}

	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return avgGoals, nil
}
