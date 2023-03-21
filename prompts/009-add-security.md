# Add security

> Let's also add proper security to the endpoints. I want to authenticate via Google. For now we'll only allow access to one email address (my own). Let's have an environment variable for that. I'll fill it in later myself.

To implement Google authentication, we will use the golang.org/x/oauth2 package for OAuth2 and the google.golang.org/api/oauth2/v2 package for accessing the Google API. First, let's add these dependencies to your project:

```sh
go get golang.org/x/oauth2
go get google.golang.org/api/oauth2/v2
```

Now, we need to set up Google OAuth2 configuration. Add the following imports to your main.go file:

```go
import (
	// ...
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
)
```

Create a function to initialize Google OAuth2 config and obtain an authenticated HTTP client:

```go
func getAuthenticatedClient(accessToken string) (*http.Client, error) {
	ctx := context.Background()

	// Use the access token to create an authenticated client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	client := oauth2.NewClient(ctx, ts)

	// Check if the token is valid
	oauth2Service, err := oauth2api.New(client)
	if err != nil {
		return nil, err
	}
	tokenInfo, err := oauth2Service.Tokeninfo().AccessToken(accessToken).Do()
	if err != nil {
		return nil, err
	}

	// Check if the email address is allowed
	allowedEmail := os.Getenv("ALLOWED_EMAIL")
	if allowedEmail == "" || allowedEmail != tokenInfo.Email {
		return nil, fmt.Errorf("unauthorized email address")
	}

	return client, nil
}
```

Now, let's create a middleware function to handle Google authentication:

```go
func googleAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request)
```

## Note

Since we took the allowed mail from the environment, I decided to manually move it to the application's config as well.

Also, ChatGPT didn't complete the `googleAuthMiddleware` function. When I had it re-generated, one import turned out not to be needed. I wonder what it would have generated on the first try.

Also, see the next prompt.
