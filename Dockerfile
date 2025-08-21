# Start from the official Golang image
FROM golang:1.24.6-alpine AS builder

WORKDIR /app

# Install git (for go mod download)
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN go build -o fiber-api main.go

# Use a minimal image for running
FROM alpine:latest
WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the built binary and .env file
COPY --from=builder /app/fiber-api .
# COPY .env .

EXPOSE 8080

CMD ["./fiber-api"]
