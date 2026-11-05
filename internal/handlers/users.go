package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"MODULE_PATH/internal/auth"
	"MODULE_PATH/internal/database"
	"MODULE_PATH/internal/respond"
)

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// userResponse is the public representation of a user. It deliberately omits
// the password hash so that credentials are never sent to clients.
type userResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func newUserResponse(user database.User) userResponse {
	return userResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
}

// CreateUser creates a new user from a JSON request body.
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	params := createUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if params.Email == "" {
		respond.Error(w, http.StatusBadRequest, "Email is required")
		return
	}

	if params.Password == "" {
		respond.Error(w, http.StatusBadRequest, "Password is required")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	user, err := h.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respond.JSON(w, http.StatusCreated, newUserResponse(user))
}

// Login validates a user's email and password.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	params := loginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respond.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respond.Error(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	respond.JSON(w, http.StatusOK, newUserResponse(user))
}
