package app

import (
	"github.com/gofiber/fiber/v2"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/config"
	healthHandler "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/health"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/logger"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/middlewares"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/router"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logger.New(cfg)

	fiberApp := fiber.New()

	mw := middlewares.New()

	health := healthHandler.NewHandler(cfg.App.Name)

	handlers := router.Handlers{
		Health: health,
	}

	r := router.New(fiberApp, mw, handlers)
	r.Register()

	log.Info().
		Str("addr", cfg.HTTP.Addr()).
		Msg("starting api-gateway")

	return fiberApp.Listen(cfg.HTTP.Addr())
}
