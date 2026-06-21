package news

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
)

type Service interface {
	List(ctx context.Context, actor domain.AuthContext, req dto.ListNewsRequest) (dto.Page[domain.News], error)
	Create(ctx context.Context, author domain.AuthContext, req dto.CreateNewsRequest) (domain.News, error)
}
