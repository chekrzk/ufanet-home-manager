package repository

import (
	"context"

	apperrors "github.com/chekrzk/ufanet-home-manager/profile-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/profile-service/internal/models"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Migrate(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).AutoMigrate(&models.Profile{}, &models.Worker{})
}

func (r *ProfileRepository) FindProfile(ctx context.Context, userID string) (models.Profile, error) {
	var profile models.Profile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	if err == nil {
		return profile, nil
	}
	if err == gorm.ErrRecordNotFound {
		return models.Profile{}, apperrors.ErrNotFound
	}
	return models.Profile{}, err
}

func (r *ProfileRepository) SaveProfile(ctx context.Context, profile *models.Profile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

func (r *ProfileRepository) CreateWorker(ctx context.Context, worker *models.Worker) error {
	return r.db.WithContext(ctx).Create(worker).Error
}

func (r *ProfileRepository) ListWorkers(ctx context.Context, houseID string) ([]models.Worker, error) {
	var workers []models.Worker
	query := r.db.WithContext(ctx).Order("created_at DESC")
	if houseID != "" {
		query = query.Where("house_id = ?", houseID)
	}
	return workers, query.Find(&workers).Error
}
