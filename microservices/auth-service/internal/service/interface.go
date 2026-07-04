package service

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByPhone(ctx context.Context, phone string) (models.User, error)
	FindByID(ctx context.Context, id string) (models.User, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
}
