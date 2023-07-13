#!/bin/bash

# Run the tests with coverage
echo "Running tests..."
go test -coverprofile=coverage.out ./...

# Generate the HTML coverage report
echo "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo "Coverage report generated as coverage.html"

