package main

import (
	"context"
	"net"
	"os/signal"
	"syscall"
	"time"

	newsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/news/v1"
	"github.com/chekrzk/ufanet-home-manager/news-service/api/interceptors"
	newsserver "github.com/chekrzk/ufanet-home-manager/news-service/api/server"
	notificationclient "github.com/chekrzk/ufanet-home-manager/news-service/infra/clients/notification_client"
	"github.com/chekrzk/ufanet-home-manager/news-service/internal/repository"
	"github.com/chekrzk/ufanet-home-manager/news-service/internal/service"
	"github.com/chekrzk/ufanet-home-manager/news-service/resources/config"
	"github.com/chekrzk/ufanet-home-manager/news-service/resources/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log := logger.New(cfg)

	db, err := gorm.Open(postgres.Open(cfg.DB.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}
	repo := repository.NewNewsRepository(db)
	migrationCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err := repo.Migrate(migrationCtx); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate database")
	}

	notificationsConn, err := grpc.NewClient(cfg.Service.NotificationsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect notifications service")
	}
	defer notificationsConn.Close()

	listener, err := net.Listen("tcp", cfg.GRPC.Addr())
	if err != nil {
		log.Fatal().Err(err).Str("addr", cfg.GRPC.Addr()).Msg("failed to listen")
	}
	newsService := service.New(repo, notificationclient.New(notificationsConn), log)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Unary(log)))
	newsv1.RegisterNewsServiceServer(grpcServer, newsserver.New(newsService))

	go func() {
		log.Info().Str("addr", cfg.GRPC.Addr()).Msg("starting news-service")
		if err := grpcServer.Serve(listener); err != nil {
			log.Error().Err(err).Msg("grpc server stopped")
			stop()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("stopping news-service")
	grpcServer.GracefulStop()
}
