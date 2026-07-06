package interceptors

import (
	"context"
	"runtime/debug"
	"time"

	authv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/auth/v1"
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
				log.Error().
					Any("panic", recovered).
					Bytes("stack", debug.Stack()).
					Str("method", info.FullMethod).
					Msg("grpc panic recovered")
				err = status.Error(codes.Internal, "internal server error")
			}
			log.Info().
				Str("method", info.FullMethod).
				Dur("duration", time.Since(start)).
				Msg("grpc request")
		}()
		if err := authorize(req); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func authorize(req any) error {
	switch req.(type) {
	case *authv1.LoginRequest, *authv1.RegisterRequest, *authv1.RefreshRequest:
		return nil
	default:
		return status.Error(codes.PermissionDenied, "access denied")
	}
}
