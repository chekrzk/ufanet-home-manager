package requests

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

func (s *Service) Create(ctx context.Context, actor domain.AuthContext, command domain.CreateRequest) (domain.Request, error) {
	s.log.Debug().Str("user_id", actor.UserID).Msg("create request via requests service")
	return s.client.Create(ctx, actor, command)
}

func (s *Service) List(ctx context.Context, actor domain.AuthContext, page domain.Pagination) (domain.Page[domain.Request], error) {
	s.log.Debug().Str("user_id", actor.UserID).Msg("list requests via requests service")
	return s.client.List(ctx, actor, page)
}

func (s *Service) Get(ctx context.Context, actor domain.AuthContext, requestID string) (domain.Request, error) {
	s.log.Debug().Str("user_id", actor.UserID).Str("request_id", requestID).Msg("get request via requests service")
	return s.client.Get(ctx, actor, requestID)
}

func (s *Service) UpdateStatus(ctx context.Context, actor domain.AuthContext, requestID string, command domain.UpdateRequestStatus) (domain.Request, error) {
	s.log.Debug().Str("user_id", actor.UserID).Str("request_id", requestID).Msg("update request status via requests service")
	return s.client.UpdateStatus(ctx, actor, requestID, command)
}

func (s *Service) AddComment(ctx context.Context, actor domain.AuthContext, requestID string, command domain.AddRequestComment) error {
	s.log.Debug().Str("user_id", actor.UserID).Str("request_id", requestID).Msg("add request comment via requests service")
	return s.client.AddComment(ctx, actor, requestID, command)
}
