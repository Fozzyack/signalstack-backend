package api

import (
	"database/sql"
	"net/http"

	"github.com/Fozzyack/signalstack-backend/internal/models"
	"github.com/Fozzyack/signalstack-backend/internal/store"
	"github.com/rs/zerolog"
)

type RequestAssignmentHandler struct {
	logger                  *zerolog.Logger
	requestAssignmentsStore store.RequestsAssignmentsStore
	userStore               store.UserStore
}

func NewRequestAssignmentHandler(logger *zerolog.Logger, requestAssignmentsStore store.RequestsAssignmentsStore, userStore store.UserStore) *RequestAssignmentHandler {
	return &RequestAssignmentHandler{
		logger:                  logger,
		requestAssignmentsStore: requestAssignmentsStore,
		userStore:               userStore,
	}
}

func (rah *RequestAssignmentHandler) GetAllRequestAssignments(w http.ResponseWriter, r *http.Request) {
	assignments, err := rah.requestAssignmentsStore.GetAllRequestAssignments(r.Context())
	if err == sql.ErrNoRows {
		rah.logger.Error().Err(err).Msg("No request assignments found")
		ErrorJSON(w, http.StatusNotFound, "Not Found")
		return
	}
	if err != nil {
		rah.logger.Error().Err(err).Msg("Could not get request assignments")
		ErrorJSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	type assignmentResponse struct {
		*models.RequestAssignment
		UserName string `json:"user_name"`
	}
	response := make([]assignmentResponse, 0, len(assignments))
	for _, assignment := range assignments {
		user, err := rah.userStore.GetUserByID(r.Context(), assignment.UserID)
		if err != nil {
			rah.logger.Error().Err(err).Str("user_id", assignment.UserID).Msg("Could not get assignment user")
			ErrorJSON(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		response = append(response, assignmentResponse{
			RequestAssignment: assignment,
			UserName:          user.Name,
		})
	}

	SendJSON(w, response)
}
