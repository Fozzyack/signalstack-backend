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

	AuthHandler        *api.AuthHandler
	HealthCheckHandler *api.HealthCheckHandler
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

	healthCheckHandler := api.NewHealthCheckHandler()
	authHandler := api.NewAuthHandler(&logger, userStore, sessionStore)

	app := &Application{
		Logger: logger,
		DB:     pgDb,

		UserStore:    userStore,
		SessionStore: sessionStore,

		AuthHandler:        authHandler,
		HealthCheckHandler: healthCheckHandler,
	}

	return app, nil

}
