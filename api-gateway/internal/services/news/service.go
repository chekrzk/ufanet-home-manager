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

// New принимает client через интерфейс, чтобы сценарии новостей можно было
// тестировать без gRPC и менять transport без переписывания сервиса.
func New(client Client, log zerolog.Logger) *Service {
	return &Service{client: client, log: log}
}

// List сохраняет проверяемый auth context рядом со сценарием получения ленты,
// чтобы downstream видел, для кого и с какими фильтрами запрошены новости.
func (s *Service) List(ctx context.Context, actor domain.AuthContext, filter domain.NewsFilter) (domain.Page[domain.News], error) {
	s.log.Debug().Str("user_id", actor.UserID).Msg("list news via news service")
	return s.client.List(ctx, actor, filter)
}

// Create передает автора вместе с командой, чтобы право на публикацию и
// аудит новости оставались частью бизнес-сценария.
func (s *Service) Create(ctx context.Context, author domain.AuthContext, command domain.CreateNews) (domain.News, error) {
	s.log.Debug().Str("user_id", author.UserID).Msg("create news via news service")
	return s.client.Create(ctx, author, command)
}
