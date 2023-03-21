# Basic predictions

> Haha, this is amazing! Thank you!
>
> Let's create another package that comes up with a prediction for the next match, based on the information that is now available in the database. Assume two inputs: the home team and the away team. Let's use the primary key values from the teams table for that.
>
> It returns a struct with home_goals and away_goals.
>
> Can you come up with a good way to make the actual prediction? We might come back to this bit later.

To create a simple prediction, we can calculate the average goals scored by each team when playing at home or away, respectively. We will create a new package called predictor to handle the prediction logic.

First, create a new directory internal/predictor and a new file internal/predictor/predictor.go:

```go
package predictor

import (
	"github.com/YOUR_USERNAME/balgpt/internal/database"
)

type Prediction struct {
	HomeGoals int
	AwayGoals int
}

type Predictor struct {
	DB *database.DB
}

func NewPredictor(db *database.DB) *Predictor {
	return &Predictor{DB: db}
}

func (p *Predictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	homeTeamAvgGoals, err := p.DB.GetTeamAvgGoals(homeTeamID, true)
	if err != nil {
		return nil, err
	}

	awayTeamAvgGoals, err := p.DB.GetTeamAvgGoals(awayTeamID, false)
	if err != nil {
		return nil, err
	}

	prediction := &Prediction{
		HomeGoals: int(homeTeamAvgGoals + 0.5),
		AwayGoals: int(awayTeamAvgGoals + 0.5),
	}

	return prediction, nil
}
```

Next, add the required method GetTeamAvgGoals to the DB struct in the internal/database/database.go file:

```go
// Add this import
import "database/sql"

// Add the following method

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
```

Now, you can use the Predictor in your main application to make predictions for the upcoming matches. Note that this is a very basic prediction method that doesn't take many factors into account. You may want to improve the prediction algorithm later by considering more factors, such as recent form, head-to-head records, or even use machine learning techniques.

### Prompt

Nice.

Now let's create an endpoint. It takes the same input as the predict function and returns the prediction in json format. It first calls the scraper, and then runs the prediction.

First, let's add the necessary packages to handle HTTP and JSON. Add these imports to your main.go file:

```go
import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/YOUR_USERNAME/balgpt/internal/database"
	"github.com/YOUR_USERNAME/balgpt/internal/predictor"
	"github.com/YOUR_USERNAME/balgpt/internal/scraper"
)
```

Now, create the endpoint handler function in your main.go file:

```go
func handlePrediction(s *scraper.Scraper, p *predictor.Predictor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		homeTeamIDStr := r.URL.Query().Get("home_team_id")
		awayTeamIDStr := r.URL.Query().Get("away_team_id")

		if homeTeamIDStr == "" || awayTeamIDStr == "" {
			http.Error(w, "Both home_team_id and away_team_id are required.", http.StatusBadRequest)
			return
		}

		homeTeamID, err := strconv.Atoi(homeTeamIDStr)
		if err != nil {
			http.Error(w, "Invalid home_team_id.", http.StatusBadRequest)
			return
		}

		awayTeamID, err := strconv.Atoi(awayTeamIDStr)
		if err != nil {
			http.Error(w, "Invalid away_team_id.", http.StatusBadRequest)
			return
		}

		err = s.Scrape()
		if err != nil {
			http.Error(w, "Error while scraping data.", http.StatusInternalServerError)
			return
		}

		prediction, err := p.Predict(homeTeamID, awayTeamID)
		if err != nil {
			http.Error(w, "Error while generating prediction.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prediction)
	}
}
```

Now, register the endpoint in your main function. Add the following code after initializing the scraper and predictor instances:

```go
http.HandleFunc("/predict", handlePrediction(scraperInstance, predictorInstance))

port := "8080"
if envPort := os.Getenv("PORT"); envPort != "" {
	port = envPort
}

fmt.Printf("Listening on port %s...\n", port)
http.ListenAndServe(":"+port, nil)
```

Now, when you run your application, it will start an HTTP server listening on port 8080 (or the port specified in the PORT environment variable). To make a prediction, send a GET request to /predict with the home_team_id and away_team_id query parameters, for example:

```bash
http://localhost:8080/predict?home_team_id=1&away_team_id=2
```

This will return a JSON response containing the predicted home and away goals. Remember that the scraper will be called before making the prediction, so it might take a few moments to get the response.


