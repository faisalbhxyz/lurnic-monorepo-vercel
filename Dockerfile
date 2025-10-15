# ---------- BUILD STAGE ----------
FROM golang:1.24-alpine AS builder

# Install git (required for go install) and gcc (optional if CGo used)
RUN apk add --no-cache git

WORKDIR /app

# Copy go module files & download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install goose CLI
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Build the Go app
RUN go build -o main .

# ---------- FINAL STAGE ----------
FROM alpine:3.20

# Working directory
WORKDIR /app

# Copy built app and goose from builder
COPY --from=builder /app/main .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy migration files
COPY ./migrations ./migrations

# Expose the port your Gin app listens on
EXPOSE 5001

# Run DB migrations first, then start app
CMD goose up && ./main
