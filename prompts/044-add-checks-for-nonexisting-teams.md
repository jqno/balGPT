# Add checks for non-existing teams

> Here's a Go function:
>
> _(snipped: AverageGoalsInLastMatches from database/database.go)_
>
> If teamID equals -1, I want it to return 0.

If you want to immediately return a value of `0` for the average goals when `teamID` is `-1`, you can add a conditional check at the beginning of the `AverageGoalsInLastMatches` function:

```go
func (db *DB) AverageGoalsInLastMatches(teamID int, numberOfMatches int) (float64, error) {
	// Return 0 immediately if teamID equals -1
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

```

By adding the `if teamID == -1` condition at the start, you ensure that the function will return early with a result of `0` without making any unnecessary database queries when the team ID is `-1`.

## Follow-up prompt

> Here's a test case for that function:
>
> _(snipped: TestAverageGoalsInLastMatches from database/database_test.go)_
>
> Write a test case for the -1 case

Sure, here's a test case for when `teamID` is `-1`:

```go
func TestAverageGoalsInLastMatchesWithNegativeTeamID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database := DB{Conn: db}

	// When teamID is -1, no database calls should be made, so we don't need to set up any expectations with sqlmock.

	avg, err := database.AverageGoalsInLastMatches(-1, 5)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, avg)
}

```

In this test, you don't need to mock any database queries or responses because the function should immediately return `0` when `teamID` is `-1`, without making any database calls. The absence of any mock expectations on the database ensures that any unintentional database interactions would cause the test to fail, which is exactly what you'd want to validate.

## Follow-up prompt

> Can we ask the mock to verify that it wasn't called?

Certainly! The `sqlmock` package provides a method to ensure that all the expectations that were set up are asserted. You can use the `mock.ExpectationsWereMet()` function at the end of your test case to check this. If any expectations (like queries, etc.) were defined but were not met, this function will return an error.

For the `-1` case, you haven't set any expectations because you don't expect any database interactions. Thus, calling `ExpectationsWereMet` will essentially verify that no database calls were made.

Here's how you can modify the test case to include this check:

```go
func TestAverageGoalsInLastMatchesWithNegativeTeamID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	database := DB{Conn: db}

	// When teamID is -1, no database calls should be made.

	avg, err := database.AverageGoalsInLastMatches(-1, 5)
	assert.NoError(t, err)
	assert.Equal(t, 0.0, avg)

	// Verify that no database interactions took place.
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

```

With the above code, if for some reason a database call does get made, the `ExpectationsWereMet` function will return an error and the test will fail, letting you know that there was an unexpected interaction with the database.

## Follow-up prompt

> Here's another database function:
>
> _(snipped: LastYearMatchScores from database/database.go)_
>
> What happens if either homeTeamID or awayTeamID is -1? If it leads to an error, please add a check.

Given the current implementation of the `LastYearMatchScores` function, if either `homeTeamID` or `awayTeamID` is `-1`, the function will proceed with the SQL query. Depending on your database setup and data, one of the following outcomes might occur:

1. If there are rows in the matches table where either home_team or away_team has a value of -1, then the function will return the score of the most recent match with that ID.
2. If no such rows exist, the function will return (0, 0, nil) because of the check for sql.ErrNoRows.

If having a `-1` value for team IDs is not valid in your domain (i.e., it doesn't make sense to have a team with an ID of `-1`), you can add a check at the beginning of the function to handle this case.

Here's the modified function with the check added:

```go
func (db *DB) LastYearMatchScores(homeTeamID, awayTeamID int) (int, int, error) {
	// Return (0, 0, nil) if either team ID is -1
	if homeTeamID == -1 || awayTeamID == -1 {
		return 0, 0, nil
	}

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

```

With the added check, if either `homeTeamID` or `awayTeamID` is `-1`, the function will immediately return `(0, 0, nil)`, preventing any unnecessary database calls.

## Follow-up prompt

> Great, that's covered.
>
> Here's another function:
>
> _(snipped: getTeamPositions from predictor/leaderboard_difference_predictor.go)_
>
> What happens here if either homeTeamID or awayTeamID equals -1? Assume that sortedLeaderboard doesn't contain any entries for -1.

If `homeTeamID` or `awayTeamID` is `-1`, and if `sortedLeaderboard` doesn't contain any entries with a `teamID` of `-1`, then the function's behavior will be as follows:

The loop will iterate through each `leaderboardEntry` in `sortedLeaderboard`. For each `leaderboardEntry`, the function will check if the entry's `teamID` matches either `homeTeamID` or `awayTeamID`.

Since `sortedLeaderboard` doesn't contain any entries with a `teamID` of `-1`, there will be no match for a `teamID` value of `-1`.

Thus:

- If homeTeamID is -1, then homePosition will remain as its default value of 0.
- If awayTeamID is -1, then awayPosition will remain as its default value of 0.

Once the loop completes:

1. If both homeTeamID and awayTeamID are -1, then the function will return (0, 0).
2. If only homeTeamID is -1 but awayTeamID is found in the leaderboard, then the function will return (0, awayPosition), where awayPosition is the found position.
3. If only awayTeamID is -1 but homeTeamID is found in the leaderboard, then the function will return (homePosition, 0), where homePosition is the found position.

To clarify this behavior, or to handle the `-1` case differently, you can add explicit checks for `-1` values at the start of the function. However, based on the current function implementation and the given assumptions, that's how the function would behave.

## Follow-up prompt

> Awesome. Let's make sure the function is never called if either team is -1. Here's the function that calls it:
>
> _(snipped: Predict from predictor/leaderboard_difference_predictor.go)_
>
> Let's update so that it does early returns when either team is -1:
>
> - If they're both -1, return 0,0.
> - If one is -1 and the other isn't, return 0 for the one and 1 for the other.

Alright, let's add those early return checks to the `Predict` function:

```go
func (l *LeaderboardDifferencePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
    // Early returns for -1 cases
    if homeTeamID == -1 && awayTeamID == -1 {
        return &Prediction{HomeGoals: 0, AwayGoals: 0}, nil
    } else if homeTeamID == -1 {
        return &Prediction{HomeGoals: 0, AwayGoals: 1}, nil
    } else if awayTeamID == -1 {
        return &Prediction{HomeGoals: 1, AwayGoals: 0}, nil
    }

    leaderboard, err := l.db.GetCurrentSeasonLeaderboard()
    if err != nil {
        return nil, err
    }

    if len(leaderboard) == 0 {
        return nil, nil
    }

    sortedLeaderboard := l.sortLeaderboard(leaderboard)
    l.logLeaderboard(sortedLeaderboard)

    homePosition, awayPosition := l.getTeamPositions(homeTeamID, awayTeamID, sortedLeaderboard)
    positionDifference := abs(homePosition - awayPosition) / 2

    if homePosition < awayPosition {
        return &Prediction{HomeGoals: positionDifference, AwayGoals: 0}, nil
    } else if homePosition > awayPosition {
        return &Prediction{HomeGoals: 0, AwayGoals: positionDifference}, nil
    } else {
        return &Prediction{HomeGoals: 0, AwayGoals: 0}, nil
    }
}

```

With the above changes:

- If both homeTeamID and awayTeamID are -1, the function will return a prediction of (0, 0).
- If only homeTeamID is -1, the function will return a prediction of (0, 1).
- If only awayTeamID is -1, the function will return a prediction of (1, 0).

These checks are at the very beginning of the function, ensuring an early return and avoiding unnecessary calculations and function calls when the provided team IDs include `-1`.

## Follow-up prompt

> Great! Here's a unit test for that function:
>
> _(snipped: TestPredictWithHomeTeamLeading from predictor/leaderboard_difference_predictor_test.go)_
>
> Write tests for the -1 cases. I'd like to verify explicitly that the database wasn't queried.

Certainly! Let's add test cases to cover the scenarios where either `homeTeamID` or `awayTeamID` is `-1`.

To make sure the database wasn't queried, we'll leverage the `AssertExpectations` method from the mocking framework. Since we didn't set any expectations for the database call, `AssertExpectations` will ensure no database calls were made.

Here are the test cases:

```go
func TestPredictWithBothTeamsMinusOne(t *testing.T) {
	mockDB := new(database_test.MockDB)
	predictor := NewLeaderboardDifferencePredictor(mockDB)

	prediction, err := predictor.Predict(-1, -1)

	assert.NoError(t, err)
	assert.Equal(t, &Prediction{HomeGoals: 0, AwayGoals: 0}, prediction)
	mockDB.AssertExpectations(t)  // Verifies that no DB calls were made
}

func TestPredictWithHomeTeamMinusOne(t *testing.T) {
	mockDB := new(database_test.MockDB)
	predictor := NewLeaderboardDifferencePredictor(mockDB)

	prediction, err := predictor.Predict(-1, 4)

	assert.NoError(t, err)
	assert.Equal(t, &Prediction{HomeGoals: 0, AwayGoals: 1}, prediction)
	mockDB.AssertExpectations(t)  // Verifies that no DB calls were made
}

func TestPredictWithAwayTeamMinusOne(t *testing.T) {
	mockDB := new(database_test.MockDB)
	predictor := NewLeaderboardDifferencePredictor(mockDB)

	prediction, err := predictor.Predict(1, -1)

	assert.NoError(t, err)
	assert.Equal(t, &Prediction{HomeGoals: 1, AwayGoals: 0}, prediction)
	mockDB.AssertExpectations(t)  // Verifies that no DB calls were made
}

```

These tests validate the `-1` scenarios and ensure that the database isn't being queried when it shouldn't be.
