package config

import (
	"wowtcp/internal/tcpserver"
	"wowtcp/pkg/challenger"
	"wowtcp/pkg/logger"

	"github.com/joeshaw/envdecode"
)

type Config struct {
	Logger     logger.Config
	Challenger challenger.Config
	Server     tcpserver.Config
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := envdecode.StrictDecode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
