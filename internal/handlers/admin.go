package handlers

import (
	"fmt"
	"net/http"

	"MODULE_PATH/internal/respond"
)

// Metrics renders the admin page with the current request count.
func (h *Handler) Metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, h.metrics.Hits())))
}

// Reset clears the metrics counter and deletes all users.
// It is only available in the dev environment.
func (h *Handler) Reset(w http.ResponseWriter, r *http.Request) {
	if h.platform != "dev" {
		respond.Error(w, http.StatusForbidden, "Reset is only allowed in dev environment")
		return
	}

	h.metrics.Reset()
	if err := h.db.DeleteUsers(r.Context()); err != nil {
		respond.Error(w, http.StatusInternalServerError, "Couldn't delete users")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Reset successful"))
}
