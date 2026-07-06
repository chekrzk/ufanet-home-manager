package notifications

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/rs/zerolog"
)

type Service struct {
	client Client
	log    zerolog.Logger
}

func New(client Client, log zerolog.Logger) *Service {
	return &Service{client: client, log: log}
}

func (s *Service) RegisterDevice(ctx context.Context, actor domain.AuthContext, device domain.RegisterDevice) error {
	s.log.Debug().Str("user_id", actor.UserID).Msg("register device via notifications service")
	return s.client.RegisterDevice(ctx, actor, device)
}

func (s *Service) UnregisterDevice(ctx context.Context, actor domain.AuthContext, device domain.UnregisterDevice) error {
	s.log.Debug().Str("user_id", actor.UserID).Msg("unregister device via notifications service")
	return s.client.UnregisterDevice(ctx, actor, device)
}

func (s *Service) List(ctx context.Context, actor domain.AuthContext, page domain.Pagination) (domain.Page[domain.Notification], error) {
	s.log.Debug().Str("user_id", actor.UserID).Msg("list notifications via notifications service")
	return s.client.List(ctx, actor, page)
}

func (s *Service) MarkRead(ctx context.Context, actor domain.AuthContext, notificationID string) error {
	s.log.Debug().Str("user_id", actor.UserID).Str("notification_id", notificationID).Msg("mark notification read via notifications service")
	return s.client.MarkRead(ctx, actor, notificationID)
}
