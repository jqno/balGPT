# Improve UX for autocompleting inputs

> Great stuff! Can you make it so it displays the team name instead of the ID in the input field, but still sends the id to the backend?

Sure! You can modify the index.html and the JavaScript code to display the team name in the input field while still sending the team ID to the backend.

1. Update the index.html file to separate the team name and ID in the datalist options, and add hidden input fields to store the team IDs:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <!-- ... -->
</head>
<body>
  <div id="signedInContent" style="display:none;">
    <h1>Soccer Prediction</h1>
    <div>
      <label for="home_team_name">Home Team:</label>
      <input list="home_teams" id="home_team_name" required>
      <input type="hidden" id="home_team_id">
      <datalist id="home_teams">
        {{range .Teams}}
        <option data-value="{{.ID}}">{{.Name}}</option>
        {{end}}
      </datalist>
    </div>
    <div>
      <label for="away_team_name">Away Team:</label>
      <input list="away_teams" id="away_team_name" required>
      <input type="hidden" id="away_team_id">
      <datalist id="away_teams">
        {{range .Teams}}
        <option data-value="{{.ID}}">{{.Name}}</option>
        {{end}}
      </datalist>
    </div>
    <!-- ... -->
  </div>
  <!-- ... -->
</body>
```

1. Update the JavaScript code to handle the input event and store the selected team ID in the hidden input fields:

```javascript
document.getElementById('home_team_name').addEventListener('input', function (e) {
  const selectedOption = e.target.list.options.namedItem(e.target.value);
  document.getElementById('home_team_id').value = selectedOption ? selectedOption.getAttribute('data-value') : '';
});

document.getElementById('away_team_name').addEventListener('input', function (e) {
  const selectedOption = e.target.list.options.namedItem(e.target.value);
  document.getElementById('away_team_id').value = selectedOption ? selectedOption.getAttribute('data-value') : '';
});
```

1. Update the API call in your JavaScript code to use the team IDs from the hidden input fields instead of the visible ones:

```javascript
document.getElementById('predict_btn').addEventListener('click', async function () {
  const homeTeamId = document.getElementById('home_team_id').value;
  const awayTeamId = document.getElementById('away_team_id').value;

  if (homeTeamId === '' || awayTeamId === '') {
    // Show an error message if the team IDs are not set
    return;
  }

  // Make the API call using the team IDs from the hidden input fields
  // ...
});
```

With these changes, the input fields will display the team names, while the team IDs are stored in hidden input fields. The API call will use the team

## Note

That last step was totally unnecessary but hey ¯\\_(ツ)_/¯
