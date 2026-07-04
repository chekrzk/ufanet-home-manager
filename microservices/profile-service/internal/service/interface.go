package service

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/profile-service/internal/models"
)

type ProfileRepository interface {
	FindProfile(ctx context.Context, userID string) (models.Profile, error)
	SaveProfile(ctx context.Context, profile *models.Profile) error
	CreateWorker(ctx context.Context, worker *models.Worker) error
	ListWorkers(ctx context.Context, houseID string) ([]models.Worker, error)
}
