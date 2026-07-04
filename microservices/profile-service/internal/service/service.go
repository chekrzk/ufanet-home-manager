package service

import (
	"context"
	stderrors "errors"
	"strings"

	apperrors "github.com/chekrzk/ufanet-home-manager/profile-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/profile-service/internal/models"
	"github.com/rs/zerolog"
)

type Service struct {
	repo ProfileRepository
	log  zerolog.Logger
}

func New(repo ProfileRepository, log zerolog.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) Me(ctx context.Context, actor models.UserContext) (models.Profile, error) {
	if strings.TrimSpace(actor.UserID) == "" {
		return models.Profile{}, apperrors.ErrInvalidArgument
	}
	profile, err := s.repo.FindProfile(ctx, actor.UserID)
	if stderrors.Is(err, apperrors.ErrNotFound) {
		return models.Profile{UserID: actor.UserID}, nil
	}
	return profile, err
}

func (s *Service) Update(ctx context.Context, cmd models.UpdateProfileCommand) (models.Profile, error) {
	if strings.TrimSpace(cmd.Actor.UserID) == "" || strings.TrimSpace(cmd.FullName) == "" {
		return models.Profile{}, apperrors.ErrInvalidArgument
	}
	profile := models.Profile{
		UserID:    cmd.Actor.UserID,
		FullName:  strings.TrimSpace(cmd.FullName),
		HouseID:   strings.TrimSpace(cmd.HouseID),
		Apartment: strings.TrimSpace(cmd.Apartment),
	}
	if err := s.repo.SaveProfile(ctx, &profile); err != nil {
		return models.Profile{}, err
	}
	return profile, nil
}

func (s *Service) AddWorker(ctx context.Context, cmd models.AddWorkerCommand) (models.Worker, error) {
	if !canManageWorkers(cmd.Actor.Role) {
		return models.Worker{}, apperrors.ErrForbidden
	}
	if strings.TrimSpace(cmd.FullName) == "" || strings.TrimSpace(cmd.Specialization) == "" {
		return models.Worker{}, apperrors.ErrInvalidArgument
	}
	worker := models.Worker{
		UserID:         strings.TrimSpace(cmd.UserID),
		FullName:       strings.TrimSpace(cmd.FullName),
		Specialization: strings.TrimSpace(cmd.Specialization),
		Phone:          strings.TrimSpace(cmd.Phone),
		HouseID:        strings.TrimSpace(cmd.HouseID),
	}
	if err := s.repo.CreateWorker(ctx, &worker); err != nil {
		return models.Worker{}, err
	}
	return worker, nil
}

func (s *Service) ListWorkers(ctx context.Context, filter models.ListWorkersFilter) ([]models.Worker, error) {
	if !canManageWorkers(filter.Actor.Role) {
		return nil, apperrors.ErrForbidden
	}
	return s.repo.ListWorkers(ctx, strings.TrimSpace(filter.HouseID))
}

func canManageWorkers(role string) bool {
	return role == "admin" || role == "manager"
}
