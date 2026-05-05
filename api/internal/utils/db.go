package utils

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() error {
	dsn := os.Getenv("GOOSE_DBSTRING")
	if dsn == "" {
		return &ConfigError{Message: "GOOSE_DBSTRING not found in environment"}
	}
	dsn = normalizeMySQLDSN(dsn)

	// GORM default SlowThreshold is 200ms; cross-region DB RTT often exceeds that without being a slow query.
	slowMS := 1000
	if v := os.Getenv("GORM_SLOW_SQL_MS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			slowMS = n
		}
	}

	cfg := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Duration(slowMS) * time.Millisecond,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	}

	database, err := gorm.Open(mysql.Open(dsn), cfg)
	if err != nil {
		return err
	}

	sqlDB, err := database.DB()
	if err != nil {
		return err
	}
	// Conservative defaults for managed/cross-region MySQL/TiDB.
	// Keep pool small to avoid exhausting remote limits and to reduce long-hanging sockets.
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}

	DB = database
	return nil
}

func normalizeMySQLDSN(dsn string) string {
	// NOTE: We intentionally avoid parsing user/pass to keep this transformation simple and safe.
	// We only touch the query string segment.
	base, q := splitQuery(dsn)
	q = ensureQueryKV(q, "tls", "skip-verify")
	q = ensureQueryKV(q, "timeout", "5s")
	q = ensureQueryKV(q, "readTimeout", "15s")
	q = ensureQueryKV(q, "writeTimeout", "15s")
	if q == "" {
		return base
	}
	return base + "?" + q
}

func splitQuery(s string) (base, query string) {
	if i := strings.Index(s, "?"); i >= 0 {
		return s[:i], s[i+1:]
	}
	return s, ""
}

func ensureQueryKV(query, key, value string) string {
	// Treat query as simple k=v&... string; we don't decode/encode because we only add safe ASCII keys/values.
	if query == "" {
		return key + "=" + value
	}
	parts := strings.Split(query, "&")
	for _, p := range parts {
		if p == key || strings.HasPrefix(p, key+"=") {
			return query
		}
	}
	return query + "&" + key + "=" + value
}

type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}
