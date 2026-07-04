package service

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/notification-service/internal/models"
)

type DeviceRepository interface {
	SaveDevice(ctx context.Context, device *models.Device) error
	DeleteDevice(ctx context.Context, userID string, token string) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event models.PublishNotificationCommand) error
}
