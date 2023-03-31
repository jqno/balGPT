# Rewrite scraper because website changed

> Good morning! I'm going to give you some context first, and then I'll ask you some questions, ok?

Good morning! Sure, please go ahead and provide the context, and then feel free to ask your questions. I'll do my best to help you.

## Follow-up prompt

> First bit of context; no questions yet.
>
> Earlier I asked you to create a scraper for a page that lists soccer match results. The page looks like this:
>
> _(snipped: what the page looked like before; see [prompts/004-scraping-the-website.md](prompt 4))_

Thank you for providing the context. I can see that this is an HTML snippet containing soccer match results. The match results are organized within different lines that are either labeled with class "even" or "odd". Each line displays the score, the competing teams and their respective logos. There are also dates in bold format, which indicate when the matches took place.

Please feel free to ask any questions related to this context, and I'll do my best to help you.

## Follow-up prompt

> Second bit of context, still no questions, so no need to reply.
>
> Here is the code that you generated for me, which works very well with the given website. It's a Go package.
>
> _(snipped: scraper.go)_

Thank you for providing the second bit of context. I see that this is the Go package code that I provided earlier for scraping soccer match results from the given website. The package includes a ScrapeData struct and functions for creating a new ScrapeData instance, scraping the data from the website, and parsing Dutch dates.

Now that I have both context snippets, feel free to ask any questions related to this code or the HTML snippet, and I'll be happy to help.

## Follow-up prompt

> Third bit of context: They changed the website! The scaping doesn't work anymore!
>
> Here's what it looks like now:
>
>    <div class="matches-panel align-left justify-center notes">
>    Zondag 19 maart 2023
>    </div>
>    <div class="matches-panel d-flex align-center justify-center Played  ">
>    <span class="fld-match">
>    <a href="https://www.fcupdate.nl/voetbalteams/nederland/vvv" class="left-team d-flex align-center justify-end">
>    <span>VVV</span>
>    <img class="lazyload" loading="lazy" data-src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/cv5ftekuejf0ryxvd5h2my47r.png" alt="VVV">
>    </a>
>    <a href="https://www.fcupdate.nl/voetbalcompetities/nederland/keuken-kampioen-divisie/programma-uitslagen/2022-2023/vvv-roda-19-03" class="score d-flex justify-center">
>    <div class="match-result">
>    <i>1</i>
>    <i class="match-result__divder">-</i>
>    <i>1</i>
>    </div>
>    </a>
>    <a href="https://www.fcupdate.nl/voetbalteams/nederland/roda" class="right-team d-flex align-center">
>    <img class="lazyload" loading="lazy" data-src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/7rmv5ns2rj7bn4otpkur0mx0x.png" alt="Roda">
>    <span>Roda</span>
>    </a>
>    </span>
>    </div>
>    <div class="matches-panel align-left justify-center notes">
>    Zondag 26 maart 2023
>    </div>
>    <div class="matches-panel d-flex align-center justify-center Played  ">
>    <span class="fld-match">
>    <a href="https://www.fcupdate.nl/voetbalteams/nederland/dordrecht" class="left-team d-flex align-center justify-end">
>    <span>Dordrecht</span>
>    <img class="lazyload" loading="lazy" data-src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/z9phg19papi9f5fd6qxzr836.png" alt="Dordrecht">
>    </a>
>    <a href="https://www.fcupdate.nl/voetbalcompetities/nederland/keuken-kampioen-divisie/programma-uitslagen/2022-2023/dordrecht-helmond-13-03" class="score d-flex justify-center">
>    <div class="match-result">
>    <i>0</i>
>    <i class="match-result__divder">-</i>
>    <i>0</i>
>    </div>
>    </a>
>    <a href="https://www.fcupdate.nl/voetbalteams/nederland/helmond" class="right-team d-flex align-center">
>    <img class="lazyload" loading="lazy" data-src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/b94pf87be9ojdgnrs5mst80un.png" alt="Helmond">
>    <span>Helmond</span>
>    </a>
>    </span>
>    </div>
>    <div class="matches-panel d-flex align-center justify-center Played  ">
>    <span class="fld-match">
>    <a href="https://www.fcupdate.nl/voetbalteams/nederland/heracles" class="left-team d-flex align-center justify-end">
>    <span><strong>Heracles</strong></span>
>    <img class="lazyload" loading="lazy" data-src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/dac758ef858jbq7pcb3gfwite.png" alt="Heracles">
>    </a>
>    <a href="https://www.fcupdate.nl/voetbalcompetities/nederland/keuken-kampioen-divisie/programma-uitslagen/2022-2023/heracles-den-bosch-13-03" class="score d-flex justify-center">
>    <div class="match-result">
>    <i>3</i>
>    <i class="match-result__divder">-</i>
>    <i>2</i>
>    </div>
>    </a>
>
> But watch out, they added upcoming matches in the same page. These cannot be scraped. It looks very similar:
>
>    <div class="matches-panel align-left justify-center notes" id="aankomende-wedstrijden">
>    Vrijdag 31 maart 2023
>    </div>
>    <div class="matches-panel d-flex align-center justify-center  Fixture ">
>    <span class="fld-match">
>    <a href="https://www.fcupdate.nl/voetbalteams/nederland/ado" class="left-team d-flex align-center justify-end">
>    <span>ADO</span>
>    <img class="lazyload" loading="lazy" data-src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/e8hiwrw6ocxm5iht6sp7jiv19.png" alt="ADO">
>     </a>
>    <a href="https://www.fcupdate.nl/voetbalcompetities/nederland/keuken-kampioen-divisie/programma-uitslagen/2022-2023/ado-telstar-31-03" class="score d-flex justify-center">
>    <span>
>    20:00
>    </span>
>    </a>
>    <a href="https://www.fcupdate.nl/voetbalteams/nederland/telstar" class="right-team d-flex align-center">
>    <img class="lazyload" loading="lazy" data-src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/zjawy6qdzk3y69v5pa3zl2cb.png" alt="Telstar">
>    <span>Telstar</span>
>    </a>
>    </span>
>    </div>
>    <div class="matches-panel d-flex align-center justify-center  Fixture ">
>    <span class="fld-match">
>    <a href="https://www.fcupdate.nl/voetbalteams/nederland/heracles" class="left-team d-flex align-center justify-end">
>    <span>Heracles</span>
>    <img class="lazyload" loading="lazy" data-src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/dac758ef858jbq7pcb3gfwite.png" alt="Heracles">
>    </a>
>    <a href="https://www.fcupdate.nl/voetbalcompetities/nederland/keuken-kampioen-divisie/programma-uitslagen/2022-2023/heracles-almere-31-03" class="score d-flex justify-center">
>    <span>
>    20:00
>    </span>
>    </a>
>    <a href="https://www.fcupdate.nl/voetbalteams/nederland/almere" class="right-team d-flex align-center">
>    <img class="lazyload" loading="lazy" data-src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/1fttb31hnskynpku8qd09yhm8.png" alt="Almere">
>    <span>Almere</span>
>    </a>
>    </span>
>    </div>
>
> What I want you to do is take the original Go code and modify it so it scrapes the new version of the website. Please stay as close to the original code as possible. Most importantly I don't want to change the function signatures, and I don't want to change the imports either if possible. Note that the dates are again in the Dutch format. When parsing, you can ignore the names of the days and use the parseDutchDate function to parse the rest.

I understand that the website has been updated and the original Go code is no longer working. I've modified the code to scrape the new version of the website, keeping the function signatures and imports unchanged. Here's the updated code:

```go
package scraper

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jqno/balGPT/internal/database"
)

type ScrapeData struct {
	DB  *database.DB
	URL string
}

func NewScrapeData(db *database.DB, url string) *ScrapeData {
	return &ScrapeData{DB: db, URL: url}
}

func (scraped *ScrapeData) Scrape() error {
	lastScrape, err := scraped.DB.GetLastScrape()
	if err != nil {
		return err
	}

	if !lastScrape.IsZero() && lastScrape.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		return nil
	}

	res, err := http.Get(scraped.URL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	currentDate := time.Time{}
	doc.Find(".matches-panel").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if selection.HasClass("align-left") && selection.HasClass("justify-center") {
			dateStr := strings.TrimSpace(selection.Text())
			dateStr = strings.TrimPrefix(dateStr, "Zondag ")
			currentDate, _ = parseDutchDate(dateStr)
			return true
		}

		if !selection.HasClass("Played") {
			return true
		}

		homeTeam := strings.TrimSpace(selection.Find(".left-team > span").Text())
		awayTeam := strings.TrimSpace(selection.Find(".right-team > span").Text())

		scoreSelection := selection.Find(".score > div > i")
		homeGoalsStr := strings.TrimSpace(scoreSelection.First().Text())
		awayGoalsStr := strings.TrimSpace(scoreSelection.Last().Text())

		homeGoals, err := strconv.Atoi(homeGoalsStr)
		if err != nil {
			return true
		}

		awayGoals, err := strconv.Atoi(awayGoalsStr)
		if err != nil {
			return true
		}

		if homeTeam == "" || awayTeam == "" {
			return true
		}

		err = scraped.DB.InsertOrUpdateMatch(homeTeam, awayTeam, homeGoals, awayGoals, currentDate)
		if err != nil {
			return true
		}

		return true
	})

	err = scraped.DB.UpdateLastScrape(time.Now())
	if err != nil {
		return err
	}

	return nil
}

func parseDutchDate(dateStr string) (time.Time, error) {
	monthNameToNumber := map[string]string{
		"januari": "01", "februari": "02", "maart": "03", "april": "04", "mei": "05", "juni": "06",
		"juli": "07", "augustus": "08", "september": "09", "oktober": "10", "november": "11", "december": "12",
	}

	for monthName, monthNumber := range monthNameToNumber {
		dateStr = strings.Replace(dateStr, monthName, monthNumber, 1)
	}

	// Parse the date string
	date, err := time.Parse("02 01 2006", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}
```

This modification should now work with the updated website structure. The code now selects elements with the `.matches-panel` class and checks for date and match elements accordingly. It also skips upcoming matches that don't have the `Played` class. The rest of the code remains largely the same, including the function signatures and imports.

## Follow-up prompt

> Thank you, this works well! However, most of the dates come back as "0001-01-01". I think this happens because of strings.TrimPrefix(dateStr, "Zondag "), which of course only works for Sundays. Please generalize this section of the Scrape() function so it works for any day, preferably without actually enumerating all possible days names, since we don't need this information anyway.
>
> Don't regenerate the entire function, you only have to give me this particular piece.

You're correct that the issue is due to the use of `strings.TrimPrefix(dateStr, "Zondag ")`. We can generalize this by using a regular expression to remove any weekday prefix from the date string. Here's the updated code snippet:

```go
import (
	// ...
	"regexp"
	// ...
)

func (scraped *ScrapeData) Scrape() error {
	// ...

	doc.Find(".matches-panel").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if selection.HasClass("align-left") && selection.HasClass("justify-center") {
			dateStr := selection.Text()
			// Remove any weekday prefix from the date string
			re := regexp.MustCompile(`^\w+\s`)
			dateStr = re.ReplaceAllString(dateStr, "")

			currentDate, _ = parseDutchDate(dateStr)
			return true
		}

		// ...
	})

	// ...
}
```

This code snippet uses the `regexp` package to create a regular expression that matches any word followed by a space at the beginning of the string. It then replaces the matched substring with an empty string, effectively removing the weekday prefix from the date string. This should work for any day of the week without enumerating all possible day names.

## Follow-up prompt

> Another small issue: dates that have a single-digit day number don't get parsed correctly, resulting in 0001-01-01 as well. Can you fix that? Please only show the code that needs to be changed.

To fix this issue, you can modify the `parseDutchDate` function to handle single-digit day numbers. You can achieve this by using the `time.Parse` function with the "2 01 2006" format string and then falling back to the "02 01 2006" format string if an error occurs. Here's the updated `parseDutchDate` function:

```go
func parseDutchDate(dateStr string) (time.Time, error) {
	// ...

	// Parse the date string with "2 01 2006" format
	date, err := time.Parse("2 01 2006", dateStr)
	if err != nil {
		// Fallback to "02 01 2006" format if there's an error
		date, err = time.Parse("02 01 2006", dateStr)
		if err != nil {
			return time.Time{}, err
		}
	}

	return date, nil
}
```

This change should allow the function to correctly parse dates with single-digit day numbers.

## Follow-up prompt

> That works! Thank you!

You're welcome! I'm glad it works for you. If you have any more questions or need further assistance, feel free to ask. Good luck with your project!
