# Remove authentication

> Good morning!
> I have an app with authentication that doesn't work. I want to remove the authentication altogether, starting with the backend. Here it is:
>
> (snipped: current version of `app.go` file)
>
> Can you remove all the authentication stuff?

and follow-up prompt:

> Nice! Now let's do the frontend. It's just a single html page:
>
> (snipped: current version of `index.html` file)

Sure, I can help you remove the authentication from the backend code. Here's the updated code without authentication:

(snipped: `app.go`, with authentication removed)

In this updated code, I have removed the `googleAuthMiddleware` function and related imports, and also removed the middleware from the `http.HandleFunc` calls for `/predict` and `/team_id`. This should disable the authentication in your backend.

To remove the authentication from the frontend, you'll need to remove the Google Sign-In button, its related functions and scripts, and also update the fetch call in the makePrediction function. Here's the updated HTML code:

(snipped: `index.html`, with authentication removed)

In this updated code, I removed the Google Sign-In button and its related functions (onSignIn, onFailure, renderButton), the Google Sign-In script tag, and the related meta tag. I also removed the accessToken variable and updated the fetch call in the makePrediction function to no longer include the Authorization header. This should disable the authentication in your frontend.

