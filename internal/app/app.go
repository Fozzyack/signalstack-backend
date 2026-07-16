package app

import (
	"os"

	"github.com/rs/zerolog"
)



type Application struct {
	Logger zerolog.Logger
}


func NewApplication() (*Application, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()


	app := &Application{
		Logger: logger,
	}
	return app, nil

}
