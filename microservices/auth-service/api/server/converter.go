package server

import (
	jwtmanager "github.com/chekrzk/ufanet-home-manager/auth-service/internal/jwt"
	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/models"
	authv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/auth/v1"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
)

func loginCommandFromProto(req *authv1.LoginRequest) models.LoginCommand {
	return models.LoginCommand{
		Phone:    req.GetPhone(),
		Password: req.GetPassword(),
	}
}

func registerCommandFromProto(req *authv1.RegisterRequest) models.RegisterCommand {
	return models.RegisterCommand{
		Phone:    req.GetPhone(),
		Password: req.GetPassword(),
		Role:     req.GetRole(),
	}
}

func tokensToProto(tokens jwtmanager.Pair) *authv1.AuthTokens {
	return &authv1.AuthTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
	}
}

func userToProto(user models.User) *commonv1.User {
	return &commonv1.User{
		Id:    user.ID,
		Phone: user.Phone,
		Role:  string(user.Role),
	}
}
