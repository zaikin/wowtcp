package tcpserver

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
	"wowtcp/pkg/challenger"
	"wowtcp/pkg/tcpio"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	ForceShutdownTimeout = 1 * time.Second
)

const (
	QuoteCommand    = "quote!"
	QuitCommand     = "quit!"
	NoncePrefix     = "nonce: "
	QuotePrefix     = "quote: "
	ChallengePrefix = "challenge: "
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
	defer func() {
		if r := recover(); r != nil {
			log := zerolog.Ctx(ctx).With().Str("method", "handleConnection").Logger()
			log.Error().Interface("recover", r).Msg("Recovered from panic")
			conn.Close()
		}
	}()

	log := zerolog.Ctx(ctx).With().Str("method", "handleConnection").Str("remoteAddr", conn.RemoteAddr().String()).Logger()
	ctx = log.WithContext(ctx)
	messages := tcpio.NewTCPReadWriter(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			message, err := messages.Read()
			if err != nil {
				log.Error().Err(err).Msg("Error reading message")
				return
			}
			log.Debug().Str("received message", message).Msg("Received message")

			switch message {
			case QuitCommand:
				log.Info().Msg("Received quit message, closing connection")
				conn.Close()
				return
			case QuoteCommand:
				if err = s.handleQuote(ctx, messages); err != nil {
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

func (s *Server) handleQuote(ctx context.Context, messages tcpio.ReadWriter) error {
	log := zerolog.Ctx(ctx).With().Str("method", "handleQuote").Logger()
	chall := s.challenger.NewChallenge("quote")

	_, err := messages.Write(ChallengePrefix + chall.GetChallengeMessage())
	if err != nil {
		err = errors.Wrap(err, "Error sending challenge")
		return err
	}

	log.Info().Str("sent challenge", chall.GetChallengeMessage()).Msg("Sent challenge to client")

	nonce, err := messages.Read()
	if err != nil {
		err = errors.Wrap(err, "Error reading nonce response")
		return err
	}

	if !strings.HasPrefix(nonce, NoncePrefix) {
		log.Error().Err(err).Str("nonce", nonce).Msg("Error parsing nonce response")
		err = errors.New("invalid nonce response format")
		return err
	}

	nonce = strings.TrimPrefix(nonce, NoncePrefix)

	log.Info().Str("received nonce", nonce).Msg("Received nonce response")

	if chall.VerifyPoW(nonce) {
		quote := s.repository.GetWoWQuote()
		if _, err = messages.Write(QuotePrefix + quote); err != nil {
			err = errors.Wrap(err, "Error sending quote")
			return err
		}
		log.Info().Str("sent quote", quote).Msg("Sent quote to client")
	} else {
		errorMessage := "Invalid nonce"
		if _, err = messages.Write(errorMessage); err != nil {
			err = errors.Wrap(err, "Error sending error message")
			return err
		}
		log.Warn().Str("sent error", errorMessage).Msg("Sent error message to client")
	}

	return nil
}
