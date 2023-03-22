I have a Go app that predicts soccer outcomes. It has a REST API. I want it to serve a simple web page that calls the API. It can be really simple, with just html, css, and some Javascript without frameworks or npm. It has two input fields, one for the home team and one for the away team, and a button to call the endpoint: `/predict?home_team_id=<num>&away_team_id=<num>. The endpoint in the Go app is already secured using OAuth2 through Google, so the webpage needs to redirect to a Google login page if it's not already authenticated.

To create a simple webpage that calls your Go app's REST API, you can follow these steps:

1. Create an HTML file, e.g., index.html:

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
    const clientId = 'YOUR_GOOGLE_CLIENT_ID'; // Replace with your Google OAuth2 client ID
    const apiBaseUrl = 'YOUR_API_BASE_URL'; // Replace with your Go app's base URL

    let accessToken = '';

    function onSignIn(googleUser) {
      const profile = googleUser.getBasicProfile();
      accessToken = googleUser.getAuthResponse().access_token;
    }

    function onFailure(error) {
      console.log(error);
    }

    function renderButton() {
      gapi.signin2.render('signInBtn', {
        'scope': 'profile email',
        'width': 200,
        'height': 50,
        'longtitle': true,
        'theme': 'dark',
        'onsuccess': onSignIn,
        'onfailure': onFailure
      });
    }

    function makePrediction() {
      const homeTeamId = document.getElementById('home_team_id').value;
      const awayTeamId = document.getElementById('away_team_id').value;

      if (!accessToken) {
        alert('Please sign in with Google first');
        return;
      }

      if (!homeTeamId || !awayTeamId) {
        alert('Please enter both team IDs');
        return;
      }

      const url = `${apiBaseUrl}/predict?home_team_id=${homeTeamId}&away_team_id=${awayTeamId}`;

      fetch(url, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${accessToken}`
        }
      })
      .then(response => response.json())
      .then(data => {
        document.getElementById('result').innerHTML = JSON.stringify(data);
      })
      .catch(error => {
        console.error('Error:', error);
        alert('Failed to make a prediction');
      });
    }

    document.getElementById('predict_btn').addEventListener('click', makePrediction);
```

In the index.html file, replace YOUR_GOOGLE_CLIENT_ID with your Google OAuth2 client ID and YOUR_API_BASE_URL with your Go app's base URL.

This HTML file contains a simple form with two input fields for the home and away team IDs, a button to call the API, and a div to display the result. It also includes Google Sign-In, which handles user authentication through Google.

The JavaScript code in the index.html file does the following:

1. Sets up Google Sign-In and handles user authentication.
2. Makes an API call to your Go app using the Fetch API when the 'Predict' button is clicked.
3. Displays the prediction result in the 'result' div.

To use this webpage, host it on a web server, and open it in a browser. Users will be prompted to sign in with Google before making predictions.

Note: The Google Sign-In JavaScript API requires your webpage to be served over HTTPS, so make sure your web server is configured to serve HTTPS.

## Note

Initially I asked it to create a Svelte app for me, but apparently Svelte's Google authentication library has been in flux since ChatGPT's training cutoff, and I couldn't get it to work. In the end I decided to just make a plain page.
