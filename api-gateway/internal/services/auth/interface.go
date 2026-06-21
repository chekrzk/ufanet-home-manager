package auth

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
)

type Service interface {
	Login(ctx context.Context, req dto.LoginRequest) (dto.AuthTokens, error)
	Register(ctx context.Context, req dto.RegisterRequest) (domain.User, error)
	Refresh(ctx context.Context, refreshToken string) (dto.AuthTokens, error)
}
