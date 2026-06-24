package news

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

func (s *Service) List(ctx context.Context, actor domain.AuthContext, filter domain.NewsFilter) (domain.Page[domain.News], error) {
	s.log.Debug().Str("user_id", actor.UserID).Msg("list news via news service")
	return s.client.List(ctx, actor, filter)
}

func (s *Service) Create(ctx context.Context, author domain.AuthContext, command domain.CreateNews) (domain.News, error) {
	s.log.Debug().Str("user_id", author.UserID).Msg("create news via news service")
	return s.client.Create(ctx, author, command)
}
