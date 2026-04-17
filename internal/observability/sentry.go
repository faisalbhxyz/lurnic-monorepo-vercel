package observability

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
)

type SentryConfig struct {
	DSN                string
	Environment        string
	Release            string
	TracesSampleRate   float64
}

func InitSentry(cfg SentryConfig) (flush func(time.Duration) bool, enabled bool) {
	if cfg.DSN == "" {
		return sentry.Flush, false
	}

	opts := sentry.ClientOptions{
		Dsn:              cfg.DSN,
		Environment:      cfg.Environment,
		Release:          cfg.Release,
		TracesSampleRate: cfg.TracesSampleRate,
		AttachStacktrace: true,
	}

	if err := sentry.Init(opts); err != nil {
		log.Printf("Sentry init failed: %v", err)
		return sentry.Flush, false
	}

	return sentry.Flush, true
}

func EnvSentryConfig(release string) SentryConfig {
	return SentryConfig{
		DSN:                os.Getenv("SENTRY_DSN"),
		Environment:        envOr("SENTRY_ENVIRONMENT", os.Getenv("GIN_MODE")),
		Release:            release,
		TracesSampleRate:   parseFloatOr("SENTRY_TRACES_SAMPLE_RATE", 0),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseFloatOr(key string, fallback float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fallback
	}
	if f < 0 {
		return 0
	}
	if f > 1 {
		return 1
	}
	return f
}

