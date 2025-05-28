# syntax=docker/dockerfile:1
FROM golang:1.24.1-alpine AS builder

# Set working directory
WORKDIR /app

# Install necessary packages
RUN apk add --no-cache git

# Copy go mod and sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the source
COPY . .

# Build the app
RUN go build -o main .

# Command to run
CMD ["/app/main"]
