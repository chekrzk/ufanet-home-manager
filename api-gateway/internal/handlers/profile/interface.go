package profile

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
)

type Service interface {
	Me(ctx context.Context, actor domain.AuthContext) (domain.User, error)
	Update(ctx context.Context, actor domain.AuthContext, command domain.UpdateProfile) (domain.User, error)
	AddWorker(ctx context.Context, actor domain.AuthContext, command domain.AddWorker) (domain.Worker, error)
	ListWorkers(ctx context.Context, actor domain.AuthContext, houseID string) ([]domain.Worker, error)
	SetWorkerAvailability(ctx context.Context, actor domain.AuthContext, command domain.SetWorkerAvailability) (domain.WorkerAvailability, error)
	ListWorkerAvailability(ctx context.Context, actor domain.AuthContext, filter domain.WorkerAvailabilityFilter) ([]domain.WorkerAvailability, error)
}
