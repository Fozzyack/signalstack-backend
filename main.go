package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Fozzyack/signalstack-backend/internal/api"
	"github.com/Fozzyack/signalstack-backend/internal/app"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	app.Logger.Info().Msg("App initialized")
	router := api.SetupRoutes(app)
	
	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 2 * time.Minute,
	}

	app.Logger.Info().Msg("Server initialized")

	app.Logger.Info().Str("port", port).Msg("Running Server")
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("Could not initialize server")
	}

}
