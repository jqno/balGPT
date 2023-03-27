# Simplify build

> I have a Go app that I build using `go build ./cmd/something/main.go`. Can I make it build with just `go build` somehow? And can I make it run as well with just `go run`?
>
> What would a GitHub actions workflow look like for such a build?

Yes, you can simplify the build and run commands by organizing your Go project structure according to the Go convention. The Go convention suggests putting the main package at the root of your project, which allows you to use `go build` and `go run` without specifying the path.

Here's how you can reorganize your project structure:

1. Move the contents of the `./cmd/something/main.go` file to a main.go file in the root of your project.
2. Update any package import paths if necessary.
3. If you have other files inside the ./cmd/something directory, you might want to move them to a separate package directory within your project, then import and use them in your main.go file.

After reorganizing your project, you can build it with just `go build` and run it with just `go run`.

To create a GitHub Actions workflow to build your Go app, you'll need to create a YAML configuration file inside the `.github/workflows` directory in your project's root. You can name the file something like `go_build.yml`. Here's an example of a GitHub Actions workflow configuration for building a Go app:

```yaml
name: Go Build

on:
  push:
    branches:
      - main
  pull_request:
    branches:
      - main

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - name: Set up Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.17

    - name: Check out code
      uses: actions/checkout@v2

    - name: Get dependencies
      run: go mod download

    - name: Build
      run: go build

    - name: Test
      run: go test -v ./...

```

This workflow does the following:

1. Triggers on push and pull_request events to the main branch.
2. Sets up a job named "build" running on the latest Ubuntu environment.
3. Sets up the Go environment with the specified Go version (1.17 in this case) using the actions/setup-go action.
4. Checks out your code using the actions/checkout action.
5. Downloads the required dependencies using go mod download.
6. Builds your Go app using go build.
7. Runs the tests using go test -v ./....

You can customize this workflow according to your needs, such as adding steps to deploy the built application or adjusting the Go version.

## Note

The `go run` story wasn't accurate, and a follow-up question (not listed) wasn't helpful. In the end, the correct invocation turned out to be `go run .`. Simple enough.

Also, ChatGPT didn't volunteer any caching logic in the GitHub workflow, which it did previously.
