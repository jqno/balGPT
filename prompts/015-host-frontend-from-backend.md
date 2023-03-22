# Host frontend from backend

> Can I have my Go app host these files, and make it fill in the Google clientId and the apiBaseUrl from an environment variable?

Yes, you can have your Go app serve these static files and fill in the clientId and apiBaseUrl from environment variables. To do this, follow these steps:

1. Create a templates folder in your Go app's directory and place the index.html file inside it.
2. Modify your Go app to serve the index.html file:

```go
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type TemplateData struct {
	ClientID  string
	ApiBaseURL string
}

func main() {
	http.HandleFunc("/", indexHandler)

	// Other routes and middleware...

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	apiBaseURL := os.Getenv("API_BASE_URL")

	if clientID == "" || apiBaseURL == "" {
		http.Error(w, "Missing environment variables", http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		ClientID:  clientID,
		ApiBaseURL: apiBaseURL,
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
```

3. Update your `index.html` file to use the template data:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Soccer Prediction</title>
  <style>
    body { font-family: Arial, sans-serif; }
    input[type="number"] { width: 50px; }
  </style>
</head>
<body>
  <h1>Soccer Prediction</h1>
  <div>
    <label for="home_team_id">Home Team ID:</label>
    <input type="number" id="home_team_id" required>
  </div>
  <div>
    <label for="away_team_id">Away Team ID:</label>
    <input type="number" id="away_team_id" required>
  </div>
  <button id="predict_btn">Predict</button>
  <div id="result"></div>

  <script>
    const clientId = '{{.ClientID}}';
    const apiBaseUrl = '{{.ApiBaseURL}}';

    // Rest of the script...

  </script>
  <script src="https://apis.google.com/js/platform.js?onload=renderButton" async defer></script>
  <meta name="google-signin-client_id" content="{{.ClientID}}">
  <div id="signInBtn"></div>
</body>
</html>
```

Now, your Go app will serve the index.html file, and the clientId and apiBaseUrl values will be filled in using the environment variables GOOGLE_CLIENT_ID and API_BASE_URL.

Remember to set these environment variables when running your Go app:

```sh
export GOOGLE_CLIENT_ID=your_google_client_id
export API_BASE_URL=your_api_base_url
go run main.go
```

With these changes, your Go app will host the HTML file and dynamically fill in the clientId and apiBaseUrl values from the environment variables. Users can now access the simple web page by visiting your Go app's base URL.

## Note

Since we already store the environment variables in the `Config` struct, I've manually added them there as well, instead of directly reading them in the `indexHandler` function.
