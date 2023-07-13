# Code coverage reporting

> Is it possible to measure code coverage in a Go project?

Yes, it is possible to measure code coverage in a Go project. The Go language has built-in support for code coverage analysis with the `go test` command.

Here's a basic walkthrough of how you might use it:

1. Run Tests with Coverage:bashCopy codego test -coverprofile=coverage.out ./...
   This command will run all tests in your Go project (the ./... specifies to look in all directories under the current one) and output a coverage profile to coverage.out.
2. Analyze Coverage:
   After running the tests, you can use the go tool cover command to analyze the output:bashCopy codego tool cover -func=coverage.out
   This will output a summary of your code's coverage, broken down by function.
3. Generate HTML Report:
   Additionally, you can generate an HTML report that shows line-by-line coverage:bashCopy codego tool cover -html=coverage.out -o coverage.html
   This will create an HTML file named coverage.html in your current directory. You can open this file in your web browser to view the coverage report.

Remember, code coverage is a useful metric to help ensure your tests are comprehensive, but it doesn't guarantee the quality of your tests. High coverage with poorly written tests could give you a false sense of security. Therefore, focus on writing meaningful tests and use coverage as one of several metrics to evaluate your testing effectiveness.

## Follow-up prompt

> Is it possible to define a coverage threshold and make it fail a build if the threshold isn't met?

The `go test` command does not natively support setting a coverage threshold that fails the build if not met. However, there are various ways to implement this functionality in your CI/CD pipeline or testing script by parsing the output of `go test`.

Here's a simple example of a bash script that would fail if code coverage is below a certain percentage:

```bash
#!/bin/bash

# Set the coverage threshold
COVERAGE_THRESHOLD=80

# Run the tests with coverage
go test -coverprofile=coverage.out ./...

# Parse the coverage percentage
COVERAGE_PERCENT=$(go tool cover -func=coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')

# Check if the coverage is below the threshold
if (( $(echo "$COVERAGE_PERCENT < $COVERAGE_THRESHOLD" | bc -l) ))
then
    echo "Code coverage of $COVERAGE_PERCENT% is below the required threshold of $COVERAGE_THRESHOLD%"
    exit 1
else
    echo "Code coverage of $COVERAGE_PERCENT% is above the required threshold of $COVERAGE_THRESHOLD%"
    exit 0
fi

```

The above script first runs the tests with coverage, then parses the total coverage percentage from the output of `go tool cover -func=coverage.out`. It then checks if this coverage percentage is below the threshold, and if so, it exits with a status of 1, indicating an error, which would typically cause a build to fail in most CI/CD environments.

Remember to make the bash script executable with `chmod +x your-script.sh`.

Note: This script relies on the `bc` command-line calculator. If `bc` isn't available in your environment, you'll need a different way to compare the floating-point numbers.

## Follow-up prompt

> OK, maybe that's not the right approach. Can I produce the coverage report on GitHub Actions, and make it accessible from there? I already have a workflow:
>
> _(snipped: .github/workflows/main.yml)_

To generate a code coverage report and make it accessible in GitHub Actions, you can add a few more steps to your workflow.

1. Generate Coverage Profile: You need to run your tests with the -coverprofile flag, which will output a coverage profile file.
2. Generate Coverage Report: Convert the coverage profile into a human-readable format. You can use the built-in go tool cover to convert it into an HTML file.
3. Upload Coverage Report as Artifact: Use the actions/upload-artifact action to upload the generated report as an artifact of your workflow. This will make the report accessible from the GitHub Actions UI.

Here's how to modify your workflow to accomplish this:

```yaml
name: CI/CD Pipeline

on:
  push:

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.20.2

      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: Install dependencies
        run: |
          go mod download

      - name: Build
        run: |
          go build

      - name: Test and Generate Coverage Profile
        run: |
          go test -coverprofile=coverage.out ./...

      - name: Generate Coverage Report
        run: |
          go tool cover -html=coverage.out -o coverage.html

      - name: Upload Coverage Report
        uses: actions/upload-artifact@v3
        with:
          name: coverage-report
          path: coverage.html
```

Now, after each run of your workflow, you can go to the GitHub Actions page for your repository, click on the specific run, and under "Artifacts" you'll find the coverage report available for download. Please note that this report will be kept for 90 days by default.

## Follow-up prompt

> Can you add a little shell script to run tests and generate the report locally? It doesn't have to open it in a browser, just create the html.

Sure, here is a simple bash script to run the tests and generate a coverage report:

```bash
#!/bin/bash

# Run the tests with coverage
echo "Running tests..."
go test -coverprofile=coverage.out ./...

# Generate the HTML coverage report
echo "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo "Coverage report generated as coverage.html"

```

You can save this script to a file, make it executable, and then run it to execute your tests and generate the report.

For example, you might save it to a file named `coverage.sh` in the root directory of your Go project. Then, make it executable:

```bash
chmod +x coverage.sh

```

And to run it:

```bash
./coverage.sh

```

This will generate a `coverage.html` file in your project's root directory, which you can view in your web browser to see the coverage report.
