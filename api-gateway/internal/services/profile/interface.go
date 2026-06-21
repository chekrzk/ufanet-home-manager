package profile

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/dto"
)

type Service interface {
	Me(ctx context.Context, actor domain.AuthContext) (domain.User, error)
	Update(ctx context.Context, actor domain.AuthContext, req dto.UpdateProfileRequest) (domain.User, error)
}
