package middleware

import (
	"net/http"
	"sync/atomic"
)

// Metrics tracks how many requests have passed through the application.
type Metrics struct {
	hits atomic.Int32
}

// NewMetrics creates an empty request counter.
func NewMetrics() *Metrics {
	return &Metrics{}
}

// Count wraps a handler so every request increments the hit counter.
func (m *Metrics) Count(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.hits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// Hits returns the current number of counted requests.
func (m *Metrics) Hits() int32 {
	return m.hits.Load()
}

// Reset sets the hit counter back to zero.
func (m *Metrics) Reset() {
	m.hits.Store(0)
}
