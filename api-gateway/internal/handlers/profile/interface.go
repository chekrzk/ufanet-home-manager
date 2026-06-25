package profile

import (
	"context"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
)

type Service interface {
	Me(ctx context.Context, actor domain.AuthContext) (domain.User, error)
	Update(ctx context.Context, actor domain.AuthContext, command domain.UpdateProfile) (domain.User, error)
}
