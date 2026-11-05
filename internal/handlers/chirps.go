package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"MODULE_PATH/internal/database"
	"MODULE_PATH/internal/respond"
)

const maxChirpBodyLength = 140

var profaneWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

type createChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

// CreateChirp creates a new chirp from a JSON request body.
func (h *Handler) CreateChirp(w http.ResponseWriter, r *http.Request) {
	params := createChirpRequest{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(params.Body) > maxChirpBodyLength {
		respond.Error(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := h.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanChirpBody(params.Body),
		UserID: params.UserID,
	})
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respond.JSON(w, http.StatusCreated, chirp)
}

// GetChirps returns every chirp, ordered by creation time.
func (h *Handler) GetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := h.db.GetChirps(r.Context())
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	// sqlc returns a nil slice when there are no rows; JSON should still be [].
	if chirps == nil {
		chirps = []database.Chirp{}
	}

	respond.JSON(w, http.StatusOK, chirps)
}

// GetChirp returns a single chirp by its ID.
func (h *Handler) GetChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respond.Error(w, http.StatusNotFound, "Chirp not found")
		return
	}

	chirp, err := h.db.GetChirp(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond.Error(w, http.StatusNotFound, "Chirp not found")
			return
		}

		respond.Error(w, http.StatusInternalServerError, "Couldn't get chirp")
		return
	}

	respond.JSON(w, http.StatusOK, chirp)
}

// cleanChirpBody replaces profane words with asterisks.
func cleanChirpBody(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		if _, ok := profaneWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
