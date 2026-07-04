package profile

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
)

type Client interface {
	Me(ctx context.Context, actor domain.AuthContext) (domain.User, error)
	Update(ctx context.Context, actor domain.AuthContext, command domain.UpdateProfile) (domain.User, error)
	AddWorker(ctx context.Context, actor domain.AuthContext, command domain.AddWorker) (domain.Worker, error)
	ListWorkers(ctx context.Context, actor domain.AuthContext, houseID string) ([]domain.Worker, error)
}
