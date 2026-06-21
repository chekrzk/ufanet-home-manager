package services

import (
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/config"
	authservice "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/auth"
	newsservice "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/news"
	notificationsservice "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/notifications"
	profileservice "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/profile"
	requestsservice "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/requests"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Services struct {
	Auth          authservice.Service
	News          newsservice.Service
	Notifications notificationsservice.Service
	Profile       profileservice.Service
	Requests      requestsservice.Service
}

type GRPCServices struct {
	Services
	conns []*grpc.ClientConn
}

func NewGRPCServices(cfg config.ServiceConfig) (*GRPCServices, error) {
	authConn, err := dial(cfg.AuthAddr)
	if err != nil {
		return nil, err
	}

	newsConn, err := dial(cfg.NewsAddr)
	if err != nil {
		_ = authConn.Close()
		return nil, err
	}

	notificationsConn, err := dial(cfg.NotificationsAddr)
	if err != nil {
		_ = authConn.Close()
		_ = newsConn.Close()
		return nil, err
	}

	requestsConn, err := dial(cfg.RequestsAddr)
	if err != nil {
		_ = authConn.Close()
		_ = newsConn.Close()
		_ = notificationsConn.Close()
		return nil, err
	}

	profileConn, err := dial(cfg.ProfileAddr)
	if err != nil {
		_ = authConn.Close()
		_ = newsConn.Close()
		_ = notificationsConn.Close()
		_ = requestsConn.Close()
		return nil, err
	}

	grpcServices := &GRPCServices{
		conns: []*grpc.ClientConn{
			authConn,
			newsConn,
			notificationsConn,
			requestsConn,
			profileConn,
		},
	}

	grpcServices.Services = Services{
		Auth:          authservice.New(authConn),
		News:          newsservice.New(newsConn),
		Notifications: notificationsservice.New(notificationsConn),
		Profile:       profileservice.New(profileConn),
		Requests:      requestsservice.New(requestsConn),
	}

	return grpcServices, nil
}

func (s *GRPCServices) Close() error {
	var lastErr error
	for _, conn := range s.conns {
		if err := conn.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func dial(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
