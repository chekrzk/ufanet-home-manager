package repository

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/news-service/internal/models"
	"gorm.io/gorm"
)

type NewsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) Migrate(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).AutoMigrate(&models.News{})
}

func (r *NewsRepository) Create(ctx context.Context, item *models.News) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *NewsRepository) List(ctx context.Context, filter models.NewsFilter) ([]models.News, int64, error) {
	page, limit := normalizePagination(filter.Pagination)
	query := r.db.WithContext(ctx).Model(&models.News{})
	if filter.DateFrom != "" {
		query = query.Where("created_at >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("created_at <= ?", filter.DateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.News
	err := query.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&items).Error
	return items, total, err
}

func normalizePagination(p models.Pagination) (int, int) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit > 100 {
		p.Limit = 20
	}
	return p.Page, p.Limit
}
