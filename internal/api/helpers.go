package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Fozzyack/signalstack-backend/internal/models"
)

func SendJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

func DecodeJSON(r *http.Request, out any) error {
	return json.NewDecoder(r.Body).Decode(out)
}

func ErrorJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func GetUserFromContext(ctx context.Context) *models.User {
	return ctx.Value("authenticated_user").(*models.User)
}
