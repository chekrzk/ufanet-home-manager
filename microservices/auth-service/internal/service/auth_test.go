package service

import (
	"context"
	"errors"
	"testing"
	"time"

	apperrors "github.com/chekrzk/ufanet-home-manager/auth-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/hasher"
	jwtmanager "github.com/chekrzk/ufanet-home-manager/auth-service/internal/jwt"
	"github.com/chekrzk/ufanet-home-manager/auth-service/internal/models"
	servicemocks "github.com/chekrzk/ufanet-home-manager/auth-service/internal/service/mocks"
	"github.com/chekrzk/ufanet-home-manager/auth-service/resources/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

func newTestService(repo UserRepository) *AuthService {
	return NewAuthService(repo, hasher.New(), jwtmanager.NewManager(config.JWTConfig{
		Secret:     "test-secret",
		Issuer:     "auth-service-test",
		AccessTTL:  time.Minute,
		RefreshTTL: time.Hour,
	}), zerolog.Nop())
}

func TestRegisterCreatesUserWithPasswordHash(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewUserRepository(t)
	svc := newTestService(repo)

	repo.On("ExistsByPhone", ctx, " 79990000001 ").Return(false, nil).Once()
	repo.On("Create", ctx, mock.MatchedBy(func(user *models.User) bool {
		user.ID = "user-1"
		return user.Phone == "79990000001" &&
			user.PasswordHash != "" &&
			user.PasswordHash != "123456" &&
			user.Role == models.RoleResident
	})).Return(nil).Once()

	user, err := svc.Register(ctx, models.RegisterCommand{Phone: " 79990000001 ", Password: "123456"})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if user.ID != "user-1" || user.Role != models.RoleResident {
		t.Fatalf("unexpected user: %+v", user)
	}
}

func TestRegisterRejectsInvalidInput(t *testing.T) {
	svc := newTestService(servicemocks.NewUserRepository(t))

	_, err := svc.Register(context.Background(), models.RegisterCommand{Phone: "", Password: "123456"})
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}

	_, err = svc.Register(context.Background(), models.RegisterCommand{Phone: "79990000001", Password: "12345"})
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument for short password, got %v", err)
	}
}

func TestRegisterRejectsExistingPhone(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewUserRepository(t)
	svc := newTestService(repo)

	repo.On("ExistsByPhone", ctx, "79990000001").Return(true, nil).Once()

	_, err := svc.Register(ctx, models.RegisterCommand{Phone: "79990000001", Password: "123456"})
	if !errors.Is(err, apperrors.ErrUserAlreadyExists) {
		t.Fatalf("expected already exists, got %v", err)
	}
}

func TestRegisterReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewUserRepository(t)
	svc := newTestService(repo)
	wantErr := errors.New("db failed")

	repo.On("ExistsByPhone", ctx, "79990000001").Return(false, wantErr).Once()

	_, err := svc.Register(ctx, models.RegisterCommand{Phone: "79990000001", Password: "123456"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

func TestLoginReturnsTokenPair(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewUserRepository(t)
	svc := newTestService(repo)
	hash, err := hasher.New().Hash("123456")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repo.On("FindByPhone", ctx, "79990000001").Return(models.User{
		ID:           "user-1",
		Phone:        "79990000001",
		PasswordHash: hash,
		Role:         models.RoleEmployee,
	}, nil).Once()

	pair, err := svc.Login(ctx, models.LoginCommand{Phone: " 79990000001 ", Password: "123456"})
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" || pair.ExpiresIn <= 0 {
		t.Fatalf("unexpected token pair: %+v", pair)
	}
}

func TestLoginRejectsBadCredentials(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewUserRepository(t)
	svc := newTestService(repo)

	repo.On("FindByPhone", ctx, "79990000001").Return(models.User{}, errors.New("not found")).Once()

	_, err := svc.Login(ctx, models.LoginCommand{Phone: "79990000001", Password: "123456"})
	if !errors.Is(err, apperrors.ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestLoginValidatesInput(t *testing.T) {
	svc := newTestService(servicemocks.NewUserRepository(t))

	_, err := svc.Login(context.Background(), models.LoginCommand{Phone: "", Password: "123456"})
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}

func TestRefreshReturnsNewTokenPair(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewUserRepository(t)
	svc := newTestService(repo)
	user := models.User{ID: "user-1", Phone: "79990000001", Role: models.RoleResident}
	pair, err := svc.tokens.NewPair(user)
	if err != nil {
		t.Fatalf("create pair: %v", err)
	}

	repo.On("FindByID", ctx, "user-1").Return(user, nil).Once()

	refreshed, err := svc.Refresh(ctx, pair.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh returned error: %v", err)
	}
	if refreshed.AccessToken == "" || refreshed.RefreshToken == "" {
		t.Fatalf("unexpected refreshed pair: %+v", refreshed)
	}
}

func TestRefreshRejectsInvalidToken(t *testing.T) {
	svc := newTestService(servicemocks.NewUserRepository(t))

	_, err := svc.Refresh(context.Background(), "invalid")
	if !errors.Is(err, apperrors.ErrUnauthorized) {
		t.Fatalf("expected unauthorized, got %v", err)
	}
}

func TestRegisterRole(t *testing.T) {
	if got := registerRole("employee"); got != models.RoleEmployee {
		t.Fatalf("registerRole(employee) = %q", got)
	}
	if got := registerRole("admin"); got != models.RoleResident {
		t.Fatalf("registerRole(admin) = %q", got)
	}
}
