package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"MODULE_PATH/internal/app"
	"MODULE_PATH/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: app.New(cfg, db),
	}

	log.Printf("Chirpy listening on %s", cfg.Address)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server stopped: %v", err)
	}
}
