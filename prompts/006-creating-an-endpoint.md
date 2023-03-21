# Creating an endpoint

> Nice.
>
> Now let's create an endpoint. It takes the same input as the predict function and returns the prediction in json format. It first calls the scraper, and then runs the prediction.

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

## Note

It got the type of `scraper.ScrapeData` wrong, even if it generated it only a few prompts before.

It suggested I add stuff to the main function, but it should have directed me to `internal/app/app.go`.

Finally, it neglected to initialize the scraper and the predictor, so I had to do that manually.
