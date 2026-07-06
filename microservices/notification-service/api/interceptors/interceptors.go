package interceptors

import (
	"context"
	"time"

	notificationsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/notifications/v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Unary(log zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error().Any("panic", recovered).Str("method", info.FullMethod).Msg("grpc panic recovered")
				err = status.Error(codes.Internal, "internal server error")
			}
			log.Info().Str("method", info.FullMethod).Dur("duration", time.Since(start)).Msg("grpc request")
		}()
		if err := authorize(req); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func authorize(req any) error {
	switch request := req.(type) {
	case *notificationsv1.RegisterDeviceRequest:
		return requireUser(request.GetUser().GetUserId(), request.GetUser().GetRole())
	case *notificationsv1.UnregisterDeviceRequest:
		return requireUser(request.GetUser().GetUserId(), request.GetUser().GetRole())
	case *notificationsv1.ListNotificationsRequest:
		return requireUser(request.GetUser().GetUserId(), request.GetUser().GetRole())
	case *notificationsv1.MarkReadRequest:
		return requireUser(request.GetUser().GetUserId(), request.GetUser().GetRole())
	case *notificationsv1.PublishNotificationRequest:
		return nil
	default:
		return status.Error(codes.PermissionDenied, "access denied")
	}
}

func requireUser(userID string, role string) error {
	if userID == "" || role == "" {
		return status.Error(codes.Unauthenticated, "authorization required")
	}
	return nil
}
