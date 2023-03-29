# I feel like an idiot

> Hi! I feel like an idiot, I committed code without checking it properly and now it doesn't work. I have an app that should use Basic Authentication to log in, but it doesn't actually check the login. Here's the front-end:
>
> _(snipped: current version of index.html)_
>
> Here's the relevant part of the backend:
>
> _(snipped: relevant part of app.go)_
>
> I think it would be good to do the following things:
>
> - Separate the frontend into a login.html and a main.html, and when the user accesses the index, the backend decides which one to show based on whether the correct credentials are provided or not.
> - The backend should return an Unauthorized HTTP return code when accessing an endpoint without proper credentials, but in that case the browser should not show a login popup dialog, instead it should redirect to the login.html page.
> - Maybe simplify a little by using just a single item in the sessionStorage instead of two.

To achieve your goals, I'll suggest the following changes to your code:

1. Create two separate HTML files: `login.html` and `main.html`.
2. Update the `indexHandler` to serve the appropriate HTML file based on the presence of correct credentials.
3. Update the front-end code to redirect to the login page when receiving an "Unauthorized" status code.
4. Simplify the `sessionStorage` usage to store a single item for authentication.

Here are the changes to your code:

1. Create `login.html`:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Login</title>
  <link rel="stylesheet" href="/static/css/styles.css">
</head>
<body>
  <div id="loginForm">
    <label for="username">Username:</label>
    <input type="text" id="username" required>
    <label for="password">Password:</label>
    <input type="password" id="password" required>
    <button id="loginBtn">Log in</button>
    <div id="loginError" style="color: red; display: none;">Invalid username or password. Please try again.</div>
  </div>
  <script>
    // ... (Update the signIn function and related code)
  </script>
</body>
</html>
```

1. Create `main.html`:

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
  <div id="signedInContent">
    <!-- ... (Move the content of the signedInContent div from your original HTML) -->
  </div>
  <script>
    // ... (Update the rest of your JavaScript code)
  </script>
</body>
</html>
```

1. Update the `indexHandler` function:

```go
func indexHandler(db *database.DB, appBaseDir string, apiBaseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		var templateFile string
		if !ok || username != a.Config.AuthUsername || password != a.Config.AuthPassword {
			templateFile = filepath.Join(appBaseDir, "templates/login.html")
		} else {
			teams, err := db.FetchTeamsFromDB()
			if err != nil {
				http.Error(w, fmt.Sprintf("Error fetching teams: %v", err), http.StatusInternalServerError)
				return
			}

			data := TemplateData{
				ApiBaseURL: apiBaseURL,
				Teams:      teams,
			}
			templateFile = filepath.Join(appBaseDir, "templates/main.html")
			tmpl, err := template.ParseFiles(templateFile)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, data)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
				return
			}
			return
		}

		tmpl, err := template.ParseFiles(templateFile)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}
```

4. Update the front-end JavaScript code in `login.html`:

```javascript
function signIn() {
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;

  if (!username || !password) {
    document.getElementById('loginError').style.display = 'block';
    return;
  }

  sessionStorage.setItem('authHeader', 'Basic ' + btoa(username + ':' + password));
  location.reload();
}

document.getElementById('loginBtn').addEventListener('click', signIn);
```

1. Update the front-end JavaScript code in main.html:

```javascript
function handleUnauthorized() {
  sessionStorage.removeItem('authHeader');
  location.reload();
}

function makeRequest(url, options) {
  const auth = sessionStorage.getItem('authHeader');
  if (!auth) {
    handleUnauthorized();
    return Promise.reject();
  }

  options = options || {};
  options.headers = options.headers || {};
  options.headers['Authorization'] = auth;

  return fetch(url, options)
    .then(response => {
      if (response.status === 401) {
        handleUnauthorized();
        return Promise.reject();
      }
      return response;
    });
}

// Update the scrapeData() and makePrediction() functions to use makeRequest() instead of fetch()
function scrapeData() {
  // ...
  makeRequest(url, {
    method: 'GET',
  })
  // ...
}

function makePrediction() {
  // ...
  makeRequest(url, {
    method: 'GET',
  })
  // ...
}

// ... (rest of the code)

```

With these changes, your app should now serve the `login.html` file when the user is not authenticated, and serve the `main.html` file when the user is authenticated. Additionally, it should redirect to the `login.html` page when an "Unauthorized" status code is received from the API, and it simplifies the usage of `sessionStorage` to store a single item for authentication.

## Follow-up prompt

> Can you give me the complete versions of main.html and login.html?

Here are the complete versions of `main.html` and `login.html`.

`main.html`:

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
  <div id="signedInContent">
    <!-- Content from the original signedInContent div -->
    <h1>Soccer Prediction</h1>
    {{if .Teams}}
      <div>
        <label for="home_team_id">Home Team ID:</label>
        <input list="home_teams" id="home_team_name" required>
        <input type="hidden" id="home_team_id">
        <datalist id="home_teams">
          {{range .Teams}}
          <option data-value="{{.ID}}">{{.Name}}</option>
          {{end}}
        </datalist>
      </div>
      <div>
        <label for="away_team_id">Away Team ID:</label>
        <input list="away_teams" id="away_team_name" required>
        <input type="hidden" id="away_team_id">
        <datalist id="away_teams">
          {{range .Teams}}
          <option data-value="{{.ID}}">{{.Name}}</option>
          {{end}}
        </datalist>
      </div>
      <button id="predict_btn">Predict</button>
      <div id="result"></div>
    {{else}}
      <div>
        <p>No teams found. Please scrape data to get team information.</p>
        <button id="scrape_btn">Scrape Data</button>
      </div>
    {{end}}
  </div>
  <script>
    const apiBaseUrl = '{{.ApiBaseURL}}';

    function handleUnauthorized() {
      sessionStorage.removeItem('authHeader');
      location.reload();
    }

    function makeRequest(url, options) {
      const auth = sessionStorage.getItem('authHeader');
      if (!auth) {
        handleUnauthorized();
        return Promise.reject();
      }

      options = options || {};
      options.headers = options.headers || {};
      options.headers['Authorization'] = auth;

      return fetch(url, options)
        .then(response => {
          if (response.status === 401) {
            handleUnauthorized();
            return Promise.reject();
          }
          return response;
        });
    }

    // Update the scrapeData() and makePrediction() functions to use makeRequest() instead of fetch()
    function scrapeData() {
      const url = `${apiBaseUrl}/scrape`;
      makeRequest(url, {
        method: 'GET',
      })
      .then(response => {
        if (response.ok) {
          alert('Data scraped successfully.');
        } else {
          alert('Failed to scrape data.');
        }
      })
      .catch(error => {
        console.error('Error:', error);
        alert('Failed to scrape data.');
      });
    }

    function makePrediction() {
      const homeTeamId = document.getElementById('home_team_id').value;
      const awayTeamId = document.getElementById('away_team_id').value;

      if (!homeTeamId || !awayTeamId) {
        alert('Please enter both team IDs');
        return;
      }

      const url = `${apiBaseUrl}/predict?home_team_id=${homeTeamId}&away_team_id=${awayTeamId}`;
      makeRequest(url, {
        method: 'GET',
      })
      .then(response => response.json())
      .then(data => {
        document.getElementById('result').innerHTML = formatPredictionResult(data);
      })
      .catch(error => {
          console.error('Error:', error);
          alert('Failed to make a prediction');
      });
    }

    function formatPredictionResult(data) {
      const homeGoals = data.HomeGoals;
      const awayGoals = data.AwayGoals;

      return `
        <h2>Prediction Result</h2>
        <table>
          <tr>
            <th>Team</th>
            <th>Goals</th>
          </tr>
          <tr>
            <td>Home Team</td>
            <td>${homeGoals}</td>
          </tr>
          <tr>
            <td>Away Team</td>
            <td>${awayGoals}</td>
          </tr>
        </table>
      `;
    }

    document.getElementById('home_team_name')?.addEventListener('input', function (e) {
      const selectedOption = Array.from(e.target.list.options).find(option => option.textContent === e.target.value);
      document.getElementById('home_team_id').value = selectedOption ? selectedOption.getAttribute('data-value') : '';
    });

    document.getElementById('away_team_name')?.addEventListener('input', function (e) {
      const selectedOption = Array.from(e.target.list.options).find(option => option.textContent === e.target.value);
      document.getElementById('away_team_id').value = selectedOption ? selectedOption.getAttribute('data-value') : '';
    });

    document.getElementById('scrape_btn')?.addEventListener('click', scrapeData);
    document.getElementById('predict_btn')?.addEventListener('click', makePrediction);
  </script>
</body>
</html>
```

`login.html`:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Soccer Prediction Login</title>
  <link rel="stylesheet" href="/static/css/styles.css">
</head>
<body>
  <div id="loginForm">
    <label for="username">Username:</label>
    <input type="text" id="username" required>
    <label for="password">Password:</label>
    <input type="password" id="password" required>
    <button id="loginBtn">Log in</button>
    <div id="loginError" style="color: red; display: none;">Invalid username or password. Please try again.</div>
  </div>

  <script>
    function signIn() {
      const username = document.getElementById('username').value;
      const password = document.getElementById('password').value;

      if (!username || !password) {
        document.getElementById('loginError').style.display = 'block';
        return;
      }

      sessionStorage.setItem('authHeader', 'Basic ' + btoa(username + ':' + password));
      location.reload();
    }

    document.getElementById('loginBtn').addEventListener('click', signIn);
  </script>
</body>
</html>
```

These are the complete versions of `main.html` and `login.html` with the changes we discussed. The `main.html` file contains the main content of the application, and the `login.html` file contains the login form. The JavaScript code in each file handles the appropriate functionality for authentication and making API requests.

## Follow-up prompt

> Regarding the backend, I need checkAuth for other endpoints as well. Can you separate that out from indexHandler and give me both functions in full?

Certainly! I'll separate the `checkAuth` function and modify the `indexHandler` function to handle the authentication in a cleaner way.

Here's the updated `checkAuth` function:

```go
func checkAuth(h http.HandlerFunc, validUsername, validPassword string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if !ok || username != validUsername || password != validPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		h(w, r)
	}
}
```

And here's the updated `indexHandler` function:

```go
func indexHandler(db *database.DB, appBaseDir string, apiBaseURL string, validUsername, validPassword string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if !ok || username != validUsername || password != validPassword {
			http.ServeFile(w, r, filepath.Join(appBaseDir, "templates/login.html"))
			return
		}

		teams, err := db.FetchTeamsFromDB()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching teams: %v", err), http.StatusInternalServerError)
			return
		}

		data := TemplateData{
			ApiBaseURL: apiBaseURL,
			Teams:      teams,
		}

		templateFile := filepath.Join(appBaseDir, "templates/main.html")
		tmpl, err := template.ParseFiles(templateFile)
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
}
```

Now, the `checkAuth` function is separate from the `indexHandler`. The `indexHandler` will serve the `login.html` file if the provided credentials are not valid, and it will serve the `main.html` file if the credentials are valid. You can use the `checkAuth` function for other endpoints as needed.

## Follow-up prompt

> Can you update login.html to show an error when the login attempt failed? It should probably validate that in in the signIn function.

To show an error when the login attempt fails, we can modify the `signIn` function in the JavaScript code of the `login.html` file to make an HTTP request to the server and check the response. If the response is 401 Unauthorized, we will display an error message.

Here's the updated `login.html` file:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Soccer Prediction - Login</title>
  <link rel="stylesheet" href="/static/css/styles.css">
</head>
<body>
  <div id="loginForm">
    <label for="username">Username:</label>
    <input type="text" id="username" required>
    <label for="password">Password:</label>
    <input type="password" id="password" required>
    <button id="loginBtn">Log in</button>
    <div id="loginError" style="color: red; display: none;">Invalid username or password. Please try again.</div>
  </div>

  <script>
    const apiBaseUrl = '{{.ApiBaseURL}}';
    async function signIn() {
      const username = document.getElementById('username').value;
      const password = document.getElementById('password').value;

      if (!username || !password) {
        document.getElementById('loginError').style.display = 'block';
        return;
      }

      const authHeader = 'Basic ' + btoa(username + ':' + password);
      const response = await fetch(apiBaseUrl, {
        method: 'GET',
        headers: {
          'Authorization': authHeader
        }
      });

      if (response.status === 401) {
        document.getElementById('loginError').style.display = 'block';
      } else if (response.ok) {
        sessionStorage.setItem('authHeader', authHeader);
        location.href = apiBaseUrl;
      } else {
        document.getElementById('loginError').style.display = 'block';
      }
    }

    document.getElementById('loginBtn').addEventListener('click', signIn);
  </script>
</body>
</html>
```

Now, the `signIn` function makes a request to the server, and if the response status is 401, it will display the error message. If the response is successful, it will store the `authHeader` in the sessionStorage and redirect the user to the main page.

## Follow-up prompt

> Let's create a login endpoint in the backend that it can call.

To create a login endpoint in the backend, we'll add a new function called `loginHandler`. This function will check the user credentials and return an appropriate response. Then, we'll add this new endpoint to the `Run()` function.

Here's the `loginHandler` function:

```go
func loginHandler(validUsername, validPassword string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()

		if !ok || username != validUsername || password != validPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized.", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Authenticated.")
	}
}
```

And here's the updated `Run()` function that includes the new `/login` endpoint:

```go
func (a *App) Run() {
	http.HandleFunc("/", indexHandler(a.DB, a.Config.AppBaseDir, a.Config.ApiBaseURL, a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/login", loginHandler(a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/predict", checkAuth(handlePrediction(a.Scraper, a.Predictor), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/scrape", checkAuth(handleScrape(a.Scraper), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/team_id", checkAuth(handleTeamID(a.DB), a.Config.AuthUsername, a.Config.AuthPassword))

	staticDir := filepath.Join(a.Config.AppBaseDir, "static")
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("Listening on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}
```

Now, you have a `/login` endpoint in the backend that checks the provided credentials and returns either an "Unauthorized" or "Authenticated" response. Make sure to update the `signIn` function in the `login.html` file to call the `/login` endpoint:

```javascript
const response = await fetch(apiBaseUrl + '/login', {
  method: 'GET',
  headers: {
    'Authorization': authHeader
  }
});
```

## Follow-up prompt

> I'm sure that if we wrap the loginHandler in checkAuth, it doesn't need to do much anymore

You're right. If you wrap the `loginHandler` function in `checkAuth`, it will already handle the authentication check. In that case, the `loginHandler` can simply return an "Authenticated" message when the credentials are valid. Here's the updated `loginHandler` function:

```go
func loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Authenticated.")
	}
}
```

And here's the updated `Run()` function with the `loginHandler` wrapped in `checkAuth`:

```go
func (a *App) Run() {
	http.HandleFunc("/", indexHandler(a.DB, a.Config.AppBaseDir, a.Config.ApiBaseURL, a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/login", checkAuth(loginHandler(), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/predict", checkAuth(handlePrediction(a.Scraper, a.Predictor), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/scrape", checkAuth(handleScrape(a.Scraper), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/team_id", checkAuth(handleTeamID(a.DB), a.Config.AuthUsername, a.Config.AuthPassword))

	staticDir := filepath.Join(a.Config.AppBaseDir, "static")
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("Listening on port %s...\n", port)
	http.ListenAndServe(":"+port, nil)
}
```

Now, the `loginHandler` function is much simpler, and the authentication logic is managed by the `checkAuth` middleware.

## Follow-up prompt

> Also, login.html is a template too. We should update indexHandler.

You're right. Since `login.html` is also a template, we should update the `indexHandler` to handle both templates, depending on whether the user is authenticated or not. Here's the updated `indexHandler` function:

```go
func indexHandler(db *database.DB, appBaseDir string, apiBaseURL string, validUsername, validPassword string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, isAuthenticated := r.BasicAuth()
		if !isAuthenticated || username != validUsername || password != validPassword {
			templateFile := filepath.Join(appBaseDir, "templates/login.html")
			tmpl, err := template.ParseFiles(templateFile)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
				return
			}
		} else {
			teams, err := db.FetchTeamsFromDB()
			if err != nil {
				http.Error(w, fmt.Sprintf("Error fetching teams: %v", err), http.StatusInternalServerError)
				return
			}

			data := TemplateData{
				ApiBaseURL: apiBaseURL,
				Teams:      teams,
			}

			templateFile := filepath.Join(appBaseDir, "templates/main.html")
			tmpl, err := template.ParseFiles(templateFile)
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
	}
}
```

Now, the `indexHandler` will check if the user is authenticated. If they are, it serves the `main.html` template; otherwise, it serves the `login.html` template.

## Follow-up prompt

> When the authentication passes, it sets the location.href to go to main.html, but it doesn't authenticate when it does so, so it returns to the login page instead.

You're correct that setting `location.href` to `main.html` doesn't include the authentication information. To resolve this issue, you can pass the authentication information as a query parameter in the URL when redirecting to the `main.html` page. Then, in the `indexHandler`, you can check for the presence of these query parameters and validate them.

Here's the updated `signIn` function in the `login.html` template:

```js
function signIn() {
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;

  if (!username || !password) {
    document.getElementById('loginError').style.display = 'block';
    return;
  }

  fetch('/login', {
    method: 'POST',
    headers: {
      'Authorization': 'Basic ' + btoa(username + ':' + password),
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ username, password }),
  })
    .then((response) => {
      if (response.ok) {
        // Redirect to main.html with the authentication information as query parameters
        window.location.href = `main.html?username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`;
      } else {
        document.getElementById('loginError').style.display = 'block';
      }
    })
    .catch((error) => {
      console.error('Error:', error);
      document.getElementById('loginError').style.display = 'block';
    });
}
```

And the updated `indexHandler` function:

```go
func indexHandler(db *database.DB, appBaseDir string, apiBaseURL string, validUsername, validPassword string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		password := r.URL.Query().Get("password")
		isAuthenticated := username == validUsername && password == validPassword

		if !isAuthenticated {
			templateFile := filepath.Join(appBaseDir, "templates/login.html")
			tmpl, err := template.ParseFiles(templateFile)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
				return
			}
		} else {
			teams, err := db.FetchTeamsFromDB()
			if err != nil {
				http.Error(w, fmt.Sprintf("Error fetching teams: %v", err), http.StatusInternalServerError)
				return
			}

			data := TemplateData{
				ApiBaseURL: apiBaseURL,
				Teams:      teams,
			}

			templateFile := filepath.Join(appBaseDir, "templates/main.html")
			tmpl, err := template.ParseFiles(templateFile)
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
	}
}
```

Now, the `indexHandler` checks the query parameters for the authentication information and validates them. If the authentication passes, it serves the `main.html` template. If not, it serves the `login.html` template.

## Note

It took a surprising amount of back-and-forth to get this stuff to work, and I'm still not sure I got it right. Actually, I know I didn't get it right, because I'm using Basic Authentication and storing stuff in the sessionStorage, which isn't right. But for now it will do.
