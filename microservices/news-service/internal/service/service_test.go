package service

import (
	"context"
	"errors"
	"testing"

	apperrors "github.com/chekrzk/ufanet-home-manager/news-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/news-service/internal/models"
	servicemocks "github.com/chekrzk/ufanet-home-manager/news-service/internal/service/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

func TestListNormalizesPagination(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewNewsRepository(t)
	svc := New(repo, nil, zerolog.Nop())

	repo.On("List", ctx, models.NewsFilter{}).Return([]models.News{{ID: "news-1"}}, int64(1), nil).Once()

	page, err := svc.List(ctx, models.NewsFilter{})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if page.Page != 1 || page.Limit != 20 || page.Total != 1 || len(page.Items) != 1 {
		t.Fatalf("unexpected page: %+v", page)
	}
}

func TestListReturnsRepositoryError(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewNewsRepository(t)
	svc := New(repo, nil, zerolog.Nop())
	wantErr := errors.New("db failed")

	repo.On("List", ctx, models.NewsFilter{}).Return([]models.News(nil), int64(0), wantErr).Once()

	_, err := svc.List(ctx, models.NewsFilter{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

func TestCreateRejectsResident(t *testing.T) {
	svc := New(servicemocks.NewNewsRepository(t), nil, zerolog.Nop())

	_, err := svc.Create(context.Background(), models.CreateNewsCommand{
		Author: models.UserContext{Role: "resident"},
		Title:  "Title",
		Body:   "Body",
	})
	if !errors.Is(err, apperrors.ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestCreateValidatesRequiredFields(t *testing.T) {
	svc := New(servicemocks.NewNewsRepository(t), nil, zerolog.Nop())

	_, err := svc.Create(context.Background(), models.CreateNewsCommand{
		Author: models.UserContext{Role: "manager"},
		Title:  " ",
		Body:   "Body",
	})
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}

func TestCreateSavesNewsAndPublishesNotification(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewNewsRepository(t)
	publisher := servicemocks.NewNotificationPublisher(t)
	svc := New(repo, publisher, zerolog.Nop())

	repo.On("Create", ctx, mock.MatchedBy(func(item *models.News) bool {
		item.ID = "news-1"
		return item.Title == "Title" &&
			item.Body == "Body" &&
			item.HouseID == "house-1" &&
			item.AuthorID == "manager-1"
	})).Return(nil).Once()
	publisher.On("Publish", ctx, mock.MatchedBy(func(event models.NotificationEvent) bool {
		return event.HouseID == "house-1" &&
			event.Type == "news.created" &&
			event.Title == "Title" &&
			event.Body == "Body" &&
			event.EntityID == "news-1"
	})).Return(nil).Once()

	item, err := svc.Create(ctx, models.CreateNewsCommand{
		Author:  models.UserContext{UserID: "manager-1", Role: "manager"},
		Title:   " Title ",
		Body:    " Body ",
		HouseID: " house-1 ",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if item.ID != "news-1" || item.Title != "Title" || item.Body != "Body" || item.HouseID != "house-1" {
		t.Fatalf("unexpected news: %+v", item)
	}
}

func TestCreateIgnoresPublisherError(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewNewsRepository(t)
	publisher := servicemocks.NewNotificationPublisher(t)
	svc := New(repo, publisher, zerolog.Nop())

	repo.On("Create", ctx, mock.AnythingOfType("*models.News")).Return(nil).Once()
	publisher.On("Publish", ctx, mock.AnythingOfType("models.NotificationEvent")).Return(errors.New("redis failed")).Once()

	_, err := svc.Create(ctx, models.CreateNewsCommand{
		Author: models.UserContext{UserID: "admin-1", Role: "admin"},
		Title:  "Title",
		Body:   "Body",
	})
	if err != nil {
		t.Fatalf("Create should ignore publisher error, got %v", err)
	}
}

func TestCanManageNews(t *testing.T) {
	cases := map[string]bool{
		"admin":    true,
		"manager":  true,
		"resident": false,
		"employee": false,
		"":         false,
	}
	for role, want := range cases {
		if got := canManageNews(role); got != want {
			t.Fatalf("canManageNews(%q) = %v, want %v", role, got, want)
		}
	}
}
