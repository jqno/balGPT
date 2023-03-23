# Restore authentication

> OK great! Now I know the app works well, I want to add some authentication back in. Let's use "Authentication: basic" headers. Please generate the backend code needed for that. The accepted credentials will be passed in via environment variables. Then update index.html. Note that I don't want to see a browser popup if possible, please do it by hiding the signedInContent div when not authenticated, and showing a signedOutContent div containing the things needed to log in.

First, let's update the backend to use Basic Authentication. We'll use environment variables to pass in the accepted username and password. Here's how to modify the main.go file:

1. Import the required packages:

```go
import (
	"encoding/base64"
	"os"
	"strings"
)
```

1. Add a new function to check the authentication credentials:

```go
func checkAuth(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization required", http.StatusUnauthorized)
		return false
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || strings.ToLower(authParts[0]) != "basic" {
		http.Error(w, "Invalid authorization header", http.StatusBadRequest)
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(authParts[1])
	if err != nil {
		http.Error(w, "Invalid base64 encoding", http.StatusBadRequest)
		return false
	}

	creds := strings.Split(string(decoded), ":")
	if len(creds) != 2 {
		http.Error(w, "Invalid credentials format", http.StatusBadRequest)
		return false
	}

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	if creds[0] != username || creds[1] != password {
		http.Error(w, "Invalid username or password", http.StatusForbidden)
		return false
	}

	return true
}
```

1. Update the handlePrediction and handleScrape functions to use checkAuth:

```go
func handlePrediction(s *scraper.ScrapeData, p *predictor.Predictor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuth(w, r) {
			return
		}
		// ...
	}
}

func handleScrape(s *scraper.ScrapeData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !checkAuth(w, r) {
			return
		}
		// ...
	}
}
```

Now let's update the frontend to use the Basic Authentication.

1. Add the following HTML to the signedOutContent div inside the index.html file:

```html
<div id="loginForm">
  <label for="username">Username:</label>
  <input type="text" id="username" required>
  <label for="password">Password:</label>
  <input type="password" id="password" required>
  <button id="loginBtn">Log in</button>
  <div id="loginError" style="color: red; display: none;">Invalid username or password. Please try again.</div>
</div>
```

1. Add the following JavaScript code inside the `<script>` tag:

```javascript
function login() {
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;

  const auth = btoa(`${username}:${password}`);
  const url = `${apiBaseUrl}/login`;

  fetch(url, {
    method: 'GET',
    headers: {
      'Authorization': `Basic ${auth}`
    }
  })
  .then(response => {
    if (response.status === 200) {
      sessionStorage.setItem('auth', auth);
      document.getElementById('signedInContent').style.display = 'block';
      document.getElementById('signedOutContent').style.display = 'none';
    } else {
      document.getElementById('loginError').style.display = 'block';
    }
  })
  .catch(error => {
    console.error('Error:', error);
    document.getElementById('loginError').style.display = 'block';
  });
}

document.getElementById('loginBtn').addEventListener('click', login);

// Check if the user is already logged in
if (sessionStorage.getItem('auth')) {
  document.getElementById('signedInContent').style.display = 'block';
  document.getElementById('signedOutContent').style.display = 'none';
}
```

3. Update the `makePrediction` and `scrapeData` functions to include the Authorization header:

```javascript
function makePrediction() {
  // ...

  const auth = sessionStorage.getItem('auth');
  if (!auth) {
    alert('Please log in first');
    return;
  }

  fetch(url, {
    method: 'GET',
    headers: {
      'Authorization': `Basic ${auth}`
    }
  })
  // ...
}

function scrapeData() {
  // ...

  const auth = sessionStorage.getItem('auth');
  if (!auth) {
    alert('Please log in first');
    return;
  }

  fetch(url, {
    method: 'GET',
    headers: {
      'Authorization': `Basic ${auth}`
    }
  })
  // ...
}
```

This will update the frontend to use the Basic Authentication. The user will be asked to provide their username and password, and the credentials will be stored in the session storage. The makePrediction and scrapeData functions will use these credentials to authenticate the requests to the backend. If the user logs out or closes the browser, the session storage will be cleared, and they will need to log in again.

## Follow-up prompt

> Awesome! But can we change the backend code a little bit? I'd like the checkAuth function to wrap around the calls to handlePrediction and handleScrape in the http.HandleFunc calls.

Sure! To achieve that, you can create a higher-order function wrapWithAuth that takes a function and returns a function with authentication check. Here's the updated backend code:

```go
func wrapWithAuth(h http.HandlerFunc, validUsername, validPassword string) http.HandlerFunc {
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

func main() {
	// ...

	// Get credentials from environment variables
	username := os.Getenv("AUTH_USERNAME")
	password := os.Getenv("AUTH_PASSWORD")

	http.HandleFunc("/predict", wrapWithAuth(handlePrediction(s, p), username, password))
	http.HandleFunc("/scrape", wrapWithAuth(handleScrape(s), username, password))

	// ...
}
```

In this code, the wrapWithAuth function takes the original function (e.g., handlePrediction or handleScrape) and the valid username and password. It returns a new function that first checks for the correct authentication credentials and then calls the original function. The http.HandleFunc calls are updated to use the wrapped functions.

## Follow-up prompt

> The divs aren't displayed correctly. Can you give the signedInContent and signedOutContent explicit initial values, and update the sessionStorage check accordingly? Also, when I'm signed in, I want the values to update accordingly as well.

Sure! I've updated the frontend code to set explicit initial values for the signedInContent and signedOutContent divs, and to update them based on the sessionStorage value. Also, I've added a function to update the divs when the user logs in or out.

```html
<!DOCTYPE html>
<html lang="en">
<!-- ... -->
<body>
  <div id="signedInContent" style="display:none;">
    <!-- ... -->
  </div>
  <div id="signedOutContent" style="display:block;">
    <h1>Please sign in</h1>
    <label for="username">Username:</label>
    <input type="text" id="username" />
    <label for="password">Password:</label>
    <input type="password" id="password" />
    <button id="signInBtn">Sign in</button>
    <div id="signInError" style="color: red; display: none;">Failed to sign in. Please try again.</div>
  </div>
  <script>
    // ...

    function updateContentVisibility() {
      const signedIn = sessionStorage.getItem('signedIn') === 'true';
      document.getElementById('signedInContent').style.display = signedIn ? 'block' : 'none';
      document.getElementById('signedOutContent').style.display = signedIn ? 'none' : 'block';
    }

    function signIn() {
      const username = document.getElementById('username').value;
      const password = document.getElementById('password').value;

      if (!username || !password) {
        document.getElementById('signInError').style.display = 'block';
        return;
      }

      sessionStorage.setItem('signedIn', 'true');
      sessionStorage.setItem('authHeader', 'Basic ' + btoa(username + ':' + password));
      updateContentVisibility();
    }

    function signOut() {
      sessionStorage.removeItem('signedIn');
      sessionStorage.removeItem('authHeader');
      updateContentVisibility();
    }

    document.getElementById('signInBtn').addEventListener('click', signIn);
    updateContentVisibility();

    // ...
  </script>
</body>
</html>
```

Now the signedOutContent div is initially set to display:block; and the signedInContent div is set to display:none;. The updateContentVisibility() function sets the display style of these divs based on the sessionStorage value. This function is called when the page loads and when the user logs in or out.

Please note that this method of storing credentials in sessionStorage is insecure and not recommended for production applications. It's only provided here for demonstration purposes. Consider using a more secure authentication solution like OAuth2 for production applications.

## Follow-up prompt

> Noted!
>
> Do I still need the original login() function with this new code?

Apologies for the confusion. You're right, you don't need the original login() function anymore. I should have removed it from the previous example. You can use the updated signIn() function in place of the original login() function. Here's the updated code without the original login() function:

```html
<!DOCTYPE html>
<html lang="en">
<!-- ... -->
<body>
  <div id="signedInContent" style="display:none;">
    <!-- ... -->
  </div>
  <div id="signedOutContent" style="display:block;">
    <h1>Please sign in</h1>
    <label for="username">Username:</label>
    <input type="text" id="username" />
    <label for="password">Password:</label>
    <input type="password" id="password" />
    <button id="signInBtn">Sign in</button>
    <div id="signInError" style="color: red; display: none;">Failed to sign in. Please try again.</div>
  </div>
  <script>
    // ...

    function updateContentVisibility() {
      const signedIn = sessionStorage.getItem('signedIn') === 'true';
      document.getElementById('signedInContent').style.display = signedIn ? 'block' : 'none';
      document.getElementById('signedOutContent').style.display = signedIn ? 'none' : 'block';
    }

    function signIn() {
      const username = document.getElementById('username').value;
      const password = document.getElementById('password').value;

      if (!username || !password) {
        document.getElementById('signInError').style.display = 'block';
        return;
      }

      sessionStorage.setItem('signedIn', 'true');
      sessionStorage.setItem('authHeader', 'Basic ' + btoa(username + ':' + password));
      updateContentVisibility();
    }

    function signOut() {
      sessionStorage.removeItem('signedIn');
      sessionStorage.removeItem('authHeader');
      updateContentVisibility();
    }

    document.getElementById('signInBtn').addEventListener('click', signIn);
    updateContentVisibility();

    // ...
  </script>
</body>
</html>
```

With this code, you can remove the original login() function, and the new signIn() function takes care of signing in and updating the visibility of the content divs accordingly.

## Note

As always, I added the new environment variables to the Config struct myself.

I was surprised to see the second version of the backend auth code was a lot smaller and simpler than the initial iteration.

Regarding the second version of the frontend code, it actually used a completely different way to store the authentication state in the sessionStorage, which wasn't reflected in the methods that actually needed to use them. I manually updated accordingly.

Also, why did ChatGPT output the same code again in the final prompt? Seemed a bit superfluous 😄

In the end I also made a small tweak to the css to make the input fields look more consistent.
