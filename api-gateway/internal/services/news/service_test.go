package news

import (
	"context"
	"errors"
	"testing"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	servicemocks "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/news/mocks"
	"github.com/rs/zerolog"
)

func TestListDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	filter := domain.NewsFilter{Pagination: domain.Pagination{Page: 2, Limit: 10}}
	want := domain.Page[domain.News]{Items: []domain.News{{ID: "news-1"}}, Page: 2, Limit: 10, Total: 1}

	client.On("List", ctx, actor, filter).Return(want, nil).Once()

	got, err := svc.List(ctx, actor, filter)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if got.Total != want.Total || got.Items[0].ID != "news-1" {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestCreateDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "manager-1", Role: "manager"}
	command := domain.CreateNews{Title: "Title", Body: "Body", HouseID: "house-1"}
	want := domain.News{ID: "news-1", Title: command.Title}

	client.On("Create", ctx, actor, command).Return(want, nil).Once()

	got, err := svc.Create(ctx, actor, command)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestCreateReturnsClientError(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "resident-1", Role: "resident"}
	command := domain.CreateNews{Title: "Title", Body: "Body"}
	wantErr := errors.New("forbidden")

	client.On("Create", ctx, actor, command).Return(domain.News{}, wantErr).Once()

	_, err := svc.Create(ctx, actor, command)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}
