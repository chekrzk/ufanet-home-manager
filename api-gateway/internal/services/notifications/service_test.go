package notifications

import (
	"context"
	"errors"
	"testing"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	servicemocks "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/notifications/mocks"
	"github.com/rs/zerolog"
)

func TestRegisterDeviceDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	device := domain.RegisterDevice{Token: "token-1", Platform: "web"}

	client.On("RegisterDevice", ctx, actor, device).Return(nil).Once()

	if err := svc.RegisterDevice(ctx, actor, device); err != nil {
		t.Fatalf("RegisterDevice returned error: %v", err)
	}
}

func TestUnregisterDeviceDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	device := domain.UnregisterDevice{Token: "token-1"}

	client.On("UnregisterDevice", ctx, actor, device).Return(nil).Once()

	if err := svc.UnregisterDevice(ctx, actor, device); err != nil {
		t.Fatalf("UnregisterDevice returned error: %v", err)
	}
}

func TestListDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	page := domain.Pagination{Page: 1, Limit: 20}
	want := domain.Page[domain.Notification]{Items: []domain.Notification{{ID: "n-1"}}, Page: 1, Limit: 20, Total: 1}

	client.On("List", ctx, actor, page).Return(want, nil).Once()

	got, err := svc.List(ctx, actor, page)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if got.Total != want.Total || got.Items[0].ID != "n-1" {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestMarkReadReturnsClientError(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	wantErr := errors.New("not found")

	client.On("MarkRead", ctx, actor, "notification-1").Return(wantErr).Once()

	err := svc.MarkRead(ctx, actor, "notification-1")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}
