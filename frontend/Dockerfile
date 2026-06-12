# Build context must be the monorepo ROOT (where `frontend/` lives). See docker-compose.yaml `web.build`.
# Keep the repo layout (`/app/frontend`) so next.config outputFileTracingRoot matches local builds.
# Flattening frontend into /app makes tracingRoot "/" and breaks standalone React at runtime.
#
# Stage 1: Build
FROM node:20-alpine AS builder
WORKDIR /app

ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_APP_KEY
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
ENV NEXT_PUBLIC_APP_KEY=${NEXT_PUBLIC_APP_KEY}

COPY package.json package-lock.json ./
COPY frontend/package*.json ./frontend/
RUN npm install --prefix frontend --force
COPY frontend/ ./frontend/

ENV NODE_OPTIONS=--max-old-space-size=6144
WORKDIR /app/frontend
RUN npm run build

# Stage 2: Production
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME=0.0.0.0

COPY --from=builder /app/frontend/.next/standalone/frontend ./
COPY --from=builder /app/frontend/.next/static ./.next/static
COPY --from=builder /app/frontend/public ./public

EXPOSE 3000
CMD ["node", "server.js"]
