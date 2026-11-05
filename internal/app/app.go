package app

import (
	"database/sql"
	"net/http"

	"MODULE_PATH/internal/config"
	"MODULE_PATH/internal/database"
	"MODULE_PATH/internal/handlers"
	"MODULE_PATH/internal/middleware"
	"MODULE_PATH/internal/router"
)

// New wires configuration, database queries, middleware, handlers, and routes
// together and returns the ready-to-serve HTTP handler.
func New(cfg config.Config, db *sql.DB) http.Handler {
	queries := database.New(db)
	metrics := middleware.NewMetrics()
	handler := handlers.New(queries, cfg.Platform, metrics)

	return router.New(handler, metrics)
}
