package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"MODULE_PATH/internal/auth"
	"MODULE_PATH/internal/respond"
)

type refreshResponse struct {
	Token string `json:"token"`
}

// Refresh exchanges a valid refresh token for a new access token.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond.Error(w, http.StatusUnauthorized, "Invalid refresh token")
			return
		}

		respond.Error(w, http.StatusInternalServerError, "Couldn't get user from refresh token")
		return
	}

	token, err := auth.MakeJWT(user.ID, h.jwtSecret, accessTokenExpiration)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "Couldn't create token")
		return
	}

	respond.JSON(w, http.StatusOK, refreshResponse{Token: token})
}

// Revoke marks the refresh token from the request header as revoked.
func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.db.RevokeRefreshToken(r.Context(), refreshToken); err != nil {
		respond.Error(w, http.StatusInternalServerError, "Couldn't revoke refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
