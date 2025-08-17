# Stage 1: Build the Go binary
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first
# This leverages Docker's layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your application source code
COPY . .

# Build the application. 
# -o /app/main specifies the output path for the binary.
# CGO_ENABLED=0 is important for creating a static binary that doesn't depend on system C libraries.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd

# ---

# Stage 2: Create the final, lightweight image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the 'builder' stage
COPY --from=builder /app/main .

# Copy other necessary assets that your application needs to run
# This includes your config, templates, static files, migrations, and JWT keys
COPY config.yaml .
COPY internal/view ./internal/view
COPY static ./static
COPY migrations ./migrations
COPY .env ./.env



# Expose the port your Go application listens on (e.g., 8080)
# Make sure this matches the port in your config.yaml or code
EXPOSE 8080

# The command to run your application when the container starts
CMD ["./main"]
