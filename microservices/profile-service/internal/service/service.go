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
	houseID := strings.TrimSpace(cmd.HouseID)
	if err := s.ensureManagerHouse(ctx, cmd.Actor, houseID); err != nil {
		return models.Worker{}, err
	}
	worker := models.Worker{
		UserID:         strings.TrimSpace(cmd.UserID),
		FullName:       strings.TrimSpace(cmd.FullName),
		Specialization: strings.TrimSpace(cmd.Specialization),
		Phone:          strings.TrimSpace(cmd.Phone),
		HouseID:        houseID,
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
	houseID := strings.TrimSpace(filter.HouseID)
	if err := s.ensureManagerHouse(ctx, filter.Actor, houseID); err != nil {
		return nil, err
	}
	return s.repo.ListWorkers(ctx, houseID)
}

func (s *Service) SetWorkerAvailability(ctx context.Context, cmd models.SetWorkerAvailabilityCommand) (models.WorkerAvailability, error) {
	if cmd.Worker.Role != "employee" {
		return models.WorkerAvailability{}, apperrors.ErrForbidden
	}
	if strings.TrimSpace(cmd.Worker.UserID) == "" || strings.TrimSpace(cmd.Specialization) == "" || strings.TrimSpace(cmd.AvailableDate) == "" || strings.TrimSpace(cmd.AvailableTime) == "" {
		return models.WorkerAvailability{}, apperrors.ErrInvalidArgument
	}
	availability := models.WorkerAvailability{
		UserID:         cmd.Worker.UserID,
		Specialization: strings.TrimSpace(cmd.Specialization),
		HouseID:        strings.TrimSpace(cmd.HouseID),
		AvailableDate:  strings.TrimSpace(cmd.AvailableDate),
		AvailableTime:  strings.TrimSpace(cmd.AvailableTime),
	}
	if err := s.repo.SaveWorkerAvailability(ctx, &availability); err != nil {
		return models.WorkerAvailability{}, err
	}
	return availability, nil
}

func (s *Service) ListWorkerAvailability(ctx context.Context, filter models.ListWorkerAvailabilityFilter) ([]models.WorkerAvailability, error) {
	if filter.Actor.Role == "resident" && strings.TrimSpace(filter.Specialization) == "" {
		return nil, apperrors.ErrInvalidArgument
	}
	if err := s.ensureManagerHouse(ctx, filter.Actor, strings.TrimSpace(filter.HouseID)); err != nil {
		return nil, err
	}
	return s.repo.ListWorkerAvailability(ctx, filter)
}

func canManageWorkers(role string) bool {
	return role == "admin" || role == "manager"
}

func (s *Service) ensureManagerHouse(ctx context.Context, actor models.UserContext, houseID string) error {
	if actor.Role != "manager" {
		return nil
	}
	if strings.TrimSpace(actor.UserID) == "" || strings.TrimSpace(houseID) == "" {
		return apperrors.ErrForbidden
	}
	profile, err := s.repo.FindProfile(ctx, actor.UserID)
	if err != nil {
		return err
	}
	if strings.TrimSpace(profile.HouseID) != houseID {
		return apperrors.ErrForbidden
	}
	return nil
}
