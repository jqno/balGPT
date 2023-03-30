# Allow unknown teams

> I'm writing an app that predicts soccer matches. Currently, it can't deal with the situation when we need to make a prediction concerting a team that isn't already known in the database. There are two parts to this; the front-end and the backend. Let's tackle the front-end first. What changes do I need to make to the code below to make it accept a named team that it doesn't know yet, either as a home team or away team, and send the value -1 to the backend when it does?
>
> _(snipped: full main.html)_

To make the front-end accept a named team that it doesn't know yet and send the value -1 to the backend, you can modify the event listeners for the 'input' event on both the home_team_name and away_team_name elements. You can set the value to -1 when the selected team is not found in the list of available teams. Here's the updated code:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Soccer Prediction</title>
  <link rel="stylesheet" href="/static/css/styles.css">
</head>
<body>
  <!-- The rest of your HTML code -->
  <script>
    // The rest of your JavaScript code

    document.getElementById('home_team_name')?.addEventListener('input', function (e) {
      const selectedOption = Array.from(e.target.list.options).find(option => option.textContent === e.target.value);
      document.getElementById('home_team_id').value = selectedOption ? selectedOption.getAttribute('data-value') : -1;
    });

    document.getElementById('away_team_name')?.addEventListener('input', function (e) {
      const selectedOption = Array.from(e.target.list.options).find(option => option.textContent === e.target.value);
      document.getElementById('away_team_id').value = selectedOption ? selectedOption.getAttribute('data-value') : -1;
    });

    // The rest of your event listeners
  </script>
</body>
</html>
```

These changes will make sure that if a team is not found in the list of available teams, the input value will be set to -1 for either the home_team_id or away_team_id, and this value will be sent to the backend when making a prediction. This allows your app to handle the case when a named team isn't already known in the database.

With this update, when the user enters a team name not found in the datalist, the front-end will send the value -1 for the corresponding team (home or away) to the backend, indicating that the team is not in the database. This will help the backend process the prediction request accordingly.

Now your front-end is able to handle unknown teams and send the appropriate value to the backend. You will need to update the backend logic to handle the case when a team has an ID of -1, and perform any necessary tasks such as fetching additional data or providing a prediction based on default values.

## Follow-up prompt

> Awesome, thank you! Now let's look at the back-end then. It consists of two functions: the handler and the predictor. The handler obviously should not fail when it encounters a -1. The predictor should not try to get the average number of goals in this case, but just predict a 0.
>
> Here's the handler:
>
> _(snipped: predictionHanlder from app.go)_
>
> Here's the predictor:
>
> _(snipped: Predict function from predictor.go)_

To modify the back-end to handle the case when the team ID is -1, you can update both the handler and predictor functions. In the handler, you can check if the team ID is -1 and set a flag accordingly. In the predictor, you can use this flag to decide whether to fetch the average goals or simply set the goals to 0.

Here's the updated handler:

```go
func handlePrediction(s *scraper.ScrapeData, p *predictor.Predictor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The rest of your handler code

		unknownHomeTeam := homeTeamID == -1
		unknownAwayTeam := awayTeamID == -1

		// The rest of your handler code

		prediction, err := p.Predict(homeTeamID, awayTeamID, unknownHomeTeam, unknownAwayTeam)
		if err != nil {
			http.Error(w, "Error while generating prediction.", http.StatusInternalServerError)
			return
		}

		// The rest of your handler code
	}
}
```

Here's the updated predictor:

```go
func (p *Predictor) Predict(homeTeamID, awayTeamID int, unknownHomeTeam, unknownAwayTeam bool) (*Prediction, error) {
	var homeTeamAvgGoals, awayTeamAvgGoals float64
	var err error

	if !unknownHomeTeam {
		homeTeamAvgGoals, err = p.DB.GetTeamAvgGoals(homeTeamID, true)
		if err != nil {
			return nil, err
		}
	}

	if !unknownAwayTeam {
		awayTeamAvgGoals, err = p.DB.GetTeamAvgGoals(awayTeamID, false)
		if err != nil {
			return nil, err
		}
	}

	prediction := &Prediction{
		HomeGoals: int(homeTeamAvgGoals + 0.5),
		AwayGoals: int(awayTeamAvgGoals + 0.5),
	}

	return prediction, nil
}
```

These changes will ensure that if a team ID is -1, the predictor will not attempt to fetch the average goals for that team and will instead set the predicted goals to 0. This should provide the desired behavior when dealing with unknown teams in your soccer match prediction app.

## Follow-up prompt

> This won't work, because the handler will throw an error when it can't find the team id in the database, which will happen when the id is -1. It should not do the query.
Also, I would prefer to keep the Predict signature the same, and just pass on the -1 as the team id. The predictor should then make its decision based on that.

Apologies for the confusion. To achieve the desired behavior, you can modify the Predictor function to handle the case when the team ID is -1 and update the GetTeamAvgGoals method to return 0 for unknown teams. Here's the updated predictor:

```go
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

Now, let's update the GetTeamAvgGoals method to return 0 for unknown teams:

```go
func (db *DB) GetTeamAvgGoals(teamID int, isHomeTeam bool) (float64, error) {
	if teamID == -1 {
		return 0, nil
	}

	// The rest of your GetTeamAvgGoals implementation
}
```

With these changes, the handler remains unchanged, and the predictor's Predict signature stays the same. The logic to handle the case when a team ID is -1 is now inside the GetTeamAvgGoals method, which returns 0 as the average number of goals for unknown teams.

This should provide the desired behavior when dealing with unknown teams in your soccer match prediction app.

## Follow-up prompt

> No, it's still wrong. I don't want to modify GetTeamAvgGoals, I want the predictor to not call it at all if a team id is -1. Also, the handler is still incorrect, because it's going to do an early return if a team id is -1, and it shouldn't do that.

Apologies for the confusion. I misunderstood your requirement earlier. Based on your clarification, here is the updated handler and predictor.

First, modify the handler to not return an error when a team ID is -1:

```go
func handlePrediction(s *scraper.ScrapeData, p *predictor.Predictor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		homeTeamIDStr := r.URL.Query().Get("home_team_id")
		awayTeamIDStr := r.URL.Query().Get("away_team_id")

		if homeTeamIDStr == "" || awayTeamIDStr == "" {
			http.Error(w, "Both home_team_id and away_team_id are required.", http.StatusBadRequest)
			return
		}

		homeTeamID, err := strconv.Atoi(homeTeamIDStr)
		if err != nil || (homeTeamID == -1 && homeTeamIDStr != "-1") {
			http.Error(w, "Invalid home_team_id.", http.StatusBadRequest)
			return
		}

		awayTeamID, err := strconv.Atoi(awayTeamIDStr)
		if err != nil || (awayTeamID == -1 && awayTeamIDStr != "-1") {
			http.Error(w, "Invalid away_team_id.", http.StatusBadRequest)
			return
		}

		// The rest of your handler code

		prediction, err := p.Predict(homeTeamID, awayTeamID)
		if err != nil {
			http.Error(w, "Error while generating prediction.", http.StatusInternalServerError)
			return
		}

		// The rest of your handler code
	}
}
```

Next, modify the predictor to not call GetTeamAvgGoals when a team ID is -1:

```go
func (p *Predictor) Predict(homeTeamID, awayTeamID int) (*Prediction, error) {
	homeTeamAvgGoals := 0.0
	awayTeamAvgGoals := 0.0
	var err error

	if homeTeamID != -1 {
		homeTeamAvgGoals, err = p.DB.GetTeamAvgGoals(homeTeamID, true)
		if err != nil {
			return nil, err
		}
	}

	if awayTeamID != -1 {
		awayTeamAvgGoals, err = p.DB.GetTeamAvgGoals(awayTeamID, false)
		if err != nil {
			return nil, err
		}
	}

	prediction := &Prediction{
		HomeGoals: int(homeTeamAvgGoals + 0.5),
		AwayGoals: int(awayTeamAvgGoals + 0.5),
	}

	return prediction, nil
}
```

Now, the `Predict` method checks if a team ID is -1 before calling `GetTeamAvgGoals`. If the team ID is -1, the method skips the call and directly assigns 0 as the average number of goals for that team.

The updated handler and predictor should now correctly handle the case when a team ID is -1, without modifying the `GetTeamAvgGoals` method.

## Follow-up prompt

> Why do we need to check that homeTeamID == -1 && homeTeamIDStr != "-1" after strconv.Atoi(homeTeamIDStr) ?

You are correct, and I apologize for the confusion. The check `(homeTeamID == -1 && homeTeamIDStr != "-1")` is redundant and unnecessary. The `strconv.Atoi` function will return an error if the conversion fails, so if `homeTeamID` is -1, it means the conversion was successful and the input string was indeed "-1".

Here is the updated handler without the unnecessary checks:

```go
func handlePrediction(s *scraper.ScrapeData, p *predictor.Predictor) http.HandlerFunc {
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

		// The rest of your handler code

		prediction, err := p.Predict(homeTeamID, awayTeamID)
		if err != nil {
			http.Error(w, "Error while generating prediction.", http.StatusInternalServerError)
			return
		}

		// The rest of your handler code
	}
}
```

Now the handler correctly processes the input without any unnecessary checks.
