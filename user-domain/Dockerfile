# Stage 1: Build binary
FROM golang:1.24.3 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go binary
RUN go build -o server ./cmd/api

# Stage 2: Run the binary in a lightweight image
FROM debian:bookworm-slim

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .

# Expose port (for example: 8080)
EXPOSE 8080

# Run the application
CMD ["./server"]
