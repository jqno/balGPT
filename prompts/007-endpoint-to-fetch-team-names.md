# Endpoint to fetch team names

> We also need an endpoint that takes a team name as a string and returns the key.

Let's create an endpoint that takes the team name as a string and returns the primary key (ID) of the team. Add a new handler function in your main.go file:

```go
func handleTeamID(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamName := r.URL.Query().Get("team_name")

		if teamName == "" {
			http.Error(w, "team_name is required.", http.StatusBadRequest)
			return
		}

		teamID, err := db.GetTeamID(teamName)
		if err != nil {
			http.Error(w, "Error while fetching team ID.", http.StatusInternalServerError)
			return
		}

		if teamID == 0 {
			http.Error(w, "Team not found.", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"team_id": teamID})
	}
}
```

Now, register the endpoint in your main function. Add the following code after initializing the scraper and predictor instances:

```go
http.HandleFunc("/team_id", handleTeamID(db))
```

Now, when you run your application, you can send a GET request to /team_id with the team_name query parameter, for example:

```bash
http://localhost:8080/team_id?team_name=TeamName
```

This will return a JSON response containing the primary key (ID) of the team, like:

```json
{
  "team_id": 1
}
```

If the team is not found, the endpoint will return a 404 status with the message "Team not found."

## Note

Again, ChatGPT told me to add the endpoint in `main.go` instead of `app.go`. It'll do this more often, I'll stop mentioning it from now on. It referred to the database as `db` instead of `a.DB`, but that was as a result of assuming that the function should go in `main.go`.

More importantly, it forgot to generate a `GetTeamID` function in the database module, so I had to add that as well.

