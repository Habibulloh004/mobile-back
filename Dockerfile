FROM golang:1.23.6-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Create uploads directory
RUN mkdir -p /app/uploads/images

# Copy the binary from the builder stage
COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./app"]