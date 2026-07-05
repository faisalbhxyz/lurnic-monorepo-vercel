#!/bin/sh
set -e

log() {
  printf '[api-entrypoint] %s\n' "$1"
}

if [ -z "${GOOSE_DBSTRING:-}" ]; then
  log 'ERROR: GOOSE_DBSTRING is not set'
  exit 1
fi

# Match api/internal/utils/db.go: ensure TLS + timeouts for managed TiDB/MySQL from Docker.
case "$GOOSE_DBSTRING" in
  *\?*)
    case "$GOOSE_DBSTRING" in
      *tls=*) ;;
      *) GOOSE_DBSTRING="${GOOSE_DBSTRING}&tls=skip-verify" ;;
    esac
    case "$GOOSE_DBSTRING" in
      *timeout=*) ;;
      *) GOOSE_DBSTRING="${GOOSE_DBSTRING}&timeout=5s" ;;
    esac
    case "$GOOSE_DBSTRING" in
      *readTimeout=*) ;;
      *) GOOSE_DBSTRING="${GOOSE_DBSTRING}&readTimeout=15s" ;;
    esac
    case "$GOOSE_DBSTRING" in
      *writeTimeout=*) ;;
      *) GOOSE_DBSTRING="${GOOSE_DBSTRING}&writeTimeout=15s" ;;
    esac
    ;;
  *)
    GOOSE_DBSTRING="${GOOSE_DBSTRING}?tls=skip-verify&timeout=5s&readTimeout=15s&writeTimeout=15s"
    ;;
esac
export GOOSE_DBSTRING

log 'running fixautoinc (non-fatal on failure)...'
if ! ./fixautoinc; then
  log 'fixautoinc failed or skipped; continuing'
fi

log 'running goose migrations...'
if ! goose -dir /app/migrations up; then
  log 'ERROR: goose migration failed'
  exit 1
fi

log 'migrations complete; starting API'
exec ./main
