package service

import (
	"context"
	"errors"
	"testing"

	apperrors "github.com/chekrzk/ufanet-home-manager/profile-service/internal/errors"
	"github.com/chekrzk/ufanet-home-manager/profile-service/internal/models"
	servicemocks "github.com/chekrzk/ufanet-home-manager/profile-service/internal/service/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

func TestMeReturnsExistingProfile(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewProfileRepository(t)
	svc := New(repo, zerolog.Nop())
	want := models.Profile{UserID: "user-1", FullName: "Ivan", HouseID: "house-1"}

	repo.On("FindProfile", ctx, "user-1").Return(want, nil).Once()

	got, err := svc.Me(ctx, models.UserContext{UserID: "user-1"})
	if err != nil {
		t.Fatalf("Me returned error: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestMeReturnsEmptyProfileWhenNotFound(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewProfileRepository(t)
	svc := New(repo, zerolog.Nop())

	repo.On("FindProfile", ctx, "user-1").Return(models.Profile{}, apperrors.ErrNotFound).Once()

	got, err := svc.Me(ctx, models.UserContext{UserID: "user-1"})
	if err != nil {
		t.Fatalf("Me returned error: %v", err)
	}
	if got.UserID != "user-1" {
		t.Fatalf("unexpected profile: %+v", got)
	}
}

func TestMeValidatesActor(t *testing.T) {
	svc := New(servicemocks.NewProfileRepository(t), zerolog.Nop())

	_, err := svc.Me(context.Background(), models.UserContext{})
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}

func TestUpdateSavesTrimmedProfile(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewProfileRepository(t)
	svc := New(repo, zerolog.Nop())

	repo.On("SaveProfile", ctx, mock.MatchedBy(func(profile *models.Profile) bool {
		return profile.UserID == "user-1" &&
			profile.FullName == "Ivan Ivanov" &&
			profile.HouseID == "house-1" &&
			profile.Apartment == "12"
	})).Return(nil).Once()

	profile, err := svc.Update(ctx, models.UpdateProfileCommand{
		Actor:     models.UserContext{UserID: "user-1"},
		FullName:  " Ivan Ivanov ",
		HouseID:   " house-1 ",
		Apartment: " 12 ",
	})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if profile.FullName != "Ivan Ivanov" || profile.HouseID != "house-1" {
		t.Fatalf("unexpected profile: %+v", profile)
	}
}

func TestUpdateValidatesRequiredFields(t *testing.T) {
	svc := New(servicemocks.NewProfileRepository(t), zerolog.Nop())

	_, err := svc.Update(context.Background(), models.UpdateProfileCommand{Actor: models.UserContext{UserID: "user-1"}})
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}

func TestAddWorkerAllowsManagerForOwnHouse(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewProfileRepository(t)
	svc := New(repo, zerolog.Nop())

	repo.On("FindProfile", ctx, "manager-1").Return(models.Profile{UserID: "manager-1", HouseID: "house-1"}, nil).Once()
	repo.On("CreateWorker", ctx, mock.MatchedBy(func(worker *models.Worker) bool {
		worker.ID = "worker-1"
		return worker.UserID == "employee-1" &&
			worker.FullName == "Worker" &&
			worker.Specialization == "electrician" &&
			worker.HouseID == "house-1"
	})).Return(nil).Once()

	worker, err := svc.AddWorker(ctx, models.AddWorkerCommand{
		Actor:          models.UserContext{UserID: "manager-1", Role: "manager"},
		UserID:         " employee-1 ",
		FullName:       " Worker ",
		Specialization: " electrician ",
		HouseID:        " house-1 ",
	})
	if err != nil {
		t.Fatalf("AddWorker returned error: %v", err)
	}
	if worker.ID != "worker-1" {
		t.Fatalf("unexpected worker: %+v", worker)
	}
}

func TestAddWorkerRejectsResident(t *testing.T) {
	svc := New(servicemocks.NewProfileRepository(t), zerolog.Nop())

	_, err := svc.AddWorker(context.Background(), models.AddWorkerCommand{
		Actor:          models.UserContext{UserID: "resident-1", Role: "resident"},
		FullName:       "Worker",
		Specialization: "electrician",
	})
	if !errors.Is(err, apperrors.ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestAddWorkerRejectsManagerFromOtherHouse(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewProfileRepository(t)
	svc := New(repo, zerolog.Nop())

	repo.On("FindProfile", ctx, "manager-1").Return(models.Profile{UserID: "manager-1", HouseID: "house-2"}, nil).Once()

	_, err := svc.AddWorker(ctx, models.AddWorkerCommand{
		Actor:          models.UserContext{UserID: "manager-1", Role: "manager"},
		FullName:       "Worker",
		Specialization: "electrician",
		HouseID:        "house-1",
	})
	if !errors.Is(err, apperrors.ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestListWorkersDelegatesToRepository(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewProfileRepository(t)
	svc := New(repo, zerolog.Nop())
	want := []models.Worker{{ID: "worker-1"}}

	repo.On("ListWorkers", ctx, "house-1").Return(want, nil).Once()

	got, err := svc.ListWorkers(ctx, models.ListWorkersFilter{Actor: models.UserContext{Role: "admin"}, HouseID: " house-1 "})
	if err != nil {
		t.Fatalf("ListWorkers returned error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "worker-1" {
		t.Fatalf("unexpected workers: %+v", got)
	}
}

func TestSetWorkerAvailabilityRequiresEmployee(t *testing.T) {
	svc := New(servicemocks.NewProfileRepository(t), zerolog.Nop())

	_, err := svc.SetWorkerAvailability(context.Background(), models.SetWorkerAvailabilityCommand{
		Worker: models.UserContext{UserID: "user-1", Role: "resident"},
	})
	if !errors.Is(err, apperrors.ErrForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestSetWorkerAvailabilitySavesTrimmedAvailability(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewProfileRepository(t)
	svc := New(repo, zerolog.Nop())

	repo.On("SaveWorkerAvailability", ctx, mock.MatchedBy(func(item *models.WorkerAvailability) bool {
		return item.UserID == "employee-1" &&
			item.Specialization == "electrician" &&
			item.HouseID == "house-1" &&
			item.AvailableDate == "2026-07-13" &&
			item.AvailableTime == "12:00"
	})).Return(nil).Once()

	got, err := svc.SetWorkerAvailability(ctx, models.SetWorkerAvailabilityCommand{
		Worker:         models.UserContext{UserID: "employee-1", Role: "employee"},
		Specialization: " electrician ",
		HouseID:        " house-1 ",
		AvailableDate:  " 2026-07-13 ",
		AvailableTime:  " 12:00 ",
	})
	if err != nil {
		t.Fatalf("SetWorkerAvailability returned error: %v", err)
	}
	if got.Specialization != "electrician" {
		t.Fatalf("unexpected availability: %+v", got)
	}
}

func TestListWorkerAvailabilityRequiresSpecializationForResident(t *testing.T) {
	svc := New(servicemocks.NewProfileRepository(t), zerolog.Nop())

	_, err := svc.ListWorkerAvailability(context.Background(), models.ListWorkerAvailabilityFilter{
		Actor: models.UserContext{UserID: "resident-1", Role: "resident"},
	})
	if !errors.Is(err, apperrors.ErrInvalidArgument) {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}

func TestListWorkerAvailabilityDelegatesToRepository(t *testing.T) {
	ctx := context.Background()
	repo := servicemocks.NewProfileRepository(t)
	svc := New(repo, zerolog.Nop())
	filter := models.ListWorkerAvailabilityFilter{
		Actor:          models.UserContext{Role: "resident"},
		Specialization: "electrician",
		HouseID:        "house-1",
	}
	want := []models.WorkerAvailability{{ID: "availability-1"}}

	repo.On("ListWorkerAvailability", ctx, filter).Return(want, nil).Once()

	got, err := svc.ListWorkerAvailability(ctx, filter)
	if err != nil {
		t.Fatalf("ListWorkerAvailability returned error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "availability-1" {
		t.Fatalf("unexpected availability: %+v", got)
	}
}

func TestCanManageWorkers(t *testing.T) {
	cases := map[string]bool{
		"admin":    true,
		"manager":  true,
		"resident": false,
		"employee": false,
	}
	for role, want := range cases {
		if got := canManageWorkers(role); got != want {
			t.Fatalf("canManageWorkers(%q) = %v, want %v", role, got, want)
		}
	}
}
