package app

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
	"github.com/chekrzk/ufanet-home-manager/profile-service/resources"
	"google.golang.org/grpc"
)

// Run поднимает profile-service после миграций, чтобы профили и работники
// имели готовую схему до первого gRPC-запроса.
func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	res, err := resources.New(ctx)
	if err != nil {
		return err
	}
	defer res.Close()

	repo := repository.NewProfileRepository(res.DB)
	migrationCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err := repo.Migrate(migrationCtx); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", res.Config.GRPC.Addr())
	if err != nil {
		return err
	}

	profileService := service.New(repo, res.Log)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Unary(res.Log)))
	profilev1.RegisterProfileServiceServer(grpcServer, profileserver.New(profileService))

	go func() {
		res.Log.Info().Str("addr", res.Config.GRPC.Addr()).Msg("starting profile-service")
		if err := grpcServer.Serve(listener); err != nil {
			res.Log.Error().Err(err).Msg("grpc server stopped")
			stop()
		}
	}()

	<-ctx.Done()
	res.Log.Info().Msg("stopping profile-service")
	grpcServer.GracefulStop()
	return nil
}
