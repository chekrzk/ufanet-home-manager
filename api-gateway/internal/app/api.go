package app

import (
	"github.com/gofiber/fiber/v2"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/config"
	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	authHandler "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/auth"
	healthHandler "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/health"
	newsHandler "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/news"
	notificationsHandler "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/notifications"
	profileHandler "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/profile"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/logger"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/middlewares"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/router"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logger.New(cfg)

	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: gwerrors.Fail,
	})

	mw := middlewares.New(cfg, log)

	svc, err := services.NewGRPCServices(cfg.Service)
	if err != nil {
		return err
	}
	defer func() {
		if err := svc.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close grpc connections")
		}
	}()

	health := healthHandler.NewHandler(cfg.App.Name)
	auth := authHandler.NewHandler(svc.Auth)
	news := newsHandler.NewHandler(svc.News)
	notifications := notificationsHandler.NewHandler(svc.Notifications)
	profile := profileHandler.NewHandler(svc.Profile)

	handlers := router.Handlers{
		Health:        health,
		Auth:          auth,
		News:          news,
		Notifications: notifications,
		Profile:       profile,
	}

	r := router.New(fiberApp, mw, handlers)
	r.Register()

	log.Info().
		Str("addr", cfg.HTTP.Addr()).
		Msg("starting api-gateway")

	return fiberApp.Listen(cfg.HTTP.Addr())
}
