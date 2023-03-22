# Autocompleting inputs

> Let's make both input fields autocompleting based on the content of the teams table in the db. We can generate the list in the backend; there's already a func fetchTeamsFromDB() ([]Team, error) that fetches the teams. How do I update the indexHandler and index.html to reflect this?

To update the indexHandler and index.html to support autocompleting team names based on the content of the teams table, you can follow these steps:

1. Modify the TemplateData struct to include the teams:

```go
type TemplateData struct {
	ClientID  string
	ApiBaseURL string
	Teams     []Team
}
```

1. Update the indexHandler to fetch the teams and pass them to the template:

```go
func indexHandler(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	apiBaseURL := os.Getenv("API_BASE_URL")

	if clientID == "" || apiBaseURL == "" {
		http.Error(w, "Missing environment variables", http.StatusInternalServerError)
		return
	}

	teams, err := fetchTeamsFromDB()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching teams: %v", err), http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		ClientID:  clientID,
		ApiBaseURL: apiBaseURL,
		Teams:     teams,
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}

3. Update the `index.html` file to include the list of teams as a JavaScript variable and use the HTML5 `datalist` element for autocompletion:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <!-- ... -->
</head>
<body>
  <div id="signedInContent" style="display:none;">
    <h1>Soccer Prediction</h1>
    <div>
      <label for="home_team_id">Home Team:</label>
      <input list="teams" id="home_team_id" required>
      <datalist id="teams">
        {{range .Teams}}
        <option value="{{.ID}}">{{.Name}}</option>
        {{end}}
      </datalist>
    </div>
    <div>
      <label for="away_team_id">Away Team:</label>
      <input list="teams" id="away_team_id" required>
      <datalist id="teams">
        {{range .Teams}}
        <option value="{{.ID}}">{{.Name}}</option>
        {{end}}
      </datalist>
    </div>
    <!-- ... -->
  </div>
  <!-- ... -->
</body>
</html>
```

In this updated index.html file, a datalist element with the ID teams is used for both input fields. The datalist is populated with the team names fetched from the database. The input elements have a list attribute with the value "teams" to link them to the datalist for autocompletion.

Now, when users start typing a team name in the input fields, they will see a list of suggestions based on the teams fetched from the database. When a user selects a suggestion, the input field will be populated with the team ID.

With these changes, the indexHandler fetches the teams using fetchTeamsFromDB() and passes them to the index.html template. The index.html file uses the datalist element for autocompleting team names based on the content of the teams table.

## Note

In my previous conversation about the frontend, ChatGPT had already created the fetchTeamsFromDB() function, which I've just recycled here, and put in `database.go`.
