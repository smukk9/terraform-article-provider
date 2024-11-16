
# Start with a small base image
FROM golang:1.22.4-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire Go application code to the container
COPY . .

# Build the Go app
RUN go build -o go-http-server

# Use a minimal image for the final build
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/go-http-server /go-http-server

# Expose the application port
EXPOSE 9999

# Run the binary
CMD ["/go-http-server"]