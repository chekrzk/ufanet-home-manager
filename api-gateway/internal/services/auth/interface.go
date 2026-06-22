package auth

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
)

type Service interface {
	Login(ctx context.Context, credentials domain.LoginCredentials) (domain.AuthTokens, error)
	Register(ctx context.Context, user domain.RegisterUser) (domain.User, error)
	Refresh(ctx context.Context, refreshToken string) (domain.AuthTokens, error)
}

type Client interface {
	Login(ctx context.Context, credentials domain.LoginCredentials) (domain.AuthTokens, error)
	Register(ctx context.Context, user domain.RegisterUser) (domain.User, error)
	Refresh(ctx context.Context, refreshToken string) (domain.AuthTokens, error)
}
