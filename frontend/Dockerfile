# Build context must be the monorepo ROOT (where `frontend/` lives). See docker-compose.yaml `web.build`.
# Do not use `docker build` with context `frontend/` unless you adjust COPY paths below.
#
# Stage 1: Build
FROM node:20-alpine AS builder
WORKDIR /app

# Baked into the client bundle — must be the public API URL including /v1 (see root .env.example).
ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_APP_KEY
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
ENV NEXT_PUBLIC_APP_KEY=${NEXT_PUBLIC_APP_KEY}

COPY frontend/package*.json ./
RUN npm install --force
COPY frontend/ .

# Coolify/small build hosts: Next "Collecting build traces" + standalone can OOM without a higher heap.
ENV NODE_OPTIONS=--max-old-space-size=6144

# Build standalone
RUN npm run build

# Stage 2: Production
FROM node:20-alpine AS runner
WORKDIR /app

# Must match Coolify “Ports” / Traefik upstream (default 3000). Bind all interfaces so the proxy can reach the process.
ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME=0.0.0.0

# Standalone output nests server.js under the app dir name (`app/` here because
# WORKDIR is /app and outputFileTracingRoot is the filesystem root). Do not copy
# the whole standalone tree or CMD looks for /app/server.js and crashes.
COPY --from=builder /app/.next/standalone/app ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

# Copy package.json (useful if you need version info or runtime scripts)
COPY --from=builder /app/package.json ./package.json

EXPOSE 3000
CMD ["node", "server.js"]