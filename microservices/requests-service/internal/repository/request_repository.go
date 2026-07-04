package repository

import (
	"context"

	apperrors "github.com/chekrzk/ufanet-home-manager/requests-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/requests-service/internal/models"
	"gorm.io/gorm"
)

type RequestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) Migrate(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).AutoMigrate(&models.MaintenanceRequest{}, &models.RequestComment{})
}

func (r *RequestRepository) Create(ctx context.Context, request *models.MaintenanceRequest) error {
	return r.db.WithContext(ctx).Create(request).Error
}

func (r *RequestRepository) List(ctx context.Context, filter models.ListRequestsFilter) ([]models.MaintenanceRequest, int64, error) {
	page, limit := normalizePagination(filter.Pagination)
	query := r.db.WithContext(ctx).Model(&models.MaintenanceRequest{})
	if filter.Actor.Role == "resident" {
		query = query.Where("user_id = ?", filter.Actor.UserID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.MaintenanceRequest
	err := query.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *RequestRepository) FindByID(ctx context.Context, id string) (models.MaintenanceRequest, error) {
	var request models.MaintenanceRequest
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&request).Error
	if err == nil {
		return request, nil
	}
	if err == gorm.ErrRecordNotFound {
		return models.MaintenanceRequest{}, apperrors.ErrNotFound
	}
	return models.MaintenanceRequest{}, err
}

func (r *RequestRepository) Save(ctx context.Context, request *models.MaintenanceRequest) error {
	return r.db.WithContext(ctx).Save(request).Error
}

func (r *RequestRepository) AddComment(ctx context.Context, comment *models.RequestComment) error {
	return r.db.WithContext(ctx).Create(comment).Error
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
