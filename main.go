package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/rs/zerolog/log"
	healthcheckServer "github.com/wisdom-oss/go-healthcheck/server"
	errorMiddleware "github.com/wisdom-oss/microservice-middlewares/v5/error"

	gographql "github.com/graph-gophers/graphql-go"

	"microservice/config"
	"microservice/graphql"
	"microservice/routes"

	"microservice/globals"
)

// the main function bootstraps the http server and handlers used for this
// microservice
func main() {
	// create a new logger for the main function
	l := log.With().Str("step", "main").Logger()
	l.Info().Msgf("starting %s service", globals.ServiceName)

	// create the healthcheck server
	hcServer := healthcheckServer.HealthcheckServer{}
	hcServer.InitWithFunc(func() error {
		// test if the database is reachable
		return globals.Db.Ping(context.Background())
	})
	err := hcServer.Start()
	if err != nil {
		l.Fatal().Err(err).Msg("unable to start healthcheck server")
	}
	go hcServer.Run()

	// create a new router
	router := chi.NewRouter()
	router.Use(config.Middlewares...)

	router.HandleFunc("/", routes.AllLocations)
	router.HandleFunc("/{stationID}", routes.SingleStation)
	router.HandleFunc("/measurements", routes.Measurements)

	// now start parsing the graphql part to allow graphql queries
	gqlSchema := gographql.MustParseSchema(graphQlSchema, &graphql.Query{}, gographql.UseFieldResolvers())
	router.Handle("/graphql", &relay.Handler{Schema: gqlSchema})
	router.NotFound(errorMiddleware.NotFoundError)

	// now boot up the service
	// Configure the HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", config.ListenAddress, globals.Environment["LISTEN_PORT"]),
		WriteTimeout: time.Second * 600,
		ReadTimeout:  time.Second * 600,
		IdleTimeout:  time.Second * 600,
		Handler:      router,
	}

	// Start the server and log errors that happen while running it
	go func() {
		if err := server.ListenAndServe(); err != nil {
			l.Fatal().Err(err).Msg("An error occurred while starting the http server")
		}
	}()

	// Set up the signal handling to allow the server to shut down gracefully

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, os.Interrupt)

	// Block further code execution until the shutdown signal was received
	l.Info().Msg("server ready to accept connections")
	<-cancelSignal

}
