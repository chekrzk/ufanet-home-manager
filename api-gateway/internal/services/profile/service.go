package profile

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	profilev1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/profile/v1"
	"google.golang.org/grpc"
)

type service struct {
	client profilev1.ProfileServiceClient
}

func New(conn *grpc.ClientConn) Service {
	return service{client: profilev1.NewProfileServiceClient(conn)}
}

func (s service) Me(ctx context.Context, actor domain.AuthContext) (domain.User, error) {
	resp, err := s.client.Me(ctx, &profilev1.MeRequest{User: userContext(actor)})
	if err != nil {
		return domain.User{}, err
	}
	return userFromProto(resp), nil
}

func (s service) Update(ctx context.Context, actor domain.AuthContext, req dto.UpdateProfileRequest) (domain.User, error) {
	resp, err := s.client.Update(ctx, &profilev1.UpdateProfileRequest{
		User:      userContext(actor),
		FullName:  req.FullName,
		Apartment: req.Apartment,
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
