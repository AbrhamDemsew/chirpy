package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"

	"MODULE_PATH/internal/auth"
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
	Body string `json:"body"`
}

// CreateChirp creates a new chirp for the user identified by the request's JWT.
func (h *Handler) CreateChirp(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(tokenString, h.jwtSecret)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

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
		UserID: userID,
	})
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respond.JSON(w, http.StatusCreated, chirp)
}

// GetChirps returns chirps ordered by creation time. When the author_id query
// parameter is provided, only chirps from that author are returned. The sort
// query parameter controls the order: "asc" (default) or "desc".
func (h *Handler) GetChirps(w http.ResponseWriter, r *http.Request) {
	var (
		chirps []database.Chirp
		err    error
	)

	authorID := r.URL.Query().Get("author_id")
	if authorID == "" {
		chirps, err = h.db.GetChirps(r.Context())
	} else {
		userID, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			respond.Error(w, http.StatusBadRequest, "Invalid author ID")
			return
		}

		chirps, err = h.db.GetChirpsByAuthor(r.Context(), userID)
	}

	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	// sqlc returns a nil slice when there are no rows; JSON should still be [].
	if chirps == nil {
		chirps = []database.Chirp{}
	}

	descending := r.URL.Query().Get("sort") == "desc"
	sort.Slice(chirps, func(i, j int) bool {
		if descending {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})

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

// DeleteChirp deletes a chirp owned by the authenticated user.
func (h *Handler) DeleteChirp(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := auth.ValidateJWT(tokenString, h.jwtSecret)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

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

	if chirp.UserID != userID {
		respond.Error(w, http.StatusForbidden, "You can only delete your own chirps")
		return
	}

	if err := h.db.DeleteChirp(r.Context(), id); err != nil {
		respond.Error(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
