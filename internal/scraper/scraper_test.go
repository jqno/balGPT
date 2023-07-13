package scraper_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jqno/balGPT/internal/database_test"
	"github.com/jqno/balGPT/internal/scraper"
	"github.com/stretchr/testify/mock"
)

func TestScrape(t *testing.T) {
	mockDB := new(database_test.MockDB)
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

const testData = `
	<div class="matches-panel align-left justify-center notes">
	Maandag 15 augustus 2022
	</div>
	<div class="matches-panel d-flex align-center justify-center Played  ">
	<span class="fld-match">
	<a href="https://www.fcupdate.nl/voetbalteams/nederland/jong-fc-utrecht" class="left-team d-flex align-center justify-end">
	<span>Jong Utrecht</span>
	<img src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/ch82doz9w9t1sw9ae3a4ffh0j.png" alt="Jong Utrecht">
	</a>
	<a href="https://www.fcupdate.nl/voetbalcompetities/nederland/keuken-kampioen-divisie/programma-uitslagen/2022-2023/jong-fc-utrecht-heracles-15-08" class="score d-flex justify-center">
	<div class="match-result">
	<i>0</i>
	<i class="match-result__divder">-</i>
	<i>3</i>
	</div>
	</a>
	<a href="https://www.fcupdate.nl/voetbalteams/nederland/heracles" class="right-team d-flex align-center">
	<img src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/dac758ef858jbq7pcb3gfwite.png" alt="Heracles">
	<span><strong>Heracles</strong></span>
	</a>
	</span>
	</div>
	<div class="matches-panel d-flex align-center justify-center Played  ">
	<span class="fld-match">
	<a href="https://www.fcupdate.nl/voetbalteams/nederland/jong-psv" class="left-team d-flex align-center justify-end">
	<span><strong>Jong PSV</strong></span>
	<img src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/2xlnfgrl6r18ftv38unwp7h4s.png" alt="Jong PSV">
	</a>
	<a href="https://www.fcupdate.nl/voetbalcompetities/nederland/keuken-kampioen-divisie/programma-uitslagen/2022-2023/jong-psv-dordrecht-15-08" class="score d-flex justify-center">
	<div class="match-result">
	<i>1</i>
	<i class="match-result__divder">-</i>
	<i>0</i>
	</div>
	</a>
	<a href="https://www.fcupdate.nl/voetbalteams/nederland/dordrecht" class="right-team d-flex align-center">
	<img src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/z9phg19papi9f5fd6qxzr836.png" alt="Dordrecht">
	<span>Dordrecht</span>
	</a>
	</span>
	</div>
	<div class="matches-panel d-flex align-center justify-center Played  ">
	<span class="fld-match">
	<a href="https://www.fcupdate.nl/voetbalteams/nederland/mvv" class="left-team d-flex align-center justify-end">
	<span><strong>MVV</strong></span>
	<img src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/bqm1bnp1bh2apyudlz1o1whr3.png" alt="MVV">
	</a>
	<a href="https://www.fcupdate.nl/voetbalcompetities/nederland/keuken-kampioen-divisie/programma-uitslagen/2022-2023/mvv-nac-15-08" class="score d-flex justify-center">
	<div class="match-result">
	<i>3</i>
	<i class="match-result__divder">-</i>
	<i>1</i>
	</div>
	</a>
	<a href="https://www.fcupdate.nl/voetbalteams/nederland/nac" class="right-team d-flex align-center">
	<img src="https://imagecdn.fcupdate.nl/?height=96&amp;width=96&amp;url=https://ext.fcupdate.nl/teams/150x150/59t7flioj4w4mpwnrwbm0m8ck.png" alt="NAC">
	<span>NAC</span>
	</a>
	</span>
	</div>
	`
