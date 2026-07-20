package app

import (
	"database/sql"
	"os"

	"github.com/Fozzyack/signalstack-backend/internal/api"
	"github.com/Fozzyack/signalstack-backend/internal/database"
	"github.com/Fozzyack/signalstack-backend/internal/store"
	"github.com/Fozzyack/signalstack-backend/migrations"
	"github.com/rs/zerolog"
)

type Application struct {
	Logger zerolog.Logger
	DB     *sql.DB

	UserStore    store.UserStore
	SessionStore store.SessionStore
	RequestStore store.RequestStore

	AuthHandler        *api.AuthHandler
	HealthCheckHandler *api.HealthCheckHandler
	RequestHandler     *api.RequestHandler
	UserHandler        *api.UserHandler
}

func NewApplication() (*Application, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	pgDb, err := database.Open()
	if err != nil {
		return nil, err
	}
	err = database.Migrate(pgDb, migrations.FS)
	if err != nil {
		return nil, err
	}

	userStore := store.NewUserStore(pgDb)
	sessionStore := store.NewSessionStore(pgDb)
	requestsStore := store.NewRequestStore(pgDb)

	healthCheckHandler := api.NewHealthCheckHandler()
	authHandler := api.NewAuthHandler(&logger, userStore, sessionStore)
	requestsHandler := api.NewRequestHandler(&logger, requestsStore)
	userHandler := api.NewUserHandler(&logger, userStore)

	app := &Application{
		Logger: logger,
		DB:     pgDb,

		UserStore:    userStore,
		SessionStore: sessionStore,
		RequestStore: requestsStore,

		AuthHandler:        authHandler,
		HealthCheckHandler: healthCheckHandler,
		RequestHandler:     requestsHandler,
		UserHandler:        userHandler,
	}

	return app, nil

}
