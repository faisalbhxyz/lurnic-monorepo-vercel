# Stage 1: Build
FROM node:20-alpine AS builder
WORKDIR /app

# Baked into the client bundle — must be the public API URL including /v1 (see root .env.example).
ARG NEXT_PUBLIC_API_URL
ARG NEXT_PUBLIC_APP_KEY
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
ENV NEXT_PUBLIC_APP_KEY=${NEXT_PUBLIC_APP_KEY}

COPY package*.json ./
RUN npm install --force
COPY . .

# Build standalone
RUN npm run build

# Stage 2: Production
FROM node:20-alpine AS runner
WORKDIR /app

# Copy standalone build
COPY --from=builder /app/.next/standalone ./   
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public

# Copy package.json (useful if you need version info or runtime scripts)
COPY --from=builder /app/package.json ./package.json

EXPOSE 3000
CMD ["node", "server.js"]