package handlers

import (
	"MODULE_PATH/internal/database"
	"MODULE_PATH/internal/middleware"
)

// Handler holds the dependencies shared by all HTTP handlers.
type Handler struct {
	db        *database.Queries
	platform  string
	jwtSecret string
	polkaKey  string
	metrics   *middleware.Metrics
}

// New creates a Handler with the dependencies it needs.
func New(db *database.Queries, platform, jwtSecret, polkaKey string, metrics *middleware.Metrics) *Handler {
	return &Handler{
		db:        db,
		platform:  platform,
		jwtSecret: jwtSecret,
		polkaKey:  polkaKey,
		metrics:   metrics,
	}
}
