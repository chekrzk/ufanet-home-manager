package profile

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

func (s *Service) Me(ctx context.Context, actor domain.AuthContext) (domain.User, error) {
	s.log.Debug().Str("user_id", actor.UserID).Msg("get profile via profile service")
	return s.client.Me(ctx, actor)
}

func (s *Service) Update(ctx context.Context, actor domain.AuthContext, command domain.UpdateProfile) (domain.User, error) {
	s.log.Debug().Str("user_id", actor.UserID).Msg("update profile via profile service")
	return s.client.Update(ctx, actor, command)
}

func (s *Service) AddWorker(ctx context.Context, actor domain.AuthContext, command domain.AddWorker) (domain.Worker, error) {
	s.log.Debug().Str("user_id", actor.UserID).Msg("add worker via profile service")
	return s.client.AddWorker(ctx, actor, command)
}

func (s *Service) ListWorkers(ctx context.Context, actor domain.AuthContext, houseID string) ([]domain.Worker, error) {
	s.log.Debug().Str("user_id", actor.UserID).Str("house_id", houseID).Msg("list workers via profile service")
	return s.client.ListWorkers(ctx, actor, houseID)
}

func (s *Service) SetWorkerAvailability(ctx context.Context, actor domain.AuthContext, command domain.SetWorkerAvailability) (domain.WorkerAvailability, error) {
	s.log.Debug().Str("user_id", actor.UserID).Str("house_id", command.HouseID).Msg("set worker availability via profile service")
	return s.client.SetWorkerAvailability(ctx, actor, command)
}

func (s *Service) ListWorkerAvailability(ctx context.Context, actor domain.AuthContext, filter domain.WorkerAvailabilityFilter) ([]domain.WorkerAvailability, error) {
	s.log.Debug().Str("user_id", actor.UserID).Str("house_id", filter.HouseID).Msg("list worker availability via profile service")
	return s.client.ListWorkerAvailability(ctx, actor, filter)
}
