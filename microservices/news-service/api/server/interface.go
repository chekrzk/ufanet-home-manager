package server

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/news-service/internal/models"
)

type NewsService interface {
	List(ctx context.Context, filter models.NewsFilter) (models.NewsPage, error)
	Create(ctx context.Context, cmd models.CreateNewsCommand) (models.News, error)
}
