# Coolify deployment (API + Next.js dashboard)

This matches what you get on **Vercel** (Next UI + Gin API under `/v1`), but on Coolify you run **two containers** behind **two domains** (or one host with path routing). Vercel uses rewrites (`vercel.json`); Coolify uses **Docker Compose** (`docker-compose.yaml`).

## Vercel vs Coolify (mental model)

| | Vercel | Coolify (this repo) |
|---|--------|----------------------|
| Frontend | Next.js build | `web` service → `frontend/Dockerfile` (standalone `node server.js`) |
| API | Serverless `/api/vercel/gin-handler` → Gin | `api` service → root `Dockerfile` → `./main` on port **5000** |
| Browser → API | Often same origin; `/v1/*` rewritten to Gin | **Public URL** must match `NEXT_PUBLIC_API_URL` (usually `https://api.<domain>/v1`) |
| Cron | `vercel.json` `crons` | Schedule HTTP GET to **`/v1/internal/cron/publish-scheduled`** with `Authorization: Bearer <CRON_SECRET>` (see below) |
| DB | External (you configure) | Bundled **`mysql`** in compose, or external DSN via `GOOSE_DBSTRING` |

## 1. Create the resource

1. Coolify → **New resource** → **Docker Compose** (not “Application” / Nixpacks-only Go — that ships **API without the dashboard**).
2. Connect the **Git** repo (full monorepo: must include `frontend/` at repo root).
3. Compose file: **`docker-compose.yaml`** (repo root).
4. **Build**: Coolify should run `docker compose build` from the compose file directory. If build fails with `lstat .../frontend`, set env **`BUILD_SOURCE_CONTEXT`** to your Git URL, e.g. `https://github.com/ORG/REPO.git#main` (see `.env.example`).

## 2. Ports & domains (Traefik)

Create **two** FQDNs (recommended):

| Service | Container port | Use for |
|---------|----------------|--------|
| **`api`** | **5000** | `https://api.example.com` — this origin + `/v1` goes into **`NEXT_PUBLIC_API_URL`** |
| **`web`** | **3000** | `https://app.example.com` — this is **`AUTH_URL`** (what users open in the browser) |

Attach each domain to the **correct** service in Coolify. Wrong mapping → **502** or login/API failures.

## 3. Required environment variables (Coolify project env)

Set these in the Coolify UI (shared across services; compose substitutes into builds and containers). Align with `.env.example`.

**Secrets (generate strong values in production)**

- **`JWT_SECRET`** — API JWT signing.
- **`AUTH_SECRET`** — NextAuth / session encryption.
- **`CRON_SECRET`** — optional but required for cron route auth if you use scheduled publishes.

**URLs (must be consistent with your domains)**

- **`NEXT_PUBLIC_API_URL`** — **Build-time** for `web`. Must end with **`/v1`**, e.g. `https://api.example.com/v1`.  
  After any change → **rebuild the `web` service** (not just restart).
- **`AUTH_URL`** — Public dashboard URL, e.g. `https://app.example.com` (scheme + host, no trailing path).

**Database**

- With **bundled MySQL**: set **`MYSQL_*`** and **`GOOSE_DBSTRING`** so credentials match (hostname **`mysql`** inside the stack).
- With **external MySQL**: remove the `mysql` service from compose (or use a prod override), set **`GOOSE_DBSTRING`** to the provider DSN reachable from **`api`**.

**Optional**

- **`NEXT_PUBLIC_APP_KEY`** — if your frontend expects it.
- **`GIN_MODE`** — `release` in production.
- **`BUILD_SOURCE_CONTEXT`** — Git URL if the build context is not a full clone.

## 4. Scheduled lessons (replace Vercel Cron)

On Vercel, cron hits `/api/vercel/gin-handler?path=internal/cron/publish-scheduled`. In Docker, Gin exposes:

`GET /v1/internal/cron/publish-scheduled`

Configure **Coolify “Scheduled Tasks”**, a system cron, or Uptime Kuma to request:

```http
GET https://api.example.com/v1/internal/cron/publish-scheduled
Authorization: Bearer <CRON_SECRET>
```

`CRON_SECRET` must match what you set in Coolify env for **`api`**.

## 5. Deploy / rebuild

- First deploy or code changes: **Deploy** with **build** so `api` and `web` images rebuild.
- After changing **`NEXT_PUBLIC_*`**: rebuild **`web`** so the client bundle embeds the new API URL.
- If the dashboard loads **without CSS** (plain HTML), see `.cursorrules`: `/_next/static/*` **404** usually means **stale/partial deploy** — full rebuild **`web`**, confirm the domain points to **`web`:3000**.

## 6. Node build memory (already in `frontend/Dockerfile`)

The builder sets `NODE_OPTIONS=--max-old-space-size=6144` to reduce OOM during Next “Collecting build traces”. If builds still fail, raise RAM on the build host or bump that value.

## 7. Local parity

```bash
cp .env.example .env
# Set JWT_SECRET, AUTH_SECRET, NEXT_PUBLIC_API_URL, AUTH_URL for your machine
docker compose up --build -d
```

API: `http://localhost:5000` · App: `http://localhost:3000`
