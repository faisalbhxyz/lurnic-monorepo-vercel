# API only (Gin). For API + Next.js dashboard together, deploy with docker-compose.yaml on Coolify (Docker Compose), not Nixpacks-only.
#
# ---------- BUILD STAGE ----------
FROM golang:1.24-alpine AS builder

# Install git (required for go install) and gcc (optional if CGo used)
RUN apk add --no-cache git

WORKDIR /app

# Copy go module files & download deps (module root is api/)
COPY api/go.mod api/go.sum ./
RUN go mod download

# Copy source code
COPY api/ .

# Install goose CLI (pin — @latest may require a newer Go than the builder image)
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.26.0

# Build the Go app
RUN go build -o main .

# ---------- FINAL STAGE ----------
FROM alpine:3.20

# tzdata (IANA zones) + wget (docker-compose healthcheck hits /health)
RUN apk add --no-cache tzdata wget
ENV TZ=UTC

# Working directory
WORKDIR /app

# Copy built app and goose from builder
COPY --from=builder /app/main .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Migration SQL (module lives under api/; do not COPY from build context — root has no migrations/)
COPY --from=builder /app/migrations ./migrations

# Expose the port your Gin app listens on (must match APP_PORT, default 5000)
EXPOSE 5000

# Run DB migrations first, then start app (explicit -dir: cwd-independent; exec for clean signals)
CMD ["sh", "-c", "goose -dir /app/migrations up && exec ./main"]
