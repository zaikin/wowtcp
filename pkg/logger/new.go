package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

func NewLogger(cfg *Config) *zerolog.Logger {
	var output io.Writer
	output = os.Stdout
	if cfg.Console {
		output = zerolog.ConsoleWriter{Out: os.Stderr}
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	if cfg.Caller {
		logger = logger.With().Caller().Logger()
	}

	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		logger.Warn().Err(err).Str("level", cfg.Level).Msg("Cannot parse logging level")
	} else {
		logger = logger.Level(level)
	}

	return &logger
}
