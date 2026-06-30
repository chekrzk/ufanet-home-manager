package service

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/notification-service/internal/models"
)

type DeviceRepository interface {
	SaveDevice(ctx context.Context, device *models.Device) error
	DeleteDevice(ctx context.Context, userID string, token string) error
	CreateNotification(ctx context.Context, notification *models.Notification) error
	ListNotifications(ctx context.Context, command models.ListNotificationsCommand) ([]models.Notification, int64, error)
	MarkRead(ctx context.Context, userID string, notificationID string) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event models.PublishNotificationCommand) error
}
