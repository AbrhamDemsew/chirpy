package router

import (
	"net/http"

	"MODULE_PATH/internal/handlers"
	"MODULE_PATH/internal/middleware"
)

// New builds the HTTP router. Think of it as the application's route table:
// handlers contain the behavior, while this file decides which handler
// belongs to which method and path.
func New(h *handlers.Handler, metrics *middleware.Metrics) http.Handler {
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/app/", metrics.Count(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Health
	mux.HandleFunc("GET /api/healthz", h.Health)

	// Users
	mux.HandleFunc("POST /api/users", h.CreateUser)
	mux.HandleFunc("PUT /api/users", h.UpdateUser)
	mux.HandleFunc("POST /api/login", h.Login)
	mux.HandleFunc("POST /api/refresh", h.Refresh)
	mux.HandleFunc("POST /api/revoke", h.Revoke)

	// Chirps
	mux.HandleFunc("GET /api/chirps", h.GetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", h.GetChirp)
	mux.HandleFunc("POST /api/chirps", h.CreateChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", h.DeleteChirp)

	// Admin
	mux.HandleFunc("GET /admin/metrics", h.Metrics)
	mux.HandleFunc("POST /admin/reset", h.Reset)

	return mux
}
