package notifications_client

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	commonv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/common/v1"
	notificationsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/notifications/v1"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Client struct {
	client notificationsv1.NotificationsServiceClient
	log    zerolog.Logger
}

func New(conn *grpc.ClientConn, log zerolog.Logger) *Client {
	return &Client{client: notificationsv1.NewNotificationsServiceClient(conn), log: log}
}

func (c *Client) RegisterDevice(ctx context.Context, actor domain.AuthContext, device domain.RegisterDevice) error {
	c.log.Debug().Str("user_id", actor.UserID).Msg("call notifications grpc register device")
	_, err := c.client.RegisterDevice(ctx, &notificationsv1.RegisterDeviceRequest{
		User:     userContext(actor),
		Token:    device.Token,
		Platform: device.Platform,
	})
	return err
}

func (c *Client) UnregisterDevice(ctx context.Context, actor domain.AuthContext, device domain.UnregisterDevice) error {
	c.log.Debug().Str("user_id", actor.UserID).Msg("call notifications grpc unregister device")
	_, err := c.client.UnregisterDevice(ctx, &notificationsv1.UnregisterDeviceRequest{
		User:  userContext(actor),
		Token: device.Token,
	})
	return err
}

func userContext(actor domain.AuthContext) *commonv1.UserContext {
	return &commonv1.UserContext{UserId: actor.UserID, Role: actor.Role}
}
