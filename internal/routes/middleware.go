package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Fozzyack/signalstack-backend/internal/app"
)

type contextKey string



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

func CheckAuth(app *app.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.TrimLeft(r.Header.Get("Authorization"), "Bearer ")
			ctx := r.Context()

			if token == "" {
				app.Logger.Error().Msg("No token provided")
				ErrorJSON(w, http.StatusUnauthorized, "No token provided")
				return
			}

			session, err := app.SessionStore.GetSessionByToken(ctx, token)
			if err != nil {
				app.Logger.Error().Err(err).Msg("Failed to get session")
				ErrorJSON(w, http.StatusUnauthorized, "Failed to get session")
				return
			}
			if session.ExpiresAt.UTC().Before(time.Now().UTC()) {
				app.Logger.Error().Msg("Session expired")
				ErrorJSON(w, http.StatusUnauthorized, "Session expired")
				return
			}

			user, err := app.UserStore.GetUserByID(ctx, session.UserID)
			if err != nil {
				app.Logger.Error().Err(err).Msg("Failed to get user")
				ErrorJSON(w, http.StatusUnauthorized, "Failed to get user")
				return
			}

			ctx = context.WithValue(ctx, "authenticated_user", user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
