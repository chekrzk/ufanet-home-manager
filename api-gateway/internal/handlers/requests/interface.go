package requests

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
)

type Service interface {
	Create(ctx context.Context, actor domain.AuthContext, command domain.CreateRequest) (domain.Request, error)
	List(ctx context.Context, actor domain.AuthContext, page domain.Pagination) (domain.Page[domain.Request], error)
	Get(ctx context.Context, actor domain.AuthContext, requestID string) (domain.Request, error)
	UpdateStatus(ctx context.Context, actor domain.AuthContext, requestID string, command domain.UpdateRequestStatus) (domain.Request, error)
	AddComment(ctx context.Context, actor domain.AuthContext, requestID string, command domain.AddRequestComment) error
}
