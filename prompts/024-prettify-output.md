# Prettify output

> Can you make the output of the prediction look a little nicer? The raw json it shows now isn't very pretty.

To make the output of the prediction look nicer, you can create a function to format the JSON data into a more human-readable format. Update the JavaScript code inside the `<script>` tag as follows:

1. Replace the following line:

```javascript
document.getElementById('result').innerHTML = JSON.stringify(data);
```

with:

```javascript
document.getElementById('result').innerHTML = formatPredictionResult(data);
```

1. Add the formatPredictionResult() function:

```javascript
function formatPredictionResult(data) {
  const homeTeam = data.home_team_name;
  const awayTeam = data.away_team_name;
  const homeWinProb = (data.home_team_win_probability * 100).toFixed(2);
  const awayWinProb = (data.away_team_win_probability * 100).toFixed(2);
  const drawProb = (data.draw_probability * 100).toFixed(2);

  return `
    <h2>Prediction Result</h2>
    <table>
      <tr>
        <th>Team</th>
        <th>Win Probability</th>
      </tr>
      <tr>
        <td>${homeTeam}</td>
        <td>${homeWinProb}%</td>
      </tr>
      <tr>
        <td>${awayTeam}</td>
        <td>${awayWinProb}%</td>
      </tr>
      <tr>
        <td>Draw</td>
        <td>${drawProb}%</td>
      </tr>
    </table>
  `;
}
```

1. (Optional) Add some CSS to style the table. Add the following code inside the <style> tag or in the linked stylesheet:

```css
table {
  width: 100%;
  border-collapse: collapse;
}

table, th, td {
  border: 1px solid black;
}

th, td {
  padding: 8px;
  text-align: left;
}

th {
  background-color: #f2f2f2;
}
```

With these changes, the prediction result will be displayed in a nicely formatted table instead of raw JSON.

## Follow-up prompt

> Hang on, here's what the json actually looks like:
>
> {"HomeGoals":2,"AwayGoals":2}

In that case, you can modify the formatPredictionResult() function to display the home and away goals. Update the JavaScript code inside the <script> tag as follows:

1. Replace the existing formatPredictionResult() function with:

```javascript
function formatPredictionResult(data) {
  const homeGoals = data.HomeGoals;
  const awayGoals = data.AwayGoals;

  return `
    <h2>Prediction Result</h2>
    <table>
      <tr>
        <th>Team</th>
        <th>Goals</th>
      </tr>
      <tr>
        <td>Home Team</td>
        <td>${homeGoals}</td>
      </tr>
      <tr>
        <td>Away Team</td>
        <td>${awayGoals}</td>
      </tr>
    </table>
  `;
}
```

With this change, the output will display the predicted goals for both home and away teams in a nicely formatted table. The CSS provided in the previous answer can still be used to style the table.
