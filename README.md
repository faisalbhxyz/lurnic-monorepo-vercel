# Dashlearn Server (Lurnic) — Local Setup

## Prerequisites
- Go 1.24+
- MySQL 8+ (or compatible)
- `goose` installed:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Run locally with a new database

### 1) Create a new MySQL database

```sql
CREATE DATABASE dashlearn_local CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2) Configure environment

```bash
cp .env.example .env
```

Then ensure `GOOSE_DBSTRING` points to your DB, e.g.:

```env
GOOSE_DBSTRING=root:@tcp(127.0.0.1:3306)/dashlearn_local?charset=utf8mb4&parseTime=True&loc=Local
```

### (Optional) Enable Sentry error monitoring

This server supports Sentry via the Gin SDK middleware. Add `SENTRY_DSN` to your `.env` to enable it:

```env
SENTRY_DSN=https://<key>@o<org>.ingest.us.sentry.io/<project>
SENTRY_ENVIRONMENT=development
SENTRY_TRACES_SAMPLE_RATE=0
```

### 3) Run migrations (create schema)

```bash
goose -dir api/migrations mysql "$GOOSE_DBSTRING" up
```

### 4) Start the API server

```bash
cd api && go run .
```

Server runs at `http://localhost:${APP_PORT}` (default `5000`).

## Multi-tenancy note
All API requests must include the `app-key` header. The server resolves it to a tenant from the `tenants` table.
