# Environment variable for scrape url

> I forget, did we define an environment variable yet for the url that we scrape?

No, we haven't defined an environment variable for the URL to be scraped yet. To do this, you can use the os package to read the environment variable. First, let's add the import statement for the os package in your main.go file:

```go
import (
	// ...
	"os"
)
```

Now, you can read the environment variable and pass it to the Scraper struct. Let's call the environment variable SCRAPER_URL. Modify the main function in your main.go file to read the environment variable and set a default value if it's not set:

```go
func main() {
	// ...

	scraperURL := os.Getenv("SCRAPER_URL")
	if scraperURL == "" {
		scraperURL = "https://www.fcupdate.nl/uitslagen/s1598/nederland-keuken-kampioen-divisie-2022-2023/"
	}

	scraper := scraper.NewScraper(db, scraperURL)
	predictor := predictor.NewPredictor(db)

	// ...
}
```

Now, when you run your application, it will use the SCRAPER_URL environment variable if it's set. Otherwise, it will use the default URL provided. To set the environment variable, you can run your application like this:

```sh
SCRAPER_URL="https://www.example.com/your-url" go run main.go
```

Or, you can set the environment variable in your system or your deployment environment, depending on your needs.

## Note

The basic logic generated was fine, but the location was just plain wrong. Also, there's no reason to have a default url. I manually rewrote it to add the url to `config.go`, where the other environment variables were also read. Then I filled in the blanks from there.
