package api

import (
	"database/sql"
	"net/http"

	"github.com/Fozzyack/signalstack-backend/internal/auth"
	"github.com/Fozzyack/signalstack-backend/internal/store"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	logger       *zerolog.Logger
	userStore    store.UserStore
	sessionStore store.SessionStore
}

func NewAuthHandler(logger *zerolog.Logger, userStore store.UserStore, sessionStore store.SessionStore) *AuthHandler {
	return &AuthHandler{
		logger:       logger,
		userStore:    userStore,
		sessionStore: sessionStore,
	}
}

func (ah *AuthHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	SendJSON(w, map[string]string{"stat": "Successful Auth"})
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ah.logger.Info().Msg("User Login Detected")
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req loginRequest
	err := DecodeJSON(r, &req)
	if err != nil {
		ah.logger.Error().Err(err).Msg("Error Decoding JSON")
		ErrorJSON(w, 400, "Invalid JSON")
		return
	}

	user, err := ah.userStore.GetUserByEmail(r.Context(), req.Email)
	if err == sql.ErrNoRows {
		ah.logger.Error().Err(err).Msg("User not Found")
		ErrorJSON(w, 404, "Incorrect Username Or Password")
		return
	}
	if err != nil {
		ah.logger.Error().Err(err).Msg("Internal Server Error in Login")
		ErrorJSON(w, 500, "Internal Server Error on Login")
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash, ah.logger) {
		ah.logger.Error().Msg("Passwords do not match")
		ErrorJSON(w, 404, "Incorrect Username Or Password")
		return
	}

	token := auth.GenerateToken()
	session, err := ah.sessionStore.CreateSession(r.Context(), user.ID, token)
	if err != nil {
		ah.logger.Error().Err(err).Msg("Error Creating Session")
		ErrorJSON(w, 500, "Internal Server Error on Login")
		return
	}

	SendJSON(w, map[string]string{"token": session.Token})

}
