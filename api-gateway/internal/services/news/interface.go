package news

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
)

type Service interface {
	List(ctx context.Context, actor domain.AuthContext, filter domain.NewsFilter) (domain.Page[domain.News], error)
	Create(ctx context.Context, author domain.AuthContext, command domain.CreateNews) (domain.News, error)
}

type Client interface {
	List(ctx context.Context, actor domain.AuthContext, filter domain.NewsFilter) (domain.Page[domain.News], error)
	Create(ctx context.Context, author domain.AuthContext, command domain.CreateNews) (domain.News, error)
}
