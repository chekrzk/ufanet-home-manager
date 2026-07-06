package service

import (
	"context"
	"strings"

	apperrors "github.com/chekrzk/ufanet-home-manager/news-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/news-service/internal/models"
	"github.com/rs/zerolog"
)

type Service struct {
	repo      NewsRepository
	publisher NotificationPublisher
	log       zerolog.Logger
}

// New принимает repo и publisher через интерфейсы, чтобы новости сохранялись
// отдельно от способа уведомления жителей.
func New(repo NewsRepository, publisher NotificationPublisher, log zerolog.Logger) *Service {
	return &Service{repo: repo, publisher: publisher, log: log}
}

// List возвращает страницу новостей с нормализованной пагинацией, чтобы клиенты
// не могли случайно запросить неограниченную ленту.
func (s *Service) List(ctx context.Context, filter models.NewsFilter) (models.NewsPage, error) {
	items, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return models.NewsPage{}, err
	}
	page, limit := filter.Pagination.Page, filter.Pagination.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return models.NewsPage{Items: items, Page: page, Limit: limit, Total: total}, nil
}

// Create проверяет административный сценарий публикации и после сохранения
// инициирует уведомление, чтобы лента и события не расходились.
func (s *Service) Create(ctx context.Context, cmd models.CreateNewsCommand) (models.News, error) {
	if !canManageNews(cmd.Author.Role) {
		return models.News{}, apperrors.ErrForbidden
	}
	if strings.TrimSpace(cmd.Title) == "" || strings.TrimSpace(cmd.Body) == "" {
		return models.News{}, apperrors.ErrInvalidArgument
	}
	item := models.News{
		Title:    strings.TrimSpace(cmd.Title),
		Body:     strings.TrimSpace(cmd.Body),
		HouseID:  strings.TrimSpace(cmd.HouseID),
		AuthorID: cmd.Author.UserID,
	}
	if err := s.repo.Create(ctx, &item); err != nil {
		return models.News{}, err
	}
	if s.publisher != nil {
		err := s.publisher.Publish(ctx, models.NotificationEvent{
			HouseID:  item.HouseID,
			Type:     "news.created",
			Title:    item.Title,
			Body:     item.Body,
			EntityID: item.ID,
		})
		if err != nil {
			s.log.Error().Err(err).Str("news_id", item.ID).Msg("failed to publish news notification")
		}
	}
	return item, nil
}

// canManageNews фиксирует роли, которым доверено публиковать новости.
func canManageNews(role string) bool {
	return role == "admin" || role == "manager"
}
