package logger

import (
	"os"

	"github.com/chekrzk/ufanet-home-manager/news-service/resources/config"
	"github.com/rs/zerolog"
)

func New(cfg *config.Config) zerolog.Logger {
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	if !cfg.Log.Pretty {
		return zerolog.New(os.Stdout).With().Timestamp().Str("service", cfg.App.Name).Str("env", cfg.App.Env).Logger()
	}
	return zerolog.New(output).With().Timestamp().Str("service", cfg.App.Name).Str("env", cfg.App.Env).Logger()
}
