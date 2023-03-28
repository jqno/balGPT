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
RUN go build -buildvcs=false -o main

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
