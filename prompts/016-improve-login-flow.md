# Improve login flow

> Let's modify index.html so that it shows only the login button and nothing else if we're not yet signed in, and that it shows the everything except the login button if we're signed in. Also add a message on the page after a failed sign-in attempt.

You can modify the index.html file to show only the login button when not signed in and hide it when the user is signed in. To achieve this, you can wrap the content you want to show or hide inside two separate divs, and then toggle their visibility based on the user's authentication status. Additionally, you can display a message after a failed sign-in attempt.

Update your index.html as follows:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Soccer Prediction</title>
  <style>
    body { font-family: Arial, sans-serif; }
    input[type="number"] { width: 50px; }
  </style>
</head>
<body>
  <div id="signedInContent" style="display:none;">
    <h1>Soccer Prediction</h1>
    <div>
      <label for="home_team_id">Home Team ID:</label>
      <input type="number" id="home_team_id" required>
    </div>
    <div>
      <label for="away_team_id">Away Team ID:</label>
      <input type="number" id="away_team_id" required>
    </div>
    <button id="predict_btn">Predict</button>
    <div id="result"></div>
  </div>
  <div id="signedOutContent">
    <h1>Please sign in</h1>
    <div id="signInBtn"></div>
    <div id="signInError" style="color: red; display: none;">Failed to sign in. Please try again.</div>
  </div>
  <script>
    const clientId = '{{.ClientID}}';
    const apiBaseUrl = '{{.ApiBaseURL}}';

    let accessToken = '';

    function onSignIn(googleUser) {
      const profile = googleUser.getBasicProfile();
      accessToken = googleUser.getAuthResponse().access_token;

      document.getElementById('signedInContent').style.display = 'block';
      document.getElementById('signedOutContent').style.display = 'none';
    }

    function onFailure(error) {
      console.log(error);
      document.getElementById('signInError').style.display = 'block';
    }

    function renderButton() {
      gapi.signin2.render('signInBtn', {
        'scope': 'profile email',
        'width': 200,
        'height': 50,
        'longtitle': true,
        'theme': 'dark',
        'onsuccess': onSignIn,
        'onfailure': onFailure
      });
    }

    // Rest of the script...

  </script>
  <script src="https://apis.google.com/js/platform.js?onload=renderButton" async defer></script>
  <meta name="google-signin-client_id" content="{{.ClientID }}">
</body>
</html>
```

With these changes:

1. The content you want to show when the user is signed in is wrapped inside a div with the ID signedInContent. By default, its display property is set to none.
2. The content you want to show when the user is not signed in is wrapped inside a div with the ID signedOutContent. This includes the Google Sign-In button and an error message for failed sign-in attempts.
3. When the user successfully signs in, the onSignIn function sets the display property of signedInContent to block and signedOutContent to none.
4. If there's an error during sign-in, the onFailure function makes the error message visible by setting the display property of signInError to block.

Now, the web page will only show the login button and an error message (if applicable) when the user is not signed in. Once the user is signed in, it will display the soccer prediction content and hide the login button.

