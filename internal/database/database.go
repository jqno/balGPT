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

func (db *DB) AverageGoalsInLastMatches(teamID int, numberOfMatches int) (float64, error) {
	if teamID == -1 {
		return 0, nil
	}

	query := `
		WITH combined AS (
			SELECT home_team AS team, home_goals AS goals, date
			FROM matches
			WHERE home_team = $1
			UNION ALL
			SELECT away_team AS team, away_goals AS goals, date
			FROM matches
			WHERE away_team = $1
		)
		SELECT AVG(goals)
		FROM (
			SELECT goals
			FROM combined
			ORDER BY date DESC
			LIMIT $2
		) last_matches;
	`

	var avgGoals float64
	err := db.Conn.QueryRow(query, teamID, numberOfMatches).Scan(&avgGoals)
	if err != nil {
		return 0, fmt.Errorf("Error fetching average goals for team %d: %v", teamID, err)
	}

	return avgGoals, nil
}

func (db *DB) LastYearMatchScores(homeTeamID, awayTeamID int) (int, int, error) {
	query := `
		SELECT home_goals, away_goals
		FROM matches
		WHERE home_team = $1 AND away_team = $2
		ORDER BY date DESC
		LIMIT 1;
	`

	var homeGoals, awayGoals int
	err := db.Conn.QueryRow(query, homeTeamID, awayTeamID).Scan(&homeGoals, &awayGoals)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil
		}
		return 0, 0, fmt.Errorf("Error fetching last match scores: %v", err)
	}

	return homeGoals, awayGoals, nil
}

func (db *DB) GetCurrentSeasonLeaderboard() (map[int]int, error) {
	seasonStart := time.Date(time.Now().Year(), time.August, 1, 0, 0, 0, 0, time.UTC)
	if time.Now().Before(seasonStart) {
		seasonStart = seasonStart.AddDate(-1, 0, 0)
	}
	query := `
		SELECT
			CASE
				WHEN home_goals > away_goals THEN home_team
				WHEN home_goals < away_goals THEN away_team
				ELSE NULL
			END AS winner,
			CASE
				WHEN home_goals = away_goals THEN home_team
				ELSE NULL
			END AS draw_team1,
			CASE
				WHEN home_goals = away_goals THEN away_team
				ELSE NULL
			END AS draw_team2
		FROM matches
		WHERE date >= $1;
	`

	rows, err := db.Conn.Query(query, seasonStart)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	points := make(map[int]int)
	for rows.Next() {
		var winner, drawTeam1, drawTeam2 sql.NullInt64
		err := rows.Scan(&winner, &drawTeam1, &drawTeam2)
		if err != nil {
			return nil, err
		}

		if winner.Valid {
			points[int(winner.Int64)] += 3
		}
		if drawTeam1.Valid {
			points[int(drawTeam1.Int64)]++
			points[int(drawTeam2.Int64)]++
		}
	}

	return points, nil
}
