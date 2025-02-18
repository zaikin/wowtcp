package config

import (
	"wowtcp/pkg/logger"

	"github.com/joeshaw/envdecode"
)

type Config struct {
	Port int `env:"APP_PORT,required"`

	Logger logger.Config
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := envdecode.StrictDecode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
