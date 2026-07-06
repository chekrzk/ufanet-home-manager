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
	return r.db.WithContext(ctx).AutoMigrate(&models.Device{}, &models.Notification{})
}

func (r *DeviceRepository) SaveDevice(ctx context.Context, device *models.Device) error {
	return r.db.WithContext(ctx).Where("token = ?", device.Token).Assign(device).FirstOrCreate(device).Error
}

func (r *DeviceRepository) DeleteDevice(ctx context.Context, userID string, token string) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND token = ?", userID, token).Delete(&models.Device{}).Error
}

func (r *DeviceRepository) CreateNotification(ctx context.Context, notification *models.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *DeviceRepository) ListNotifications(ctx context.Context, command models.ListNotificationsCommand) ([]models.Notification, int64, error) {
	page, limit := normalizePagination(command.Pagination)
	query := r.db.WithContext(ctx).Model(&models.Notification{}).Where("user_id = ? OR user_id IS NULL", command.User.UserID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.Notification
	err := query.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *DeviceRepository) MarkRead(ctx context.Context, userID string, notificationID string) error {
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ? AND (user_id = ? OR user_id IS NULL)", notificationID, userID).
		Update("read", true).Error
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
