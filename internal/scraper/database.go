package scraper

import "time"

type DB interface {
	GetLastScrape() (time.Time, error)
	InsertOrUpdateMatch(homeTeam, awayTeam string, homeGoals, awayGoals int, date time.Time) error
	UpdateLastScrape(t time.Time) error
}
