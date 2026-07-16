package routes

import (
	"github.com/Fozzyack/signalstack-backend/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	return r
}
