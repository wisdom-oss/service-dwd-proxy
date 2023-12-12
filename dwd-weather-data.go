package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
	wisdomMiddleware "github.com/wisdom-oss/microservice-middlewares/v3"

	"github.com/wisdom-oss/service-dwd-proxy/globals"
	"github.com/wisdom-oss/service-dwd-proxy/helpers"
	"github.com/wisdom-oss/service-dwd-proxy/routes"
)

// the main function bootstraps the http server and handlers used for this
// microservice
func main() {
	// create a new logger for the main function
	l := log.With().Str("part", "http-server").Logger()
	l.Info().Msgf("starting %s service", globals.ServiceName)

	// create a new router
	router := chi.NewRouter()
	// add some middlewares to the router to allow identifying requests
	router.Use(wisdomMiddleware.ErrorHandler(globals.ServiceName, globals.Errors))
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(httplog.Handler(l))
	compressor := chiMiddleware.NewCompressor(5)
	compressor.SetEncoder("br", func(w io.Writer, level int) io.Writer {
		return brotli.NewWriterLevel(w, level)
	})
	router.Use(compressor.Handler)
	// now add the authorization middleware to the router
	router.Use(wisdomMiddleware.Authorization(globals.AuthorizationConfiguration, globals.ServiceName))
	// now mount the admin router
	router.Get("/", routes.DiscoverMetadata)
	router.Get("/{stationID}", routes.StationInformation)
	router.Get("/{stationID}/{dataType}", routes.DataTypeInformation)
	router.Get("/{stationID}/{dataType}/{resolution}", routes.TimeSeries)

	// now boot up the service
	// Configure the HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", globals.Environment["LISTEN_PORT"]),
		WriteTimeout: time.Second * 600,
		ReadTimeout:  time.Second * 600,
		IdleTimeout:  time.Second * 600,
		Handler:      router,
	}

	// now run an initial discovery of all dwd stations if no station list is
	// currently present
	stationListEntries, err := globals.RedisClient.Exists(context.Background(), "dwd-station-list").Result()
	if err != nil {
		l.Error().Err(err).Msg("unable to verify the existence of the station list")
	}
	if stationListEntries == 0 {
		log.Warn().Msg("no station list currently present. running initial update")
		helpers.RunDiscovery()
	}

	// Start the server and log errors that happen while running it
	go func() {
		l.Info().Msg("starting http server for requests")
		if err := server.ListenAndServe(); err != nil {
			l.Fatal().Err(err).Msg("An error occurred while starting the http server")
		}
	}()

	// Set up the signal handling to allow the server to shut down gracefully

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, os.Interrupt)

	// now create a ticker that is used to periodically renew the station
	// discovery
	discoveryUpdate := time.Tick(5 * 60 * time.Second)

serviceLoop:
	for {
		select {
		case <-discoveryUpdate:
			log.Info().Msg("updating station data")
			helpers.RunDiscovery()
			log.Info().Msg("finished station update")
		case <-cancelSignal:
			break serviceLoop
		}
	}

	log.Info().Msg("shutting down service after STOPSIGNAL")

}
