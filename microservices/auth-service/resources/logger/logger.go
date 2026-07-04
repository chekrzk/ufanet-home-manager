package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/chekrzk/ufanet-home-manager/auth-service/resources/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New(cfg *config.Config) zerolog.Logger {
	level := parseLevel(cfg.Log.Level)
	zerolog.SetGlobalLevel(level)

	var output io.Writer = os.Stdout
	if cfg.Log.Pretty {
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
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
