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

// New принимает абстрактный auth client, чтобы gateway-сценарии не зависели
// от конкретной gRPC-реализации.
func New(client Client, log zerolog.Logger) *Service {
	return &Service{client: client, log: log}
}

// Login оставляет HTTP-слой тонким: сценарий авторизации проходит через
// доменную модель и единую точку вызова auth-service.
func (s *Service) Login(ctx context.Context, credentials domain.LoginCredentials) (domain.AuthTokens, error) {
	s.log.Debug().Msg("login via auth service")
	return s.client.Login(ctx, credentials)
}

// Register отделяет регистрацию пользователя от транспорта, чтобы handler
// не знал, каким upstream-сервисом она фактически выполняется.
func (s *Service) Register(ctx context.Context, user domain.RegisterUser) (domain.User, error) {
	s.log.Debug().Msg("register via auth service")
	return s.client.Register(ctx, user)
}

// Refresh держит обновление токенов в auth-сценарии, а не размазывает его
// между HTTP handler и gRPC client.
func (s *Service) Refresh(ctx context.Context, refreshToken string) (domain.AuthTokens, error) {
	s.log.Debug().Msg("refresh tokens via auth service")
	return s.client.Refresh(ctx, refreshToken)
}
