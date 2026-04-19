package main

import (
	"dashlearn/pkg/server"
	"fmt"
	"log"
	"os"
	"time"
	_ "time/tzdata" // embed IANA tz DB so time.LoadLocation works without OS zoneinfo (Alpine/scratch/minimal hosts)

	"github.com/getsentry/sentry-go"
)

var Version = "v1.0.24"

func main() {
	fmt.Println("🚀 DashLearn Server Starting... Version:", Version)

	// Optional .env is loaded inside server.NewEngine (also used by Vercel handler); Coolify/Docker use injected env only.

	debugRoutesEnabled := os.Getenv("ENABLE_DEBUG_ROUTES") == "true"

	router, flush, err := server.NewEngine(Version)
	if err != nil {
		if debugRoutesEnabled {
			log.Fatalf("Failed to start server: %v", err)
		}
		log.Fatal("Failed to start server:", err)
	}

	defer flush(2 * time.Second)
	defer func() {
		if r := recover(); r != nil {
			sentry.CurrentHub().Recover(r)
			flush(2 * time.Second)
			panic(r)
		}
	}()

	router.Run(":" + os.Getenv("APP_PORT"))
}
