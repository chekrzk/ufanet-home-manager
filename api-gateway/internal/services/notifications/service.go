package notifications

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/rs/zerolog"
)

type service struct {
	client Client
	log    zerolog.Logger
}

func New(client Client, log zerolog.Logger) Service {
	return service{client: client, log: log}
}

func (s service) RegisterDevice(ctx context.Context, actor domain.AuthContext, device domain.RegisterDevice) error {
	s.log.Debug().Str("user_id", actor.UserID).Msg("register device via notifications service")
	return s.client.RegisterDevice(ctx, actor, device)
}

func (s service) UnregisterDevice(ctx context.Context, actor domain.AuthContext, device domain.UnregisterDevice) error {
	s.log.Debug().Str("user_id", actor.UserID).Msg("unregister device via notifications service")
	return s.client.UnregisterDevice(ctx, actor, device)
}
