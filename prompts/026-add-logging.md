# Add logging

> I forgot, we should probably also add some logging. Login attempts should be logged, and predictions as well.

Sure, let's add logging for login attempts and predictions. We'll use the standard Go logging package, "log", to accomplish this. I'll provide you with the code snippets to include logging in the appropriate places.

1. Logging login attempts:

You can add logging for both successful and failed login attempts in the checkAuth middleware function. Here's the updated code for the middleware:

(note: ChatGPT started hallucinating the wrong checkAuth function, so I wrote something myself here)

1. Logging predictions:

In the predictHandler function, you can add the following line after calling the predict function:

```go
log.Printf("Prediction for home_team_id=%d, away_team_id=%d: %d - %d", homeTeamID, awayTeamID, prediction.HomeGoals, prediction.AwayGoals)
```

Replace the existing predictHandler with the updated one:

```go
func predictHandler(w http.ResponseWriter, r *http.Request) {
    homeTeamID, _ := strconv.Atoi(r.URL.Query().Get("home_team_id"))
    awayTeamID, _ := strconv.Atoi(r.URL.Query().Get("away_team_id"))

    prediction, err := predict(homeTeamID, awayTeamID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Prediction for home_team_id=%d, away_team_id=%d: %d - %d", homeTeamID, awayTeamID, prediction.HomeGoals, prediction.AwayGoals)

    err = json.NewEncoder(w).Encode(prediction)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
```

With these changes, login attempts and predictions will be logged to the standard output. You can adjust the log format or output destination by configuring the log package according to your preferences.
