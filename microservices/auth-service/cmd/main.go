package main

import (
	"context"
	"net"
	"os/signal"
	"syscall"
	"time"

	"github.com/chekrzk/ufanet-home-manager/auth-service/api/interceptors"
	authserver "github.com/chekrzk/ufanet-home-manager/auth-service/api/server"
	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/hasher"
	jwtmanager "github.com/chekrzk/ufanet-home-manager/auth-service/internal/jwt"
	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/repository"
	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/service"
	"github.com/chekrzk/ufanet-home-manager/auth-service/resources/config"
	"github.com/chekrzk/ufanet-home-manager/auth-service/resources/logger"
	authv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/auth/v1"
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

	userRepo := repository.NewUserRepository(db)
	migrationCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err := userRepo.Migrate(migrationCtx); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate database")
	}

	authService := service.NewAuthService(
		userRepo,
		hasher.New(),
		jwtmanager.NewManager(cfg.JWT),
		log,
	)

	listener, err := net.Listen("tcp", cfg.GRPC.Addr())
	if err != nil {
		log.Fatal().Err(err).Str("addr", cfg.GRPC.Addr()).Msg("failed to listen")
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Unary(log)))
	authv1.RegisterAuthServiceServer(grpcServer, authserver.New(authService))

	go func() {
		log.Info().Str("addr", cfg.GRPC.Addr()).Msg("starting auth-service")
		if err := grpcServer.Serve(listener); err != nil {
			log.Error().Err(err).Msg("grpc server stopped")
			stop()
		}
	}()

	<-ctx.Done()
	log.Info().Msg("stopping auth-service")
	grpcServer.GracefulStop()
}
