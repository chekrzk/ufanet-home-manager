package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	servicemocks "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/services/profile/mocks"
	"github.com/rs/zerolog"
)

func TestMeDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	want := domain.User{ID: "user-1", FullName: "Ivan"}

	client.On("Me", ctx, actor).Return(want, nil).Once()

	got, err := svc.Me(ctx, actor)
	if err != nil {
		t.Fatalf("Me returned error: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestUpdateDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	command := domain.UpdateProfile{FullName: "Ivan", HouseID: "house-1"}
	want := domain.User{ID: "user-1", FullName: "Ivan", HouseID: "house-1"}

	client.On("Update", ctx, actor, command).Return(want, nil).Once()

	got, err := svc.Update(ctx, actor, command)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestAddWorkerDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "manager-1", Role: "manager"}
	command := domain.AddWorker{FullName: "Worker", Specialization: "electrician"}
	want := domain.Worker{ID: "worker-1", FullName: "Worker"}

	client.On("AddWorker", ctx, actor, command).Return(want, nil).Once()

	got, err := svc.AddWorker(ctx, actor, command)
	if err != nil {
		t.Fatalf("AddWorker returned error: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestListWorkersDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "manager-1", Role: "manager"}
	want := []domain.Worker{{ID: "worker-1"}}

	client.On("ListWorkers", ctx, actor, "house-1").Return(want, nil).Once()

	got, err := svc.ListWorkers(ctx, actor, "house-1")
	if err != nil {
		t.Fatalf("ListWorkers returned error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "worker-1" {
		t.Fatalf("unexpected workers: %+v", got)
	}
}

func TestSetWorkerAvailabilityDelegatesToClient(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "employee-1", Role: "employee"}
	command := domain.SetWorkerAvailability{Specialization: "electrician", HouseID: "house-1", AvailableDate: "2026-07-13"}
	want := domain.WorkerAvailability{ID: "availability-1", UserID: "employee-1"}

	client.On("SetWorkerAvailability", ctx, actor, command).Return(want, nil).Once()

	got, err := svc.SetWorkerAvailability(ctx, actor, command)
	if err != nil {
		t.Fatalf("SetWorkerAvailability returned error: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestListWorkerAvailabilityReturnsClientError(t *testing.T) {
	ctx := context.Background()
	client := servicemocks.NewClient(t)
	svc := New(client, zerolog.Nop())
	actor := domain.AuthContext{UserID: "user-1", Role: "resident"}
	filter := domain.WorkerAvailabilityFilter{HouseID: "house-1", Specialization: "electrician"}
	wantErr := errors.New("profile unavailable")

	client.On("ListWorkerAvailability", ctx, actor, filter).Return([]domain.WorkerAvailability(nil), wantErr).Once()

	_, err := svc.ListWorkerAvailability(ctx, actor, filter)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}
