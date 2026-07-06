package profile

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gwerrors "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/errors"
	handlermocks "github.com/chekrzk/ufanet-home-manager/api-gateway/internal/handlers/profile/mocks"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/constant"
	"github.com/chekrzk/ufanet-home-manager/api-gateway/internal/models/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

func newTestApp(handler *Handler) *fiber.App {
	app := fiber.New(fiber.Config{ErrorHandler: func(c *fiber.Ctx, err error) error {
		return gwerrors.Fail(c, err)
	}})
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(constant.CtxUserID, "user-1")
		c.Locals(constant.CtxRole, "manager")
		return c.Next()
	})
	app.Get("/profile", handler.Me)
	app.Put("/profile", handler.Update)
	app.Post("/profile/workers", handler.AddWorker)
	app.Get("/profile/workers", handler.ListWorkers)
	app.Post("/profile/workers/availability", handler.SetWorkerAvailability)
	app.Get("/profile/workers/availability", handler.ListWorkerAvailability)
	return app
}

func jsonRequest(method, path, body string) *http.Request {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestMe(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "manager"}

	service.On("Me", mock.Anything, actor).Return(domain.User{ID: "user-1"}, nil).Once()

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/profile", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
}

func TestUpdate(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "manager"}
	command := domain.UpdateProfile{FullName: "Ivan", HouseID: "house-1", Apartment: "12"}

	service.On("Update", mock.Anything, actor, command).Return(domain.User{ID: "user-1", FullName: "Ivan"}, nil).Once()

	resp, err := app.Test(jsonRequest(http.MethodPut, "/profile", `{"full_name":"Ivan","house_id":"house-1","apartment":"12"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
}

func TestAddWorker(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "manager"}
	command := domain.AddWorker{UserID: "employee-1", FullName: "Worker", Specialization: "electrician", Phone: "79990000001", HouseID: "house-1"}

	service.On("AddWorker", mock.Anything, actor, command).Return(domain.Worker{ID: "worker-1"}, nil).Once()

	resp, err := app.Test(jsonRequest(http.MethodPost, "/profile/workers", `{"user_id":"employee-1","full_name":"Worker","specialization":"electrician","phone":"79990000001","house_id":"house-1"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusCreated)
	}
}

func TestListWorkers(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "manager"}

	service.On("ListWorkers", mock.Anything, actor, "house-1").Return([]domain.Worker{{ID: "worker-1"}}, nil).Once()

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/profile/workers?house_id=house-1", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
}

func TestSetWorkerAvailability(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "manager"}
	command := domain.SetWorkerAvailability{Specialization: "electrician", HouseID: "house-1", AvailableDate: "2026-07-13", AvailableTime: "12:00"}

	service.On("SetWorkerAvailability", mock.Anything, actor, command).Return(domain.WorkerAvailability{ID: "availability-1"}, nil).Once()

	resp, err := app.Test(jsonRequest(http.MethodPost, "/profile/workers/availability", `{"specialization":"electrician","house_id":"house-1","available_date":"2026-07-13","available_time":"12:00"}`))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusCreated)
	}
}

func TestListWorkerAvailability(t *testing.T) {
	service := handlermocks.NewService(t)
	app := newTestApp(NewHandler(service))
	actor := domain.AuthContext{UserID: "user-1", Role: "manager"}
	filter := domain.WorkerAvailabilityFilter{HouseID: "house-1", Specialization: "electrician", AvailableDate: "2026-07-13"}

	service.On("ListWorkerAvailability", mock.Anything, actor, filter).Return([]domain.WorkerAvailability{{ID: "availability-1"}}, nil).Once()

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/profile/workers/availability?house_id=house-1&specialization=electrician&available_date=2026-07-13", nil))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, fiber.StatusOK)
	}
}
