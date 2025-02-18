package tcpserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	ForceShutdownTimeout = 5 * time.Second
)

type Server struct {
	listener net.Listener
}

func NewServer(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	return &Server{
		listener: listener,
	}, nil
}

func (s *Server) Start() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Error().Err(err).Msg("Error accepting connection")
			continue
		}

		go handleConnection(conn)
	}
}

func (s *Server) Shutdown(ctx context.Context) {
	s.listener.Close()
	select {
	case <-ctx.Done():
		log.Info().Msg("Shutdown complete")
	case <-time.After(ForceShutdownTimeout):
		log.Warn().Msg("Shutdown timeout reached, forcing close")
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Info().Str("client", conn.RemoteAddr().String()).Msg("Client connected")

	// TODO: to service

	buf := make([]byte, 1024) //TODO: messages types
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Error().Err(err).Msg("Error reading from connection")
			return
		}
		if n == 0 {
			return
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Error().Err(err).Msg("Error writing to connection")
			return
		}
	}
}
