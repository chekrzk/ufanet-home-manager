package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/config"
	authclient "github.com/chekrzk/ufanet-home-manager/api-gateway/infra/clients/auth_client"
	newsclient "github.com/chekrzk/ufanet-home-manager/api-gateway/infra/clients/news_client"
	notificationsclient "github.com/chekrzk/ufanet-home-manager/api-gateway/infra/clients/notifications_client"
	profileclient "github.com/chekrzk/ufanet-home-manager/api-gateway/infra/clients/profile_client"
	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	authHandler "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/auth"
	healthHandler "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/health"
	newsHandler "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/news"
	notificationsHandler "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/notifications"
	profileHandler "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/profile"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/logger"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/middlewares"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/router"
	authservice "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/auth"
	newsservice "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/news"
	notificationsservice "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/notifications"
	profileservice "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/profile"
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

	conns, err := serviceConnections(cfg.Service)
	if err != nil {
		return err
	}
	defer func() {
		closeConnections(log, conns.all()...)
	}()

	health := healthHandler.NewHandler(cfg.App.Name)
	authClient := authclient.New(conns.auth, log)
	newsClient := newsclient.New(conns.news, log)
	notificationsClient := notificationsclient.New(conns.notifications, log)
	profileClient := profileclient.New(conns.profile, log)

	auth := authHandler.NewHandler(authservice.New(authClient, log))
	news := newsHandler.NewHandler(newsservice.New(newsClient, log))
	notifications := notificationsHandler.NewHandler(notificationsservice.New(notificationsClient, log))
	profile := profileHandler.NewHandler(profileservice.New(profileClient, log))

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

type grpcConnections struct {
	auth          *grpc.ClientConn
	news          *grpc.ClientConn
	notifications *grpc.ClientConn
	requests      *grpc.ClientConn
	profile       *grpc.ClientConn
}

func (c grpcConnections) all() []*grpc.ClientConn {
	return []*grpc.ClientConn{
		c.auth,
		c.news,
		c.notifications,
		c.requests,
		c.profile,
	}
}

func serviceConnections(cfg config.ServiceConfig) (grpcConnections, error) {
	var conns grpcConnections

	var err error
	conns.auth, err = dial(cfg.AuthAddr)
	if err != nil {
		return conns, err
	}

	conns.news, err = dial(cfg.NewsAddr)
	if err != nil {
		closeConnections(zerolog.Nop(), conns.auth)
		return conns, err
	}

	conns.notifications, err = dial(cfg.NotificationsAddr)
	if err != nil {
		closeConnections(zerolog.Nop(), conns.auth, conns.news)
		return conns, err
	}

	conns.requests, err = dial(cfg.RequestsAddr)
	if err != nil {
		closeConnections(zerolog.Nop(), conns.auth, conns.news, conns.notifications)
		return conns, err
	}

	conns.profile, err = dial(cfg.ProfileAddr)
	if err != nil {
		closeConnections(zerolog.Nop(), conns.auth, conns.news, conns.notifications, conns.requests)
		return conns, err
	}

	return conns, nil
}

func dial(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func closeConnections(log zerolog.Logger, conns ...*grpc.ClientConn) {
	for _, conn := range conns {
		if conn == nil {
			continue
		}
		if err := conn.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close grpc connection")
		}
	}
}
