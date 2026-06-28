package main

import (
	"context"
	"net"
	"os/signal"
	"syscall"
	"time"

	notificationsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/notifications/v1"
	"github.com/chekrzk/ufanet-home-manager/notification-service/api/interceptors"
	notificationserver "github.com/chekrzk/ufanet-home-manager/notification-service/api/server"
	redisstream "github.com/chekrzk/ufanet-home-manager/notification-service/infra/streams/redis_stream"
	"github.com/chekrzk/ufanet-home-manager/notification-service/internal/repository"
	"github.com/chekrzk/ufanet-home-manager/notification-service/internal/service"
	"github.com/chekrzk/ufanet-home-manager/notification-service/resources/config"
	"github.com/chekrzk/ufanet-home-manager/notification-service/resources/logger"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
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
	repo := repository.NewDeviceRepository(db)
	migrationCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err := repo.Migrate(migrationCtx); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate database")
	}

	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr})
	defer redisClient.Close()

	listener, err := net.Listen("tcp", cfg.GRPC.Addr())
	if err != nil {
		log.Fatal().Err(err).Str("addr", cfg.GRPC.Addr()).Msg("failed to listen")
	}
	notificationService := service.New(repo, redisstream.New(redisClient, cfg.Redis.Stream), log)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Unary(log)))
	notificationsv1.RegisterNotificationsServiceServer(grpcServer, notificationserver.New(notificationService))

	go func() {
		log.Info().Str("addr", cfg.GRPC.Addr()).Msg("starting notification-service")
		if err := grpcServer.Serve(listener); err != nil {
			log.Error().Err(err).Msg("grpc server stopped")
			stop()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("stopping notification-service")
	grpcServer.GracefulStop()
}
