package app

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
	"github.com/chekrzk/ufanet-home-manager/auth-service/resources"
	authv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/auth/v1"
	"google.golang.org/grpc"
)

// Run собирает auth-service как gRPC-приложение: ресурсы создаются один раз,
// миграции выполняются до приема трафика, а shutdown закрывает сервер корректно.
func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	res, err := resources.New(ctx)
	if err != nil {
		return err
	}
	defer res.Close()

	userRepo := repository.NewUserRepository(res.DB)
	migrationCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err := userRepo.Migrate(migrationCtx); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", res.Config.GRPC.Addr())
	if err != nil {
		return err
	}

	authService := service.NewAuthService(userRepo, hasher.New(), jwtmanager.NewManager(res.Config.JWT), res.Log)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Unary(res.Log)))
	authv1.RegisterAuthServiceServer(grpcServer, authserver.New(authService))

	go func() {
		res.Log.Info().Str("addr", res.Config.GRPC.Addr()).Msg("starting auth-service")
		if err := grpcServer.Serve(listener); err != nil {
			res.Log.Error().Err(err).Msg("grpc server stopped")
			stop()
		}
	}()

	<-ctx.Done()
	res.Log.Info().Msg("stopping auth-service")
	grpcServer.GracefulStop()
	return nil
}
