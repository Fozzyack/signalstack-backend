package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/Fozzyack/signalstack-backend/internal/store"
	"github.com/rs/zerolog"
)

type RequestHandler struct {
	logger       *zerolog.Logger
	requestStore store.RequestStore
}

func NewRequestHandler(logger *zerolog.Logger, requestStore store.RequestStore) *RequestHandler {
	return &RequestHandler{
		logger:       logger,
		requestStore: requestStore,
	}
}

func (rh *RequestHandler) GetAllRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := rh.requestStore.GetRequests(context.Background())
	if err != sql.ErrNoRows {
		rh.logger.Error().Err(err).Msg("Could not get Requests")
		ErrorJSON(w, 404, "No Requests Found")
		return
	}
	if err != nil {
		rh.logger.Error().Err(err).Msg("Could not get Requests")
		ErrorJSON(w, 500, "Internal Server Errror")
		return
	}

	SendJSON(w, requests)
}
