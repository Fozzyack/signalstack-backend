package routes

import (
	"github.com/Fozzyack/signalstack-backend/internal/app"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/health", app.HealthCheckHandler.HealthCheck)

	r.Post("/auth/login", app.AuthHandler.Login)
	r.With(CheckAuth(app)).Group(func(r chi.Router) {

		r.Get("/auth/check", app.AuthHandler.CheckAuth)
		r.Get("/requests", app.RequestHandler.GetAllRequests)
		r.Get("/request-assignments", app.RequestAssignmentHandler.GetAllRequestAssignments)

		r.Get("/users/me", app.UserHandler.GetLoggedInUser)
		r.Put("/users/me", app.UserHandler.UpdateUser)

	})

	return r
}
