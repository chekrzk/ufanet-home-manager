package main

import (
	"context"
	"net"
	"os/signal"
	"syscall"
	"time"

	profilev1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/profile/v1"
	"github.com/chekrzk/ufanet-home-manager/profile-service/api/interceptors"
	profileserver "github.com/chekrzk/ufanet-home-manager/profile-service/api/server"
	"github.com/chekrzk/ufanet-home-manager/profile-service/internal/repository"
	"github.com/chekrzk/ufanet-home-manager/profile-service/internal/service"
	"github.com/chekrzk/ufanet-home-manager/profile-service/resources/config"
	"github.com/chekrzk/ufanet-home-manager/profile-service/resources/logger"
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
	repo := repository.NewProfileRepository(db)
	migrationCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err := repo.Migrate(migrationCtx); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate database")
	}

	listener, err := net.Listen("tcp", cfg.GRPC.Addr())
	if err != nil {
		log.Fatal().Err(err).Str("addr", cfg.GRPC.Addr()).Msg("failed to listen")
	}
	profileService := service.New(repo, log)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Unary(log)))
	profilev1.RegisterProfileServiceServer(grpcServer, profileserver.New(profileService))

	go func() {
		log.Info().Str("addr", cfg.GRPC.Addr()).Msg("starting profile-service")
		if err := grpcServer.Serve(listener); err != nil {
			log.Error().Err(err).Msg("grpc server stopped")
			stop()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("stopping profile-service")
	grpcServer.GracefulStop()
}
