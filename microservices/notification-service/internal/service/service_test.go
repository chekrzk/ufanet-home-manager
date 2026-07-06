package service

import (
	"context"
	"errors"
	"testing"

	apperrors "github.com/chekrzk/ufanet-home-manager/notification-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/notification-service/internal/models"
	servicemocks "github.com/chekrzk/ufanet-home-manager/notification-service/internal/service/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

func TestRegisterDeviceSavesTrimmedDevice(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewDeviceRepository(t)
	svc := New(repo, servicemocks.NewEventPublisher(t), zerolog.Nop())

	repo.On("SaveDevice", ctx, mock.MatchedBy(func(device *models.Device) bool {
		return device.UserID == "user-1" &&
			device.Token == "token-1" &&
			device.Platform == "web"
	})).Return(nil).Once()

	err := svc.RegisterDevice(ctx, models.RegisterDeviceCommand{
		User:     models.UserContext{UserID: "user-1"},
		Token:    " token-1 ",
		Platform: " web ",
	})
	if err != nil {
		t.Fatalf("RegisterDevice returned error: %v", err)
	}
}

func TestRegisterDeviceValidatesInput(t *testing.T) {
	svc := New(servicemocks.NewDeviceRepository(t), servicemocks.NewEventPublisher(t), zerolog.Nop())

	err := svc.RegisterDevice(context.Background(), models.RegisterDeviceCommand{User: models.UserContext{UserID: "user-1"}})
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}

func TestUnregisterDeviceDeletesToken(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewDeviceRepository(t)
	svc := New(repo, servicemocks.NewEventPublisher(t), zerolog.Nop())

	repo.On("DeleteDevice", ctx, "user-1", "token-1").Return(nil).Once()

	err := svc.UnregisterDevice(ctx, models.UnregisterDeviceCommand{
		User:  models.UserContext{UserID: "user-1"},
		Token: " token-1 ",
	})
	if err != nil {
		t.Fatalf("UnregisterDevice returned error: %v", err)
	}
}

func TestPublishCreatesNotificationAndPublishesEvent(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewDeviceRepository(t)
	publisher := servicemocks.NewEventPublisher(t)
	svc := New(repo, publisher, zerolog.Nop())

	cmd := models.PublishNotificationCommand{
		UserID:   " user-1 ",
		HouseID:  " house-1 ",
		Type:     " request.created ",
		Title:    " Title ",
		Body:     " Body ",
		EntityID: " entity-1 ",
	}
	repo.On("CreateNotification", ctx, mock.MatchedBy(func(notification *models.Notification) bool {
		return notification.UserID != nil &&
			*notification.UserID == "user-1" &&
			notification.HouseID == "house-1" &&
			notification.Type == "request.created" &&
			notification.Title == "Title" &&
			notification.Body == "Body" &&
			notification.EntityID == "entity-1"
	})).Return(nil).Once()
	publisher.On("Publish", ctx, cmd).Return(nil).Once()

	if err := svc.Publish(ctx, cmd); err != nil {
		t.Fatalf("Publish returned error: %v", err)
	}
}

func TestPublishAllowsHouseNotificationWithoutUser(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewDeviceRepository(t)
	publisher := servicemocks.NewEventPublisher(t)
	svc := New(repo, publisher, zerolog.Nop())
	cmd := models.PublishNotificationCommand{HouseID: "house-1", Type: "news.created", Title: "Title"}

	repo.On("CreateNotification", ctx, mock.MatchedBy(func(notification *models.Notification) bool {
		return notification.UserID == nil && notification.HouseID == "house-1"
	})).Return(nil).Once()
	publisher.On("Publish", ctx, cmd).Return(nil).Once()

	if err := svc.Publish(ctx, cmd); err != nil {
		t.Fatalf("Publish returned error: %v", err)
	}
}

func TestPublishValidatesRequiredFields(t *testing.T) {
	svc := New(servicemocks.NewDeviceRepository(t), servicemocks.NewEventPublisher(t), zerolog.Nop())

	err := svc.Publish(context.Background(), models.PublishNotificationCommand{Type: "", Title: "Title"})
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}

func TestListNormalizesPagination(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewDeviceRepository(t)
	svc := New(repo, servicemocks.NewEventPublisher(t), zerolog.Nop())
	cmd := models.ListNotificationsCommand{User: models.UserContext{UserID: "user-1"}}

	repo.On("ListNotifications", ctx, cmd).Return([]models.Notification{{ID: "n-1"}}, int64(1), nil).Once()

	page, err := svc.List(ctx, cmd)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if page.Page != 1 || page.Limit != 20 || page.Total != 1 || len(page.Items) != 1 {
		t.Fatalf("unexpected page: %+v", page)
	}
}

func TestListValidatesUser(t *testing.T) {
	svc := New(servicemocks.NewDeviceRepository(t), servicemocks.NewEventPublisher(t), zerolog.Nop())

	_, err := svc.List(context.Background(), models.ListNotificationsCommand{})
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}

func TestMarkReadDelegatesToRepository(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewDeviceRepository(t)
	svc := New(repo, servicemocks.NewEventPublisher(t), zerolog.Nop())

	repo.On("MarkRead", ctx, "user-1", "notification-1").Return(nil).Once()

	if err := svc.MarkRead(ctx, models.UserContext{UserID: "user-1"}, "notification-1"); err != nil {
		t.Fatalf("MarkRead returned error: %v", err)
	}
}

func TestMarkReadValidatesInput(t *testing.T) {
	svc := New(servicemocks.NewDeviceRepository(t), servicemocks.NewEventPublisher(t), zerolog.Nop())

	err := svc.MarkRead(context.Background(), models.UserContext{UserID: ""}, "notification-1")
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}
