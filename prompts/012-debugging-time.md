# Debugging time!

At this point I noticed there were a few bugs, so I went into an interactive session to resolve them. I won't reproduce the conversation here, but I'll mention the bugs themselves:

- If no scrapes had been performed yet, the `stats` table would be empty and the `GetLastScrape` function would return with a `sql.ErrNoRows` error. Instead, it should return an empty time and no error. Also, the `Scrape()` function should check for the zero case and continue.
- Dates were scraped incorrectly and as a result, were never found.
- Once the dates were found, they were parsed with the format `"2006-02-01"`, but they are actually formatted as `"01 februari 2006"`, with the Dutch month names. We needed to add a new date parsing function.

One interesting observation was that ChatGPT came up with a slightly complicated function that required new imports. In a follow up, I asked it if I needed to `go get` these imports. It replied yes, gave me the commands to do so, and then continued with "However, upon further review, there is a simpler solution that does not require additional dependencies. You can [...]" followed by a much simpler implementation that, indeed, didn't need any additional dependencies.
