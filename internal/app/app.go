package app

import (
	"database/sql"
	"os"

	"github.com/Fozzyack/signalstack-backend/internal/database"
	"github.com/Fozzyack/signalstack-backend/migrations"
	"github.com/rs/zerolog"
)



type Application struct {
	Logger zerolog.Logger
	DB *sql.DB
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

	app := &Application{
		Logger: logger,
		DB: pgDb,
	}
	return app, nil

}
