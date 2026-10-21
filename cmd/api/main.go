package main

import (
	"log"
	"net/http"

	"finance-tracker/internal/app"
	"finance-tracker/internal/config"
	"finance-tracker/internal/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dbPool, err := db.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	router := app.NewRouter()

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	log.Printf("server started on :%s", cfg.AppPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
