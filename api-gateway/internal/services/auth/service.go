package auth

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/rs/zerolog"
)

type Service struct {
	client Client
	log    zerolog.Logger
}

func New(client Client, log zerolog.Logger) *Service {
	return &Service{client: client, log: log}
}

func (s *Service) Login(ctx context.Context, credentials domain.LoginCredentials) (domain.AuthTokens, error) {
	s.log.Debug().Msg("login via auth service")
	return s.client.Login(ctx, credentials)
}

func (s *Service) Register(ctx context.Context, user domain.RegisterUser) (domain.User, error) {
	s.log.Debug().Msg("register via auth service")
	return s.client.Register(ctx, user)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (domain.AuthTokens, error) {
	s.log.Debug().Msg("refresh tokens via auth service")
	return s.client.Refresh(ctx, refreshToken)
}
