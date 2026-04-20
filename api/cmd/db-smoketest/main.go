package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("GOOSE_DBSTRING")
	if dsn == "" {
		log.Fatal("GOOSE_DBSTRING not set")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("sql open: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	var one int
	if err := db.QueryRow("SELECT 1").Scan(&one); err != nil {
		log.Fatalf("select 1: %v", err)
	}

	fmt.Println("db ok:", one)
}

