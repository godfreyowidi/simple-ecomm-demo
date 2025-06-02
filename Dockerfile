# Build stage
FROM golang:1.24 AS builder


WORKDIR /app

# Copy go.mod and go.sum first for caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o savanna-app ./main.go


# Final stage: minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/savanna-app .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./savanna-app"]
