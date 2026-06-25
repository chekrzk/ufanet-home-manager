package repository

import (
	"context"

	apperrors "github.com/chekrzk/ufanet-home-manager/auth-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Migrate(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).AutoMigrate(&models.User{})
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err == nil {
		return user, nil
	}
	if err == gorm.ErrRecordNotFound {
		return models.User{}, apperrors.ErrUserNotFound
	}
	return models.User{}, err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err == nil {
		return user, nil
	}
	if err == gorm.ErrRecordNotFound {
		return models.User{}, apperrors.ErrUserNotFound
	}
	return models.User{}, err
}

func (r *UserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("phone = ?", phone).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
