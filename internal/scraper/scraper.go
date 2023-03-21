package scraper

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jqno/balGPT/internal/database"
)

const URL = "https://www.fcupdate.nl/uitslagen/s1598/nederland-keuken-kampioen-divisie-2022-2023/"

type ScrapeData struct {
	DB *database.DB
}

func NewScrapeData(db *database.DB) *ScrapeData {
	return &ScrapeData{DB: db}
}

func (scraped *ScrapeData) Scrape() error {
	lastScrape, err := scraped.DB.GetLastScrape()
	if err != nil {
		return err
	}

	if lastScrape.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		return nil
	}

	res, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	currentDate := ""
	doc.Find("b, .line").EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if date := selection.Find("b").Text(); date != "" {
			currentDate = date
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
