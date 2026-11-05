package respond

import (
	"encoding/json"
	"net/http"
)

// JSON writes payload as a JSON response with the given status code.
func JSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(data)
}

// Error writes a standard {"error": "..."} JSON response.
func Error(w http.ResponseWriter, code int, msg string) {
	JSON(w, code, struct {
		Error string `json:"error"`
	}{Error: msg})
}
