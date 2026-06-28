package server

import (
	"context"

	jwtmanager "github.com/chekrzk/ufanet-home-manager/auth-service/internal/jwt"
	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/models"
)

type AuthService interface {
	Login(ctx context.Context, cmd models.LoginCommand) (jwtmanager.Pair, error)
	Register(ctx context.Context, cmd models.RegisterCommand) (models.User, error)
	Refresh(ctx context.Context, refreshToken string) (jwtmanager.Pair, error)
}
