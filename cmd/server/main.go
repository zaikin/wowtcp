package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wowtcp/internal/config"
	"wowtcp/internal/repository"
	"wowtcp/internal/tcpserver"
	"wowtcp/pkg/challenger"
	"wowtcp/pkg/logger"

	"github.com/rs/zerolog/log"
)

const (
	DefaultShutdownTimeout = 1 * time.Second
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}

	logger := logger.NewLogger(&cfg.Logger)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(logger.WithContext(ctx))

	repo := repository.NewInMemoryRepository()
	challenger := challenger.NewHashcashChallenger(&cfg.Challenger)

	server, err := tcpserver.NewServer(cfg.Server, repo, challenger)
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

		server.Shutdown(ctx)
		cancel()

		close(doneChan)
	}()

	log.Info().Msg("Starting TCP server")
	server.Start(ctx)

	<-doneChan
	log.Info().Msg("Server shut down gracefully")
}
