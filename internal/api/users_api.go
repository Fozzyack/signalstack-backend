package api

import (
	"net/http"

	"github.com/Fozzyack/signalstack-backend/internal/store"
	"github.com/rs/zerolog"
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
