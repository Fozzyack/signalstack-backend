package api

import (
	"database/sql"
	"net/http"

	"github.com/Fozzyack/signalstack-backend/internal/models"
	"github.com/Fozzyack/signalstack-backend/internal/store"
	"github.com/rs/zerolog"
)

type RequestHandler struct {
	logger       *zerolog.Logger
	requestStore store.RequestStore
	requestAssignmentStore store.RequestsAssignmentsStore
}

func NewRequestHandler(logger *zerolog.Logger, requestStore store.RequestStore, requestAssignmentsStore store.RequestsAssignmentsStore) *RequestHandler {
	return &RequestHandler{
		logger:       logger,
		requestStore: requestStore,
		requestAssignmentStore: requestAssignmentsStore,
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

func (rh *RequestHandler) GetAllUserRequests(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	requestAssignments, err := rh.requestAssignmentStore.GetRequestAssignmentsByUserID(r.Context(), user.ID)
	if err == sql.ErrNoRows {
		rh.logger.Error().Err(err).Msg("No Request Assignments found")
		ErrorJSON(w, 404, "Not Found")
		return
	}
	if err != nil {
		rh.logger.Error().Err(err).Msg("Could not get Request Assignments")
		ErrorJSON(w, 500, "Internal Server Error")
		return
	}

	requests := make([]*models.Request, 0)
	for _, assignment := range requestAssignments {
		request, err := rh.requestStore.GetRequestById(r.Context(), assignment.RequestID)
		if err != nil {
			rh.logger.Error().Err(err).Msg("Could not get Request")
			ErrorJSON(w, 500, "Internal Server Error")
			return
		}
		requests = append(requests, request)
	}
	SendJSON(w, requests)
}





















