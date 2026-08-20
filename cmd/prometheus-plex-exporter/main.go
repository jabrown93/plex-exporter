package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/grafana/plexporter/pkg/metrics"
	"github.com/grafana/plexporter/pkg/plex"
)

const (
	MetricsServerAddr = ":9000"
)

var (
	log = slog.New(slog.NewTextHandler(os.Stderr, nil))
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	serverAddress := os.Getenv("PLEX_SERVER")
	if serverAddress == "" {
		log.Error("PLEX_SERVER environment variable must be specified")
		os.Exit(1)
	}

	plexToken := os.Getenv("PLEX_TOKEN")
	if plexToken == "" {
		log.Error("PLEX_TOKEN environment variable must be specified")
		os.Exit(1)
	}

	server, err := plex.NewServer(serverAddress, plexToken)
	if err != nil {
		log.Error("cannot initialize connection to plex server", "error", err)
		os.Exit(1)
	}

	metrics.Register(server)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	metricsServer := http.Server{
		Addr:         MetricsServerAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info("starting metrics server on " + MetricsServerAddr)
		err = metricsServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Error("cannot start metrics server", "error", err)
		}
	}()

	exitCode := 0
	err = server.Listen(ctx, log)
	if err != nil {
		log.Error("cannot listen to plex server events", "error", err)
		exitCode = 1
	}

	log.Debug("shutting down metrics server")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()
	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		log.Error("cannot gracefully shutdown metrics server", "error", err)
	}

	os.Exit(exitCode)
}
