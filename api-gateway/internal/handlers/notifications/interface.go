package notifications

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
)

type Service interface {
	RegisterDevice(ctx context.Context, actor domain.AuthContext, device domain.RegisterDevice) error
	UnregisterDevice(ctx context.Context, actor domain.AuthContext, device domain.UnregisterDevice) error
	List(ctx context.Context, actor domain.AuthContext, page domain.Pagination) (domain.Page[domain.Notification], error)
	MarkRead(ctx context.Context, actor domain.AuthContext, notificationID string) error
}
