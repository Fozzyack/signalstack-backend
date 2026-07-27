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

func (rh *RequestHandler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	type newRequestType struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Name        string `json:"clientName"`
		Email       string `json:"clientEmail"`
	}

	newRequest := newRequestType{}
	DecodeJSON(r, &newRequest)
	rh.logger.Info().Msg("Request Created")
	request, err := rh.requestStore.CreateRequest(r.Context(), newRequest.Title, newRequest.Description, newRequest.Name, newRequest.Email)
	if err != nil {
		rh.logger.Error().Err(err).Msg("Could not create Request")
		ErrorJSON(w, 500, "Internal Server Error")
		return
	}

	// Extra logic here for
	// Sending email
	// Adding tags

	SendJSON(w, request)
}





















