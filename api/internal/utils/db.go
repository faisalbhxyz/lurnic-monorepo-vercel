package utils

import (
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
	dsn := os.Getenv("GOOSE_DBSTRING")
	if dsn == "" {
		return &ConfigError{Message: "GOOSE_DBSTRING not found in environment"}
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

	database, err := gorm.Open(mysql.Open(dsn), cfg)
	if err != nil {
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
