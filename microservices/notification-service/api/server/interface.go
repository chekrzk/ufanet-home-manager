package server

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/notification-service/internal/models"
)

type NotificationService interface {
	RegisterDevice(ctx context.Context, cmd models.RegisterDeviceCommand) error
	UnregisterDevice(ctx context.Context, cmd models.UnregisterDeviceCommand) error
	Publish(ctx context.Context, cmd models.PublishNotificationCommand) error
}
