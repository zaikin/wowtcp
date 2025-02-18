package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wowtcp/internal/config"
	"wowtcp/internal/tcpserver"
	"wowtcp/pkg/logger"

	"github.com/rs/zerolog/log"
)

const (
	DefaultShutdownTimeout = 5 * time.Second
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}

	log := logger.NewLogger(&cfg.Logger)

	server, err := tcpserver.NewServer(cfg.Port)
	if err != nil {
		log.Error().Err(err).Msg("Error starting TCP server")
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	doneChan := make(chan struct{})

	go func() {
		<-sigChan
		log.Info().Msg("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
		defer cancel()

		server.Shutdown(ctx)

		close(doneChan)
	}()

	log.Info().Msg("Starting TCP server")
	server.Start()

	<-doneChan
	log.Info().Msg("Server shut down gracefully")
}
