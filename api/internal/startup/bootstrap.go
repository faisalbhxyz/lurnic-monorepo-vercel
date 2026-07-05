package startup

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"

	"dashlearn/internal/dsn"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
)

// Bootstrap normalizes GOOSE_DBSTRING, repairs TiDB AUTO_INCREMENT quirks, and runs SQL migrations.
func Bootstrap() error {
	norm, err := dsn.Normalize(os.Getenv("GOOSE_DBSTRING"))
	if err != nil {
		return err
	}
	os.Setenv("GOOSE_DBSTRING", norm)

	if err := runFixAutoInc(); err != nil {
		log.Printf("[startup] fixautoinc warning: %v", err)
	}

	return runMigrations(norm)
}

func runFixAutoInc() error {
	if _, err := os.Stat("./fixautoinc"); err != nil {
		return nil
	}
	cmd := exec.Command("./fixautoinc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	return cmd.Run()
}

func runMigrations(dsnStr string) error {
	dir := migrationDir()
	log.Printf("[startup] running goose migrations from %s", dir)

	db, err := sql.Open("mysql", dsnStr)
	if err != nil {
		return fmt.Errorf("open migration DB: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("migration DB ping: %w", err)
	}

	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}

	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	log.Printf("[startup] migrations complete")
	return nil
}

func migrationDir() string {
	if dir := os.Getenv("GOOSE_MIGRATION_DIR"); dir != "" {
		return dir
	}
	if _, err := os.Stat("/app/migrations"); err == nil {
		return "/app/migrations"
	}
	return "migrations"
}
