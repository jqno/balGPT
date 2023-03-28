# Access static files

> I have some files that the app needs access to, like migrations scripts in `./db/migrations`, static assets in `./static`, and website templates in `./templates`. How do I make sure they are also deployed correctly to fly.io?

To ensure that your migration scripts, static assets, and website templates are deployed correctly to Fly.io, you can create a Dockerfile in your project's root directory. This Dockerfile will define the container image used by Fly.io to run your application and will include the necessary files.

Here's an example `Dockerfile` for your Go app:

```Dockerfile
# Use the official Golang image as the base image
FROM golang:1.20.2 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o main

# Use a minimal base image for the final stage
FROM gcr.io/distroless/base-debian10

# Copy the built binary from the builder stage
COPY --from=builder /app/main /app/main

# Copy the migrations, static assets, and templates
COPY --from=builder /app/db/migrations /app/db/migrations
COPY --from=builder /app/static /app/static
COPY --from=builder /app/templates /app/templates

# Set the working directory
WORKDIR /app

# Start the application
CMD ["/app/main"]
```

This `Dockerfile` uses a multi-stage build process:

1. The first stage uses the official Golang image as the base image and builds your application.

2. The second stage uses a minimal base image and copies the built binary and the required files (migrations, static assets, and templates) from the builder stage.

Next, update your `fly.toml` file to use Docker for building and deploying your app. Replace the `[build]` section (or add it, if you removed it earlier) with the following:

```toml
[build]
  image = "registry.fly.io/your-app-name:deployment-${app.version}"
  dockerfile = "Dockerfile"
```

Replace `"your-app-name"` with the unique name of your app.

Now, when you deploy your app to Fly.io, it will use the specified `Dockerfile` to build the container image, including the required files. Your app should have access to the migration scripts, static assets, and templates at runtime.

Remember to adjust your app's code to look for these files in the `/app` directory within the container (e.g., `/app/db/migrations`, `/app/static`, and `/app/templates`).

## Follow-up prompt

> If I adjust the app's code to do that, I probably need to move files around locally as well to match. How do I do that?

You don't need to move files around locally. Instead, you can make your app configurable to use the correct paths based on the environment it's running in. You can use an environment variable to specify the base directory for your app, which will help you set the correct paths for your migrations, static assets, and templates.

First, add a new environment variable in your `fly.toml` file:

```toml
[env]
  # Your other environment variables...
  APP_BASE_DIR = "/app"
```

Then, in your Go app, read the `APP_BASE_DIR` environment variable and use it to construct the paths for your migrations, static assets, and templates. For example:

```go
package main

import (
  "os"
  "path/filepath"
)

func main() {
  // Read the APP_BASE_DIR environment variable
  appBaseDir := os.Getenv("APP_BASE_DIR")
  if appBaseDir == "" {
    // If the variable is not set, use the current directory
    appBaseDir = "."
  }

  // Construct the paths for your migrations, static assets, and templates
  migrationsDir := filepath.Join(appBaseDir, "db", "migrations")
  staticDir := filepath.Join(appBaseDir, "static")
  templatesDir := filepath.Join(appBaseDir, "templates")

  // Your application logic...
}
```

Now, your app will use the correct paths based on the environment it's running in. When running locally, the paths will point to your local directories (assuming the `APP_BASE_DIR` environment variable is not set). When running in the Fly.io container, the paths will point to the directories within the container as specified in the `fly.toml` file.

## Follow-up prompt

> I have this code:
>
>   fs := http.FileServer(http.Dir("static"))
>   http.Handle("/static/", http.StripPrefix("/static/", fs))
>
> And this code:
>
>   _, b, _, _ := runtime.Caller(0)
>   basepath := filepath.Dir(b)
>   migrationsPath := filepath.Join(basepath, "../../db/migrations")
>
> How do I adapt them to use the appBaseDir? Especially the second one looks like it should be simplified.

You can adapt the code to use the `appBaseDir` as follows:

First, read the `APP_BASE_DIR` environment variable and define your base paths:

```go
package main

import (
	"os"
	"path/filepath"
)

func main() {
	// Read the APP_BASE_DIR environment variable
	appBaseDir := os.Getenv("APP_BASE_DIR")
	if appBaseDir == "" {
		// If the variable is not set, use the current directory
		appBaseDir = "."
	}

	// Construct the paths for your migrations and static assets
	staticDir := filepath.Join(appBaseDir, "static")
	migrationsDir := filepath.Join(appBaseDir, "db", "migrations")

	// Your application logic...
}
```

Next, update the code that serves the static assets using the staticDir:

```go
	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
```

Finally, update the code that sets the migrations path to use the migrationsDir:

```go
	migrationsPath := migrationsDir
```

By using the appBaseDir, you've simplified the logic for determining the base path of your application, making it more readable and maintainable. Now, your code will work correctly both locally and when deployed to Fly.io.

## Note

The `image` in the `[build]` section of `fly.toml` was plain wrong, and unneeded, so I removed it.

More insidious was the fact that the GitHub Actions deployment workflow no longer worked, because Go expects a `.git` directory to be present so it can use the git commit hash for tagging. The error message suggests a fix: add a flag that disables this tagging. We'll just go ahead and do that.
