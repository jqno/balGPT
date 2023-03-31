package scraper

import (
	"net/http"
	"strconv"
	"regexp"
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
			re := regexp.MustCompile(`^\w+\s`)
			dateStr = re.ReplaceAllString(dateStr, "")
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
	date, err := time.Parse("2 01 2006", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}
