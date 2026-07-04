package repository

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/notification-service/internal/models"
	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Migrate(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).AutoMigrate(&models.Device{})
}

func (r *DeviceRepository) SaveDevice(ctx context.Context, device *models.Device) error {
	return r.db.WithContext(ctx).Where("token = ?", device.Token).Assign(device).FirstOrCreate(device).Error
}

func (r *DeviceRepository) DeleteDevice(ctx context.Context, userID string, token string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND token = ?", userID, token).Delete(&models.Device{}).Error
}
