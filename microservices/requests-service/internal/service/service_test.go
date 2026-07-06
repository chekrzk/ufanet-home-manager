package service

import (
	"context"
	"errors"
	"testing"

	apperrors "github.com/chekrzk/ufanet-home-manager/requests-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/requests-service/internal/models"
	servicemocks "github.com/chekrzk/ufanet-home-manager/requests-service/internal/service/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

func TestCreateSavesRequestAndPublishesToUser(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewRequestRepository(t)
	publisher := servicemocks.NewNotificationPublisher(t)
	svc := New(repo, publisher, zerolog.Nop())

	repo.On("Create", ctx, mock.MatchedBy(func(request *models.MaintenanceRequest) bool {
		request.ID = "request-1"
		return request.UserID == "user-1" &&
			request.Category == "electrician" &&
			request.Description == "No light" &&
			request.Status == models.RequestStatusNew &&
			request.AssignedTo == nil &&
			request.PreferredDate == "2026-07-13" &&
			request.Address == "Prospekt Oktyabrya, 107"
	})).Return(nil).Once()
	publisher.On("Publish", ctx, mock.MatchedBy(func(event models.NotificationEvent) bool {
		return event.UserID == "user-1" &&
			event.Type == "request.created" &&
			event.EntityID == "request-1"
	})).Return(nil).Once()

	request, err := svc.Create(ctx, models.CreateRequestCommand{
		User:          models.UserContext{UserID: "user-1", Role: "resident"},
		Category:      " electrician ",
		Description:   " No light ",
		PreferredDate: " 2026-07-13 ",
		Address:       " Prospekt Oktyabrya, 107 ",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if request.ID != "request-1" || request.AssignedTo != nil {
		t.Fatalf("unexpected request: %+v", request)
	}
}

func TestCreatePublishesAssignmentForWorker(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewRequestRepository(t)
	publisher := servicemocks.NewNotificationPublisher(t)
	svc := New(repo, publisher, zerolog.Nop())

	repo.On("Create", ctx, mock.AnythingOfType("*models.MaintenanceRequest")).Run(func(args mock.Arguments) {
		args.Get(1).(*models.MaintenanceRequest).ID = "request-1"
	}).Return(nil).Once()
	publisher.On("Publish", ctx, mock.MatchedBy(func(event models.NotificationEvent) bool {
		return event.UserID == "user-1" && event.Type == "request.created"
	})).Return(nil).Once()
	publisher.On("Publish", ctx, mock.MatchedBy(func(event models.NotificationEvent) bool {
		return event.UserID == "worker-1" && event.Type == "request.assigned"
	})).Return(nil).Once()

	request, err := svc.Create(ctx, models.CreateRequestCommand{
		User:             models.UserContext{UserID: "user-1"},
		Category:         "electrician",
		Description:      "No light",
		AssignedWorkerID: " worker-1 ",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if request.AssignedTo == nil || *request.AssignedTo != "worker-1" {
		t.Fatalf("unexpected assignment: %+v", request.AssignedTo)
	}
}

func TestCreateValidatesInput(t *testing.T) {
	svc := New(servicemocks.NewRequestRepository(t), nil, zerolog.Nop())

	_, err := svc.Create(context.Background(), models.CreateRequestCommand{User: models.UserContext{UserID: "user-1"}})
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}

func TestListNormalizesPagination(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewRequestRepository(t)
	svc := New(repo, nil, zerolog.Nop())
	filter := models.ListRequestsFilter{Pagination: models.Pagination{Page: -1, Limit: 500}}

	repo.On("List", ctx, filter).Return([]models.MaintenanceRequest{{ID: "request-1"}}, int64(1), nil).Once()

	page, err := svc.List(ctx, filter)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if page.Page != 1 || page.Limit != 20 || page.Total != 1 || len(page.Items) != 1 {
		t.Fatalf("unexpected page: %+v", page)
	}
}

func TestGetRejectsResidentAccessToOtherRequest(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewRequestRepository(t)
	svc := New(repo, nil, zerolog.Nop())

	repo.On("FindByID", ctx, "request-1").Return(models.MaintenanceRequest{ID: "request-1", UserID: "owner-1"}, nil).Once()

	_, err := svc.Get(ctx, models.GetRequestCommand{
		Actor:     models.UserContext{UserID: "resident-1", Role: "resident"},
		RequestID: "request-1",
	})
	if !errors.Is(err, apperrors.ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestGetAllowsOwner(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewRequestRepository(t)
	svc := New(repo, nil, zerolog.Nop())
	want := models.MaintenanceRequest{ID: "request-1", UserID: "owner-1"}

	repo.On("FindByID", ctx, "request-1").Return(want, nil).Once()

	got, err := svc.Get(ctx, models.GetRequestCommand{
		Actor:     models.UserContext{UserID: "owner-1", Role: "resident"},
		RequestID: "request-1",
	})
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if got.ID != want.ID {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestUpdateStatusSavesTimestampsAndPublishes(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewRequestRepository(t)
	publisher := servicemocks.NewNotificationPublisher(t)
	svc := New(repo, publisher, zerolog.Nop())
	current := models.MaintenanceRequest{ID: "request-1", UserID: "user-1", Status: models.RequestStatusNew}

	repo.On("FindByID", ctx, "request-1").Return(current, nil).Once()
	repo.On("Save", ctx, mock.MatchedBy(func(request *models.MaintenanceRequest) bool {
		return request.ID == "request-1" &&
			request.Status == models.RequestStatusInProgress &&
			request.AcceptedAt != nil &&
			request.AssignedTo != nil &&
			*request.AssignedTo == "worker-1"
	})).Return(nil).Once()
	publisher.On("Publish", ctx, mock.MatchedBy(func(event models.NotificationEvent) bool {
		return event.UserID == "user-1" &&
			event.Type == "request.status_changed" &&
			event.EntityID == "request-1"
	})).Return(nil).Once()

	got, err := svc.UpdateStatus(ctx, models.UpdateRequestStatusCommand{
		Actor:      models.UserContext{UserID: "manager-1", Role: "manager"},
		RequestID:  "request-1",
		Status:     models.RequestStatusInProgress,
		AssignedTo: " worker-1 ",
	})
	if err != nil {
		t.Fatalf("UpdateStatus returned error: %v", err)
	}
	if got.AcceptedAt == nil {
		t.Fatalf("accepted timestamp was not set")
	}
}

func TestUpdateStatusRejectsInvalidRoleOrStatus(t *testing.T) {
	svc := New(servicemocks.NewRequestRepository(t), nil, zerolog.Nop())

	_, err := svc.UpdateStatus(context.Background(), models.UpdateRequestStatusCommand{
		Actor:  models.UserContext{Role: "resident"},
		Status: models.RequestStatusDone,
	})
	if !errors.Is(err, apperrors.ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestUpdateStatusRejectsUnknownStatus(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewRequestRepository(t)
	svc := New(repo, nil, zerolog.Nop())

	repo.On("FindByID", ctx, "request-1").Return(models.MaintenanceRequest{ID: "request-1", UserID: "user-1"}, nil).Once()

	_, err := svc.UpdateStatus(ctx, models.UpdateRequestStatusCommand{
		Actor:     models.UserContext{Role: "admin"},
		RequestID: "request-1",
		Status:    "unknown",
	})
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}

func TestAddCommentSavesTrimmedComment(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewRequestRepository(t)
	svc := New(repo, nil, zerolog.Nop())

	repo.On("FindByID", ctx, "request-1").Return(models.MaintenanceRequest{ID: "request-1", UserID: "user-1"}, nil).Once()
	repo.On("AddComment", ctx, mock.MatchedBy(func(comment *models.RequestComment) bool {
		return comment.RequestID == "request-1" &&
			comment.UserID == "user-1" &&
			comment.Text == "Text"
	})).Return(nil).Once()

	err := svc.AddComment(ctx, models.AddRequestCommentCommand{
		Actor:     models.UserContext{UserID: "user-1", Role: "resident"},
		RequestID: "request-1",
		Text:      " Text ",
	})
	if err != nil {
		t.Fatalf("AddComment returned error: %v", err)
	}
}

func TestAddCommentRejectsOtherResident(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewRequestRepository(t)
	svc := New(repo, nil, zerolog.Nop())

	repo.On("FindByID", ctx, "request-1").Return(models.MaintenanceRequest{ID: "request-1", UserID: "owner-1"}, nil).Once()

	err := svc.AddComment(ctx, models.AddRequestCommentCommand{
		Actor:     models.UserContext{UserID: "resident-1", Role: "resident"},
		RequestID: "request-1",
		Text:      "Text",
	})
	if !errors.Is(err, apperrors.ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestAddCommentValidatesInput(t *testing.T) {
	svc := New(servicemocks.NewRequestRepository(t), nil, zerolog.Nop())

	err := svc.AddComment(context.Background(), models.AddRequestCommentCommand{Actor: models.UserContext{UserID: "user-1"}})
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}

func TestHelpers(t *testing.T) {
	if !canManageRequests("admin") || !canManageRequests("manager") || !canManageRequests("employee") {
		t.Fatalf("expected admin, manager and employee to manage requests")
	}
	if canManageRequests("resident") {
		t.Fatalf("resident must not manage requests")
	}
	if !validStatus(models.RequestStatusNew) || !validStatus(models.RequestStatusInProgress) || !validStatus(models.RequestStatusDone) || !validStatus(models.RequestStatusCanceled) {
		t.Fatalf("known request statuses must be valid")
	}
	if validStatus("unknown") {
		t.Fatalf("unknown request status must be invalid")
	}

	request := models.MaintenanceRequest{}
	assignRequest(&request, " worker-1 ")
	if request.AssignedTo == nil || *request.AssignedTo != "worker-1" {
		t.Fatalf("assignment was not set: %+v", request.AssignedTo)
	}
	assignRequest(&request, " ")
	if request.AssignedTo != nil {
		t.Fatalf("assignment was not cleared: %+v", request.AssignedTo)
	}
}
