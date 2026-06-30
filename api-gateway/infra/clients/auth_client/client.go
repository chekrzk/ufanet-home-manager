package auth_client

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	authv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/auth/v1"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Client struct {
	client authv1.AuthServiceClient
	log    zerolog.Logger
}

func New(conn *grpc.ClientConn, log zerolog.Logger) *Client {
	return &Client{client: authv1.NewAuthServiceClient(conn), log: log}
}

func (c *Client) Login(ctx context.Context, credentials domain.LoginCredentials) (domain.AuthTokens, error) {
	c.log.Debug().Msg("call auth grpc login")
	resp, err := c.client.Login(ctx, &authv1.LoginRequest{
		Phone:    credentials.Phone,
		Password: credentials.Password,
	})
	if err != nil {
		return domain.AuthTokens{}, err
	}
	return domain.AuthTokens{
		AccessToken:  resp.GetAccessToken(),
		RefreshToken: resp.GetRefreshToken(),
		ExpiresIn:    resp.GetExpiresIn(),
	}, nil
}

func (c *Client) Register(ctx context.Context, user domain.RegisterUser) (domain.User, error) {
	c.log.Debug().Msg("call auth grpc register")
	resp, err := c.client.Register(ctx, &authv1.RegisterRequest{
		Phone:    user.Phone,
		Password: user.Password,
		FullName: user.FullName,
		Role:     user.Role,
	})
	if err != nil {
		return domain.User{}, err
	}
	return userFromProto(resp), nil
}

func (c *Client) Refresh(ctx context.Context, refreshToken string) (domain.AuthTokens, error) {
	c.log.Debug().Msg("call auth grpc refresh")
	resp, err := c.client.Refresh(ctx, &authv1.RefreshRequest{RefreshToken: refreshToken})
	if err != nil {
		return domain.AuthTokens{}, err
	}
	return domain.AuthTokens{
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
