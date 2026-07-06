package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	servicemocks "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/auth/mocks"
	"github.com/rs/zerolog"
)

func TestLoginDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	credentials := domain.LoginCredentials{Phone: "79990000001", Password: "123456"}
	want := domain.AuthTokens{AccessToken: "access", RefreshToken: "refresh", ExpiresIn: 60}

	client.On("Login", ctx, credentials).Return(want, nil).Once()

	got, err := svc.Login(ctx, credentials)
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestRegisterDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	user := domain.RegisterUser{Phone: "79990000001", Password: "123456", Role: "resident"}
	want := domain.User{ID: "user-1", Phone: user.Phone, Role: "resident"}

	client.On("Register", ctx, user).Return(want, nil).Once()

	got, err := svc.Register(ctx, user)
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestRefreshReturnsClientError(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	wantErr := errors.New("auth unavailable")

	client.On("Refresh", ctx, "refresh-token").Return(domain.AuthTokens{}, wantErr).Once()

	_, err := svc.Refresh(ctx, "refresh-token")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}
