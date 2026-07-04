package service

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/requests-service/internal/models"
)

type RequestRepository interface {
	Create(ctx context.Context, request *models.MaintenanceRequest) error
	List(ctx context.Context, filter models.ListRequestsFilter) ([]models.MaintenanceRequest, int64, error)
	FindByID(ctx context.Context, id string) (models.MaintenanceRequest, error)
	Save(ctx context.Context, request *models.MaintenanceRequest) error
	AddComment(ctx context.Context, comment *models.RequestComment) error
}

type NotificationPublisher interface {
	Publish(ctx context.Context, event models.NotificationEvent) error
}
