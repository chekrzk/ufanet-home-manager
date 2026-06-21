package notifications

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	notificationsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/notifications/v1"
	"google.golang.org/grpc"
)

type service struct {
	client notificationsv1.NotificationsServiceClient
}

func New(conn *grpc.ClientConn) Service {
	return service{client: notificationsv1.NewNotificationsServiceClient(conn)}
}

func (s service) RegisterDevice(ctx context.Context, actor domain.AuthContext, req dto.RegisterDeviceRequest) error {
	_, err := s.client.RegisterDevice(ctx, &notificationsv1.RegisterDeviceRequest{
		User:     userContext(actor),
		Token:    req.Token,
		Platform: req.Platform,
	})
	return err
}

func (s service) UnregisterDevice(ctx context.Context, actor domain.AuthContext, req dto.UnregisterDeviceRequest) error {
	_, err := s.client.UnregisterDevice(ctx, &notificationsv1.UnregisterDeviceRequest{
		User:  userContext(actor),
		Token: req.Token,
	})
	return err
}

func userContext(actor domain.AuthContext) *commonv1.UserContext {
	return &commonv1.UserContext{UserId: actor.UserID, Role: actor.Role}
}
