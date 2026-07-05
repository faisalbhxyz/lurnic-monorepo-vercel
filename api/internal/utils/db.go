package utils

import (
	"context"
	"dashlearn/internal/dsn"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() error {
	raw := os.Getenv("GOOSE_DBSTRING")
	if raw == "" {
		return &ConfigError{Message: "GOOSE_DBSTRING not found in environment"}
	}
	dsnStr, err := dsn.Normalize(raw)
	if err != nil {
		return err
	}

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

	database, err := gorm.Open(mysql.Open(dsnStr), cfg)
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

type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return e.Message
}
