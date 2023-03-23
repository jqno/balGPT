# Fixing a bug

> Can you fix a small bug for me? When I select teams, the home_team_id and away_team_id aren't actually updated, so I always get the "Please enter both team IDs" error.

I see the issue. It's because the `input` event listener is trying to find the `option` element by its name, but it should be looking for the `option` by its text content. Here's the updated JavaScript code that should fix the issue:

```javascript
document.getElementById('home_team_name').addEventListener('input', function (e) {
  const selectedOption = Array.from(e.target.list.options).find(option => option.textContent === e.target.value);
  document.getElementById('home_team_id').value = selectedOption ? selectedOption.getAttribute('data-value') : '';
});

document.getElementById('away_team_name').addEventListener('input', function (e) {
  const selectedOption = Array.from(e.target.list.options).find(option => option.textContent === e.target.value);
  document.getElementById('away_team_id').value = selectedOption ? selectedOption.getAttribute('data-value') : '';
});
```

In this updated code, I replaced `e.target.list.options.namedItem(e.target.value)` with `Array.from(e.target.list.options).find(option => option.textContent === e.target.value)` to find the `option` element by its text content instead of its name. This should correctly update the `home_team_id` and `away_team_id` when selecting teams.
