package requests

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
)

type Service interface {
	Create(ctx context.Context, actor domain.AuthContext, req dto.CreateRequestRequest) (domain.Request, error)
	List(ctx context.Context, actor domain.AuthContext, page dto.Pagination) (dto.Page[domain.Request], error)
	Get(ctx context.Context, actor domain.AuthContext, requestID string) (domain.Request, error)
	UpdateStatus(ctx context.Context, actor domain.AuthContext, requestID string, req dto.UpdateRequestStatusRequest) (domain.Request, error)
	AddComment(ctx context.Context, actor domain.AuthContext, requestID string, req dto.AddRequestCommentRequest) error
}
