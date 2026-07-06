package app

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
	"github.com/chekrzk/ufanet-home-manager/notification-service/resources"
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

	repo := repository.NewDeviceRepository(res.DB)
	migrationCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	if err := repo.Migrate(migrationCtx); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", res.Config.GRPC.Addr())
	if err != nil {
		return err
	}

	notificationService := service.New(repo, redisstream.New(res.Redis, res.Config.Redis.Stream), res.Log)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptors.Unary(res.Log)))
	notificationsv1.RegisterNotificationsServiceServer(grpcServer, notificationserver.New(notificationService))

	go func() {
		res.Log.Info().Str("addr", res.Config.GRPC.Addr()).Msg("starting notification-service")
		if err := grpcServer.Serve(listener); err != nil {
			res.Log.Error().Err(err).Msg("grpc server stopped")
			stop()
		}
	}()

	<-ctx.Done()
	res.Log.Info().Msg("stopping notification-service")
	grpcServer.GracefulStop()
	return nil
}
