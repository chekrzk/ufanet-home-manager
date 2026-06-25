package profile_client

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	profilev1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/profile/v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Client struct {
	client profilev1.ProfileServiceClient
	log    zerolog.Logger
}

func New(conn *grpc.ClientConn, log zerolog.Logger) *Client {
	return &Client{client: profilev1.NewProfileServiceClient(conn), log: log}
}

func (c *Client) Me(ctx context.Context, actor domain.AuthContext) (domain.User, error) {
	c.log.Debug().Str("user_id", actor.UserID).Msg("call profile grpc me")
	resp, err := c.client.Me(ctx, &profilev1.MeRequest{User: userContext(actor)})
	if err != nil {
		return domain.User{}, err
	}
	return userFromProto(resp), nil
}

func (c *Client) Update(ctx context.Context, actor domain.AuthContext, command domain.UpdateProfile) (domain.User, error) {
	c.log.Debug().Str("user_id", actor.UserID).Msg("call profile grpc update")
	resp, err := c.client.Update(ctx, &profilev1.UpdateProfileRequest{
		User:      userContext(actor),
		FullName:  command.FullName,
		Apartment: command.Apartment,
	})
	if err != nil {
		return domain.User{}, err
	}
	return userFromProto(resp), nil
}

func userContext(actor domain.AuthContext) *commonv1.UserContext {
	return &commonv1.UserContext{UserId: actor.UserID, Role: actor.Role}
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
