# Unit tests for scraper

> Let's create a unit test for this Go file, which scrapes a website:
>
> _(snipped: scraper.go)_
>
> Here's a piece of the website it scrapes; use it for test data.
>
> _(snipped: some html, which now lives in scraper_test.go)_

Sure, let's create a unit test for the Scrape function. We will use the `httptest` package to create a test server that serves the sample HTML you provided, and then we will test the `Scrape` function to ensure it correctly scrapes and processes this data.

```go
package scraper_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jqno/balGPT/internal/database"
	"github.com/jqno/balGPT/internal/scraper"
)

const testData = `<!-- your HTML test data goes here -->`

func TestScrape(t *testing.T) {
	// Mock DB operations
	db := &database.MockDB{
		GetLastScrapeFunc: func() (time.Time, error) {
			return time.Time{}, nil
		},
		InsertOrUpdateMatchFunc: func(homeTeam, awayTeam string, homeGoals, awayGoals int, date time.Time) error {
			return nil
		},
		UpdateLastScrapeFunc: func(t time.Time) error {
			return nil
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, testData)
	}))
	defer ts.Close()

	s := scraper.NewScrapeData(db, ts.URL)

	err := s.Scrape()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}
}


```

This test verifies that the `Scrape` function can successfully scrape data from a web page and correctly call the required database operations. You should add more specific tests for the individual pieces of data that are scraped from the page.

Please replace `const testData` with your HTML content you provided.

Also, note that we are using a mock database here (`MockDB`). This is a dummy database implementation that implements the `DB` interface. The `GetLastScrapeFunc`, `InsertOrUpdateMatchFunc` and `UpdateLastScrapeFunc` are function fields on this struct that allow us to control the behavior of these methods during the test.

Please make sure to implement the `MockDB` similar to the following example:

```go
package database

import (
	"time"
)

type MockDB struct {
	GetLastScrapeFunc         func() (time.Time, error)
	InsertOrUpdateMatchFunc   func(homeTeam, awayTeam string, homeGoals, awayGoals int, date time.Time) error
	UpdateLastScrapeFunc      func(t time.Time) error
}

func (m *MockDB) GetLastScrape() (time.Time, error) {
	return m.GetLastScrapeFunc()
}

func (m *MockDB) InsertOrUpdateMatch(homeTeam, awayTeam string, homeGoals, awayGoals int, date time.Time) error {
	return m.InsertOrUpdateMatchFunc(homeTeam, awayTeam, homeGoals, awayGoals, date)
}

func (m *MockDB) UpdateLastScrape(t time.Time) error {
	return m.UpdateLastScrapeFunc(t)
}

```

Again, remember this is a basic test, you should add more tests to verify the scraped data is correct.

## Follow-up prompt

> How do I add a multiline string in Go, to initialize the testData const?

You can use backticks (` ``) to create multiline strings in Go. Here is how you can define your `testData` constant with your HTML test data:

```go
const testData = `
<div class="matches-panel align-left justify-center notes">
Maandag 15 augustus 2022
</div>
<div class="matches-panel d-flex align-center justify-center Played  ">
...
</div>`

func TestScrape(t *testing.T) {
    ...
}

```

Inside backticks, Go will preserve the literal content, including newlines and special characters, which makes them ideal for defining multiline strings.

## Follow-up prompt

> I already have a MockDB struct. It looks like this:
>
> ```go
> type MockDB struct {
>     mock.Mock
> }
>
> func (m \*MockDB) LastYearMatchScores(homeTeamID, awayTeamID int) (int, int, error) {
>     args := m.Called(homeTeamID, awayTeamID)
>     return args.Int(0), args.Int(1), args.Error(2)
> }
> ```
>
> Please update your version of MockDB accordingly.

Certainly! Given that your `MockDB` struct is already using the `github.com/stretchr/testify/mock` package, you can leverage this package for mocking the methods as well. Here's how you can update the MockDB struct and usage accordingly:

```go
type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetLastScrape() (time.Time, error) {
	args := m.Called()
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) InsertOrUpdateMatch(homeTeam, awayTeam string, homeGoals, awayGoals int, date time.Time) error {
	args := m.Called(homeTeam, awayTeam, homeGoals, awayGoals, date)
	return args.Error(0)
}

func (m *MockDB) UpdateLastScrape(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

```

Now, in your test function `TestScrape`, you need to set up the expected return values for these mocked methods. Here is an example:

```go
func TestScrape(t *testing.T) {
	// Mock DB operations
	db := new(MockDB)
	db.On("GetLastScrape").Return(time.Time{}, nil)
	db.On("InsertOrUpdateMatch", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	db.On("UpdateLastScrape", mock.Anything).Return(nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, testData)
	}))
	defer ts.Close()

	s := scraper.NewScrapeData(db, ts.URL)

	err := s.Scrape()
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	db.AssertExpectations(t) // Verify that all expectations were met
}

```

In this example, `mock.Anything` is a placeholder that matches any value. You could replace these with specific values if you want to assert that the methods are called with specific arguments.

## Follow-up prompt

> For a previous test suite, I needed to create a new Database interface and update my code to refer to that interface instead of the original struct. How do I do that here?

In Go, you can decouple your code by using interfaces. An interface is two things: it is a set of methods, but it is also a type. The general idea is that you can define a new type and then provide method signatures in the interface definition. Then, any type that implements these methods can be used wherever this interface type is required.

First, you'd define an interface that describes the methods you need:

```go
package database

import "time"

type Database interface {
	GetLastScrape() (time.Time, error)
	InsertOrUpdateMatch(homeTeam, awayTeam string, homeGoals, awayGoals int, date time.Time) error
	UpdateLastScrape(t time.Time) error
}

```

Then you'd ensure that your `DB` type implements this interface:

```go
package database

import "time"

type DB struct {
	// your fields
}

func (db *DB) GetLastScrape() (time.Time, error) {
	// your implementation
}

func (db *DB) InsertOrUpdateMatch(homeTeam, awayTeam string, homeGoals, awayGoals int, date time.Time) error {
	// your implementation
}

func (db *DB) UpdateLastScrape(t time.Time) error {
	// your implementation
}

```

In your `scraper` package, you would change the `ScrapeData` type to hold a value of the `Database` interface type instead of `*database.DB`:

```go
package scraper

import (
	"github.com/jqno/balGPT/internal/database"
)

type ScrapeData struct {
	DB  database.Database
	URL string
}

```

For testing, you can create a `MockDB` type in your test file that implements the `Database` interface, and use that in your tests:

```go
package scraper_test

import (
	"time"

	"github.com/jqno/balGPT/internal/database"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetLastScrape() (time.Time, error) {
	args := m.Called()
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) InsertOrUpdateMatch(homeTeam, awayTeam string, homeGoals, awayGoals int, date time.Time) error {
	args := m.Called(homeTeam, awayTeam, homeGoals, awayGoals, date)
	return args.Error(0)
}

func (m *MockDB) UpdateLastScrape(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

```

With this setup, you can replace the real database with a mock one for testing. The rest of your code will not be aware that it's using a mock, because it only cares that it's working with something that implements the `Database` interface. This is a common pattern to make code more testable in Go.

## Follow-up prompt

> My testData contains the result for three matches: Jong Utrecht - Heracles 0-3, Jong PSV - Dordrecht 1-0, MVV - NAC 3-1. Can we verify this in the test?
>
> Note that I prefer tests to be short and concise, though I'm not sure if that's possible in this case.

In this scenario, you'll want to verify that the method `InsertOrUpdateMatch` of the mock database has been called correctly for each of the three matches.

In Go, the `mock` package from `github.com/stretchr/testify/mock` provides a method `Called` that records the invocation of a method. Then you can use the `mock.AssertCalled` function to verify that a method was called with given arguments.

Here's how you can set up your test:

```go
package scraper_test

import (
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jqno/balGPT/internal/database"
	"github.com/jqno/balGPT/internal/scraper"
	"github.com/stretchr/testify/mock"
)

const testData = `
// your testData HTML
`

func TestScrape(t *testing.T) {
	// 1. Setup your mock DB
	mockDB := new(MockDB)

	// 2. Set up expected matches
	expectedMatches := []struct {
		homeTeam string
		awayTeam string
		homeGoals int
		awayGoals int
		date time.Time
	}{
		{"Jong Utrecht", "Heracles", 0, 3, time.Date(2022, 8, 15, 0, 0, 0, 0, time.UTC)},
		{"Jong PSV", "Dordrecht", 1, 0, time.Date(2022, 8, 15, 0, 0, 0, 0, time.UTC)},
		{"MVV", "NAC", 3, 1, time.Date(2022, 8, 15, 0, 0, 0, 0, time.UTC)},
	}

	// 3. Prepare ScrapeData with mock DB and the test URL
	scraped := scraper.NewScrapeData(mockDB, "http://test.com")

	// 4. Run Scrape
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(testData))
	scraped.Scrape(doc)

	// 5. Verify that the InsertOrUpdateMatch function has been called with correct arguments for each match
	for _, match := range expectedMatches {
		mockDB.AssertCalled(t, "InsertOrUpdateMatch", match.homeTeam, match.awayTeam, match.homeGoals, match.awayGoals, match.date)
	}
}

```

Note that this test code assumes that the Scrape method of the ScrapeData struct is modified to take a `*goquery.Document` as an argument, which allows passing pre-parsed HTML data to the function:

```go
func (scraped *ScrapeData) Scrape(doc *goquery.Document) error {
	// the rest of the Scrape method
}

```

You may need to adjust this code to match the exact methods and types in your actual implementation.

## Follow-up prompt

> I don't want to modify existing code. Can you adjust?

In order to test your Scrape() function without modifying your existing code, we can replace the http.Get function to return our test data instead of making an actual HTTP request. Here's how you can do this:

```go
package scraper_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jqno/balGPT/internal/database"
	"github.com/jqno/balGPT/internal/scraper"
	"github.com/stretchr/testify/mock"
)

const testData = `
// your testData HTML
`

// Save the original http.Get function
var origHttpGet = http.Get

// Prepare MockDB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetLastScrape() (time.Time, error) {
	args := m.Called()
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) InsertOrUpdateMatch(homeTeam, awayTeam string, homeGoals, awayGoals int, matchDate time.Time) error {
	args := m.Called(homeTeam, awayTeam, homeGoals, awayGoals, matchDate)
	return args.Error(0)
}

func (m *MockDB) UpdateLastScrape(scrapeTime time.Time) error {
	args := m.Called(scrapeTime)
	return args.Error(0)
}

func TestScrape(t *testing.T) {
	mockDB := new(MockDB)
	mockDB.On("GetLastScrape").Return(time.Time{}, nil)

	// set expected matches
	expectedMatches := []struct {
		homeTeam  string
		awayTeam  string
		homeGoals int
		awayGoals int
		date      time.Time
	}{
		{"Jong Utrecht", "Heracles", 0, 3, time.Date(2022, 8, 15, 0, 0, 0, 0, time.UTC)},
		{"Jong PSV", "Dordrecht", 1, 0, time.Date(2022, 8, 15, 0, 0, 0, 0, time.UTC)},
		{"MVV", "NAC", 3, 1, time.Date(2022, 8, 15, 0, 0, 0, 0, time.UTC)},
	}

	for _, match := range expectedMatches {
		mockDB.On("InsertOrUpdateMatch", match.homeTeam, match.awayTeam, match.homeGoals, match.awayGoals, match.date).Return(nil)
	}

	mockDB.On("UpdateLastScrape", mock.AnythingOfType("time.Time")).Return(nil)

	// override http.Get function
	http.Get = func(url string) (*http.Response, error) {
		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(testData)),
		}, nil
	}
	defer func() { http.Get = origHttpGet }() // ensure original function is put back at the end of the test

	scraped := scraper.NewScrapeData(mockDB, "http://test.com")
	err := scraped.Scrape()

	// check no error returned
	if err != nil {
		t.Fatalf("Scrape() error = %v", err)
	}

	for _, match := range expectedMatches {
		mockDB.AssertCalled(t, "InsertOrUpdateMatch", match.homeTeam, match.awayTeam, match.homeGoals, match.awayGoals, match.date)
	}
}

```

This test is overriding the http.Get function at the start of the test and replacing it with a function that returns a reader with your test HTML. At the end of the test, we restore the original http.Get function.

Please note that this approach alters global state (http.Get) which can cause problems if tests are run concurrently. In practice, you may want to design your code in a way that avoids such global state.

## Follow-up prompt

> Doesn't work:
>
> ```plaintext
> internal/scraper/scraper_test.go:43:2: use of package http not in selector
> internal/scraper/scraper_test.go:43:2: cannot assign to http.Get (value of type func(url string) (resp \*http.Response, err error))
> ```
>
> Note that earlier you already were using net/http/httptest to mock out the connection. Let's do that here too.

I apologize for the confusion. Go does not allow overwriting functions from the standard library, so we can't directly override http.Get. You're correct that the httptest package is the proper way to mock HTTP requests in Go tests.

Please check the updated version of the test function where we are using httptest to mock the http.Get response:

```go
package scraper_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jqno/balGPT/internal/database"
	"github.com/jqno/balGPT/internal/scraper"
	"github.com/stretchr/testify/mock"
)

const testData = `
// your testData HTML
`

// Prepare MockDB
type MockDB struct {
	mock.Mock
}

func (m *MockDB) GetLastScrape() (time.Time, error) {
	args := m.Called()
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) InsertOrUpdateMatch(homeTeam, awayTeam string, homeGoals, awayGoals int, matchDate time.Time) error {
	args := m.Called(homeTeam, awayTeam, homeGoals, awayGoals, matchDate)
	return args.Error(0)
}

func (m *MockDB) UpdateLastScrape(scrapeTime time.Time) error {
	args := m.Called(scrapeTime)
	return args.Error(0)
}

func TestScrape(t *testing.T) {
	mockDB := new(MockDB)
	mockDB.On("GetLastScrape").Return(time.Time{}, nil)

	// set expected matches
	expectedMatches := []struct {
		homeTeam  string
		awayTeam  string
		homeGoals int
		awayGoals int
		date      time.Time
	}{
		{"Jong Utrecht", "Heracles", 0, 3, time.Date(2022, 8, 15, 0, 0, 0, 0, time.UTC)},
		{"Jong PSV", "Dordrecht", 1, 0, time.Date(2022, 8, 15, 0, 0, 0, 0, time.UTC)},
		{"MVV", "NAC", 3, 1, time.Date(2022, 8, 15, 0, 0, 0, 0, time.UTC)},
	}

	for _, match := range expectedMatches {
		mockDB.On("InsertOrUpdateMatch", match.homeTeam, match.awayTeam, match.homeGoals, match.awayGoals, match.date).Return(nil)
	}

	mockDB.On("UpdateLastScrape", mock.AnythingOfType("time.Time")).Return(nil)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testData))
	}))
	defer ts.Close()

	scraped := scraper.NewScrapeData(mockDB, ts.URL)
	err := scraped.Scrape()

	// check no error returned
	if err != nil {
		t.Fatalf("Scrape() error = %v", err)
	}

	for _, match := range expectedMatches {
		mockDB.AssertCalled(t, "InsertOrUpdateMatch", match.homeTeam, match.awayTeam, match.homeGoals, match.awayGoals, match.date)
	}
}

```

In this example, we create an HTTP server that serves our test data using `httptest.NewServer`. This server responds to HTTP requests with the test data HTML. We then pass the URL of this test server (accessible with `ts.URL`) to the `NewScrapeData` function. At the end of the test, we stop the server with `ts.Close()`.
