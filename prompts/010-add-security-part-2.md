# Adding security, part 2

> BTW, what do we do with the googleAuthMiddleware function? How do we connect it to the endpoints?

To use the googleAuthMiddleware function with your endpoints, you need to wrap your handler functions with the middleware function. In your main.go file, where you set up the routes, you can wrap the handlers using the middleware like this:

```go
package main

import (
	// ...
)

func main() {
	// ...
	
	http.HandleFunc("/get_team_id", googleAuthMiddleware(getTeamIDHandler))
	http.HandleFunc("/predict", googleAuthMiddleware(predictHandler))
	
	// ...
}
```

By wrapping your handler functions with the googleAuthMiddleware, you're ensuring that the middleware function will be called before your handler functions. If the request passes the authentication check, your handler functions will be executed; otherwise, the middleware will send an unauthorized response.

## Note

Turns out that the return type for `googleAuthMiddleware` should not be `http.Handler`, but `http.HandlerFunc`.

