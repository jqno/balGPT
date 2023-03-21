# Scraping the website

> Amazing, thank you!
>
> What are units of code called in Go? Modules? Classes? Files?
>
> I want to create one that scrapes a website for match results. Whenever it encounters a new team, it should be added to the teams table, and whenever it finds a new score, it should be added to matches. The match dates are in the sub headings in the table.
>
> When it's done, it should update last_scrape. Also, when last_scrape is set to today, it should not scrape at all.
>
> The website to scrape is at https://www.fcupdate.nl/uitslagen/s1598/nederland-keuken-kampioen-divisie-2022-2023/ This url should be set as an external property because it changes every year.
>
> Here's part of the page's source. The beginning and end of the file change occasionally, but the middle containing the divs with the scores has been stable for many years.
>
> <b>13 maart 2023</b>                            <div class="line even">
>                 <div class="center score"><span class="big"><a href="/voetbaluitslagen/364619/keuken-kampioen-divisie-almere-city-fc-vvv-venlo/" title="Wedstrijdverloop: Almere City FC - VVV-Venlo">1 - 3</a></span></div>
>                 <div class="float-left club"><a href="/voetbalclubs/383/almere-city-fc/">Almere City FC<img src="https://adn2.fcupdate.nl/clubs/small/383.png" width="16" height="16"></a></div>
>                 <div class="float-right club"><a href="/voetbalclubs/15/vvv-venlo/"><img src="https://adn2.fcupdate.nl/clubs/small/15.png" width="16" height="16">VVV-Venlo</a></div>
>               </div>
>                                                                     <div class="line odd">
>                 <div class="center score"><span class="big"><a href="/voetbaluitslagen/364620/keuken-kampioen-divisie-fc-eindhoven-ado-den-haag/" title="Wedstrijdverloop: FC Eindhoven - ADO Den Haag">1 - 1</a></span></div>
>                 <div class="float-left club"><a href="/voetbalclubs/3/fc-eindhoven/">FC Eindhoven<img src="https://adn2.fcupdate.nl/clubs/small/3.png" width="16" height="16"></a></div>
>                 <div class="float-right club"><a href="/voetbalclubs/27/ado-den-haag/"><img src="https://adn2.fcupdate.nl/clubs/small/27.png" width="16" height="16">ADO Den Haag</a></div>
>               </div>
>                                                                     <div class="line even">
>                 <div class="center score"><span class="big"><a href="/voetbaluitslagen/364621/keuken-kampioen-divisie-de-graafschap-pec-zwolle/" title="Wedstrijdverloop: De Graafschap - PEC Zwolle">0 - 1</a></span></div>
>                 <div class="float-left club"><a href="/voetbalclubs/20/de-graafschap/">De Graafschap<img src="https://adn2.fcupdate.nl/clubs/small/20.png" width="16" height="16"></a></div>
>                 <div class="float-right club"><a href="/voetbalclubs/114/pec-zwolle/"><img src="https://adn2.fcupdate.nl/clubs/small/114.png" width="16" height="16">PEC Zwolle</a></div>
>               </div>
>                                                                     <div class="line odd">
>                 <div class="center score"><span class="big"><a href="/voetbaluitslagen/364622/keuken-kampioen-divisie-nac-breda-jong-ajax/" title="Wedstrijdverloop: NAC Breda - Jong Ajax">2 - 6</a></span></div>
>                 <div class="float-left club"><a href="/voetbalclubs/28/nac-breda/">NAC Breda<img src="https://adn2.fcupdate.nl/clubs/small/28.png" width="16" height="16"></a></div>
>                 <div class="float-right club"><a href="/voetbalclubs/430/jong-ajax/"><img src="https://adn2.fcupdate.nl/clubs/small/430.png" width="16" height="16">Jong Ajax</a></div>
>               </div>
>                                                                     <div class="line even">
>                 <div class="center score"><span class="big"><a href="/voetbaluitslagen/364623/keuken-kampioen-divisie-roda-jc-jong-az/" title="Wedstrijdverloop: Roda JC - Jong AZ">2 - 0</a></span></div>
>                 <div class="float-left club"><a href="/voetbalclubs/22/roda-jc/">Roda JC<img src="https://adn2.fcupdate.nl/clubs/small/22.png" width="16" height="16"></a></div>
>                 <div class="float-right club"><a href="/voetbalclubs/32219/jong-az/"><img src="https://adn2.fcupdate.nl/clubs/small/32219.png" width="16" height="16">Jong AZ</a></div>
>               </div>
>                                                                     <div class="line odd">
>                 <div class="center score"><span class="big"><a href="/voetbaluitslagen/364624/keuken-kampioen-divisie-jong-fc-utrecht-jong-psv/" title="Wedstrijdverloop: Jong FC Utrecht - Jong PSV">0 - 1</a></span></div>
>                 <div class="float-left club"><a href="/voetbalclubs/31857/jong-fc-utrecht/">Jong FC Utrecht<img src="https://adn2.fcupdate.nl/clubs/small/31857.png" width="16" height="16"></a></div>
>                 <div class="float-right club"><a href="/voetbalclubs/28079/jong-psv/"><img src="https://adn2.fcupdate.nl/clubs/small/28079.png" width="16" height="16">Jong PSV</a></div>
>               </div>
>                                         <b>10 maart 2023</b>                            <div class="line even">
>                 <div class="center score"><span class="big"><a href="/voetbaluitslagen/364609/keuken-kampioen-divisie-heracles-almelo-jong-az/" title="Wedstrijdverloop: Heracles Almelo - Jong AZ">2 - 0</a></span></div>
>                 <div class="float-left club"><a href="/voetbalclubs/18/heracles-almelo/">Heracles Almelo<img src="https://adn2.fcupdate.nl/clubs/small/18.png" width="16" height="16"></a></div>
>                 <div class="float-right club"><a href="/voetbalclubs/32219/jong-az/"><img src="https://adn2.fcupdate.nl/clubs/small/32219.png" width="16" height="16">Jong AZ</a></div>
>               </div>
>                                                                     <div class="line odd">
>                 <div class="center score"><span class="big">Afg.</span></div>
>                 <div class="float-left club"><a href="/voetbalclubs/383/almere-city-fc/">Almere City FC<img src="https://adn2.fcupdate.nl/clubs/small/383.png" width="16" height="16"></a></div>
>                 <div class="float-right club"><a href="/voetbalclubs/31857/jong-fc-utrecht/"><img src="https://adn2.fcupdate.nl/clubs/small/31857.png" width="16" height="16">Jong FC Utrecht</a></div>
>               </div>


In Go, units of code are called packages. A package is a collection of Go source files in the same directory with the same package name at the top of each file. You can create a new package for scraping the website and handling the data.

First, add the required libraries for scraping:

```sh
$ go get -u github.com/PuerkitoBio/goquery
```

Now, create a new package called scraper:

internal/scraper/scraper.go:

```go
package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/YOUR_USERNAME/balgpt/internal/database"
)

const URL = "https://www.fcupdate.nl/uitslagen/s1598/nederland-keuken-kampioen-divisie-2022-2023/"

type ScrapeData struct {
	DB *database.DB
}

func NewScrapeData(db *database.DB) *ScrapeData {
	return &ScrapeData{DB: db}
}

func (s *ScrapeData) Scrape() error {
	lastScrape, err := s.DB.GetLastScrape()
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
	doc.Find("b, .line").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if date := s.Find("b").Text(); date != "" {
			currentDate = date
			return true
		}

		score := s.Find(".center.score > span.big > a").Text()
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

		homeTeam := strings.TrimSpace(s.Find(".float-left.club > a").Text())
		awayTeam := strings.TrimSpace(s.Find(".float-right.club > a").Text())

		if homeTeam == "" || awayTeam == "" {
			return true
		}

		err = s.DB.InsertOrUpdateMatch(homeTeam, awayTeam, homeGoals, awayGoals, currentDate)
		if err != nil {
			return true
		}

		return true
	})

	err = s.DB.UpdateLastScrape(time.Now())
	if err != nil {
		return err
	}

	return nil
}
```

Next, add the required methods to the DB struct in the internal/database/database.go file:

```go
// Add this import
import "time"

// Add the following methods

func (db *DB) GetLastScrape() (time.Time, error) {
	var lastScrape time.Time
	err := db.Conn.QueryRow("SELECT last_scrape FROM stats ORDER BY id DESC LIMIT 1").Scan(&lastScrape)
	if err != nil {
		return time.Time{}, err
	}

	return lastScrape, nil
}

func (db *DB) InsertOrUpdateMatch(homeTeam, awayTeam string, homeGoals, awayGoals int, dateStr string) error {
```

## Note

Somehow, as you can see, I didn't get the body for the `InsertOrUpdateMatch` function. The function `UpdateLastScrape` was also missing. I'm not sure if the tool I used to save the conversation dropped them for some reason, or if ChatGPT didn't actually generate them. I had no access to the ChatGPT history at the time of writing, so I had to ask for alternative implementations in a separate conversation.

Also, there was a small bug in the scraper code: there were two variables named `s`, of different types, and they clashed at one point. I gave them different names.
