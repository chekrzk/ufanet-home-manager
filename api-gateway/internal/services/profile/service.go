package profile

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

func (s service) Me(ctx context.Context, actor domain.AuthContext) (domain.User, error) {
	s.log.Debug().Str("user_id", actor.UserID).Msg("get profile via profile service")
	return s.client.Me(ctx, actor)
}

func (s service) Update(ctx context.Context, actor domain.AuthContext, command domain.UpdateProfile) (domain.User, error) {
	s.log.Debug().Str("user_id", actor.UserID).Msg("update profile via profile service")
	return s.client.Update(ctx, actor, command)
}
