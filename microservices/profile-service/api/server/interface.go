package server

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/profile-service/internal/models"
)

type ProfileService interface {
	Me(ctx context.Context, actor models.UserContext) (models.Profile, error)
	Update(ctx context.Context, cmd models.UpdateProfileCommand) (models.Profile, error)
	AddWorker(ctx context.Context, cmd models.AddWorkerCommand) (models.Worker, error)
	ListWorkers(ctx context.Context, filter models.ListWorkersFilter) ([]models.Worker, error)
	SetWorkerAvailability(ctx context.Context, cmd models.SetWorkerAvailabilityCommand) (models.WorkerAvailability, error)
	ListWorkerAvailability(ctx context.Context, filter models.ListWorkerAvailabilityFilter) ([]models.WorkerAvailability, error)
}
