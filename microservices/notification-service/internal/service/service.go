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
	if err := s.publisher.Publish(ctx, cmd); err != nil {
		return err
	}
	s.log.Info().Str("type", cmd.Type).Str("entity_id", cmd.EntityID).Msg("notification event published")
	return nil
}
