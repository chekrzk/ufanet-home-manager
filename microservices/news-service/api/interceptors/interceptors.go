package interceptors

import (
	"context"
	"time"

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
		return handler(ctx, req)
	}
}
