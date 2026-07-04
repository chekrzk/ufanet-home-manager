package server

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/requests-service/internal/models"
)

type RequestsService interface {
	Create(ctx context.Context, cmd models.CreateRequestCommand) (models.MaintenanceRequest, error)
	List(ctx context.Context, filter models.ListRequestsFilter) (models.RequestsPage, error)
	Get(ctx context.Context, cmd models.GetRequestCommand) (models.MaintenanceRequest, error)
	UpdateStatus(ctx context.Context, cmd models.UpdateRequestStatusCommand) (models.MaintenanceRequest, error)
	AddComment(ctx context.Context, cmd models.AddRequestCommentCommand) error
}
