package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/config"
)

func New(cfg *config.Config) zerolog.Logger {
	level := parseLevel(cfg.Log.Level)

	zerolog.SetGlobalLevel(level)

	var output io.Writer

	if cfg.Log.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	} else {
		output = os.Stdout
	}

	logger := zerolog.New(output).
		With().
		Timestamp().
		Str("service", cfg.App.Name).
		Str("env", cfg.App.Env).
		Logger()

	log.Logger = logger

	return logger
}

func parseLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
