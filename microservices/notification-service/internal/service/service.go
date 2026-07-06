package service

import (
	"context"
	"strings"

	apperrors "github.com/chekrzk/ufanet-home-manager/notification-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/notification-service/internal/models"
	"github.com/rs/zerolog"
)

type Service struct {
	devices   DeviceRepository
	publisher EventPublisher
	log       zerolog.Logger
}

func New(devices DeviceRepository, publisher EventPublisher, log zerolog.Logger) *Service {
	return &Service{devices: devices, publisher: publisher, log: log}
}

func (s *Service) RegisterDevice(ctx context.Context, cmd models.RegisterDeviceCommand) error {
	if strings.TrimSpace(cmd.User.UserID) == "" || strings.TrimSpace(cmd.Token) == "" || strings.TrimSpace(cmd.Platform) == "" {
		return apperrors.ErrInvalidArgument
	}
	return s.devices.SaveDevice(ctx, &models.Device{
		UserID:   cmd.User.UserID,
		Token:    strings.TrimSpace(cmd.Token),
		Platform: strings.TrimSpace(cmd.Platform),
	})
}

func (s *Service) UnregisterDevice(ctx context.Context, cmd models.UnregisterDeviceCommand) error {
	if strings.TrimSpace(cmd.User.UserID) == "" || strings.TrimSpace(cmd.Token) == "" {
		return apperrors.ErrInvalidArgument
	}
	return s.devices.DeleteDevice(ctx, cmd.User.UserID, strings.TrimSpace(cmd.Token))
}

func (s *Service) Publish(ctx context.Context, cmd models.PublishNotificationCommand) error {
	if strings.TrimSpace(cmd.Type) == "" || strings.TrimSpace(cmd.Title) == "" {
		return apperrors.ErrInvalidArgument
	}
	userID := strings.TrimSpace(cmd.UserID)
	notification := models.Notification{
		HouseID:  strings.TrimSpace(cmd.HouseID),
		Type:     strings.TrimSpace(cmd.Type),
		Title:    strings.TrimSpace(cmd.Title),
		Body:     strings.TrimSpace(cmd.Body),
		EntityID: strings.TrimSpace(cmd.EntityID),
	}
	if userID != "" {
		notification.UserID = &userID
	}
	if err := s.devices.CreateNotification(ctx, &notification); err != nil {
		return err
	}
	if err := s.publisher.Publish(ctx, cmd); err != nil {
		return err
	}
	s.log.Info().Str("type", cmd.Type).Str("entity_id", cmd.EntityID).Msg("notification event published")
	return nil
}

func (s *Service) List(ctx context.Context, cmd models.ListNotificationsCommand) (models.NotificationsPage, error) {
	if strings.TrimSpace(cmd.User.UserID) == "" {
		return models.NotificationsPage{}, apperrors.ErrInvalidArgument
	}
	items, total, err := s.devices.ListNotifications(ctx, cmd)
	if err != nil {
		return models.NotificationsPage{}, err
	}
	page, limit := cmd.Pagination.Page, cmd.Pagination.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return models.NotificationsPage{Items: items, Page: page, Limit: limit, Total: total}, nil
}

func (s *Service) MarkRead(ctx context.Context, user models.UserContext, notificationID string) error {
	if strings.TrimSpace(user.UserID) == "" || strings.TrimSpace(notificationID) == "" {
		return apperrors.ErrInvalidArgument
	}
	return s.devices.MarkRead(ctx, user.UserID, notificationID)
}
