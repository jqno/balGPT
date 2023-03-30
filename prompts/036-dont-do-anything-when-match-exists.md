# Don't do anything when match exists

> How do I make sure that nothing happens if the match exists, instead of an update?
>
> _(snipped: InsertOrUpdateMatch database function)_

If you want to make sure that nothing happens if the match exists, you can remove the update section from the code. The modified function would look like this:

```go
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
```

Now, if the match already exists, the function will simply return `nil` and not perform any updates.

## Note

OK it was probably overkill to ask this at all, but I'm trying to make a point here so I have to ask all the things :).
