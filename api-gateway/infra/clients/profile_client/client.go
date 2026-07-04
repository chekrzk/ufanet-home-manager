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
		HouseId:   command.HouseID,
		Apartment: command.Apartment,
	})
	if err != nil {
		return domain.User{}, err
	}
	return userFromProto(resp), nil
}

func (c *Client) AddWorker(ctx context.Context, actor domain.AuthContext, command domain.AddWorker) (domain.Worker, error) {
	c.log.Debug().Str("user_id", actor.UserID).Msg("call profile grpc add worker")
	resp, err := c.client.AddWorker(ctx, &profilev1.AddWorkerRequest{
		Actor:          userContext(actor),
		UserId:         command.UserID,
		FullName:       command.FullName,
		Specialization: command.Specialization,
		Phone:          command.Phone,
		HouseId:        command.HouseID,
	})
	if err != nil {
		return domain.Worker{}, err
	}
	return workerFromProto(resp), nil
}

func (c *Client) ListWorkers(ctx context.Context, actor domain.AuthContext, houseID string) ([]domain.Worker, error) {
	c.log.Debug().Str("user_id", actor.UserID).Str("house_id", houseID).Msg("call profile grpc list workers")
	resp, err := c.client.ListWorkers(ctx, &profilev1.ListWorkersRequest{
		Actor:   userContext(actor),
		HouseId: houseID,
	})
	if err != nil {
		return nil, err
	}
	items := make([]domain.Worker, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		items = append(items, workerFromProto(item))
	}
	return items, nil
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

func workerFromProto(worker *commonv1.Worker) domain.Worker {
	if worker == nil {
		return domain.Worker{}
	}
	return domain.Worker{
		ID:             worker.GetId(),
		UserID:         worker.GetUserId(),
		FullName:       worker.GetFullName(),
		Specialization: worker.GetSpecialization(),
		Phone:          worker.GetPhone(),
		HouseID:        worker.GetHouseId(),
	}
}
