# Handle empty database

> There's another small issue. When we run the app for the first time, the database is still empty and there are no teams. The Teams variable (which is filled in as a template from the Go backend) will be empty. In that case, I would like to display a single button that calls the 'scrape' endpoint. Can you generate the html for that? Also, can you generate code for the backend? The prediction code looks like this, you can adapt it as needed:
>
> (snipped `handlePrediction` function from `app.go`)

First, let's create an HTML button that will call the `/scrape` endpoint when clicked. Add the following button inside the div with the id `signedInContent`:

```html
<button id="scrape_btn" style="display:none;">Scrape Data</button>
```

Now, let's create a JavaScript function that sends a request to the `/scrape` endpoint:

```javascript
function scrapeData() {
  const url = `${apiBaseUrl}/scrape`;

  fetch(url, {
    method: 'GET',
  })
  .then(response => {
    if (response.ok) {
      alert('Data scraped successfully.');
      location.reload();
    } else {
      alert('Failed to scrape data.');
    }
  })
  .catch(error => {
    console.error('Error:', error);
    alert('Failed to scrape data.');
  });
}

document.getElementById('scrape_btn').addEventListener('click', scrapeData);
```

Make sure to add this JavaScript code to your existing JavaScript code in the HTML file.

Now, let's create a new handler function for the `/scrape` endpoint in your Go backend:

```go
func handleScrape(s *scraper.ScrapeData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := s.Scrape()
		if err != nil {
			http.Error(w, "Error while scraping data.", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
```

Finally, let's register the new endpoint in the `Run()` method of your App:

```go
http.HandleFunc("/scrape", handleScrape(a.Scraper))
```

Now, when you click the "Scrape Data" button, it will send a request to the `/scrape` endpoint, which will scrape the data and update the database. After scraping is completed, the page will be reloaded to display the updated teams.

## Follow-up prompt

> Can you update the front-end so that the scrape button is only ever shown if there are no teams, and that the existing widgets are only ever shown if there are teams?

Sure! You can use a conditional statement in the template to show or hide the buttons and widgets depending on whether there are any teams. Update your HTML code as follows:

```html
<div id="signedInContent" style="display:none;">
  <h1>Soccer Prediction</h1>
  {{if .Teams}}
    <div>
      <!-- Your existing widgets for selecting teams and making predictions -->
    </div>
  {{else}}
    <div>
      <p>No teams found. Please scrape data to get team information.</p>
      <button id="scrape_btn">Scrape Data</button>
    </div>
  {{end}}
</div>
```

By using the `{{if .Teams}}` conditional, we ensure that the team selection widgets and the predict button are only shown if there are teams in the .Teams variable. If there are no teams, a message and the "Scrape Data" button are displayed instead.

You can remove the `style="display:none;"` attribute from the `<button id="scrape_btn">` tag, as it is now conditionally rendered within the `{{if .Teams}}` block.

Also, you don't need the `location.reload();` line in the `scrapeData()` JavaScript function anymore, as the page will automatically display the correct content based on the presence of teams. Remove this line from your JavaScript code:

```javascript
location.reload();
```

Next, you need to add an event listener for the "Scrape Data" button to call the `scrapeData()` function when the button is clicked. Update your JavaScript code as follows:

```javascript
document.getElementById('scrape_btn')?.addEventListener('click', scrapeData);
```

Make sure to place this event listener after the scrapeData() function definition. The `?.` operator is used to ensure that the event listener is only added if the button with the ID `'scrape_btn'` exists. If there are teams, the button will not be in the DOM, and the event listener will not be added.

With these changes, your frontend should now display the appropriate content based on whether there are any teams in the database. When there are no teams, the "Scrape Data" button will be shown, and when there are teams, the team selection widgets and the "Predict" button will be displayed.

## Note

It actually added back a `style="display:none;"` attribute on the main `div`, so the UI was completely gone. I fixed it manually.

There were some other places where I had to add `?.` as well for the page to work, which I also did manually.
