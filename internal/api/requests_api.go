package api

import (
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
	requests, err := rh.requestStore.GetRequests(r.Context())
	if err == sql.ErrNoRows {
		rh.logger.Error().Err(err).Msg("No Requests found")
		ErrorJSON(w, 404, "Not Found")
		return
	}
	if err != nil {
		rh.logger.Error().Err(err).Msg("Could not get Requests")
		ErrorJSON(w, 500, "Internal Server Error")
		return
	}

	SendJSON(w, requests)
}
