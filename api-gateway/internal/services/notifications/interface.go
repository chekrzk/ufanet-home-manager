package notifications

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
)

type Service interface {
	RegisterDevice(ctx context.Context, actor domain.AuthContext, req dto.RegisterDeviceRequest) error
	UnregisterDevice(ctx context.Context, actor domain.AuthContext, req dto.UnregisterDeviceRequest) error
}
