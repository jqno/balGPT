# Add tests

> Let's add some unit tests. Especially the predictor and the scraper, but also for the security. And perhaps also an integration test?

To write unit tests in Go, you'll create _test.go files containing your test functions. Let's start by creating tests for the scraper and predictor packages.

1. Create a file named scraper_test.go in the scraper package.
2. Import necessary packages and create a test function for the scraper:

```go
package scraper

import (
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TestScraper(t *testing.T) {
	// Set up the test environment
	dbUrl := os.Getenv("TEST_DATABASE_URL")
	if dbUrl == "" {
		t.Fatal("TEST_DATABASE_URL not set")
	}

	db, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		t.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	// Prepare the database
	// ...

	// Initialize the scraper
	scraperURL := "https://www.fcupdate.nl/uitslagen/s1598/nederland-keuken-kampioen-divisie-2022-2023/"
	s := NewScraper(db, scraperURL)

	// Run the scraper
	err = s.Scrape()
	if err != nil {
		t.Fatalf("Scraper failed: %v", err)
	}

	// Test the expected outcomes
	// ...
}
```

1. Create a file named predictor_test.go in the predictor package.
2. Import necessary packages and create a test function for the predictor:

```go
package predictor

import (
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

func TestPredictor(t *testing.T) {
	// Set up the test environment
	dbUrl := os.Getenv("TEST_DATABASE_URL")
	if dbUrl == "" {
		t.Fatal("TEST_DATABASE_URL not set")
	}

	db, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		t.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	// Prepare the database
	// ...

	// Initialize the predictor
	p := NewPredictor(db)

	// Run the predictor and test the outcomes
	// ...
}
```

For testing the security middleware, you can create a test function in the main_test.go file. You can mock the Google OAuth2 token validation process to test the allowed and disallowed email addresses.

For integration tests, you can write test functions that interact with your running application through HTTP requests, simulating a real user interaction. You can use the net/http/httptest package to send requests and check responses from your application.

Remember to adjust the // ... parts with the actual tests you want to perform, like checking if the scraper fetched the correct data or if the predictor returned a reasonable result.

To run the tests, you can use the go test command followed by the package name or ./... to run tests for all packages:

```sh
go test ./...
```

## Note

These tests were so very incomplete and tied to the environment that I didn't even bother incorporating them. I should have paid more attention while I was generating prompts, and asked for better tests. Maybe another day.
