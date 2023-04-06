# Leaderboard predictor

> I want to create another predictor, which calculates the difference between the two teams on the leader board, and awards the difference in position divided by 2 as the number of goals for the team highest on the board, and zero to the other.
>
> The leaderboard is defined as follows. For each match, the winning team gets 3 points and the losing team 0; it doesn't matter if they're home or away. If it's a draw, both teams get 1 point. The team in 1st place on the leaderboard is the team with the largest number of points; the team in 2nd place has the second largest number of points, and so on. Only the matches in the current season are counted. Seasons start on August 1st; that's when the leaderboard is reset.
>
> Note that it's possible that no leaderboard exists yet because no matches have yet been played since August 1st. In that case, the predictor should return nil.
>
> Please write this predictor, and the database query or queries required to calculate it.

Let's start by writing a function in the `database.go` file to get the leaderboard for the current season. We'll then create a new predictor called `LeaderboardDifferencePredictor` that uses this leaderboard to calculate the prediction.

First, add the following function to your `database.go` file:

```go
// database/database.go
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
```

This function calculates the current season leaderboard by fetching all matches from the current season and computing the points for each team based on the match results. If no matches have been played yet, the function will return an empty map.

Now, let's create the `LeaderboardDifferencePredictor`. Create a new file `leaderboard_difference_predictor.go` inside the `predictor` directory:

```go
// predictor/leaderboard_difference_predictor.go
package predictor

import (
	"database"
	"sort"
)

type LeaderboardDifferencePredictor struct {
	db *database.DB
}

func NewLeaderboardDifferencePredictor(db *database.DB) *LeaderboardDifferencePredictor {
	return &LeaderboardDifferencePredictor{db: db}
}

func (l *LeaderboardDifferencePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	leaderboard, err := l.db.GetCurrentSeasonLeaderboard()
	if err != nil {
		return nil, err
	}

	if len(leaderboard) == 0 {
		return nil, nil
	}

	homePoints := leaderboard[homeTeamID]
	awayPoints := leaderboard[awayTeamID]

	positionDifference := abs(homePoints - awayPoints) / 2

	if homePoints > awayPoints {
		return &Prediction{HomeGoals: positionDifference, AwayGoals: 0}, nil
	} else if homePoints < awayPoints {
		return &Prediction{HomeGoals: 0, AwayGoals: positionDifference}, nil
	} else {
		return &Prediction{HomeGoals: 0, AwayGoals: 0}, nil
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
```

This function `abs` calculates the absolute value of an integer, which is used in the `Predict` function to compute the position difference between the home team and the away team.

Now, you can use the `LeaderboardDifferencePredictor` in your application.

## Follow-up prompt

> Awesome!
>
> Please update the Predict function to also print the leaderboard to the log, in descending order of number of points on the board.

To print the leaderboard in descending order of points, we can modify the `Predict` function of the `LeaderboardDifferencePredictor` as follows:

First, we need to import the "log" package at the beginning of the `leaderboard_difference_predictor.go` file:

```go
import (
	"database"
	"log"
	"sort"
)
```

Then, update the `Predict` function to print the sorted leaderboard:

```go
func (l *LeaderboardDifferencePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	leaderboard, err := l.db.GetCurrentSeasonLeaderboard()
	if err != nil {
		return nil, err
	}

	if len(leaderboard) == 0 {
		return nil, nil
	}

	type leaderboardEntry struct {
		teamID int
		points int
	}

	// Sort the leaderboard in descending order of points
	sortedLeaderboard := make([]leaderboardEntry, 0, len(leaderboard))
	for teamID, points := range leaderboard {
		sortedLeaderboard = append(sortedLeaderboard, leaderboardEntry{teamID: teamID, points: points})
	}

	sort.Slice(sortedLeaderboard, func(i, j int) bool {
		return sortedLeaderboard[i].points > sortedLeaderboard[j].points
	})

	// Print the sorted leaderboard
	log.Println("Current Season Leaderboard:")
	for _, entry := range sortedLeaderboard {
		log.Printf("Team ID: %d, Points: %d\n", entry.teamID, entry.points)
	}

	homePoints := leaderboard[homeTeamID]
	awayPoints := leaderboard[awayTeamID]

	positionDifference := abs(homePoints - awayPoints) / 2

	if homePoints > awayPoints {
		return &Prediction{HomeGoals: positionDifference, AwayGoals: 0}, nil
	} else if homePoints < awayPoints {
		return &Prediction{HomeGoals: 0, AwayGoals: positionDifference}, nil
	} else {
		return &Prediction{HomeGoals: 0, AwayGoals: 0}, nil
	}
}
```

This modification first creates a sorted slice of the leaderboard entries and then prints the sorted leaderboard to the log. Note that this approach assumes that the team IDs are unique.

Now, when you call the `Predict` function of the `LeaderboardDifferencePredictor`, it will also print the current season leaderboard in descending order of points.

## Follow-up prompt

> You know what? Let's extract logging the leaderboard to a separate, private function.

Sure, we can create a separate private function for logging the leaderboard. Update the `leaderboard_difference_predictor.go` file as follows:

1. Create the private function logLeaderboard:

```go
func (l *LeaderboardDifferencePredictor) logLeaderboard(leaderboard map[int]int) {
	type leaderboardEntry struct {
		teamID int
		points int
	}

	sortedLeaderboard := make([]leaderboardEntry, 0, len(leaderboard))
	for teamID, points := range leaderboard {
		sortedLeaderboard = append(sortedLeaderboard, leaderboardEntry{teamID: teamID, points: points})
	}

	sort.Slice(sortedLeaderboard, func(i, j int) bool {
		return sortedLeaderboard[i].points > sortedLeaderboard[j].points
	})

	log.Println("Current Season Leaderboard:")
	for _, entry := range sortedLeaderboard {
		log.Printf("Team ID: %d, Points: %d\n", entry.teamID, entry.points)
	}
}
```

1. Update the Predict function to call this new private function:

```go
func (l *LeaderboardDifferencePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	leaderboard, err := l.db.GetCurrentSeasonLeaderboard()
	if err != nil {
		return nil, err
	}

	if len(leaderboard) == 0 {
		return nil, nil
	}

	l.logLeaderboard(leaderboard)

	homePoints := leaderboard[homeTeamID]
	awayPoints := leaderboard[awayTeamID]

	positionDifference := abs(homePoints - awayPoints) / 2

	if homePoints > awayPoints {
		return &Prediction{HomeGoals: positionDifference, AwayGoals: 0}, nil
	} else if homePoints < awayPoints {
		return &Prediction{HomeGoals: 0, AwayGoals: positionDifference}, nil
	} else {
		return &Prediction{HomeGoals: 0, AwayGoals: 0}, nil
	}
}
```

Now, the `logLeaderboard` function is responsible for logging the leaderboard, and the `Predict` function simply calls it. This makes the code more modular and easier to maintain.

## Follow-up prompt

> This is really great, and the leaderboard was correct on the first try!

I'm glad to hear that the leaderboard works correctly and that you're happy with the solution! If you have any more questions or need further assistance with your Go app or anything else, please don't hesitate to ask. I'm here to help!
