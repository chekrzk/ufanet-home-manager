package app

import (
	"github.com/chekrzk/ufanet-home-manager/api-gateway/config"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/logger"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/router"
)

// Run собирает gateway в одном месте, чтобы main оставался точкой запуска,
// а детали конфигурации, DI и HTTP-сервера не протекали наружу.
func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := logger.New(cfg)

	container, err := newContainer(cfg, log)
	if err != nil {
		return err
	}
	defer container.Close()

	r := router.New(container.app, container.mw, container.handlers)
	r.Register()

	log.Info().
		Str("addr", cfg.HTTP.Addr()).
		Msg("starting api-gateway")

	return container.app.Listen(cfg.HTTP.Addr())
}
