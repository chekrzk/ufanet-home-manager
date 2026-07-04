package notification_client

import (
	"context"

	notificationsv1 "github.com/chekrzk/ufanet-home-manager/contracts/gen/go/notifications/v1"
	"github.com/chekrzk/ufanet-home-manager/requests-service/internal/models"
	"google.golang.org/grpc"
)

type Client struct {
	client notificationsv1.NotificationsServiceClient
}

func New(conn grpc.ClientConnInterface) *Client {
	return &Client{client: notificationsv1.NewNotificationsServiceClient(conn)}
}

func (c *Client) Publish(ctx context.Context, event models.NotificationEvent) error {
	_, err := c.client.Publish(ctx, &notificationsv1.PublishNotificationRequest{
		UserId:   event.UserID,
		Type:     event.Type,
		Title:    event.Title,
		Body:     event.Body,
		EntityId: event.EntityID,
	})
	return err
}
