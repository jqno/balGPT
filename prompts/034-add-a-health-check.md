# Add a health check

> I have a Go app with a Postgres database deployed to Fly.io. Can we give it a health check, that also verifies the connection with the database?
>
> Here's part of the app's Run() function, let's add it there:
>
>	http.HandleFunc("/", indexHandler(a.DB, a.Config.AppBaseDir, a.Config.ApiBaseURL, a.Config.AuthUsername, a.Config.AuthPassword))
>	http.HandleFunc("/login", checkAuth(loginHandler(), a.Config.AuthUsername, a.Config.AuthPassword))
>	http.HandleFunc("/predict", checkAuth(handlePrediction(a.Scraper, a.Predictor), a.Config.AuthUsername, a.Config.AuthPassword))
>	http.HandleFunc("/scrape", checkAuth(handleScrape(a.Scraper), a.Config.AuthUsername, a.Config.AuthPassword))
>	http.HandleFunc("/team_id", checkAuth(handleTeamID(a.DB), a.Config.AuthUsername, a.Config.AuthPassword))


To add a health check that verifies the connection with the database, you can create a new handler function that checks the connection to the Postgres database and returns an appropriate HTTP status code based on the connection status. Let's create a new handler function called `healthCheckHandler` and register it with the http package:

1. First, create the healthCheckHandler function:

```go
func healthCheckHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check the database connection
		if err := db.Ping(); err != nil {
			// If there is an error, return a 500 Internal Server Error status code
			http.Error(w, "Database connection failed", http.StatusInternalServerError)
			return
		}

		// If everything is fine, return a 200 OK status code
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
```

1. Then, register the healthCheckHandler in your Run() function:

```go
	http.HandleFunc("/", indexHandler(a.DB, a.Config.AppBaseDir, a.Config.ApiBaseURL, a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/login", checkAuth(loginHandler(), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/predict", checkAuth(handlePrediction(a.Scraper, a.Predictor), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/scrape", checkAuth(handleScrape(a.Scraper), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/team_id", checkAuth(handleTeamID(a.DB), a.Config.AuthUsername, a.Config.AuthPassword))
	http.HandleFunc("/health", healthCheckHandler(a.DB)) // Add this line
```

Now, you can use the `/health` endpoint to perform a health check that also verifies the connection with the database. If the connection to the database is working, the endpoint will return a 200 OK status code with the message "OK". If the connection fails, it will return a 500 Internal Server Error status code with the message "Database connection failed".

## Follow-up prompt

> Can we tell Fly.io about our health endpoint?

Yes, you can configure Fly.io to use your health check endpoint for monitoring the health of your application. Fly.io supports custom health checks via the `fly.toml` configuration file in your project.

To set up a custom health check for your Go app running on Fly.io, follow these steps:

1. Open the `fly.toml` file in your project directory. If you don't have one, create a new file named `fly.toml`.
2. Add or update the `[services]` section in the `fly.toml` file to include a custom health check configuration pointing to your `/health` endpoint:

```toml
app = "your-app-name"

[[services]]
  internal_port = 8080
  protocol = "tcp"

  [services.concurrency]
    hard_limit = 25
    soft_limit = 20

  [[services.ports]]
    handlers = ["http"]
    port = "80"

  [[services.ports]]
    handlers = ["tls", "http"]
    port = "443"

  [[services.tcp_checks]]
    interval = 10000 # in milliseconds, 10000 ms = 10 seconds
    timeout = 2000 # in milliseconds, 2000 ms = 2 seconds
    type = "http"
    method = "GET"
    path = "/health" # The path to your health check endpoint
    protocol = "http"
    port = "80"
    tls_skip_verify = false
    status = 200 # Expected HTTP status code for a healthy instance
```

Replace `your-app-name` with the actual name of your Fly.io application. Adjust the `internal_port`, `protocol`, `ports`, and other settings according to your app's configuration.

The example above sets up an HTTP health check on port 80 at the `/health` path. It checks the health of your application every 10 seconds (as defined by the `interval` field), with a 2-second timeout (as defined by the `timeout` field). A 200 OK status code is expected for a healthy instance.

1. Save the changes to the fly.toml file and redeploy your application:

```sh
$ flyctl deploy
```

After deployment, Fly.io will use the custom health check to monitor your application. If the health check fails, Fly.io can automatically restart or replace unhealthy instances to maintain application reliability.

