package app

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
	"github.com/chekrzk/ufanet-home-manager/news-service/resources"
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

	repo := repository.NewNewsRepository(res.DB)
	migrationCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err := repo.Migrate(migrationCtx); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", res.Config.GRPC.Addr())
	if err != nil {
		return err
	}

	newsService := service.New(repo, notificationclient.New(res.NotificationsConn), res.Log)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Unary(res.Log)))
	newsv1.RegisterNewsServiceServer(grpcServer, newsserver.New(newsService))

	go func() {
		res.Log.Info().Str("addr", res.Config.GRPC.Addr()).Msg("starting news-service")
		if err := grpcServer.Serve(listener); err != nil {
			res.Log.Error().Err(err).Msg("grpc server stopped")
			stop()
		}
	}()

	<-ctx.Done()
	res.Log.Info().Msg("stopping news-service")
	grpcServer.GracefulStop()
	return nil
}
