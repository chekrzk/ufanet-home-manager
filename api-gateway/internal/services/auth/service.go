package auth

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
	authv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/auth/v1"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	"google.golang.org/grpc"
)

type service struct {
	client authv1.AuthServiceClient
}

func New(conn *grpc.ClientConn) Service {
	return service{client: authv1.NewAuthServiceClient(conn)}
}

func (s service) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthTokens, error) {
	resp, err := s.client.Login(ctx, &authv1.LoginRequest{
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return dto.AuthTokens{}, err
	}
	return dto.AuthTokens{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
		ExpiresIn:    resp.GetExpiresIn(),
	}, nil
}

func (s service) Register(ctx context.Context, req dto.RegisterRequest) (domain.User, error) {
	resp, err := s.client.Register(ctx, &authv1.RegisterRequest{
		Phone:     req.Phone,
		Password:  req.Password,
		FullName:  req.FullName,
		HouseId:   req.HouseID,
		Apartment: req.Apartment,
	})
	if err != nil {
		return domain.User{}, err
	}
	return userFromProto(resp), nil
}

func (s service) Refresh(ctx context.Context, refreshToken string) (dto.AuthTokens, error) {
	resp, err := s.client.Refresh(ctx, &authv1.RefreshRequest{RefreshToken: refreshToken})
	if err != nil {
		return dto.AuthTokens{}, err
	}
	return dto.AuthTokens{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
		ExpiresIn:    resp.GetExpiresIn(),
	}, nil
}

func userFromProto(user *commonv1.User) domain.User {
	if user == nil {
		return domain.User{}
	}
	return domain.User{
		ID:        user.GetId(),
		Phone:     user.GetPhone(),
		FullName:  user.GetFullName(),
		Role:      user.GetRole(),
		HouseID:   user.GetHouseId(),
		Apartment: user.GetApartment(),
	}
}
