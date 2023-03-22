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
	doc.Find("b, .line").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if selection.Is("b") {
			currentDate, _ = parseDutchDate(selection.Text())
			return true
		}

		score := selection.Find(".center.score > span.big > a").Text()
		scores := strings.Split(score, " - ")

		if len(scores) != 2 {
			return true
		}

		homeGoals, err := strconv.Atoi(strings.TrimSpace(scores[0]))
		if err != nil {
			return true
		}

		awayGoals, err := strconv.Atoi(strings.TrimSpace(scores[1]))
		if err != nil {
			return true
		}

		homeTeam := strings.TrimSpace(selection.Find(".float-left.club > a").Text())
		awayTeam := strings.TrimSpace(selection.Find(".float-right.club > a").Text())

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
