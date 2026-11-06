package handlers

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"MODULE_PATH/internal/auth"
	"MODULE_PATH/internal/respond"
)

type polkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

// PolkaWebhook handles subscription events sent by Polka. Only user.upgraded
// events are relevant to Chirpy; every other event is acknowledged and ignored.
func (h *Handler) PolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || subtle.ConstantTimeCompare([]byte(apiKey), []byte(h.polkaKey)) != 1 {
		respond.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	params := polkaWebhookRequest{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respond.Error(w, http.StatusNotFound, "User not found")
		return
	}

	_, err = h.db.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond.Error(w, http.StatusNotFound, "User not found")
			return
		}

		respond.Error(w, http.StatusInternalServerError, "Couldn't upgrade user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
