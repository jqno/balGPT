Let's add a GitHub actions workflow!

Sure, I can help you set up a GitHub Actions workflow to build, test, and deploy your application. Create a new file in your repository under .github/workflows/main.yml and add the following content:

```yaml
name: CI/CD Pipeline

on:
  push:
    branches:
      - main

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v2

      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.17

      - name: Cache Go modules
        uses: actions/cache@v2
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Install dependencies
        run: |
          cd server
          go mod download

      - name: Build and Test
        run: |
          cd server
          go build -v
          go test -v
```

This GitHub Actions workflow is configured to:

1. Trigger on push events to the main branch.
2. Build and test your Go application using the latest version of Go (1.17).
3. Cache Go modules for faster subsequent builds.

Make sure to replace cd server with the appropriate path to your Go application's directory if it's different.

Once you've committed the workflow file, GitHub Actions will automatically run the workflow on every push to the main branch, building, and testing your application.

## Note

- The current Go version used in this project was actually 1.20.
- The commands to build and test were incorrect.
- Actually, we don't even have tests yet at this point, but that's because I decided not to use them -- see [prompt 11](prompts/011-add-tests.md).
