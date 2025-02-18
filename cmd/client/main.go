package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
	"wowtcp/pkg/challenger"
	"wowtcp/pkg/logger"
	tcpmessage "wowtcp/pkg/tcpMessage"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func handleMessage(ctx context.Context, conn net.Conn, message string) {
	log := zerolog.Ctx(ctx).With().Str("component", "client").Logger()
	switch {
	case message == "quit!":
		log.Info().Msg("Received quit message, closing connection")
		conn.Close()
		return
	case strings.HasPrefix(message, "challenge: "):
		handleChallenge(ctx, conn, message)
	case strings.HasPrefix(message, "quote: "):
		log.Info().Str("quote", message).Msg("Received quote")
	default:
		log.Warn().Str("message", message).Msg("Unknown message received")
	}
}

func handleChallenge(ctx context.Context, conn net.Conn, message string) {
	log := zerolog.Ctx(ctx).With().Str("component", "handleChallenge").Logger()
	log.Info().Str("challengeReq", message).Msg("Received challenge request")

	chall := challenger.HashcashChallenge{}

	if err := chall.ParseChallengeMessage(message); err != nil {
		log.Error().Err(err).Msg("Error parsing challenge")
		return
	}

	log.Info().Msg("Start solve")

	nonce := chall.SolvePoW()
	log.Info().Str("nonce", nonce).Msg("Solved PoW")

	if _, err := tcpmessage.Write(conn, fmt.Sprintf("%s\n", nonce)); err != nil {
		log.Error().Err(err).Msg("Error sending nonce")
		return
	}
	log.Info().Str("nonce", nonce).Msg("Sent nonce to server")
}

func main() {
	host := flag.String("host", "", "Host to connect to")
	port := flag.String("port", "", "Port to connect to")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	log.Debug().Str("host", *host).Str("port", *port).Bool("debug", *debug).Msg("Starting client parameters")

	loggerCfg := logger.Config{Console: true, Caller: true, Level: "info"}
	if *debug {
		loggerCfg.Level = "debug"
	}

	logger := logger.NewLogger(&loggerCfg)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(logger.WithContext(ctx))
	defer cancel()
	log := logger

	log.Info().Msg("Starting client")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to server")
	}
	defer conn.Close()

	log.Info().Msg("Connected to server")

	for {
		if _, err = tcpmessage.Write(conn, "quote!\n"); err != nil {
			log.Fatal().Err(err).Msg("Error sending quote request")
		}
		log.Info().Msg("Sent quote request")

		message, err := tcpmessage.Read(conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Info().Msg("Connection closed by server")
				break
			}
			log.Error().Err(err).Msg("Error reading message")
			continue
		}
		handleMessage(ctx, conn, message)
		time.Sleep(5 * time.Second)
	}
}
