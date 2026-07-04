package service

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/news-service/internal/models"
)

type NewsRepository interface {
	Create(ctx context.Context, item *models.News) error
	List(ctx context.Context, filter models.NewsFilter) ([]models.News, int64, error)
}

type NotificationPublisher interface {
	Publish(ctx context.Context, event models.NotificationEvent) error
}
