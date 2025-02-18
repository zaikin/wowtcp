package tcpserver

import (
	"context"
	"fmt"
	"net"
	"time"
	"wowtcp/pkg/challenger"
	tcpmessage "wowtcp/pkg/tcpMessage"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	ForceShutdownTimeout = 1 * time.Second
)

//go:generate mockery --name=Repository --output=./mocks --outpkg=mocks
type Repository interface {
	GetWoWQuote() string
}

//go:generate mockery --name=Challenger --output=./mocks --outpkg=mocks
type Challenger interface {
	NewChallenge(resourceType string) challenger.Challenge
}

type Server struct {
	listener   net.Listener
	repository Repository
	challenger Challenger
}

func NewServer(cfg Config, repo Repository, challenger Challenger) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, err
	}

	return &Server{
		listener:   listener,
		repository: repo,
		challenger: challenger,
	}, nil
}

func (s *Server) Start(ctx context.Context) {
	logger := zerolog.Ctx(ctx).With().Str("component", "server").Logger()
	ctx = logger.WithContext(ctx)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				logger.Warn().Msg("Listener closed, stopping server")
				return
			}
			logger.Error().Err(err).Msg("Error accepting connection")
			continue
		}

		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	logger := zerolog.Ctx(ctx).With().Str("component", "server").Logger()
	s.listener.Close()
	select {
	case <-ctx.Done():
		logger.Info().Msg("Shutdown complete")
	case <-time.After(ForceShutdownTimeout):
		logger.Warn().Msg("Shutdown timeout reached, forcing close")
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	log := zerolog.Ctx(ctx).With().Str("method", "handleConnection").Str("remoteAddr", conn.RemoteAddr().String()).Logger()
	ctx = log.WithContext(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			message, err := tcpmessage.Read(conn)
			if err != nil {
				log.Error().Err(err).Msg("Error reading message")
				return
			}
			log.Debug().Str("received message", message).Msg("Received message")

			switch message {
			case "quit!":
				log.Info().Msg("Received quit message, closing connection")
				conn.Close()
				return
			case "quote!":
				if err = s.handleQuote(ctx, conn); err != nil {
					log.Warn().Err(err).Msg("Error handling quote")
					conn.Close()
					return
				}
			default:
				log.Debug().Str("message", message).Msg("Unknown message received")
			}
		}
	}
}

func (s *Server) handleQuote(ctx context.Context, conn net.Conn) error {
	log := zerolog.Ctx(ctx).With().Str("method", "handleQuote").Logger()
	chall := s.challenger.NewChallenge("quote")

	_, err := tcpmessage.Write(conn, chall.GetChallengeMessage()+"\n")
	if err != nil {
		err = errors.Wrap(err, "Error sending challenge")
		return err
	}

	log.Info().Str("sent challenge", chall.GetChallengeMessage()).Msg("Sent challenge to client")

	nonce, err := tcpmessage.Read(conn)
	if err != nil {
		err = errors.Wrap(err, "Error reading nonce response")
		return err
	}

	log.Info().Str("received nonce", nonce).Msg("Received nonce response")

	if chall.VerifyPoW(nonce) {
		quote := s.repository.GetWoWQuote()
		if _, err := tcpmessage.Write(conn, fmt.Sprintf("quote: %s\n", quote)); err != nil {
			err := errors.Wrap(err, "Error sending quote")
			return err
		}
		log.Info().Str("sent quote", quote).Msg("Sent quote to client")
	} else {
		errorMessage := "Invalid nonce"
		if _, err := tcpmessage.Write(conn, errorMessage+"\n"); err != nil {
			err = errors.Wrap(err, "Error sending error message")
			return err
		}
		log.Warn().Str("sent error", errorMessage).Msg("Sent error message to client")
	}

	return nil
}
