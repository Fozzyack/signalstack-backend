package api

import (
	"context"
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

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ah.logger.Info().Msg("User Login Detected")
	type loginRequest struct {
		email    string
		password string
	}
	var req loginRequest
	DecodeJSON(r, req)

	user, err := ah.userStore.GetUserByEmail(context.Background(), req.email)
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

	if !auth.CheckPassword(req.password, user.PasswordHash) {
		ah.logger.Error().Err(err).Msg("Passwords do not match")
		ErrorJSON(w, 404, "Incorrect Username Or Password")
		return
	}

	token := auth.GenerateToken()
	session, err := ah.sessionStore.CreateSession(context.Background(), user.ID, token)
	if err != nil {
		ah.logger.Error().Err(err).Msg("Error Creating Session")
		ErrorJSON(w, 500, "Internal Server Error on Login")
		return
	}

	SendJSON(w, map[string]string{"token": session.Token})

}
