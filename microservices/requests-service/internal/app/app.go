package app

import (
	"context"
	"net"
	"os/signal"
	"syscall"
	"time"

	requestsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/requests/v1"
	"github.com/chekrzk/ufanet-home-manager/requests-service/api/interceptors"
	requestsserver "github.com/chekrzk/ufanet-home-manager/requests-service/api/server"
	notificationclient "github.com/chekrzk/ufanet-home-manager/requests-service/infra/clients/notification_client"
	"github.com/chekrzk/ufanet-home-manager/requests-service/internal/repository"
	"github.com/chekrzk/ufanet-home-manager/requests-service/internal/service"
	"github.com/chekrzk/ufanet-home-manager/requests-service/resources"
	"google.golang.org/grpc"
)

func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	res, err := resources.New(ctx)
	if err != nil {
		return err
	}
	defer res.Close()

	repo := repository.NewRequestRepository(res.DB)
	migrationCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err := repo.Migrate(migrationCtx); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", res.Config.GRPC.Addr())
	if err != nil {
		return err
	}

	requestsService := service.New(repo, notificationclient.New(res.NotificationsConn), res.Log)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Unary(res.Log)))
	requestsv1.RegisterRequestsServiceServer(grpcServer, requestsserver.New(requestsService))

	go func() {
		res.Log.Info().Str("addr", res.Config.GRPC.Addr()).Msg("starting requests-service")
		if err := grpcServer.Serve(listener); err != nil {
			res.Log.Error().Err(err).Msg("grpc server stopped")
			stop()
		}
	}()

	<-ctx.Done()
	res.Log.Info().Msg("stopping requests-service")
	grpcServer.GracefulStop()
	return nil
}
