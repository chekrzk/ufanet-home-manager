package requests

import (
	"context"
	"errors"
	"testing"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	servicemocks "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/requests/mocks"
	"github.com/rs/zerolog"
)

func TestCreateDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	command := domain.CreateRequest{Category: "electrician", Description: "No light"}
	want := domain.Request{ID: "request-1", UserID: "user-1", Category: "electrician"}

	client.On("Create", ctx, actor, command).Return(want, nil).Once()

	got, err := svc.Create(ctx, actor, command)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestListDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	page := domain.Pagination{Page: 1, Limit: 20}
	want := domain.Page[domain.Request]{Items: []domain.Request{{ID: "request-1"}}, Page: 1, Limit: 20, Total: 1}

	client.On("List", ctx, actor, page).Return(want, nil).Once()

	got, err := svc.List(ctx, actor, page)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if got.Total != want.Total || got.Items[0].ID != "request-1" {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestGetDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	want := domain.Request{ID: "request-1", UserID: "user-1"}

	client.On("Get", ctx, actor, "request-1").Return(want, nil).Once()

	got, err := svc.Get(ctx, actor, "request-1")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestUpdateStatusDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "manager-1", Role: "manager"}
	command := domain.UpdateRequestStatus{Status: "in_progress", AssignedTo: "worker-1"}
	want := domain.Request{ID: "request-1", Status: "in_progress"}

	client.On("UpdateStatus", ctx, actor, "request-1", command).Return(want, nil).Once()

	got, err := svc.UpdateStatus(ctx, actor, "request-1", command)
	if err != nil {
		t.Fatalf("UpdateStatus returned error: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestAddCommentReturnsClientError(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	command := domain.AddRequestComment{Text: "Need help"}
	wantErr := errors.New("request not found")

	client.On("AddComment", ctx, actor, "request-1", command).Return(wantErr).Once()

	err := svc.AddComment(ctx, actor, "request-1", command)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}
