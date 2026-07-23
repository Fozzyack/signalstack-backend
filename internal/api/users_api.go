package api

import (
	"net/http"

	"github.com/Fozzyack/signalstack-backend/internal/auth"
	"github.com/Fozzyack/signalstack-backend/internal/store"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	logger    *zerolog.Logger
	userStore store.UserStore
}

func NewUserHandler(logger *zerolog.Logger, userStore store.UserStore) *UserHandler {
	return &UserHandler{
		logger:    logger,
		userStore: userStore,
	}
}

func (uh *UserHandler) GetLoggedInUser(w http.ResponseWriter, r *http.Request) {

	type userResponse struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	user := GetUserFromContext(r.Context())

	responsePayload := userResponse{
		Name:  user.Name,
		Email: user.Email,
	}

	SendJSON(w, responsePayload)
}

func (uh *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	type updateUserRequest struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		CurrentPassword string `json:"currentPassword"`
	}

	user := GetUserFromContext(r.Context())

	var req updateUserRequest
	err := DecodeJSON(r, &req)
	if err != nil {
		uh.logger.Error().Err(err).Msg("Error Decoding JSON")
		ErrorJSON(w, 400, "Invalid JSON")
		return
	}

	if !auth.CheckPassword(req.CurrentPassword, user.PasswordHash, uh.logger) {
		uh.logger.Error().Msg("Passwords do not match")
		ErrorJSON(w, 404, "Incorrect Password")
		return
	}

	passwordHash := user.PasswordHash
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			uh.logger.Error().Err(err).Msg("Could not hash password")
			ErrorJSON(w, 500, "Internal Server Error")
			return
		}
		passwordHash = string(hashedPassword)
	}

	updatedUser, err := uh.userStore.UpdateUser(
		r.Context(),
		user.ID,
		req.Name,
		req.Email,
		passwordHash,
	)
	if err != nil {
		uh.logger.Error().Err(err).Msg("Could not update user")
		ErrorJSON(w, 500, "Internal Server Error")
		return
	}

	SendJSON(w, updatedUser)
}
