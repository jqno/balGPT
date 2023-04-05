# Improve predictors

> I want to predict soccer matches in a Go app. I have a table with historical data in a database. I want to create a number of predictors which we'll define later. They all follow this structure:
>
> 	type Prediction struct {
> 		HomeGoals int
> 		AwayGoals int
> 	}
>
> 	type Predictor struct {
> 		DB *database.DB
> 	}
>
> 	func NewPredictor(db *database.DB) *Predictor {
> 		return &Predictor{DB: db}
> 	}
>
> 	func (p *Predictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
> 	    // ...
> 	}
>
> I want to start with the predictor that combines a bunch of other predictors and combines their results by taking the median of the HomeGoals and the AwayGoals respectively. Can you write that for me?

Yes, I can help you create a composite predictor that takes the median of HomeGoals and AwayGoals predictions from a list of other predictors. Here's the code:

```go
package main

import (
	"database"
	"errors"
	"sort"
)

type CompositePredictor struct {
	Predictor
	predictors []Predictor
}

func NewCompositePredictor(db *database.DB, predictors []Predictor) *CompositePredictor {
	return &CompositePredictor{Predictor{DB: db}, predictors}
}

func (c *CompositePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	if len(c.predictors) == 0 {
		return nil, errors.New("No predictors provided")
	}

	homeGoals := make([]int, 0, len(c.predictors))
	awayGoals := make([]int, 0, len(c.predictors))

	for _, predictor := range c.predictors {
		prediction, err := predictor.Predict(homeTeamID, awayTeamID)
		if err != nil {
			return nil, err
		}
		homeGoals = append(homeGoals, prediction.HomeGoals)
		awayGoals = append(awayGoals, prediction.AwayGoals)
	}

	sort.Ints(homeGoals)
	sort.Ints(awayGoals)

	medianHomeGoals := homeGoals[len(homeGoals)/2]
	medianAwayGoals := awayGoals[len(awayGoals)/2]

	return &Prediction{HomeGoals: medianHomeGoals, AwayGoals: medianAwayGoals}, nil
}
```

Now you can use the `NewCompositePredictor` function to create a composite predictor that takes a list of other predictors and combines their results by taking the median of the HomeGoals and the AwayGoals respectively.

## Follow-up prompt

> Can you modify the code so that the CompositePredictor uses the same interface as the individual predictors? Also, there will be more than 2 predictors -- I don't know how many. So please adjust the median calculation for that.

Sure, I can modify the code to make the `CompositePredictor` implement the same interface as the individual predictors. I'll also adjust the median calculation to handle any number of predictors.

```go
package main

import (
	"database"
	"errors"
	"sort"
)

type Prediction struct {
	HomeGoals int
	AwayGoals int
}

type Predictor interface {
	Predict(homeTeamID, awayTeamID int) (*Prediction, error)
}

type BasePredictor struct {
	DB *database.DB
}

func NewPredictor(db *database.DB) *BasePredictor {
	return &BasePredictor{DB: db}
}

type CompositePredictor struct {
	BasePredictor
	predictors []Predictor
}

func NewCompositePredictor(db *database.DB, predictors []Predictor) *CompositePredictor {
	return &CompositePredictor{BasePredictor{DB: db}, predictors}
}

func (c *CompositePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	if len(c.predictors) == 0 {
		return nil, errors.New("No predictors provided")
	}

	homeGoals := make([]int, 0, len(c.predictors))
	awayGoals := make([]int, 0, len(c.predictors))

	for _, predictor := range c.predictors {
		prediction, err := predictor.Predict(homeTeamID, awayTeamID)
		if err != nil {
			return nil, err
		}
		homeGoals = append(homeGoals, prediction.HomeGoals)
		awayGoals = append(awayGoals, prediction.AwayGoals)
	}

	sort.Ints(homeGoals)
	sort.Ints(awayGoals)

	medianHomeGoals := calculateMedian(homeGoals)
	medianAwayGoals := calculateMedian(awayGoals)

	return &Prediction{HomeGoals: medianHomeGoals, AwayGoals: medianAwayGoals}, nil
}

func calculateMedian(sortedValues []int) int {
	middle := len(sortedValues) / 2

	if len(sortedValues)%2 == 0 {
		return (sortedValues[middle-1] + sortedValues[middle]) / 2
	}

	return sortedValues[middle]
}
```

Now the `CompositePredictor` uses the same `Predictor` interface as individual predictors. The median calculation has also been adjusted to work with any number of predictors.

## Follow-up prompt

> Let's define a package predictor and have all the predictors in their own file. The CompositePredictor should be in a separate file from the generic stuff.

I understand that you want to organize the code into a `predictor` package and have separate files for the generic and composite predictors. Here's how you can structure your package:

Create a directory named `predictor`, and inside the directory, create three files: `predictor.go`, `base_predictor.go`, and `composite_predictor.go`.

1. predictor.go: This file will contain the Predictor interface and the Prediction struct.

```go
// predictor/predictor.go
package predictor

type Prediction struct {
	HomeGoals int
	AwayGoals int
}

type Predictor interface {
	Predict(homeTeamID, awayTeamID int) (*Prediction, error)
}
```

1. base_predictor.go: This file will contain the BasePredictor struct and its constructor.

```go
// predictor/base_predictor.go
package predictor

import "database"

type BasePredictor struct {
	DB *database.DB
}

func NewBasePredictor(db *database.DB) *BasePredictor {
	return &BasePredictor{DB: db}
}
```

1. composite_predictor.go: This file will contain the CompositePredictor struct, its constructor, and the Predict method.

```go
// predictor/composite_predictor.go
package predictor

import (
	"errors"
	"sort"
)

type CompositePredictor struct {
	BasePredictor
	predictors []Predictor
}

func NewCompositePredictor(db *database.DB, predictors []Predictor) *CompositePredictor {
	return &CompositePredictor{BasePredictor{DB: db}, predictors}
}

func (c *CompositePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	if len(c.predictors) == 0 {
		return nil, errors.New("No predictors provided")
	}

	homeGoals := make([]int, 0, len(c.predictors))
	awayGoals := make([]int, 0, len(c.predictors))

	for _, predictor := range c.predictors {
		prediction, err := predictor.Predict(homeTeamID, awayTeamID)
		if err != nil {
			return nil, err
		}
		homeGoals = append(homeGoals, prediction.HomeGoals)
		awayGoals = append(awayGoals, prediction.AwayGoals)
	}

	sort.Ints(homeGoals)
	sort.Ints(awayGoals)

	medianHomeGoals := calculateMedian(homeGoals)
	medianAwayGoals := calculateMedian(awayGoals)

	return &Prediction{HomeGoals: medianHomeGoals, AwayGoals: medianAwayGoals}, nil
}

func calculateMedian(sortedValues []int) int {
	middle := len(sortedValues) / 2

	if len(sortedValues)%2 == 0 {
		return (sortedValues[middle-1] + sortedValues[middle]) / 2
	}

	return sortedValues[middle]
}
```

Now you have a `predictor` package with the `Predictor` interface and `Prediction` struct defined in `predictor.go`, the `BasePredictor` struct in `base_predictor.go`, and the `CompositePredictor` struct in `composite_predictor.go`.

## Follow-up prompt

> Explain the purpose of the BasePredictor.

The `BasePredictor` serves as a foundation for building various specific predictor implementations. Its primary purpose is to store common components that all predictors might need, such as a reference to the database. In this case, the `BasePredictor` holds a reference to the `DB` object, which allows access to the historical data.

Using a `BasePredictor` struct is helpful because it promotes code reuse and provides a more organized structure. Instead of repeating the same code for initializing and accessing the database in each predictor, you can embed the `BasePredictor` in your specific predictor implementations, inheriting its properties and methods.

However, note that the current code uses the `BasePredictor` only in the `CompositePredictor`. To utilize its full potential, you should embed it in all of your other predictor implementations as well. This way, when you create a new predictor, you can simply embed the `BasePredictor` and focus on the unique logic for that predictor, knowing that the common components, like the database reference, are already handled.

In summary, the `BasePredictor` serves as a base structure for specific predictor implementations, holding common properties (such as the database reference) and potentially common methods that all predictors might need. This helps improve code organization, reusability, and maintainability.

## Follow-up prompt

> Awesome, thank you.
>
> Let's create our first, super-simple predictor: the home advantage predictor, which always predicts a score of 1-0 for the home team.

Sure, here's the code for the `HomeAdvantagePredictor`:

1. Create a new file home_advantage_predictor.go inside the predictor directory.

```go
// predictor/home_advantage_predictor.go
package predictor

type HomeAdvantagePredictor struct {
	BasePredictor
}

func NewHomeAdvantagePredictor(db *database.DB) *HomeAdvantagePredictor {
	return &HomeAdvantagePredictor{BasePredictor{DB: db}}
}

func (h *HomeAdvantagePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	return &Prediction{HomeGoals: 1, AwayGoals: 0}, nil
}
```

The `HomeAdvantagePredictor` struct embeds the `BasePredictor` and implements the `Predict` method, always returning a score of 1-0 for the home team. You can create a new `HomeAdvantagePredictor` using the `NewHomeAdvantagePredictor` constructor.

Now, you can use the `HomeAdvantagePredictor` in combination with other predictors inside the `CompositePredictor`.

## Follow-up prompt

> Since both the HomeAdvantagePredictor and the CompositePredictor don't need a database, and there won't be other shared infrastructure, let's remove the BasePredictor altogether and pass the database directly to those predictors that need it in the future.

Sure, I can remove the `BasePredictor` and update the `HomeAdvantagePredictor` and `CompositePredictor` accordingly. Here's the updated code:

1. Remove the base_predictor.go file from the predictor directory.
2. Update the composite_predictor.go file:

```go
// predictor/composite_predictor.go
package predictor

import (
	"errors"
	"sort"
)

type CompositePredictor struct {
	predictors []Predictor
}

func NewCompositePredictor(predictors []Predictor) *CompositePredictor {
	return &CompositePredictor{predictors: predictors}
}

func (c *CompositePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	if len(c.predictors) == 0 {
		return nil, errors.New("No predictors provided")
	}

	homeGoals := make([]int, 0, len(c.predictors))
	awayGoals := make([]int, 0, len(c.predictors))

	for _, predictor := range c.predictors {
		prediction, err := predictor.Predict(homeTeamID, awayTeamID)
		if err != nil {
			return nil, err
		}
		homeGoals = append(homeGoals, prediction.HomeGoals)
		awayGoals = append(awayGoals, prediction.AwayGoals)
	}

	sort.Ints(homeGoals)
	sort.Ints(awayGoals)

	medianHomeGoals := calculateMedian(homeGoals)
	medianAwayGoals := calculateMedian(awayGoals)

	return &Prediction{HomeGoals: medianHomeGoals, AwayGoals: medianAwayGoals}, nil
}

func calculateMedian(sortedValues []int) int {
	middle := len(sortedValues) / 2

	if len(sortedValues)%2 == 0 {
		return (sortedValues[middle-1] + sortedValues[middle]) / 2
	}

	return sortedValues[middle]
}
```

1. Update the home_advantage_predictor.go file:

```go
// predictor/home_advantage_predictor.go
package predictor

type HomeAdvantagePredictor struct {
}

func NewHomeAdvantagePredictor() *HomeAdvantagePredictor {
	return &HomeAdvantagePredictor{}
}

func (h *HomeAdvantagePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	return &Prediction{HomeGoals: 1, AwayGoals: 0}, nil
}
```

With these changes, the `BasePredictor` has been removed, and the `CompositePredictor` and `HomeAdvantagePredictor` no longer use or require a reference to the database. You can pass the database directly to future predictors that need it.

## Follow-up prompt

> Awesome!
>
> In my app, I have a predictor of type `*predictor.Predictor`. How do I initialize it as a CompositePredictor with a HomeAdvantagePredictor, and how do I call it?

Since `Predictor` is an interface, you can't directly create an instance of type `*predictor.Predictor`. Instead, you should initialize a variable of type `predictor.Predictor` and assign it an instance of `CompositePredictor` containing a `HomeAdvantagePredictor`. Then you can call the `Predict` method using the variable.

Here's an example:

```go
package main

import (
	"fmt"
	"your_project_name/predictor" // Replace "your_project_name" with the name of your project
)

func main() {
	// Create a HomeAdvantagePredictor
	homeAdvantagePredictor := predictor.NewHomeAdvantagePredictor()

	// Initialize a CompositePredictor with a HomeAdvantagePredictor
	predictors := []predictor.Predictor{homeAdvantagePredictor}
	compositePredictor := predictor.NewCompositePredictor(predictors)

	// Assign the CompositePredictor instance to a variable of type predictor.Predictor
	var myPredictor predictor.Predictor = compositePredictor

	// Call the Predict method using the myPredictor variable
	homeTeamID := 1
	awayTeamID := 2
	prediction, err := myPredictor.Predict(homeTeamID, awayTeamID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Predicted score: %d-%d\n", prediction.HomeGoals, prediction.AwayGoals)
}
```

This example demonstrates how to initialize a `CompositePredictor` with a `HomeAdvantagePredictor`, assign it to a variable of type `predictor.Predictor`, and call the `Predict` method to get a prediction. Replace `"your_project_name"` with the name of your project in the import statement to properly import the `predictor` package.

## Follow-up prompt

> Let's create another predictor, that predicts for each given team the average of the number of goals scored in the past 8 matches. Note that the team could have been either the home team or the away team in those matches.
>
> It needs to access a Postgres database for that. Here's what the corresponding table looks like:
>
>   CREATE TABLE matches (
>       id SERIAL PRIMARY KEY,
>       home_team INTEGER REFERENCES teams(id),
>       away_team INTEGER REFERENCES teams(id),
>       home_goals INTEGER,
>       away_goals INTEGER,
>       date DATE NOT NULL
>   );
>
> I already have a database.go file with this struct:
>
>     type DB struct {
>         Conn *sql.DB
>     }
>
> You can create the functions you need, but they have to be of this form:
>
>   func (db *DB) FunctionName() (returntypes, error) {
>       // ...
>   }
> 
> Don't create new structs and be careful that only goals scored by the given teamID are queried.

I understand that you want to create a new predictor that predicts the average number of goals scored by each team in their past 8 matches. I'll add the required functions to access the Postgres database and create a new predictor that uses these functions.

First, let's add the required functions in your `database.go` file:

```go
// database/database.go
package database

import (
	"database/sql"
	"fmt"
)

type DB struct {
	Conn *sql.DB
}

func (db *DB) AverageGoalsInLastMatches(teamID int, numberOfMatches int) (float64, error) {
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

Now, let's create the new predictor called `AverageGoalsPredictor` in the `predictor` package:

1. Create a new file average_goals_predictor.go inside the predictor directory.

```go
// predictor/average_goals_predictor.go
package predictor

import (
	"database"
	"math"
)

type AverageGoalsPredictor struct {
	db *database.DB
}

func NewAverageGoalsPredictor(db *database.DB) *AverageGoalsPredictor {
	return &AverageGoalsPredictor{db: db}
}

func (a *AverageGoalsPredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	homeAvgGoals, err := a.db.AverageGoalsInLastMatches(homeTeamID, 8)
	if err != nil {
		return nil, err
	}

	awayAvgGoals, err := a.db.AverageGoalsInLastMatches(awayTeamID, 8)
	if err != nil {
		return nil, err
	}

	homeGoals := int(math.Round(homeAvgGoals))
	awayGoals := int(math.Round(awayAvgGoals))

	return &Prediction{HomeGoals: homeGoals, AwayGoals: awayGoals}, nil
}
```

Now you have a new predictor called `AverageGoalsPredictor` that calculates the average number of goals scored by each team in their past 8 matches and predicts the outcome based on these averages. You can use this predictor in combination with other predictors inside the `CompositePredictor`.

## Follow-up prompt

> The next predictor I have in mind, might not return a prediction at all. Please update the CompositePredictor accordingly. Does idiomatic Go use something like an Optional or Maybe type? If so, use it and update the function signatures.

In Go, the idiomatic way to represent an optional value is to use a pointer and return `nil` when there's no value to return. This is already the approach used in the existing predictor code, so we don't need to change the function signatures. To update the `CompositePredictor` to handle predictors that might not return a prediction, we just need to modify the `Predict` method to check for `nil` predictions and handle them accordingly. Here's the updated `composite_predictor.go` file:

```go
// predictor/composite_predictor.go
package predictor

import (
	"errors"
	"sort"
)

type CompositePredictor struct {
	predictors []Predictor
}

func NewCompositePredictor(predictors []Predictor) *CompositePredictor {
	return &CompositePredictor{predictors: predictors}
}

func (c *CompositePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	if len(c.predictors) == 0 {
		return nil, errors.New("No predictors provided")
	}

	homeGoals := make([]int, 0, len(c.predictors))
	awayGoals := make([]int, 0, len(c.predictors))

	for _, predictor := range c.predictors {
		prediction, err := predictor.Predict(homeTeamID, awayTeamID)
		if err != nil {
			return nil, err
		}
		if prediction != nil {
			homeGoals = append(homeGoals, prediction.HomeGoals)
			awayGoals = append(awayGoals, prediction.AwayGoals)
		}
	}

	if len(homeGoals) == 0 || len(awayGoals) == 0 {
		return nil, errors.New("No predictions available")
	}

	sort.Ints(homeGoals)
	sort.Ints(awayGoals)

	medianHomeGoals := calculateMedian(homeGoals)
	medianAwayGoals := calculateMedian(awayGoals)

	return &Prediction{HomeGoals: medianHomeGoals, AwayGoals: medianAwayGoals}, nil
}

func calculateMedian(sortedValues []int) int {
	middle := len(sortedValues) / 2

	if len(sortedValues)%2 == 0 {
		return (sortedValues[middle-1] + sortedValues[middle]) / 2
	}

	return sortedValues[middle]
}
```

Now the `CompositePredictor` checks if the prediction returned by a predictor is `nil` and only appends the home and away goals to the respective slices if it's not. If there are no predictions available after iterating through all predictors, the `Predict` method returns an error indicating that there are no predictions available. This approach allows you to use predictors that might not return a prediction, and the `CompositePredictor` will still function correctly.

## Follow-up prompt

> Awesome!
>
> Now let's write the predictor I was talking about: for a given home team and away team, it returns the scores of the match that the same teams played last year, where the home team was the home team and the away team was the away team. Note that we don't know exactly when that match was played. We do know it was the most recent match between the two, though. Note that such a match might not exist in the database.

To create this predictor, we'll first add a function in your `database.go` file to fetch the scores of the most recent match between the two teams played last year:

```go
// database/database.go
package database

import (
	"database/sql"
	"fmt"
	"time"
)

type DB struct {
	Conn *sql.DB
}

// ...

func (db *DB) LastYearMatchScores(homeTeamID, awayTeamID int) (*Prediction, error) {
	query := `
		SELECT home_goals, away_goals
		FROM matches
		WHERE home_team = $1 AND away_team = $2
		AND date > (NOW() - INTERVAL '1 year')
		ORDER BY date DESC
		LIMIT 1;
	`

	var homeGoals, awayGoals int
	err := db.Conn.QueryRow(query, homeTeamID, awayTeamID).Scan(&homeGoals, &awayGoals)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Error fetching last year's match scores: %v", err)
	}

	return &Prediction{HomeGoals: homeGoals, AwayGoals: awayGoals}, nil
}
```

Now, let's create the new predictor called `LastYearMatchPredictor` in the `predictor` package:

1. Create a new file last_year_match_predictor.go inside the predictor directory.

```go
// predictor/last_year_match_predictor.go
package predictor

import (
	"database"
)

type LastYearMatchPredictor struct {
	db *database.DB
}

func NewLastYearMatchPredictor(db *database.DB) *LastYearMatchPredictor {
	return &LastYearMatchPredictor{db: db}
}

func (l *LastYearMatchPredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	prediction, err := l.db.LastYearMatchScores(homeTeamID, awayTeamID)
	if err != nil {
		return nil, err
	}

	if prediction == nil {
		return nil, nil
	}

	return &Prediction{HomeGoals: prediction.HomeGoals, AwayGoals: prediction.AwayGoals}, nil
}

```

Now you have a new predictor called `LastYearMatchPredictor` that returns the scores of the most recent match between the two teams played last year, where the home team was the home team and the away team was the away team. The predictor returns `nil` if no such match exists in the database. You can use this predictor in combination with other predictors inside the `CompositePredictor`.

## Follow-up prompt

> Last year's match could have been played more than a year ago, depending on when exactly in the season it was played. Let's lose the AND condition, I think the ORDER and LIMIT clauses are sufficient.

Sure, we can modify the `LastYearMatchScores` function in the `database.go` file to remove the `AND` condition related to the date. Here's the updated function:

```go
// database/database.go
package database

import (
	"database/sql"
	"fmt"
)

type DB struct {
	Conn *sql.DB
}

// ...

func (db *DB) LastYearMatchScores(homeTeamID, awayTeamID int) (*Prediction, error) {
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
			return nil, nil
		}
		return nil, fmt.Errorf("Error fetching last match scores: %v", err)
	}

	return &Prediction{HomeGoals: homeGoals, AwayGoals: awayGoals}, nil
}
```

This change will make the `LastYearMatchPredictor` return the scores of the most recent match between the two teams, without any constraint on the date. The `ORDER BY` and `LIMIT` clauses will ensure you get the latest match between the teams.

## Follow-up prompt

> Also, database.go doesn't (and shouldn't!) know about the Prediction struct. Let's return two ints instead.

You're right. It's better to decouple the `database` package from the `predictor` package. We can modify the `LastYearMatchScores` function in the `database.go` file to return two `int` values instead of a `*Prediction`. Here's the updated function:

```go
// database/database.go
package database

import (
	"database/sql"
	"fmt"
)

type DB struct {
	Conn *sql.DB
}

// ...

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
```

Now, we should update the `Predict` method in the `LastYearMatchPredictor` to handle the returned `int` values:

```go
// predictor/last_year_match_predictor.go
package predictor

import (
	"database"
)

type LastYearMatchPredictor struct {
	db *database.DB
}

func NewLastYearMatchPredictor(db *database.DB) *LastYearMatchPredictor {
	return &LastYearMatchPredictor{db: db}
}

func (l *LastYearMatchPredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	homeGoals, awayGoals, err := l.db.LastYearMatchScores(homeTeamID, awayTeamID)
	if err != nil {
		return nil, err
	}

	if homeGoals == 0 && awayGoals == 0 {
		return nil, nil
	}

	return &Prediction{HomeGoals: homeGoals, AwayGoals: awayGoals}, nil
}
```

Now the `LastYearMatchScores` function in the `database` package returns two `int` values and the `LastYearMatchPredictor` handles them correctly, keeping the `database` and `predictor` packages decoupled.

## Follow-up prompt

> Great! Now let's add a predictor that does the same thing as the LastYearMatchPredictor, but where the home team and the away team are flipped. Let's try to reuse as much code from LastYearMatchPredictor as possible.

To create a predictor that returns the scores of the most recent match between the two teams with flipped home and away teams, we can reuse the existing `LastYearMatchScores` function in the `database.go` file and create a new predictor called `FlippedLastYearMatchPredictor`. Here's how to do it:

1. Create a new file flipped_last_year_match_predictor.go inside the predictor directory.

```go
// predictor/flipped_last_year_match_predictor.go
package predictor

import (
	"database"
)

type FlippedLastYearMatchPredictor struct {
	db *database.DB
}

func NewFlippedLastYearMatchPredictor(db *database.DB) *FlippedLastYearMatchPredictor {
	return &FlippedLastYearMatchPredictor{db: db}
}

func (f *FlippedLastYearMatchPredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	awayGoals, homeGoals, err := f.db.LastYearMatchScores(awayTeamID, homeTeamID)
	if err != nil {
		return nil, err
	}

	if homeGoals == 0 && awayGoals == 0 {
		return nil, nil
	}

	return &Prediction{HomeGoals: homeGoals, AwayGoals: awayGoals}, nil
}
```

This new predictor is very similar to the `LastYearMatchPredictor`, but it flips the home and away team IDs when calling the `LastYearMatchScores` function. Consequently, the predictor returns the scores of the most recent match between the two teams with flipped home and away teams. You can use this predictor in combination with other predictors inside the `CompositePredictor`.

## Follow-up prompt

> Can we modify the existing LastYearMatchPredictor a little bit and do something clever with two constructors, instead of adding a FlippedLastYearMatchPredictor that duplicates a lot of code?

Yes, we can modify the existing `LastYearMatchPredictor` to handle both cases by introducing a field in the struct to control the flipping behavior. Here's the modified `LastYearMatchPredictor`:

```go
// predictor/last_year_match_predictor.go
package predictor

import (
	"database"
)

type LastYearMatchPredictor struct {
	db          *database.DB
	flippedTeams bool
}

func NewLastYearMatchPredictor(db *database.DB) *LastYearMatchPredictor {
	return &LastYearMatchPredictor{db: db, flippedTeams: false}
}

func NewFlippedLastYearMatchPredictor(db *database.DB) *LastYearMatchPredictor {
	return &LastYearMatchPredictor{db: db, flippedTeams: true}
}

func (l *LastYearMatchPredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	var team1ID, team2ID int

	if l.flippedTeams {
		team1ID = awayTeamID
		team2ID = homeTeamID
	} else {
		team1ID = homeTeamID
		team2ID = awayTeamID
	}

	team1Goals, team2Goals, err := l.db.LastYearMatchScores(team1ID, team2ID)
	if err != nil {
		return nil, err
	}

	if team1Goals == 0 && team2Goals == 0 {
		return nil, nil
	}

	if l.flippedTeams {
		return &Prediction{HomeGoals: team2Goals, AwayGoals: team1Goals}, nil
	}
	return &Prediction{HomeGoals: team1Goals, AwayGoals: team2Goals}, nil
}
```

We added a new field `flippedTeams` to the `LastYearMatchPredictor` struct, which controls whether the home and away teams are flipped or not. We also introduced a new constructor `NewFlippedLastYearMatchPredictor` that sets this field to `true`. The `Predict` method has been updated to handle both cases based on the value of the `flippedTeams` field.

Now, you can create a regular `LastYearMatchPredictor` with `NewLastYearMatchPredictor(db)` and a flipped version with `NewFlippedLastYearMatchPredictor(db)`. This approach avoids code duplication and handles both cases within a single predictor.

## Follow-up prompt

> Can the Predict function in this approach be simplified further?

Yes, we can simplify the `Predict` function further by removing the conditional checks for `flippedTeams` inside the function and only swap the team IDs before calling `LastYearMatchScores` based on the value of `flippedTeams`. Here's the updated `Predict` function:

```go
func (l *LastYearMatchPredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	if l.flippedTeams {
		homeTeamID, awayTeamID = awayTeamID, homeTeamID
	}

	homeGoals, awayGoals, err := l.db.LastYearMatchScores(homeTeamID, awayTeamID)
	if err != nil {
		return nil, err
	}

	if homeGoals == 0 && awayGoals == 0 {
		return nil, nil
	}

	return &Prediction{HomeGoals: homeGoals, AwayGoals: awayGoals}, nil
}
```

This implementation is cleaner, as it swaps the team IDs before calling `LastYearMatchScores` based on the `flippedTeams` flag and doesn't require further conditional checks in the function. The resulting prediction is directly returned without additional modifications.

## Follow-up prompt

> We should probably flip the values back though, otherwise the goals will be assigned to the wrong teams.

Apologies for the confusion. You're right; we should flip the goals back when the `flippedTeams` flag is true. Here's the corrected `Predict` function:

```go
func (l *LastYearMatchPredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	if l.flippedTeams {
		homeTeamID, awayTeamID = awayTeamID, homeTeamID
	}

	homeGoals, awayGoals, err := l.db.LastYearMatchScores(homeTeamID, awayTeamID)
	if err != nil {
		return nil, err
	}

	if homeGoals == 0 && awayGoals == 0 {
		return nil, nil
	}

	if l.flippedTeams {
		homeGoals, awayGoals = awayGoals, homeGoals
	}

	return &Prediction{HomeGoals: homeGoals, AwayGoals: awayGoals}, nil
}
```

Now, the function correctly flips the goals back when the `flippedTeams` flag is true, ensuring that the goals are assigned to the correct teams in the resulting `Prediction` object.

## Follow-up prompt

> Great. One last thing, 0-0 is a perfectly valid soccer match result, so we shouldn't return nils in that case.

Apologies for the oversight. Let's modify the `Predict` function to handle the 0-0 case correctly by checking if the returned error is `sql.ErrNoRows`. Here's the updated `Predict` function:

```go
func (l *LastYearMatchPredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	if l.flippedTeams {
		homeTeamID, awayTeamID = awayTeamID, homeTeamID
	}

	homeGoals, awayGoals, err := l.db.LastYearMatchScores(homeTeamID, awayTeamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if l.flippedTeams {
		homeGoals, awayGoals = awayGoals, homeGoals
	}

	return &Prediction{HomeGoals: homeGoals, AwayGoals: awayGoals}, nil
}
```

Now, the function will only return `nil, nil` if there is no match found in the database for the given team IDs. In case of a 0-0 match result, a `Prediction` object with `HomeGoals` and `AwayGoals` set to 0 will be returned correctly.

## Follow-up prompt

> Yeah!
>
> How do I update the initialization for the CompositePredictor with all of the predictors we've defined so far? Please write the initialization in as few lines as possible.

You've reached the current usage cap for GPT-4. You can continue with the default model now, or try again after 3:30 PM. Learn moreUse default model

## Continuing in a new conversation with GPT-3.5

> Can I inline the predictors variable in this bit of Go code?
>
> 	predictors := []predictor.Predictor{
> 		predictor.NewHomeAdvantagePredictor(),
> 		predictor.NewAverageGoalsPredictor(db),
> 		predictor.NewLastYearMatchPredictor(db),
> 		predictor.NewFlippedLastYearMatchPredictor(db),
> 	}
> 	predictor := predictor.NewCompositePredictor(predictors)

Yes, you can inline the `predictors` variable in the Go code. Here is an example of how to do it:

```scss
predictor := predictor.NewCompositePredictor([]predictor.Predictor{
    predictor.NewHomeAdvantagePredictor(),
    predictor.NewAverageGoalsPredictor(db),
    predictor.NewLastYearMatchPredictor(db),
    predictor.NewFlippedLastYearMatchPredictor(db),
})
```

This will create a `CompositePredictor` instance using a slice of `Predictor` objects as an argument to the constructor, without the need for an intermediate variable.

## Follow-up prompt

> Can I get rid of the []predictor.Predictor bit?

Yes, you can get rid of the `[]predictor.Predictor` type declaration by using a variadic function instead of a slice.

Here's how you can modify the code to use a variadic function:

```css
predictor := predictor.NewCompositePredictor(
	predictor.NewHomeAdvantagePredictor(),
	predictor.NewAverageGoalsPredictor(db),
	predictor.NewLastYearMatchPredictor(db),
	predictor.NewFlippedLastYearMatchPredictor(db),
)
```

With this change, the `NewCompositePredictor` function is defined to accept zero or more arguments of type `predictor.Predictor`. This allows you to call the function without explicitly creating a slice and passing it as an argument.

Note that if the `NewCompositePredictor` function is defined to only accept a slice, then you won't be able to use this approach.

## Follow-up prompt

> NewCompositePredictor is currently defined like this:
>
> 	func NewCompositePredictor(predictors []Predictor) *CompositePredictor {
> 		return &CompositePredictor{predictors: predictors}
> 	}
>
> How do I make it variadic?

To make the `NewCompositePredictor` function variadic, you need to change its parameter list to accept zero or more arguments of type `Predictor`. Here's how you can modify the function to be variadic:

```go
func NewCompositePredictor(predictors ...Predictor) *CompositePredictor {
	return &CompositePredictor{predictors: predictors}
}
```

With this change, the `NewCompositePredictor` function will accept zero or more arguments of type `Predictor`. You can call the function with any number of arguments of type `Predictor`, separated by commas. For example:

```css
predictor := predictor.NewCompositePredictor(
	predictor.NewHomeAdvantagePredictor(),
	predictor.NewAverageGoalsPredictor(db),
	predictor.NewLastYearMatchPredictor(db),
	predictor.NewFlippedLastYearMatchPredictor(db),
)
```

In this example, the `NewCompositePredictor` function is called with four arguments of type `Predictor`, each separated by a comma.

## Follow-up prompt

> Can you add logging for each predictor output in this function? Also, is there a way to print the predictor's runtime type name along with it?
>
> _(snipped: current implementation of composite_predictor.go's Predict funtion)_

Sure, here's the modified function with logging added for each predictor's output, including its runtime type name:

```go
func (c *CompositePredictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	if len(c.predictors) == 0 {
		return nil, errors.New("No predictors provided")
	}

	homeGoals := make([]int, 0, len(c.predictors))
	awayGoals := make([]int, 0, len(c.predictors))

	for _, predictor := range c.predictors {
		predictionStartTime := time.Now()
		prediction, err := predictor.Predict(homeTeamID, awayTeamID)
		predictionEndTime := time.Now()

		if err != nil {
			return nil, err
		}
		if prediction != nil {
			log.Printf("Prediction from %T: %+v, runtime: %v", predictor, prediction, predictionEndTime.Sub(predictionStartTime))
			homeGoals = append(homeGoals, prediction.HomeGoals)
			awayGoals = append(awayGoals, prediction.AwayGoals)
		}
	}

	if len(homeGoals) == 0 || len(awayGoals) == 0 {
		return nil, errors.New("No predictions available")
	}

	sort.Ints(homeGoals)
	sort.Ints(awayGoals)

	medianHomeGoals := calculateMedian(homeGoals)
	medianAwayGoals := calculateMedian(awayGoals)

	return &Prediction{HomeGoals: medianHomeGoals, AwayGoals: medianAwayGoals}, nil
}
```

This implementation uses the Go `log` package to print a message for each predictor's output, including the runtime type name and the elapsed time for the prediction. Note that this implementation assumes that the predictor's type is directly accessible from the predictor's variable, which is usually the case in Go, but may not be true for all languages or implementations.
